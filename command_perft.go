package main

import (
	"fmt"
	"time"
)

type CommandPerft struct {
	CommandIO

	FEN   string `arg:"" help:"FEN position"`
	Depth int    `default:"1" help:"Depth to search."`
}

func (cmd CommandPerft) Run() error {
	game, err := NewGame(cmd.FEN)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Out(), "FEN: %s\n\n", game.Board().FEN())

	start := time.Now()
	nodes := cmd.perft(game, cmd.Depth)

	fmt.Fprintf(cmd.Out(), "\nFound %d nodes in %s\n", nodes, time.Since(start))

	return nil
}

func (cmd CommandPerft) perft(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0

	for _, move := range game.Board().GenerateMoves(MoveGeneratorOptions{}) {
		game.MakeMove(move)
		n := cmd.perft(game, depth-1)
		nodes += n

		if depth == cmd.Depth {
			fmt.Fprintf(cmd.Out(), "%s: %d\n", move, n)
		}

		game.UnmakeMove()
	}

	return nodes
}
