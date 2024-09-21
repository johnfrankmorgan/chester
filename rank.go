package main

import "iter"

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

func RankFromString(s string) (Rank, bool) {
	if len(s) != 1 || s[0] < '1' || s[0] > '8' {
		return 0, false
	}

	return Rank(s[0] - '1'), true
}

func (r Rank) String() string {
	if r >= RankFirst && r <= RankLast {
		return string('1' + byte(r))
	}

	return repr(r)
}

func Ranks() iter.Seq[Rank] {
	return func(yield func(Rank) bool) {
		for rank := RankFirst; rank <= RankLast; rank++ {
			if !yield(rank) {
				break
			}
		}
	}
}

func RanksReversed() iter.Seq[Rank] {
	return func(yield func(Rank) bool) {
		for rank := RankLast; rank >= RankFirst; rank-- {
			if !yield(rank) {
				break
			}
		}
	}
}
