package main

import (
	"strings"
)

type Piece uint8

const (
	EmptySquare Piece = 0

	BlackPawn   = _Black | Piece(Pawn)
	BlackKnight = _Black | Piece(Knight)
	BlackBishop = _Black | Piece(Bishop)
	BlackRook   = _Black | Piece(Rook)
	BlackQueen  = _Black | Piece(Queen)
	BlackKing   = _Black | Piece(King)

	WhitePawn   = _White | Piece(Pawn)
	WhiteKnight = _White | Piece(Knight)
	WhiteBishop = _White | Piece(Bishop)
	WhiteRook   = _White | Piece(Rook)
	WhiteQueen  = _White | Piece(Queen)
	WhiteKing   = _White | Piece(King)

	_Black = Piece(Black << 3)
	_White = Piece(White << 3)
)

func NewPiece(c Color, k PieceType) Piece {
	return Piece(c<<3) | Piece(k)
}

func PieceFromString(s string) (Piece, bool) {
	switch s {
	case "p":
		return BlackPawn, true

	case "n":
		return BlackKnight, true

	case "b":
		return BlackBishop, true

	case "r":
		return BlackRook, true

	case "q":
		return BlackQueen, true

	case "k":
		return BlackKing, true

	case "P":
		return WhitePawn, true

	case "N":
		return WhiteKnight, true

	case "B":
		return WhiteBishop, true

	case "R":
		return WhiteRook, true

	case "Q":
		return WhiteQueen, true

	case "K":
		return WhiteKing, true

	default:
		return 0, false
	}
}

func (p Piece) String() string {
	s := p.Type().String()

	if len(s) != 1 {
		return repr(p)
	}

	if p.Color() == White {
		s = strings.ToUpper(s)
	}

	return s
}

func (p Piece) Type() PieceType {
	return PieceType(p & 0b111)
}

func (p Piece) Color() Color {
	return Color(p >> 3)
}

type PieceType uint8

const (
	Pawn PieceType = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King

	PieceTypeCount = 6
)

func (t PieceType) String() string {
	types := [...]string{"p", "n", "b", "r", "q", "k"}

	if t != 0 && int(t) <= len(types) {
		return types[t-1]
	}

	return repr(t)
}
