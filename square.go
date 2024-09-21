package main

import (
	"fmt"
	"iter"
)

type Square int8

const (
	SquareA1 Square = iota
	SquareB1
	SquareC1
	SquareD1
	SquareE1
	SquareF1
	SquareG1
	SquareH1

	SquareA2
	SquareB2
	SquareC2
	SquareD2
	SquareE2
	SquareF2
	SquareG2
	SquareH2

	SquareA3
	SquareB3
	SquareC3
	SquareD3
	SquareE3
	SquareF3
	SquareG3
	SquareH3

	SquareA4
	SquareB4
	SquareC4
	SquareD4
	SquareE4
	SquareF4
	SquareG4
	SquareH4

	SquareA5
	SquareB5
	SquareC5
	SquareD5
	SquareE5
	SquareF5
	SquareG5
	SquareH5

	SquareA6
	SquareB6
	SquareC6
	SquareD6
	SquareE6
	SquareF6
	SquareG6
	SquareH6

	SquareA7
	SquareB7
	SquareC7
	SquareD7
	SquareE7
	SquareF7
	SquareG7
	SquareH7

	SquareA8
	SquareB8
	SquareC8
	SquareD8
	SquareE8
	SquareF8
	SquareG8
	SquareH8

	SquareFirst = SquareA1
	SquareLast  = SquareH8

	SquareCount = FileCount * RankCount
)

func NewSquare(file File, rank Rank) Square {
	return Square(rank*8) + Square(file)
}

func SquareFromString(s string) (Square, bool) {
	if len(s) != 2 {
		return 0, false
	}

	file, ok := FileFromString(s[:1])
	if !ok {
		return 0, false
	}

	rank, ok := RankFromString(s[1:])
	if !ok {
		return 0, false
	}

	return NewSquare(file, rank), true
}

func (s Square) String() string {
	if s.Valid() {
		return fmt.Sprintf("%s%s", s.File(), s.Rank())
	}

	return repr(s)
}

func (s Square) File() File {
	return File(s % FileCount)
}

func (s Square) Rank() Rank {
	return Rank(s / FileCount)
}

func (s Square) Bitboard() Bitboard {
	return 1 << s
}

func (s Square) Valid() bool {
	return s >= SquareFirst && s <= SquareLast
}

func Squares() iter.Seq[Square] {
	return func(yield func(Square) bool) {
		for square := SquareFirst; square <= SquareLast; square++ {
			if !yield(square) {
				break
			}
		}
	}
}

var SquaresToEdge = func() func(Square, Direction) Square {
	var lookup [SquareCount][DirectionCount]Square

	for src := range Squares() {
		file := Square(src.File())
		rank := Square(src.Rank())

		north := RankCount - rank - 1
		south := rank
		east := FileCount - file - 1
		west := file

		lookup[src] = [DirectionCount]Square{
			North:     north,
			South:     south,
			East:      east,
			West:      west,
			NorthEast: min(north, east),
			SouthEast: min(south, east),
			NorthWest: min(north, west),
			SouthWest: min(south, west),
		}
	}

	return func(src Square, dir Direction) Square {
		return lookup[src][dir]
	}
}()
