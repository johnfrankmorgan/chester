package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommandGenerateMagics(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CommandGenerateMagicsTest{})
}

type CommandGenerateMagicsTest struct {
	suite.Suite
}

func (t *CommandGenerateMagicsTest) TearDownTest() {
	PanicIfError(
		Magic.LoadDefault(),
	)
}

func (t *CommandGenerateMagicsTest) TestRun() {
	for _, test := range []struct {
		name   string
		skip   bool
		cmd    CommandGenerateMagics
		assert func(*CommandGenerateMagicsTest, CommandGenerateMagics, error)
	}{
		{
			name: "file creation fails",
			cmd:  CommandGenerateMagics{},
			assert: func(t *CommandGenerateMagicsTest, cmd CommandGenerateMagics, err error) {
				t.Assert().ErrorIs(err, os.ErrNotExist)
			},
		},
		{
			name: "diagonal generation",
			cmd: CommandGenerateMagics{
				Output:   "/tmp/chester.diagonal.gob",
				Diagonal: true,
			},
			assert: func(t *CommandGenerateMagicsTest, cmd CommandGenerateMagics, err error) {
				t.Assert().NoError(err)

				magics := Must(os.Open(cmd.Output))
				defer magics.Close()

				PanicIfError(Magic.Load(magics))

				moves := Magic.Diagonal(SquareH8, SquareG7.Bitboard()|SquareH7.Bitboard()|SquareG8.Bitboard())
				t.Assert().Equal(SquareG7.Bitboard(), moves)
			},
		},
		{
			name: "orthogonal generation",
			skip: testing.Short(),
			cmd: CommandGenerateMagics{
				Output:     "/tmp/chester.orthogonal.gob",
				Orthogonal: true,
			},
			assert: func(t *CommandGenerateMagicsTest, cmd CommandGenerateMagics, err error) {
				t.Assert().NoError(err)

				magics := Must(os.Open(cmd.Output))
				defer magics.Close()

				PanicIfError(Magic.Load(magics))

				moves := Magic.Orthogonal(SquareH8, SquareG7.Bitboard()|SquareH7.Bitboard()|SquareG8.Bitboard())
				t.Assert().Equal(SquareG8.Bitboard()|SquareH7.Bitboard(), moves)
			},
		},
		{
			name: "king generation",
			cmd: CommandGenerateMagics{
				Output: "/tmp/chester.king.gob",
				King:   true,
			},
			assert: func(t *CommandGenerateMagicsTest, cmd CommandGenerateMagics, err error) {
				t.Assert().NoError(err)

				magics := Must(os.Open(cmd.Output))
				defer magics.Close()

				PanicIfError(Magic.Load(magics))

				expected := SquareD4.Bitboard() |
					SquareD5.Bitboard() |
					SquareE5.Bitboard() |
					SquareF5.Bitboard() |
					SquareF4.Bitboard() |
					SquareF3.Bitboard() |
					SquareE3.Bitboard() |
					SquareD3.Bitboard()

				t.Assert().Equal(expected, Magic.King(SquareE4))
			},
		},
		{
			name: "knight generation",
			cmd: CommandGenerateMagics{
				Output: "/tmp/chester.knight.gob",
				Knight: true,
			},
			assert: func(t *CommandGenerateMagicsTest, cmd CommandGenerateMagics, err error) {
				t.Assert().NoError(err)

				magics := Must(os.Open(cmd.Output))
				defer magics.Close()

				PanicIfError(Magic.Load(magics))

				t.Assert().Equal(SquareA3.Bitboard()|SquareC3.Bitboard()|SquareD2.Bitboard(), Magic.Knight(SquareB1))
			},
		},
	} {
		t.Run(test.name, func() {
			if test.skip {
				t.T().SkipNow()
			}

			defer os.Remove(test.cmd.Output)

			err := test.cmd.Run()

			test.assert(t, test.cmd, err)
		})
	}
}
