package main

import (
	"math/bits"
	"strings"
)

type Bitboard uint64

func (b Bitboard) String() string {
	s := strings.Builder{}

	for rank := RankLast; rank >= RankFirst; rank-- {
		if rank < RankLast {
			s.WriteByte('\n')
		}

		s.WriteString(rank.String())
		s.WriteByte(' ')

		for file := FileFirst; file <= FileLast; file++ {
			if file > FileFirst {
				s.WriteByte(' ')
			}

			if b.IsSet(NewSquare(file, rank).Bitboard()) {
				s.WriteByte('X')
			} else {
				s.WriteByte('.')
			}
		}
	}

	s.WriteString("\n  a b c d e f g h")

	return s.String()
}

func (b Bitboard) IsSet(bits Bitboard) bool {
	return b&bits == bits
}

func (b Bitboard) AnySet(bits Bitboard) bool {
	return b&bits > 0
}

func (b *Bitboard) Set(bits Bitboard) {
	*b |= bits
}

func (b Bitboard) SetCount() int {
	return bits.OnesCount64(uint64(b))
}

func (b *Bitboard) Clear(bits Bitboard) {
	*b &= ^bits
}

func (b *Bitboard) PopLSB() Bitboard {
	lsb := Bitboard(bits.TrailingZeros64(uint64(*b)))

	*b &= *b - 1

	return lsb
}
