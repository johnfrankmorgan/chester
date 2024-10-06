package main

import (
	"context"
	"log/slog"
	"math/rand"
	"time"
)

type SearchContext struct {
	context.Context

	Game     *Game
	BestMove Move

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
		slog.Debug("starting iteration", "depth", sctx.Depth)

		if sctx.Err() != nil {
			slog.Warn("search aborted")
			break
		}

		eval := search(sctx, sctx.Depth, -EvalInf, EvalInf)
		slog.Debug("completed iteration", "depth", sctx.Depth, "eval", eval, "bestmove", sctx.BestMove)
	}

	if sctx.BestMove.IsZero() {
		moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})
		sctx.BestMove = moves[rand.Intn(len(moves))]

		slog.Warn("failed to find best move, selected random move", "move", sctx.BestMove)
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

	if depth == 0 {
		sctx.Nodes++

		return quiesce(sctx, alpha, beta)
	}

	moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})

	if len(moves) == 0 {
		if sctx.Game.Board().Attacks.Checks > 0 {
			return -(EvalMate - Eval(sctx.Depth-depth))
		}

		return 0
	}

	OrderMoves(sctx.Game.Board(), sctx.BestMove, moves)

	for _, move := range moves {
		if depth == sctx.Depth {
			sctx.CurrentMove = move
		}

		sctx.Game.MakeMove(move)

		extension := 0
		if sctx.Extensions < SearchMaxExtensions {
			if sctx.Game.Board().Attacks.Checks > 0 {
				slog.Debug("extending search", "reason", "check", "depth", depth, "move", move)

				extension = 1
				sctx.Extensions++
			} else if move.Flags&MoveFlagPromoteAny != 0 {
				slog.Debug("extending search", "reason", "promotion", "depth", depth, "move", move)

				extension = 1
				sctx.Extensions++
			}
		}

		eval := -search(sctx, depth-1+extension, -beta, -alpha)
		sctx.Game.UnmakeMove()

		if eval >= beta {
			return beta
		} else if eval > alpha {
			alpha = eval

			if depth == sctx.Depth && sctx.Err() == nil {
				sctx.BestMove = move

				if n, ok := eval.MateIn(); ok {
					slog.Debug("mate", "in", n, "move", move)
				}
			}
		}
	}

	return alpha
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
