package main

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
)

type Searcher struct {
	game *Game
	opts SearchOptions

	transpositions *Transpositions
	repetitions    *Repetitions
}

type SearchOptions struct {
	Depth          int
	OpeningBook    bool
	Transpositions bool
}

const (
	_SearcherMaxExtensions      = 3
	_SearcherInitialAlpha       = -100000000
	_SearcherInitialBeta        = -_SearcherInitialAlpha
	_SearcherMateEvaluation     = 1000000
	_SearcherTranspositionsSize = 4096 * 4096 * 4
	_SearcherMaxDepth           = 16
)

var _tt *Transpositions

func NewSearcher(game *Game, opts SearchOptions) *Searcher {
	if _tt == nil {
		_tt = NewTranspositions(_SearcherTranspositionsSize, opts.Transpositions)
	}

	s := &Searcher{
		game: game,
		opts: opts,

		transpositions: _tt,
		repetitions:    NewRepetitions(),
	}

	for _, b := range s.game.Boards() {
		s.repetitions.Push(b.Zobrist(), b.Moves.Half == 1)
	}

	return s
}

func (s *Searcher) Search(ctx context.Context) (Move, int) {
	if s.opts.OpeningBook && s.game.Board().Moves.Full/2 < Precomputed.OpeningBook.Depth {
		fen := s.game.Board().FEN()

		slog.Info("trying opening book")

		if openings, ok := Precomputed.OpeningBook.Moves[fen]; ok && len(openings) > 0 {
			opening := openings[rand.Intn(len(openings))]

			slog.Info("using book move",
				"move", opening.Move,
				"eco", Precomputed.OpeningBook.ECOs[opening.ECO],
				"name", Precomputed.OpeningBook.Names[opening.Name])

			return opening.Move, 0
		}
	}

	best := struct {
		move Move
		eval int
	}{
		eval: math.MinInt,
	}

	slog.Info("running iterative deepening search", "max-depth", s.opts.Depth)

loop:
	for depth := 1; depth <= _SearcherMaxDepth; depth++ {
		slog.Info("starting search", "depth", depth)

		move := Move{}
		eval := s._search(ctx, depth, 0, _SearcherInitialAlpha, _SearcherInitialBeta, 0, &move)

		select {
		case <-ctx.Done():
			// FIXME: don't discard search results on cancellation
			slog.Info("search cancelled", "reason", ctx.Err())
			break loop

		default:
		}

		if eval > best.eval || (best.move.IsZero() && !move.IsZero()) {
			best.move = move
			best.eval = eval
		}
	}

	if best.move.IsZero() {
		best.move = s.game.Board().GenerateMoves(MoveGeneratorOptions{})[0]

		slog.Warn("failed to find best move. using first generated move", "move", best.move)
	}

	if abs(best.eval) > _SearcherMateEvaluation-_SearcherMaxDepth {
		slog.Info("mate", "moves", _SearcherMateEvaluation-abs(best.eval))
	}

	slog.Info("best move", "move", best.move, "eval", best.eval)

	return best.move, best.eval
}

func (s *Searcher) _search(ctx context.Context, depth, root, alpha, beta, extensions int, best *Move) int {
	select {
	case <-ctx.Done():
		return 0

	default:
	}

	slog.Debug("searching", "depth", depth, "root", root, "alpha", alpha, "beta", beta)

	if root > 0 {
		if s.game.Board().Moves.Half > 100 || s.repetitions.Contains(s.game.Board().Zobrist()) {
			slog.Debug("draw", "depth", depth, "root", root, "alpha", alpha, "beta", beta)
			return 0
		}

		s.repetitions.Push(s.game.Board().Zobrist(), s.game.Board().Moves.Half == 1)
		defer s.repetitions.Pop()

		alpha = max(alpha, -_SearcherMateEvaluation+root)
		beta = min(beta, _SearcherMateEvaluation-root)

		if alpha >= beta {
			return alpha
		}
	}

	if t, ok := s.transpositions.Get(s.game.Board(), depth, alpha, beta); ok {
		slog.Debug("found cached evaluation", "transposition", t)

		if root == 0 {
			slog.Info("using cached best move", "transposition", t, "move", t.Move)
			*best = t.Move
		}

		return t.Evaluation
	}

	if depth == 0 {
		return s._quiesce(ctx, alpha, beta)
	}

	hash := Move{}

	if move, ok := s.transpositions.Move(s.game.Board()); ok {
		hash = move
	}

	moves := s.game.Board().GenerateMoves(MoveGeneratorOptions{HashMove: hash})

	if len(moves) == 0 {
		if s.game.Board().Attacks.Checks.Check {
			slog.Debug("mate", "depth", depth, "root", root)

			return -_SearcherMateEvaluation + root
		}

		return 0
	}

	bound := TranspositionUpper

	position := struct {
		best Move
	}{}

	for _, move := range moves {
		s.game.MakeMove(move)

		extension := 0

		if extensions <= _SearcherMaxExtensions {
			if s.game.Board().Attacks.Checks.Check {
				slog.Debug("extending search", "depth", depth, "root", root, "alpha", alpha, "beta", beta, "extensions", extensions+1, "reason", "check")
				extension = 1
			} else if s.game.Board().Bitboards.Pieces[PiecePawn].IsSet(move.To.Bitboard()) && (move.To.Rank() == Rank2 || move.To.Rank() == Rank7) {
				slog.Debug("extending search", "depth", depth, "root", root, "alpha", alpha, "beta", beta, "extensions", extensions+1, "reason", "pre-promotion pawn")
				extension = 1
			} else if s.game.Board().Bitboards.Pieces[PiecePawn].IsSet(move.To.Bitboard()) && (move.To.Rank() == Rank1 || move.To.Rank() == Rank8) {
				slog.Debug("extending search", "depth", depth, "root", root, "alpha", alpha, "beta", beta, "extensions", extensions+1, "reason", "promoting pawn")
				extension = 1
			}
		}

		eval := -s._search(ctx, depth-1+extension, root+1, -beta, -alpha, extensions+extension, best)

		if root == 0 {
			slog.Info("evaluated move", "move", move, "depth", depth, "eval", eval)
		}

		s.game.UnmakeMove()

		select {
		case <-ctx.Done():
			return 0

		default:
		}

		if eval >= beta {
			s.transpositions.Store(s.game.Board(), depth, beta, TranspositionLower, move)
			return beta
		}

		if eval > alpha {
			bound = TranspositionExact
			alpha = eval

			if root == 0 {
				*best = move
			}

			position.best = move
		}
	}

	s.transpositions.Store(s.game.Board(), depth, alpha, bound, position.best)

	return alpha
}

func (s *Searcher) _quiesce(ctx context.Context, alpha, beta int) int {
	select {
	case <-ctx.Done():
		return 0

	default:
	}

	eval := Evaluate(s.game.Board())

	if eval >= beta {
		return beta
	}

	if eval > alpha {
		alpha = eval
	}

	for _, move := range s.game.Board().GenerateMoves(MoveGeneratorOptions{CapturesOnly: true}) {
		s.game.MakeMove(move)
		eval = -s._quiesce(ctx, -beta, -alpha)
		s.game.UnmakeMove()

		if eval >= beta {
			return beta
		}

		if eval > alpha {
			alpha = eval
		}
	}

	return alpha
}
