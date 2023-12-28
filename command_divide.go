package main

import "fmt"

type CommandDivide struct {
	CommandIO

	FEN   string `arg:"" help:"FEN position"`
	Depth int    `default:"1" help:"Depth to search."`
}

func (cmd CommandDivide) Run() error {
	game, err := NewGame(cmd.FEN)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Out(), "FEN: %s\n\n", game.Board().FEN())

	nodes := cmd.divide(game, cmd.Depth)

	fmt.Fprintf(cmd.Out(), "\nNodes: %d\n", nodes)

	return nil
}

func (cmd CommandDivide) divide(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0

	for _, move := range game.Board().GenerateMoves(MoveGeneratorOptions{}) {
		game.MakeMove(move)
		n := cmd.divide(game, depth-1)
		nodes += n

		if depth == cmd.Depth {
			fmt.Fprintf(cmd.Out(), "%s: %d\n", move.UCI(), n)
		}

		game.UnmakeMove()
	}

	return nodes
}
