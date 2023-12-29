package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

type CommandGenerateOpenings struct {
	CommandIO

	File  *os.File `arg:"" help:"CSV file containing openings."`
	Depth int      `default:"10" help:"Max depth to generate openings for."`
}

func (cmd CommandGenerateOpenings) Run() error {
	defer cmd.File.Close()

	csv := csv.NewReader(cmd.File)

	book := OpeningBook{
		Depth: cmd.Depth,
		Moves: make(map[string][]OpeningMove),
	}

	slog.Info("reading openings")

	for n := 0; ; n++ {
		row, err := csv.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		if n == 0 {
			continue
		}

		slog.Info("storing", "eco", row[0], "name", row[1])

		book.ECOs = append(book.ECOs, row[0])
		book.Names = append(book.Names, row[1])

		eco := len(book.ECOs) - 1
		name := len(book.Names) - 1
		moves := strings.Fields(row[2])

		game := must(NewGame(BoardStartPositionFEN))

		for i := 0; i < min(len(moves), cmd.Depth); i++ {
			move := must(NewUCIMove(game.Board(), moves[i]))

			fen := game.Board().FEN()

			book.Moves[fen] = append(book.Moves[fen], OpeningMove{
				ECO:  eco,
				Name: name,
				Move: move,
			})

			game.MakeMove(move)
		}
	}

	out := json.NewEncoder(cmd.Out())
	out.SetIndent("", "  ")

	return out.Encode(book)
}
