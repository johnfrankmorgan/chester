package main

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBitboard(t *testing.T) {
	t.Parallel()

	suite.Run(t, &BitboardTest{})
}

type BitboardTest struct {
	suite.Suite
}

func (t *BitboardTest) TestString() {
	rank := func(rank Rank) Bitboard {
		return 0b11111111 << Bitboard(rank*8)
	}

	for i, test := range []struct {
		bitboard Bitboard
		expected []string
	}{
		{rank(Rank8) | rank(Rank6) | rank(Rank2) | rank(Rank1), []string{
			"8 X X X X X X X X",
			"7 . . . . . . . .",
			"6 X X X X X X X X",
			"5 . . . . . . . .",
			"4 . . . . . . . .",
			"3 . . . . . . . .",
			"2 X X X X X X X X",
			"1 X X X X X X X X",
			"  a b c d e f g h",
		}},
		{SquareA1.Bitboard(), []string{
			"8 . . . . . . . .",
			"7 . . . . . . . .",
			"6 . . . . . . . .",
			"5 . . . . . . . .",
			"4 . . . . . . . .",
			"3 . . . . . . . .",
			"2 . . . . . . . .",
			"1 X . . . . . . .",
			"  a b c d e f g h",
		}},
		{SquareF5.Bitboard() | SquareD3.Bitboard(), []string{
			"8 . . . . . . . .",
			"7 . . . . . . . .",
			"6 . . . . . . . .",
			"5 . . . . . X . .",
			"4 . . . . . . . .",
			"3 . . . X . . . .",
			"2 . . . . . . . .",
			"1 . . . . . . . .",
			"  a b c d e f g h",
		}},
	} {
		t.Run(strconv.Itoa(i), func() {
			expected := strings.Join(test.expected, "\n")

			t.Assert().Equal(expected, test.bitboard.String())
		})
	}
}

func (t *BitboardTest) TestIsSet() {
	b := Bitboard(0b1110)

	t.Assert().True(b.IsSet(0b1110))
	t.Assert().True(b.IsSet(0b1010))
	t.Assert().False(b.IsSet(0b1001))
}

func (t *BitboardTest) TestAnySet() {
	b := Bitboard(0b1111)

	t.Assert().True(b.AnySet(0b1000))
	t.Assert().True(b.AnySet(0b0001))
	t.Assert().True(b.AnySet(0b1111))
	t.Assert().False(b.AnySet(0b10000))
	t.Assert().False(b.AnySet(0))
}

func (t *BitboardTest) TestSet() {
	b := Bitboard(0b1000)
	b.Set(0b101)
	b.Set(0b010)

	t.Assert().EqualValues(0b1111, b)
}

func (t *BitboardTest) TestSetCount() {
	for _, test := range []struct {
		bitboard Bitboard
		expected int
	}{
		{0b111, 3},
		{0b100100000001, 3},
		{0, 0},
		{1, 1},
		{math.MaxUint64, 64},
	} {
		t.Assert().Equal(test.expected, test.bitboard.SetCount())
	}
}

func (t *BitboardTest) TestClear() {
	b := Bitboard(0b1010)
	b.Clear(0b10)

	t.Assert().EqualValues(0b1000, b)
}

func (t *BitboardTest) TestPopLSB() {
	b := Bitboard(0b1111)

	lsb := b.PopLSB()
	t.Assert().EqualValues(0, lsb)
	t.Assert().EqualValues(0b1110, b)

	lsb = b.PopLSB()
	t.Assert().EqualValues(1, lsb)
	t.Assert().EqualValues(0b1100, b)

	lsb = b.PopLSB()
	t.Assert().EqualValues(2, lsb)
	t.Assert().EqualValues(0b1000, b)

	lsb = b.PopLSB()
	t.Assert().EqualValues(3, lsb)
	t.Assert().EqualValues(0, b)
}
