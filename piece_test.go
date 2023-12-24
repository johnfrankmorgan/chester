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
		{NewPiece(ColorBlack, PieceQueen), "q"},
		{NewPiece(ColorWhite, PieceKing), "K"},
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
		{NewPiece(ColorBlack, PieceQueen), ColorBlack},
		{NewPiece(ColorWhite, PieceKing), ColorWhite},
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
		{PieceBlackPawn, PiecePawn},
		{PieceWhitePawn, PiecePawn},
		{NewPiece(ColorBlack, PieceQueen), PieceQueen},
		{NewPiece(ColorWhite, PieceKing), PieceKing},
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
		{PieceBlackPawn, PiecePawn, true},
		{PieceWhitePawn, PiecePawn, true},
		{NewPiece(ColorBlack, PieceQueen), PieceQueen, true},
		{NewPiece(ColorWhite, PieceKing), PieceKing, true},
		{NewPiece(ColorWhite, PieceKing), PieceQueen, false},
		{PieceBlackKnight, PieceBishop, false},
	} {
		t.Run(test.piece.String(), func() {
			t.Assert().Equal(test.expected, test.piece.Is(test.kind))
		})
	}
}
