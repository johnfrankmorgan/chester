package main

import (
	"fmt"
	"io"
)

func Perft(game *Game, depth int, divide io.Writer) uint64 {
	return perft(game, depth, divide, true)
}

func perft(game *Game, depth int, divide io.Writer, root bool) uint64 {
	if depth == 0 {
		return 1
	}

	total := uint64(0)

	for _, move := range GenerateMoves(game.Board(), MoveGenerationOptions{}) {
		game.MakeMove(move)

		nodes := perft(game, depth-1, divide, false)

		if root && divide != nil {
			fmt.Fprintf(divide, "%5s: %d\n", move, nodes)
		}

		total += nodes

		game.UnmakeMove()
	}

	return total
}
