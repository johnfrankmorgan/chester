package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type UCI struct {
	OpeningBook         bool          `help:"Use opening book to play opening moves" default:"true" negatable:""`
	OpeningBookMoves    int           `help:"Number of moves to play from the opening book" default:"20"`
	DefaultMoveTime     time.Duration `help:"Default time to spend calculating the best move" default:"250ms" env:"CHESTER_DEFAULT_MOVE_TIME"`
	DefaultInfoInterval time.Duration `help:"Default interval to send info messages" default:"500ms"`

	stdin  io.Reader
	stdout io.Writer

	quit  bool
	debug bool

	game *Game
	tt   *TranspositionTable
	sctx *SearchContext
	stop func()
}

func (uci *UCI) Run(ctx context.Context) error {
	uci.debug = true

	uci.stdin = os.Stdin
	uci.stdout = os.Stdout

	slog.Info("starting uci engine")

	uci.tt = NewTranspositionTable()

	return uci.run(ctx)
}

func (uci *UCI) run(ctx context.Context) error {
	input := bufio.NewScanner(uci.stdin)

	for !uci.quit && input.Scan() {
		if ctx.Err() != nil {
			return nil
		}

		cmd := UCICommandFromString(input.Text())
		if len(cmd) == 0 {
			slog.Warn("empty command")
			continue
		}

		slog.Debug("received command", "command", cmd)

		uci.handle(ctx, cmd)
	}

	if err := input.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func (uci *UCI) handle(ctx context.Context, cmd UCICommand) {
	switch name := cmd.Name(); name {
	case "uci":
		uci.send("uciok")

	case "isready":
		uci.send("readyok")

	case "ucinewgame":
		// no-op for now

	case "debug":
		uci.debug = cmd.BoolArg("on")

	case "position":
		uci.position(cmd)

	case "print":
		uci.send(uci.game.Board())

	case "go":
		if uci.game == nil {
			slog.Warn("no position set")
			return
		}

		if uci.sctx != nil {
			slog.Warn("search already in progress")
			return
		}

		if depth, ok := cmd.IntArg("perft"); ok {
			slog.Debug("running perft", "depth", depth)

			divide := io.Writer(nil)
			if uci.debug {
				divide = uci.stdout
			}

			uci.send("nodes:", Perft(uci.game, depth, divide))
			return
		}

		if cmd.BoolArg("infinite") {
			ctx, cancel := context.WithCancel(ctx)
			uci.stop = cancel

			slog.Debug("starting infinite search")
			go uci.search(ctx)
		} else {
			timeout := uci.DefaultMoveTime

			if mtime, ok := cmd.IntArg("movetime"); ok {
				timeout = time.Duration(mtime) * time.Millisecond
			} else {
				slog.Info("movetime not provided, using default")
			}

			ctx, cancel := context.WithTimeout(ctx, timeout)
			uci.stop = cancel

			slog.Debug("starting movetime search", "timeout", timeout)
			go uci.search(ctx)
		}

	case "stop":
		if uci.sctx == nil {
			slog.Warn("attempted to stop without a search in progress")
			return
		}

		uci.stop()

	case "setoption", "ponderhit":
		slog.Warn("not implemented", "command", name)

	case "quit":
		if uci.stop != nil {
			uci.stop()
		}

		uci.quit = true

	default:
		slog.Warn("invalid command", "command", cmd)
	}
}

func (uci *UCI) position(cmd UCICommand) {
	fen := BoardStartPos

	if !cmd.BoolArg("startpos") {
		if len(cmd) < 3 {
			slog.Warn("invalid position command", "command", cmd)
			return
		}

		if cmd[1] != "fen" {
			slog.Warn("missing fen", "command", cmd)
			return
		}

		end := slices.Index(cmd, "moves")
		if end == -1 {
			end = len(cmd)
		}

		fen = strings.Join(cmd[2:end], " ")
	}

	slog.Debug("setting position", "fen", fen)

	game, err := GameFromFEN(fen)
	if err != nil {
		slog.Warn("invalid position", "error", err)
		return
	}

	if moves := slices.Index(cmd, "moves"); moves != -1 {
		for _, move := range cmd[moves+1:] {
			if !game.MakeUCIMove(move) {
				slog.Warn("invalid move", "move", move)
				return
			}
		}
	}

	uci.game = game
}

func (uci *UCI) search(ctx context.Context) {
	uci.sctx = &SearchContext{
		Context: ctx,
		Game:    uci.game,
		TT:      uci.tt,
	}

	go func() {
		ticker := time.NewTicker(uci.DefaultInfoInterval)
		defer ticker.Stop()

		for range ticker.C {
			if ctx.Err() != nil {
				break
			}

			uci.info()
		}
	}()

	go func() {
		if uci.OpeningBook && len(uci.sctx.Game.Moves()) < uci.OpeningBookMoves {
			slog.Debug("trying book move")

			move := RandomOpeningMove(uci.sctx.Game.Moves()...)
			if move != nil {
				slog.Info("using book move", "move", move)
				uci.send("bestmove", move)

				uci.stop()
				uci.sctx = nil
				uci.stop = nil

				return
			}
		}

		defer func() {
			uci.stop()
			uci.info()
			uci.send("bestmove", uci.sctx.Best)

			uci.sctx = nil
			uci.stop = nil
		}()

		Search(uci.sctx)
	}()

	<-ctx.Done()
}

func (uci *UCI) info() {
	if !uci.debug {
		return
	}

	if uci.sctx.CurrentMove.IsZero() {
		return
	}

	uci.send(
		"info",
		"time", time.Since(uci.sctx.Start).Milliseconds(),
		"depth", uci.sctx.Depth,
		"nodes", uci.sctx.Nodes,
		"currmove", uci.sctx.CurrentMove,
	)
}

func (uci *UCI) send(msg ...any) {
	slog.Debug("sending", "msg", msg)

	if _, err := fmt.Fprintln(uci.stdout, msg...); err != nil {
		slog.Warn("failed to send", "error", err)
	}
}

type UCICommand []string

func UCICommandFromString(cmd string) UCICommand {
	return strings.Fields(cmd)
}

func (cmd UCICommand) String() string {
	return strings.Join(cmd, " ")
}

func (cmd UCICommand) Name() string {
	return cmd[0]
}

func (cmd UCICommand) Arg(name string) (string, bool) {
	index := slices.Index(cmd, name)
	if index == -1 {
		return "", false
	}

	if index < len(cmd)-1 {
		return cmd[index+1], true
	}

	return "", false
}

func (cmd UCICommand) IntArg(name string) (int, bool) {
	arg, ok := cmd.Arg(name)
	if !ok {
		return 0, false
	}

	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0, false
	}

	return value, true
}

func (cmd UCICommand) BoolArg(name string) bool {
	return slices.Contains(cmd, name)
}
