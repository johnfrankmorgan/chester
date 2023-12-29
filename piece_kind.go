package main

type PieceKind uint8

const (
	PieceNone PieceKind = iota

	PiecePawn
	PieceKnight
	PieceBishop
	PieceRook
	PieceQueen
	PieceKing

	PieceCount = 6

	_PieceCountIncludingNone = PieceCount + 1
)

var _PieceKindValues = [...]int{
	PiecePawn:   100,
	PieceKnight: 320,
	PieceBishop: 330,
	PieceRook:   500,
	PieceQueen:  900,
	PieceKing:   20000,
}

func (pk PieceKind) String() string {
	switch pk {
	case PiecePawn:
		return "p"

	case PieceKnight:
		return "n"

	case PieceBishop:
		return "b"

	case PieceRook:
		return "r"

	case PieceQueen:
		return "q"

	case PieceKing:
		return "k"

	default:
		return istr(pk)
	}
}

func (pk PieceKind) Value() int {
	return _PieceKindValues[pk]
}
