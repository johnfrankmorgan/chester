package main

import (
	"math"
	"math/bits"
	"strings"
)

type Bitboard uint64

const (
	BitboardFileA Bitboard = 0x101010101010101
	BitboardFileB          = BitboardFileA << 1
	BitboardFileC          = BitboardFileB << 1
	BitboardFileD          = BitboardFileC << 1
	BitboardFileE          = BitboardFileD << 1
	BitboardFileF          = BitboardFileE << 1
	BitboardFileG          = BitboardFileF << 1
	BitboardFileH          = BitboardFileG << 1

	BitboardRank1 Bitboard = 0b11111111
	BitboardRank2          = BitboardRank1 << 8
	BitboardRank3          = BitboardRank2 << 8
	BitboardRank4          = BitboardRank3 << 8
	BitboardRank5          = BitboardRank4 << 8
	BitboardRank6          = BitboardRank5 << 8
	BitboardRank7          = BitboardRank6 << 8
	BitboardRank8          = BitboardRank7 << 8

	BitboardAll Bitboard = math.MaxUint64
)

var (
	BitboardFiles = [FileCount]Bitboard{
		FileA: BitboardFileA,
		FileB: BitboardFileB,
		FileC: BitboardFileC,
		FileD: BitboardFileD,
		FileE: BitboardFileE,
		FileF: BitboardFileF,
		FileG: BitboardFileG,
		FileH: BitboardFileH,
	}

	BitboardRanks = [RankCount]Bitboard{
		Rank1: BitboardRank1,
		Rank2: BitboardRank2,
		Rank3: BitboardRank3,
		Rank4: BitboardRank4,
		Rank5: BitboardRank5,
		Rank6: BitboardRank6,
		Rank7: BitboardRank7,
		Rank8: BitboardRank8,
	}

	BitboardsCastle = ColorTable[struct {
		Kingside struct {
			Attackers, Blockers Bitboard
		}

		Queenside struct {
			Attackers, Blockers Bitboard
		}
	}]{
		ColorWhite: {
			Kingside: struct{ Attackers, Blockers Bitboard }{
				Attackers: SquareF1.Bitboard() | SquareG1.Bitboard(),
				Blockers:  SquareF1.Bitboard() | SquareG1.Bitboard(),
			},

			Queenside: struct{ Attackers, Blockers Bitboard }{
				Attackers: SquareC1.Bitboard() | SquareD1.Bitboard(),
				Blockers:  SquareB1.Bitboard() | SquareC1.Bitboard() | SquareD1.Bitboard(),
			},
		},

		ColorBlack: {
			Kingside: struct{ Attackers, Blockers Bitboard }{
				Attackers: SquareF8.Bitboard() | SquareG8.Bitboard(),
				Blockers:  SquareF8.Bitboard() | SquareG8.Bitboard(),
			},

			Queenside: struct{ Attackers, Blockers Bitboard }{
				Attackers: SquareC8.Bitboard() | SquareD8.Bitboard(),
				Blockers:  SquareB8.Bitboard() | SquareC8.Bitboard() | SquareD8.Bitboard(),
			},
		},
	}
)

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
