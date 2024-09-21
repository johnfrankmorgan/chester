package main

import (
	"context"
	"math/rand"
	"time"
)

type SearchContext struct {
	Game     *Game
	BestMove Move

	Start       time.Time
	Depth       int
	Nodes       int
	CurrentMove Move
}

func Search(ctx context.Context, sctx *SearchContext) {
	sctx.Start = time.Now()

	moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})

	sctx.BestMove = moves[rand.Intn(len(moves))]
}
