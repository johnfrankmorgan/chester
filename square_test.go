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

func (t *SquareTest) TestAlignMask() {
	for src := SquareFirst; src < SquareLast; src++ {
		for dst := SquareFirst; dst < SquareLast; dst++ {
			mask := src.AlignMask(dst)

			if mask == 0 {
				continue
			}

			t.Run(src.String()+dst.String(), func() {
				t.Assert().True(mask.IsSet(src.Bitboard()))
				t.Assert().True(mask.IsSet(dst.Bitboard()))
			})
		}
	}

	for _, test := range []struct {
		src      Square
		dst      Square
		expected Bitboard
	}{
		{SquareA3, SquareD3, BitboardRank3},
		{SquareE2, SquareE7, BitboardFileE},
	} {
		t.Run(test.src.String()+test.dst.String(), func() {
			t.Assert().Equal(test.expected, test.src.AlignMask(test.dst))
		})
	}
}

func TestFile(t *testing.T) {
	t.Parallel()

	suite.Run(t, &FileTest{})
}

type FileTest struct {
	suite.Suite
}

func (t *FileTest) TestString() {
	for _, test := range []struct {
		file     File
		expected string
	}{
		{FileA, "a"},
		{FileB, "b"},
		{FileC, "c"},
		{FileD, "d"},
		{FileE, "e"},
		{FileF, "f"},
		{FileG, "g"},
		{FileH, "h"},
		{10, "main.File(10)"},
		{-100, "main.File(-100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.file.String())
		})
	}
}

func (t *FileTest) TestValid() {
	for _, test := range []struct {
		file     File
		expected bool
	}{
		{FileA, true},
		{FileB, true},
		{FileC, true},
		{FileD, true},
		{FileE, true},
		{FileF, true},
		{FileG, true},
		{FileH, true},
		{10, false},
		{-100, false},
	} {
		t.Run(test.file.String(), func() {
			t.Assert().Equal(test.expected, test.file.Valid())
		})
	}
}

func TestRank(t *testing.T) {
	t.Parallel()

	suite.Run(t, &RankTest{})
}

type RankTest struct {
	suite.Suite
}

func (t *RankTest) TestString() {
	for _, test := range []struct {
		file     Rank
		expected string
	}{
		{Rank1, "1"},
		{Rank2, "2"},
		{Rank3, "3"},
		{Rank4, "4"},
		{Rank5, "5"},
		{Rank6, "6"},
		{Rank7, "7"},
		{Rank8, "8"},
		{10, "main.Rank(10)"},
		{-100, "main.Rank(-100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.file.String())
		})
	}
}

func (t *RankTest) TestValid() {
	for _, test := range []struct {
		file     Rank
		expected bool
	}{
		{Rank1, true},
		{Rank2, true},
		{Rank3, true},
		{Rank4, true},
		{Rank5, true},
		{Rank6, true},
		{Rank7, true},
		{Rank8, true},
		{10, false},
		{-100, false},
	} {
		t.Run(test.file.String(), func() {
			t.Assert().Equal(test.expected, test.file.Valid())
		})
	}
}
