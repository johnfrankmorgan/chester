package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	ucilexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "FEN", Pattern: `(?i)[1-8pnbrqk/]+ [wb] [kq-]+ [a-f1-8-]+ [0-9]+ [0-9]+`},
		{Name: "Integer", Pattern: `[0-9]+`},
		{Name: "String", Pattern: `[a-zA-Z_0-9]+`},
		{Name: "Whitespace", Pattern: `\s+`},
	})

	uciparser = participle.MustBuild[UCICommand](
		participle.Lexer(ucilexer),
		participle.Elide("Whitespace"),
	)
)

var (
	ErrUCI     = errors.New("uci")
	ErrUCIQuit = fmt.Errorf("%w: quit", ErrUCI)
)

type UCI struct {
	input  *bufio.Scanner
	output io.Writer

	game  *Game
	debug bool
}

func NewUCI(input io.Reader, output io.Writer) *UCI {
	return &UCI{
		input:  bufio.NewScanner(input),
		output: output,
	}
}

func (uci *UCI) Run() error {
	if err := uci.run(); err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("%w: %w", ErrUCIQuit, err)
		} else if !errors.Is(err, ErrUCI) {
			err = fmt.Errorf("%w: %w", ErrUCI, err)
		}

		return err
	}

	return nil
}

func (uci *UCI) run() error {
	uci.input.Scan()

	if err := uci.input.Err(); err != nil {
		return err
	}

	raw := uci.input.Text()

	slog.Info("received uci command", "command", raw)

	cmd, err := uciparser.ParseString("", raw)
	if err != nil {
		return err
	}

	return cmd.Execute(uci)
}

func (uci *UCI) Respond(format string, args ...any) error {
	response := fmt.Sprintf(strings.TrimSpace(format), args...)

	slog.Info("sending uci response", "response", response)

	_, err := fmt.Fprintln(uci.output, response)

	return err
}

type UCICommand struct {
	UCI       *UCICommandUCI       `parser:"( @@"`
	Debug     *UCICommandDebug     `parser:"| @@"`
	IsReady   *UCICommandIsReady   `parser:"| @@"`
	SetOption *UCICommandSetOption `parser:"| @@"`
	Position  *UCICommandPosition  `parser:"| @@"`
	Go        *UCICommandGo        `parser:"| @@"`
	Stop      *UCICommandStop      `parser:"| @@"`
	Quit      *UCICommandQuit      `parser:"| @@)"`
}

func (cmd UCICommand) Execute(uci *UCI) error {
	switch {
	case cmd.UCI != nil:
		return cmd.UCI.Execute(uci)

	case cmd.Debug != nil:
		return cmd.Debug.Execute(uci)

	case cmd.IsReady != nil:
		return cmd.IsReady.Execute(uci)

	case cmd.SetOption != nil:
		return cmd.SetOption.Execute(uci)

	case cmd.Position != nil:
		return cmd.Position.Execute(uci)

	case cmd.Go != nil:
		return cmd.Go.Execute(uci)

	case cmd.Stop != nil:
		return cmd.Stop.Execute(uci)

	case cmd.Quit != nil:
		return cmd.Quit.Execute(uci)
	}

	return fmt.Errorf("%w: invalid command", ErrUCI)
}

type UCICommandUCI struct {
	Prefix string `parser:"@'uci'"`
}

func (cmd UCICommandUCI) Execute(uci *UCI) error {
	return uci.Respond("uciok")
}

type UCICommandDebug struct {
	Prefix string `parser:"@'debug'"`
	State  string `parser:"@('on' | 'off')"`
}

func (cmd UCICommandDebug) Execute(uci *UCI) error {
	uci.debug = cmd.State == "on"

	return nil
}

type UCICommandIsReady struct {
	Prefix string `parser:"@'isready'"`
}

func (cmd UCICommandIsReady) Execute(uci *UCI) error {
	return uci.Respond("readyok")
}

type UCICommandSetOption struct {
	Prefix string `parser:"@'setoption'"`
	Name   string `parser:"'name' @String"`
	Value  string `parser:"'value' @(String | Integer)"`
}

func (cmd UCICommandSetOption) Execute(uci *UCI) error {
	return fmt.Errorf("%w: unsupported option: %s", ErrUCI, cmd.Name)
}

type UCICommandPosition struct {
	Prefix   string   `parser:"@'position'"`
	StartPos bool     `parser:"( @'startpos'"`
	FEN      string   `parser:"| 'fen' @FEN)"`
	Moves    []string `parser:"('moves' @String*)?"`
}

func (cmd UCICommandPosition) Execute(uci *UCI) error {
	panic("todo")
}

type UCICommandGo struct {
	Prefix string `parser:"@'go'"`
}

func (cmd UCICommandGo) Execute(uci *UCI) error {
	panic("todo")
}

type UCICommandStop struct {
	Prefix string `parser:"@'stop'"`
}

func (cmd UCICommandStop) Execute(uci *UCI) error {
	panic("todo")
}

type UCICommandQuit struct {
	Prefix string `parser:"@'quit'"`
}

func (cmd UCICommandQuit) Execute(uci *UCI) error {
	return ErrUCIQuit
}
