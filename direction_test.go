package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDirection(t *testing.T) {
	t.Parallel()

	suite.Run(t, &DirectionTest{})
}

type DirectionTest struct {
	suite.Suite
}

func (t *DirectionTest) TestString() {
	for _, test := range []struct {
		dir      Direction
		expected string
	}{
		{DirectionNorth, "north"},
		{DirectionSouth, "south"},
		{DirectionEast, "east"},
		{DirectionWest, "west"},
		{DirectionNorthEast, "north east"},
		{DirectionSouthWest, "south west"},
		{DirectionNorthWest, "north west"},
		{DirectionSouthEast, "south east"},
		{100, "main.Direction(100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.dir.String())
		})
	}
}

func (t *DirectionTest) TestOffset() {
	for _, test := range []struct {
		dir      Direction
		expected Square
	}{
		{DirectionNorth, 8},
		{DirectionSouth, -8},
		{DirectionEast, 1},
		{DirectionWest, -1},
		{DirectionNorthEast, 9},
		{DirectionSouthWest, -9},
		{DirectionNorthWest, 7},
		{DirectionSouthEast, -7},
	} {
		t.Run(test.dir.String(), func() {
			t.Assert().Equal(test.expected, test.dir.Offset())
		})
	}
}

func (t *DirectionTest) TestOpposite() {
	for _, test := range []struct {
		dir      Direction
		expected Direction
	}{
		{DirectionNorth, DirectionSouth},
		{DirectionSouth, DirectionNorth},
		{DirectionEast, DirectionWest},
		{DirectionWest, DirectionEast},
		{DirectionNorthEast, DirectionSouthWest},
		{DirectionSouthWest, DirectionNorthEast},
		{DirectionNorthWest, DirectionSouthEast},
		{DirectionSouthEast, DirectionNorthWest},
	} {
		t.Run(test.dir.String(), func() {
			t.Assert().Equal(test.expected, test.dir.Opposite())
		})
	}

}

func (t *DirectionTest) TestToEdge() {
	for _, test := range []struct {
		dir      Direction
		src      Square
		expected Square
	}{
		{DirectionNorth, SquareA1, 7},
		{DirectionNorthEast, SquareA1, 7},
		{DirectionSouth, SquareB2, 1},
		{DirectionNorth, SquareH8, 0},
		{DirectionNorthEast, SquareH8, 0},
		{DirectionWest, SquareE4, 4},
		{DirectionEast, SquareE4, 3},
	} {
		t.Run(fmt.Sprintf("%s(%s)", test.src, test.dir), func() {
			t.Assert().Equal(test.expected, test.dir.ToEdge(test.src))
		})
	}
}

func (t *DirectionTest) TestIsDiagonal() {
	for _, test := range []struct {
		dir      Direction
		expected bool
	}{
		{DirectionNorth, false},
		{DirectionSouth, false},
		{DirectionEast, false},
		{DirectionWest, false},
		{DirectionNorthEast, true},
		{DirectionSouthWest, true},
		{DirectionNorthWest, true},
		{DirectionSouthEast, true},
	} {
		t.Run(test.dir.String(), func() {
			t.Assert().Equal(test.expected, test.dir.IsDiagonal())
		})
	}
}

func (t *DirectionTest) TestMask() {
	for _, test := range []struct {
		dir      Direction
		src      Square
		expected Bitboard
	}{
		{
			dir:      DirectionNorth,
			src:      SquareA1,
			expected: BitboardFileA,
		},
		{
			dir: DirectionNorthEast,
			src: SquareA1,
			expected: SquareA1.Bitboard() |
				SquareB2.Bitboard() |
				SquareC3.Bitboard() |
				SquareD4.Bitboard() |
				SquareE5.Bitboard() |
				SquareF6.Bitboard() |
				SquareG7.Bitboard() |
				SquareH8.Bitboard(),
		},
		{
			dir:      DirectionSouth,
			src:      SquareB2,
			expected: BitboardFileB,
		},
		{
			dir:      DirectionNorth,
			src:      SquareH8,
			expected: BitboardFileH,
		},
		{
			dir:      DirectionWest,
			src:      SquareE4,
			expected: BitboardRank4,
		},
		{
			dir:      DirectionEast,
			src:      SquareE4,
			expected: BitboardRank4,
		},
		{
			dir: DirectionNorthWest,
			src: SquareE2,
			expected: SquareF1.Bitboard() |
				SquareE2.Bitboard() |
				SquareD3.Bitboard() |
				SquareC4.Bitboard() |
				SquareB5.Bitboard() |
				SquareA6.Bitboard(),
		},
	} {
		t.Run(fmt.Sprintf("%s(%s)", test.src, test.dir), func() {
			t.Assert().Equal(test.expected, test.dir.Mask(test.src))
		})
	}
}
