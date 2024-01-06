package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMagics(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MagicsTest{})
}

type MagicsTest struct {
	suite.Suite
}

func (t *MagicsTest) TestDiagonalMoves() {
	for _, test := range []struct {
		src      Square
		blockers Bitboard
		expected Bitboard
	}{
		{
			src:      SquareA1,
			blockers: SquareE5.Bitboard() | SquareH8.Bitboard(),
			expected: SquareB2.Bitboard() | SquareC3.Bitboard() | SquareD4.Bitboard() | SquareE5.Bitboard(),
		},
		{
			src:      SquareA2,
			blockers: SquareE5.Bitboard() | SquareH8.Bitboard(),
			expected: SquareB1.Bitboard() |
				SquareB3.Bitboard() |
				SquareC4.Bitboard() |
				SquareD5.Bitboard() |
				SquareE6.Bitboard() |
				SquareF7.Bitboard() |
				SquareG8.Bitboard(),
		},
		{
			src:      SquareE4,
			blockers: SquareA1.Bitboard(),
			expected: SquareD5.Bitboard() |
				SquareC6.Bitboard() |
				SquareB7.Bitboard() |
				SquareA8.Bitboard() |
				SquareD3.Bitboard() |
				SquareC2.Bitboard() |
				SquareB1.Bitboard() |
				SquareF5.Bitboard() |
				SquareG6.Bitboard() |
				SquareH7.Bitboard() |
				SquareF3.Bitboard() |
				SquareG2.Bitboard() |
				SquareH1.Bitboard(),
		},
	} {
		t.Run(test.src.String(), func() {
			moves := Magic.DiagonalMoves(test.src, test.blockers)

			t.Assert().Equal(test.expected, moves)
		})
	}
}

func (t *MagicsTest) TestOrthogonalMoves() {
	for _, test := range []struct {
		src      Square
		blockers Bitboard
		expected Bitboard
	}{
		{
			src:      SquareB2,
			blockers: SquareB5.Bitboard() | SquareH2.Bitboard(),
			expected: SquareB3.Bitboard() |
				SquareB4.Bitboard() |
				SquareB5.Bitboard() |
				SquareA2.Bitboard() |
				SquareB1.Bitboard() |
				SquareC2.Bitboard() |
				SquareD2.Bitboard() |
				SquareE2.Bitboard() |
				SquareF2.Bitboard() |
				SquareG2.Bitboard() |
				SquareH2.Bitboard(),
		},
		{
			src: SquareD4,
			expected: SquareA4.Bitboard() |
				SquareB4.Bitboard() |
				SquareC4.Bitboard() |
				SquareE4.Bitboard() |
				SquareF4.Bitboard() |
				SquareG4.Bitboard() |
				SquareH4.Bitboard() |
				SquareD1.Bitboard() |
				SquareD2.Bitboard() |
				SquareD3.Bitboard() |
				SquareD5.Bitboard() |
				SquareD6.Bitboard() |
				SquareD7.Bitboard() |
				SquareD8.Bitboard(),
		},
	} {
		t.Run(test.src.String(), func() {
			moves := Magic.OrthogonalMoves(test.src, test.blockers)

			t.Assert().Equal(test.expected, moves)
		})
	}
}

func (t *MagicsTest) TestKingMoves() {
	for _, test := range []struct {
		src      Square
		expected Bitboard
	}{
		{
			src: SquareB2,
			expected: SquareA3.Bitboard() | SquareB3.Bitboard() | SquareC3.Bitboard() |
				SquareA2.Bitboard() | SquareC2.Bitboard() |
				SquareA1.Bitboard() | SquareB1.Bitboard() | SquareC1.Bitboard(),
		},
		{
			src: SquareH8,
			expected: SquareG8.Bitboard() |
				SquareG7.Bitboard() | SquareH7.Bitboard(),
		},
	} {
		t.Run(test.src.String(), func() {
			moves := Magic.KingMoves(test.src)

			t.Assert().Equal(test.expected, moves)
		})
	}
}

func (t *MagicsTest) TestKnightMoves() {
	for _, test := range []struct {
		src      Square
		expected Bitboard
	}{
		{
			src: SquareG1,
			expected: SquareH3.Bitboard() |
				SquareF3.Bitboard() |
				SquareE2.Bitboard(),
		},
		{
			src: SquareD4,
			expected: SquareC6.Bitboard() |
				SquareE6.Bitboard() |
				SquareF5.Bitboard() |
				SquareF3.Bitboard() |
				SquareE2.Bitboard() |
				SquareC2.Bitboard() |
				SquareB3.Bitboard() |
				SquareB5.Bitboard(),
		},
	} {
		t.Run(test.src.String(), func() {
			moves := Magic.KnightMoves(test.src)

			t.Assert().Equal(test.expected, moves)
		})
	}
}

func (t *MagicsTest) TestPawnAttacks() {
	for _, test := range []struct {
		src      Square
		color    Color
		expected Bitboard
	}{
		{
			src:      SquareG2,
			color:    ColorWhite,
			expected: SquareF3.Bitboard() | SquareH3.Bitboard(),
		},
		{
			src:      SquareH7,
			color:    ColorBlack,
			expected: SquareG6.Bitboard(),
		},
		{
			src:      SquareB7,
			color:    ColorWhite,
			expected: SquareA8.Bitboard() | SquareC8.Bitboard(),
		},
	} {
		t.Run(test.src.String(), func() {
			moves := Magic.PawnAttacks(test.color, test.src)

			t.Assert().Equal(test.expected, moves)
		})
	}
}
