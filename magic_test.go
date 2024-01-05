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

func (t *MagicsTest) TestDiagonal() {
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
			moves := Magic.Diagonal(test.src, test.blockers)

			t.Assert().Equal(test.expected, moves)
		})
	}
}

func (t *MagicsTest) TestOrthogonal() {
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
			moves := Magic.Orthogonal(test.src, test.blockers)

			t.Assert().Equal(test.expected, moves)
		})
	}
}
