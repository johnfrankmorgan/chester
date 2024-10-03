package main

import (
	"context"
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
}

func Search(sctx *SearchContext) {
	sctx.Start = time.Now()

	moves := GenerateMoves(sctx.Game.Board(), MoveGenerationOptions{})

	sctx.BestMove = moves[rand.Intn(len(moves))]
}
