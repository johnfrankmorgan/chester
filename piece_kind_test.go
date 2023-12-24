package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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
		{PiecePawn, "p"},
		{PieceKnight, "n"},
		{PieceBishop, "b"},
		{PieceRook, "r"},
		{PieceQueen, "q"},
		{PieceKing, "k"},
		{10, "main.PieceKind(10)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.kind.String())
		})
	}
}
