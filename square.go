package main

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

var (
	_SquareAlignMasks [SquareCount][SquareCount]Bitboard
)

type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH

	FileFirst = FileA
	FileLast  = FileH

	FileCount = 8
)

type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8

	RankFirst = Rank1
	RankLast  = Rank8

	RankCount = 8
)

func init() {
	for src := SquareFirst; src <= SquareLast; src++ {
		for _, dir := range Directions {
			mask := dir.Mask(src)
			alignment := mask

			for mask > 0 {
				dst := Square(mask.PopLSB())

				_SquareAlignMasks[src][dst] = alignment
			}
		}
	}
}

func NewSquare(file File, rank Rank) Square {
	return Square(rank*FileCount) + Square(file)
}

func (s Square) String() string {
	if s.Valid() {
		return s.File().String() + s.Rank().String()
	}

	return UnknownNumeric(s)
}

func (s Square) Valid() bool {
	return s >= SquareFirst && s <= SquareLast
}

func (s Square) File() File {
	return File(s & 0b111)
}

func (s Square) Rank() Rank {
	return Rank(s >> 3)
}

func (s Square) Bitboard() Bitboard {
	return 1 << Bitboard(s)
}

func (s Square) AlignMask(other Square) Bitboard {
	return _SquareAlignMasks[s][other]
}

func (f File) String() string {
	if f.Valid() {
		return string(byte(f) + 'a')
	}

	return UnknownNumeric(f)
}

func (f File) Valid() bool {
	return f >= FileFirst && f <= FileLast
}

func (r Rank) String() string {
	if r.Valid() {
		return string(byte(r) + '1')
	}

	return UnknownNumeric(r)
}

func (r Rank) Valid() bool {
	return r >= RankFirst && r <= RankLast
}
