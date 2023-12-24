package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSquare(t *testing.T) {
	t.Parallel()

	suite.Run(t, &SquareTest{})
}

type SquareTest struct {
	suite.Suite
}

func (t *SquareTest) TestString() {
	for _, test := range []struct {
		square   Square
		expected string
	}{
		{NewSquare(FileA, Rank3), "a3"},
		{SquareA1, "a1"},
		{SquareB5, "b5"},
		{SquareH8, "h8"},
		{-1, "main.Square(-1)"},
		{100, "main.Square(100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.square.String())
		})
	}
}

func (t *SquareTest) TestValid() {
	for _, test := range []struct {
		square   Square
		expected bool
	}{
		{NewSquare(FileA, Rank3), true},
		{SquareA1, true},
		{SquareB5, true},
		{SquareH8, true},
		{-10, false},
		{100, false},
	} {
		t.Run(test.square.String(), func() {
			t.Assert().Equal(test.expected, test.square.Valid())
		})
	}
}

func (t *SquareTest) TestFile() {
	for _, test := range []struct {
		square   Square
		expected File
	}{
		{NewSquare(FileA, Rank3), FileA},
		{SquareA1, FileA},
		{SquareB5, FileB},
		{SquareH8, FileH},
	} {
		t.Run(test.square.String(), func() {
			t.Assert().Equal(test.expected, test.square.File())
		})
	}
}

func (t *SquareTest) TestRank() {
	for _, test := range []struct {
		square   Square
		expected Rank
	}{
		{NewSquare(FileA, Rank3), Rank3},
		{SquareA1, Rank1},
		{SquareB5, Rank5},
		{SquareH8, Rank8},
	} {
		t.Run(test.square.String(), func() {
			t.Assert().Equal(test.expected, test.square.Rank())
		})
	}
}

func (t *SquareTest) TestCoord() {
	for _, test := range []struct {
		square Square
		file   File
		rank   Rank
	}{
		{NewSquare(FileA, Rank3), FileA, Rank3},
		{SquareA1, FileA, Rank1},
		{SquareB5, FileB, Rank5},
		{SquareH8, FileH, Rank8},
	} {
		t.Run(test.square.String(), func() {
			file, rank := test.square.Coord()

			t.Assert().Equal(test.file, file)
			t.Assert().Equal(test.rank, rank)
		})
	}
}

func (t *SquareTest) TestBitboard() {
	for _, test := range []struct {
		square   Square
		expected Bitboard
	}{
		{SquareA1, 0b00001},
		{SquareB1, 0b00010},
		{SquareC1, 0b00100},
		{SquareD1, 0b01000},
		{SquareE1, 0b10000},
		{SquareH8, 1 << 63},
	} {
		t.Run(test.square.String(), func() {
			t.Assert().Equal(test.expected, test.square.Bitboard())
		})
	}
}
