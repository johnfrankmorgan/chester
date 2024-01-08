package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPiece(t *testing.T) {
	t.Parallel()

	suite.Run(t, &PieceTest{})
}

type PieceTest struct {
	suite.Suite
}

func (t *PieceTest) TestString() {
	for _, test := range []struct {
		piece    Piece
		expected string
	}{
		{PieceBlackPawn, "p"},
		{PieceWhitePawn, "P"},
		{NewPiece(ColorBlack, PieceKindQueen), "q"},
		{NewPiece(ColorWhite, PieceKindKing), "K"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.piece.String())
		})
	}
}

func (t *PieceTest) TestColor() {
	for _, test := range []struct {
		piece    Piece
		expected Color
	}{
		{PieceBlackPawn, ColorBlack},
		{PieceWhitePawn, ColorWhite},
		{NewPiece(ColorBlack, PieceKindQueen), ColorBlack},
		{NewPiece(ColorWhite, PieceKindKing), ColorWhite},
	} {
		t.Run(test.piece.String(), func() {
			t.Assert().Equal(test.expected, test.piece.Color())
		})
	}
}

func (t *PieceTest) TestKind() {
	for _, test := range []struct {
		piece    Piece
		expected PieceKind
	}{
		{PieceBlackPawn, PieceKindPawn},
		{PieceWhitePawn, PieceKindPawn},
		{NewPiece(ColorBlack, PieceKindQueen), PieceKindQueen},
		{NewPiece(ColorWhite, PieceKindKing), PieceKindKing},
	} {
		t.Run(test.piece.String(), func() {
			t.Assert().Equal(test.expected, test.piece.Kind())
		})
	}
}

func (t *PieceTest) TestIs() {
	for _, test := range []struct {
		piece    Piece
		kind     PieceKind
		expected bool
	}{
		{PieceBlackPawn, PieceKindPawn, true},
		{PieceWhitePawn, PieceKindPawn, true},
		{NewPiece(ColorBlack, PieceKindQueen), PieceKindQueen, true},
		{NewPiece(ColorWhite, PieceKindKing), PieceKindKing, true},
		{NewPiece(ColorWhite, PieceKindKing), PieceKindQueen, false},
		{PieceBlackKnight, PieceKindBishop, false},
	} {
		t.Run(test.piece.String(), func() {
			t.Assert().Equal(test.expected, test.piece.Is(test.kind))
		})
	}
}

func TestPieceKind(t *testing.T) {
	t.Parallel()

	suite.Run(t, &PieceKindTest{})
}

type PieceKindTest struct {
	suite.Suite
}

func (t *PieceKindTest) TestString() {
	for _, test := range []struct {
		kind     PieceKind
		expected string
	}{
		{PieceKindPawn, "p"},
		{PieceKindKnight, "n"},
		{PieceKindBishop, "b"},
		{PieceKindRook, "r"},
		{PieceKindQueen, "q"},
		{PieceKindKing, "k"},
		{10, "main.PieceKind(10)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.kind.String())
		})
	}
}

func TestColor(t *testing.T) {
	t.Parallel()

	suite.Run(t, &ColorTest{})
}

type ColorTest struct {
	suite.Suite
}

func (t *ColorTest) TestString() {
	for _, test := range []struct {
		color    Color
		expected string
	}{
		{ColorWhite, "w"},
		{ColorBlack, "b"},
		{10, "main.Color(10)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.color.String())
		})
	}
}

func (t *ColorTest) TestValid() {
	for _, test := range []struct {
		color    Color
		expected bool
	}{
		{ColorWhite, true},
		{ColorBlack, true},
		{10, false},
	} {
		t.Run(test.color.String(), func() {
			t.Assert().Equal(test.expected, test.color.Valid())
		})
	}
}

func (t *ColorTest) TestOpponent() {
	for _, test := range []struct {
		color    Color
		expected Color
	}{
		{ColorWhite, ColorBlack},
		{ColorBlack, ColorWhite},
	} {
		t.Run(test.color.String(), func() {
			t.Assert().Equal(test.expected, test.color.Opponent())
		})
	}
}
