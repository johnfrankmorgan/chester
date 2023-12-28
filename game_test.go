package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestGame(t *testing.T) {
	t.Parallel()

	suite.Run(t, &GameTest{})
}

type GameTest struct {
	suite.Suite
}

func (t *GameTest) TestNewGame() {
	t.Run("invalid fens result in error", func() {
		game, err := NewGame("")

		t.Assert().ErrorIs(err, ErrInvalidFEN)
		t.Assert().Nil(game)
	})

	t.Run("valid fens result in new game", func() {
		game, err := NewGame(BoardStartPositionFEN)

		t.Assert().NoError(err)
		t.Assert().NotNil(game)
	})
}

func (t *GameTest) TestBoard() {
	board := must(NewBoard(BoardStartPositionFEN))

	game := &Game{_boards: []Board{board}}

	t.Assert().Equal(board, *game.Board())
}

func (t *GameTest) TestMakeMove() {
	game := must(NewGame(BoardStartPositionFEN))

	game.MakeMove(NewMove(SquareE2, SquareE4, MoveFlagsDoublePawnPush))

	t.Assert().Equal(2, len(game._boards))
}

func (t *GameTest) TestUnmakeMove() {
	game := must(NewGame(BoardStartPositionFEN))

	board := *game.Board()

	game.MakeMove(NewMove(SquareE2, SquareE4, MoveFlagsDoublePawnPush))
	game.MakeMove(NewMove(SquareE2, SquareE4, MoveFlagsDoublePawnPush))
	game.MakeMove(NewMove(SquareE2, SquareE4, MoveFlagsDoublePawnPush))

	game.UnmakeMove()
	game.UnmakeMove()
	game.UnmakeMove()

	t.Assert().Equal(1, len(game._boards))
	t.Assert().Equal(board, *game.Board())
}
