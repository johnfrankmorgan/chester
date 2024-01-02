package piece

import "github.com/johnfrankmorgan/chester/internal/chester/util"

type Kind uint8

const (
	None Kind = iota

	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	Count = 7 // includes "None"
)

func (k Kind) String() string {
	switch k {
	case Pawn:
		return "p"

	case Knight:
		return "n"

	case Bishop:
		return "b"

	case Rook:
		return "r"

	case Queen:
		return "q"

	case King:
		return "k"
	}

	return util.UnknownNumeric(k)
}
