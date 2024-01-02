package square

import "github.com/johnfrankmorgan/chester/internal/chester/util"

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

func (r Rank) String() string {
	if r.Valid() {
		return string('1' + byte(r))
	}

	return util.UnknownNumeric(r)
}

func (r Rank) Valid() bool {
	return r >= RankFirst && r <= RankLast
}
