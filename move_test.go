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
