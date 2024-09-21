package main

import (
	"iter"
	"math/bits"
	"strings"
)

type Bitboard uint64

func BitboardForFile(file File) Bitboard {
	return 0x101010101010101 << file
}

func BitboardForRank(rank Rank) Bitboard {
	return 0b11111111 << (rank * FileCount)
}

var BitboardInDirection = func() func(Square, Direction) Bitboard {
	var lookup [SquareCount][DirectionCount]Bitboard

	for src := range Squares() {
		for dir := range Directions() {
			for mul := Square(1); mul <= SquaresToEdge(src, dir); mul++ {
				dst := src + dir.Offset()*mul

				lookup[src][dir] = lookup[src][dir].Occupy(dst)
			}
		}
	}

	return func(src Square, dir Direction) Bitboard {
		return lookup[src][dir]
	}
}()

var BitboardAlignedAlong = func() func(Square, Square) Bitboard {
	var lookup [SquareCount][SquareCount]Bitboard

	for src := range Squares() {
		for dir := range Directions() {
			mask := src.Bitboard().
				Set(BitboardInDirection(src, dir)).
				Set(BitboardInDirection(src, dir.Opposite()))

			for dst := range mask.Occupied() {
				lookup[src][dst] = mask
			}
		}
	}

	return func(src, dst Square) Bitboard {
		return lookup[src][dst]
	}
}()

func (b Bitboard) String() string {
	s := strings.Builder{}

	for rank := range RanksReversed() {
		if rank < RankLast {
			s.WriteByte('\n')
		}

		s.WriteString(rank.String())
		s.WriteByte(' ')

		for file := range Files() {
			if file > FileFirst {
				s.WriteByte(' ')
			}

			if b.IsOccupied(NewSquare(file, rank)) {
				s.WriteByte('X')
			} else {
				s.WriteByte('.')
			}
		}
	}

	s.WriteString("\n  a b c d e f g h")

	return s.String()
}

func (b Bitboard) And(other Bitboard) Bitboard {
	return b & other
}

func (b Bitboard) Set(bits Bitboard) Bitboard {
	return b | bits
}

func (b Bitboard) IsSet(bits Bitboard) bool {
	return b&bits == bits
}

func (b Bitboard) AnySet(bits Bitboard) bool {
	return b&bits != 0
}

func (b Bitboard) Clear(bits Bitboard) Bitboard {
	return b &^ bits
}

func (b Bitboard) Occupy(sq Square) Bitboard {
	return b.Set(sq.Bitboard())
}

func (b Bitboard) IsOccupied(sq Square) bool {
	return b.IsSet(sq.Bitboard())
}

func (b Bitboard) Unoccupy(sq Square) Bitboard {
	return b.Clear(sq.Bitboard())
}

func (b Bitboard) Occupied() iter.Seq[Square] {
	return func(yield func(Square) bool) {
		sq := Square(0)

		for b != 0 {
			b, sq = b.PopLSB()

			if !yield(sq) {
				break
			}
		}
	}
}

func (b Bitboard) OnesCount() int {
	return bits.OnesCount64(uint64(b))
}

func (b Bitboard) PopLSB() (Bitboard, Square) {
	lsb := Bitboard(bits.TrailingZeros64(uint64(b)))

	b &= b - 1

	return b, Square(lsb)
}
