package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommandDivide(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CommandDivideTest{})
}

type CommandDivideTest struct {
	suite.Suite

	cmd CommandDivide
	out bytes.Buffer
}

func (t *CommandDivideTest) SetupTest() {
	t.cmd = CommandDivide{
		FEN:   BoardStartPositionFEN,
		Depth: 2,
	}

	t.cmd.SetOut(wcloser{&t.out})
}

func (t *CommandDivideTest) TestRun() {
	t.Run("valid fen", func() {
		err := t.cmd.Run()

		expected := []string{
			"FEN: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"",
			"b1c3: 20",
			"b1a3: 20",
			"g1h3: 20",
			"g1f3: 20",
			"a2a3: 20",
			"a2a4: 20",
			"b2b3: 20",
			"b2b4: 20",
			"c2c3: 20",
			"c2c4: 20",
			"d2d3: 20",
			"d2d4: 20",
			"e2e3: 20",
			"e2e4: 20",
			"f2f3: 20",
			"f2f4: 20",
			"g2g3: 20",
			"g2g4: 20",
			"h2h3: 20",
			"h2h4: 20",
			"",
			"Nodes: 400",
			"",
		}

		t.Assert().NoError(err)
		t.Assert().Equal(strings.Join(expected, "\n"), t.out.String())
	})

	t.Run("invalid fen", func() {
		t.cmd.FEN = ""

		t.Assert().Error(t.cmd.Run())
	})
}
