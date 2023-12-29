package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMove(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MoveTest{})
}

type MoveTest struct {
	suite.Suite
}

func (t *MoveTest) TestNewUCIMove() {
	for _, test := range []struct {
		board    Board
		move     string
		err      string
		expected Move
	}{
		{
			board:    must(NewBoard(BoardStartPositionFEN)),
			move:     "a2a4",
			expected: NewMove(SquareA2, SquareA4, MoveFlagsDoublePawnPush),
		},
		{
			board:    must(NewBoard(BoardStartPositionFEN)),
			move:     "a2a3q",
			expected: NewMove(SquareA2, SquareA3, MoveFlagsPromoteToQueen),
		},
		{
			board:    must(NewBoard(BoardStartPositionFEN)),
			move:     "a2a3r",
			expected: NewMove(SquareA2, SquareA3, MoveFlagsPromoteToRook),
		},
		{
			board:    must(NewBoard(BoardStartPositionFEN)),
			move:     "a2a3b",
			expected: NewMove(SquareA2, SquareA3, MoveFlagsPromoteToBishop),
		},
		{
			board:    must(NewBoard(BoardStartPositionFEN)),
			move:     "a2a3n",
			expected: NewMove(SquareA2, SquareA3, MoveFlagsPromoteToKnight),
		},
		{
			board: SetupTestBoard([SquareCount]Piece{
				SquareE8: PieceBlackKing,
			}, nil),
			move:     "e8g8",
			expected: NewMove(SquareE8, SquareG8, MoveFlagsCastleKingside),
		},
		{
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
			}, nil),
			move:     "e1c1",
			expected: NewMove(SquareE1, SquareC1, MoveFlagsCastleQueenside),
		},
		{
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteRook,
				SquareA1: PieceBlackPawn,
			}, nil),
			move:     "e1a1",
			expected: NewMove(SquareE1, SquareA1, MoveFlagsCapture),
		},
		{
			board: SetupTestBoard([SquareCount]Piece{
				SquareE5: PieceWhitePawn,
			}, nil),
			move:     "e5d6",
			expected: NewMove(SquareE5, SquareD6, MoveFlagsCapture, MoveFlagsCaptureEnPassant),
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "a",
			err:   "invalid move: a",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "a2sdfssl",
			err:   "invalid move: a2sdfssl",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "s1d2",
			err:   "invalid source file: s",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "asd2",
			err:   "invalid source rank: s",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "a122",
			err:   "invalid destination file: 2",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "a1d9",
			err:   "invalid destination rank: 9",
		},
		{
			board: must(NewBoard(BoardStartPositionFEN)),
			move:  "a1d8k",
			err:   "invalid promotion: k",
		},
	} {
		t.Run(test.move, func() {
			move, err := NewUCIMove(&test.board, test.move)

			if test.err != "" {
				t.Assert().ErrorIs(err, ErrUCI)
				t.Assert().ErrorContains(err, test.err)
			} else {
				t.Assert().NoError(err)
			}

			t.Assert().Equal(test.expected, move)
		})
	}
}

func (t *MoveTest) TestString() {
	for _, test := range []struct {
		move     Move
		expected string
	}{
		{NewMove(SquareA1, SquareD1), "a1d1"},
		{NewMove(SquareE2, SquareE4), "e2e4"},
		{NewMove(SquareE4, SquareD5, MoveFlagsCapture), "e4d5 (c)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.move.String())
		})
	}
}

func (t *MoveTest) TestUCI() {
	for _, test := range []struct {
		move     Move
		expected string
	}{
		{NewMove(SquareA1, SquareD1), "a1d1"},
		{NewMove(SquareE2, SquareE4), "e2e4"},
		{NewMove(SquareE4, SquareD5, MoveFlagsPromoteToQueen), "e4d5q"},
		{NewMove(SquareE4, SquareD5, MoveFlagsPromoteToRook), "e4d5r"},
		{NewMove(SquareE4, SquareD5, MoveFlagsPromoteToBishop), "e4d5b"},
		{NewMove(SquareE4, SquareD5, MoveFlagsPromoteToKnight), "e4d5n"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.move.UCI())
		})
	}
}

func (t *MoveTest) TestPromotion() {
	for _, test := range []struct {
		move     Move
		expected PieceKind
	}{
		{NewMove(SquareA3, SquareA4), PieceNone},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToQueen), PieceQueen},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToRook), PieceRook},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToBishop), PieceBishop},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToKnight), PieceKnight},
	} {
		t.Run(test.move.String(), func() {
			t.Assert().Equal(test.expected, test.move.Promotion())
		})
	}
}

func (t *MoveTest) TestIsZero() {
	for _, test := range []struct {
		move     Move
		expected bool
	}{
		{Move{}, true},
		{NewMove(SquareA1, SquareA2), false},
	} {
		t.Run(test.move.String(), func() {
			t.Assert().Equal(test.expected, test.move.IsZero())
		})
	}
}
