package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"log/slog"
)

type CommandUCI struct {
	CommandIO

	Depth          int           `default:"4" help:"Depth to search to."`
	OpeningBook    bool          `default:"true" negatable:"true" help:"Enable opening book."`
	Transpositions bool          `default:"true" negatable:"true" help:"Enable transposition cache."`
	ThinkTime      time.Duration `default:"600ms" help:"Thinking time."`

	sc        *bufio.Scanner
	game      *Game
	searching bool
	cancel    func()
}

var (
	ErrUCI     = errors.New("uci")
	ErrUCIQuit = fmt.Errorf("%w: quit", ErrUCI)
)

func (cmd *CommandUCI) Run() error {
	for {
		if err := cmd.run(); err != nil {
			if errors.Is(err, ErrUCIQuit) {
				return nil
			}

			slog.Error("uci failure", "error", err)

			return err
		}
	}
}

func (cmd *CommandUCI) recv() ([]string, error) {
	if cmd.sc == nil {
		cmd.sc = bufio.NewScanner(cmd.In())
	}

	if !cmd.sc.Scan() {
		return nil, fmt.Errorf("%w: failed to read command", ErrUCI)
	}

	command := cmd.sc.Text()

	if err := cmd.sc.Err(); err != nil {
		return nil, fmt.Errorf("%w: failed to read command: %w", ErrUCI, err)
	}

	return strings.Fields(strings.TrimSpace(command)), nil
}

func (cmd *CommandUCI) send(format string, args ...any) error {
	response := fmt.Sprintf(strings.TrimSpace(format), args...)

	slog.Info("sending response", "response", response)

	if _, err := fmt.Fprintln(cmd.Out(), response); err != nil {
		return fmt.Errorf("%w: failed to send response: %w", ErrUCI, err)
	}

	return nil
}

func (cmd *CommandUCI) run() error {
	command, err := cmd.recv()
	if err != nil {
		return err
	}

	slog.Info("received command", "command", command)

	if len(command) == 0 {
		slog.Warn("empty command")

		return nil
	}

	switch command[0] {
	case "quit":
		return ErrUCIQuit

	case "uci":
		return cmd.send("uciok")

	case "isready":
		return cmd.send("readyok")

	case "ucinewgame":
		cmd.game = nil
		return nil

	case "position":
		return cmd.RunPosition(command)

	case "go":
		return cmd.RunGo(command)

	case "stop":
		return cmd.RunStop(command)

	default:
		slog.Error("unknown command", "command", command)
		return nil
	}
}

func (cmd *CommandUCI) RunPosition(command []string) error {
	if len(command) == 1 {
		return fmt.Errorf("%w: invalid position: %s", ErrUCI, command)
	}

	command = command[1:]

	switch command[0] {
	case "startpos":
		cmd.game = must(NewGame(BoardStartPositionFEN))
		command = command[1:]

	case "fen":
		command = command[1:]

		if len(command) < 6 {
			return fmt.Errorf("%w: invalid fen: %s", ErrUCI, command)
		}

		fen := strings.Join(command[:6], " ")

		game, err := NewGame(fen)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUCI, err)
		}

		cmd.game = game

		command = command[6:]
	}

	if len(command) == 0 {
		return nil
	}

	if command[0] != "moves" {
		return fmt.Errorf("%w: invalid position: %s", ErrUCI, command)
	}

	command = command[1:]

	for _, move := range command {
		move, err := NewUCIMove(cmd.game.Board(), move)
		if err != nil {
			return err
		}

		slog.Info("making move from uci command", "move", move)

		cmd.game.MakeMove(move)
	}

	return nil
}

func (cmd *CommandUCI) RunGo(command []string) error {
	slog.Info("searching", "timeout", cmd.ThinkTime, "command", command)

	go func() {
		cmd.searching = true

		ctx, cancel := context.WithTimeout(context.Background(), cmd.ThinkTime)
		defer cancel()

		cmd.cancel = cancel

		move, _ := NewSearcher(cmd.game, SearchOptions{
			Depth:          cmd.Depth,
			OpeningBook:    cmd.OpeningBook,
			Transpositions: cmd.Transpositions,
		}).Search(ctx)

		check(cmd.send("bestmove %s", move.UCI()))
		cmd.searching = false
		cmd.cancel = nil
	}()

	return nil
}

func (cmd *CommandUCI) RunStop(command []string) error {
	if !cmd.searching {
		slog.Error("can't stop, not searching", "command", command)
		return nil
	}

	cmd.cancel()

	return nil
}
