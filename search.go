package main

import (
	"context"
	"log/slog"
	"time"
)

type SearchContext struct {
	context.Context

	Game *Game
	TT   *TranspositionTable

	Best Move

	Start       time.Time
	Depth       int
	Nodes       int
	CurrentMove Move

	Extensions int
}

const (
	SearchMaxDepth      = 32
	SearchMaxExtensions = 3
)

func Search(sctx *SearchContext) {
	sctx.Start = time.Now()

	for sctx.Depth = 1; sctx.Depth <= SearchMaxDepth; sctx.Depth++ {
		start := time.Now()

		slog.Debug("starting iteration", "depth", sctx.Depth)

		if sctx.Err() != nil {
			slog.Warn("search aborted")
			break
		}

		eval := search(sctx, sctx.Depth, -EvalInf, EvalInf)
		slog.Debug("completed iteration", "depth", sctx.Depth, "eval", eval, "bestmove", sctx.Best)

		if n, ok := eval.MateIn(); ok {
			slog.Debug("mate", "in", n, "move", sctx.Best)
			break
		}

		deadline, ok := sctx.Deadline()
		if !ok {
			continue
		}

		if d := time.Since(start); time.Now().Add(d * 2).After(deadline) {
			slog.Warn("unlikely to have time for another iteration")
			break
		}
	}

	if sctx.Best.IsZero() {
		moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})
		sctx.Best = moves[0]

		slog.Warn("failed to find best move, selected first", "move", sctx.Best)
	}
}

func search(sctx *SearchContext, depth int, alpha, beta Eval) Eval {
	if sctx.Err() != nil {
		return 0
	}

	if depth != sctx.Depth {
		if sctx.Game.Board().Moves.Half >= 100 {
			return 0
		}

		for b := range sctx.Game.Boards() {
			if b != sctx.Game.Board() && b.Zobrist == sctx.Game.Board().Zobrist {
				return 0
			}
		}
	}

	if t, ok := sctx.TT.Get(sctx.Game.Board().Zobrist); ok && t.Depth >= depth {
		if sctx.Extensions == 0 && depth == sctx.Depth {
			sctx.Best = t.Best
		}

		switch t.Bound {
		case BoundBeta:
			if t.Eval >= beta {
				return beta
			}

		case BoundAlpha:
			if t.Eval <= alpha {
				return alpha
			}

		case BoundExact:
			return t.Eval
		}
	}

	if depth == 0 {
		sctx.Nodes++

		return quiesce(sctx, alpha, beta)
	}

	moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})
	trans := Transposition{
		Key:   sctx.Game.Board().Zobrist,
		Depth: depth,
		Bound: BoundAlpha,
	}

	if len(moves) == 0 {
		if sctx.Game.Board().Attacks.Checks > 0 {
			return -(EvalMate - Eval(sctx.Depth-depth))
		}

		return 0
	}

	OrderMoves(sctx.Game.Board(), sctx.Best, moves)

	for _, move := range moves {
		if depth == sctx.Depth {
			sctx.CurrentMove = move
		}

		sctx.Game.MakeMove(move)

		extension := 0
		if sctx.Extensions < SearchMaxExtensions {
			if sctx.Game.Board().Attacks.Checks > 0 {
				extension = 1
			} else if move.Flags&MoveFlagPromoteAny != 0 {
				extension = 1
			}
		}

		sctx.Extensions += extension

		eval := -search(sctx, depth-1+extension, -beta, -alpha)
		sctx.Game.UnmakeMove()

		sctx.Extensions -= extension

		if eval >= beta {
			if sctx.Err() == nil {
				sctx.TT.Store(Transposition{
					Key:   sctx.Game.Board().Zobrist,
					Eval:  beta,
					Bound: BoundBeta,
					Best:  move,
					Depth: depth,
				})
			}

			return beta
		} else if eval > alpha {
			if sctx.Err() == nil {
				trans.Best = move
				trans.Bound = BoundExact
			}

			alpha = eval
		}
	}

	trans.Eval = alpha

	if sctx.Err() == nil {
		sctx.TT.Store(trans)

		if sctx.Extensions == 0 && depth == sctx.Depth {
			sctx.Best = trans.Best
		}
	}

	return trans.Eval
}

func quiesce(sctx *SearchContext, alpha, beta Eval) Eval {
	if eval := Evaluate(sctx.Game.Board()); eval >= beta {
		return eval
	} else if eval > alpha {
		alpha = eval
	}

	moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{CapturesOnly: true})
	OrderMoves(sctx.Game.Board(), Move{}, moves)

	for _, move := range moves {
		sctx.Game.MakeMove(move)
		eval := -quiesce(sctx, -beta, -alpha)
		sctx.Game.UnmakeMove()

		if eval >= beta {
			return beta
		} else if eval > alpha {
			alpha = eval
		}
	}

	return alpha
}
