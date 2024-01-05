package main

import "strings"

type Piece uint8

const (
	PieceNone Piece = 0

	PieceBlackPawn   = _PieceBlack | Piece(PieceKindPawn)
	PieceBlackKnight = _PieceBlack | Piece(PieceKindKnight)
	PieceBlackBishop = _PieceBlack | Piece(PieceKindBishop)
	PieceBlackRook   = _PieceBlack | Piece(PieceKindRook)
	PieceBlackQueen  = _PieceBlack | Piece(PieceKindQueen)
	PieceBlackKing   = _PieceBlack | Piece(PieceKindKing)

	PieceWhitePawn   = _PieceWhite | Piece(PieceKindPawn)
	PieceWhiteKnight = _PieceWhite | Piece(PieceKindKnight)
	PieceWhiteBishop = _PieceWhite | Piece(PieceKindBishop)
	PieceWhiteRook   = _PieceWhite | Piece(PieceKindRook)
	PieceWhiteQueen  = _PieceWhite | Piece(PieceKindQueen)
	PieceWhiteKing   = _PieceWhite | Piece(PieceKindKing)

	_PieceBlack = Piece(ColorBlack << 3)
	_PieceWhite = Piece(ColorWhite << 3)

	_PieceColorMask = _PieceBlack | _PieceWhite
)

type PieceKind uint8

const (
	PieceKindNone PieceKind = iota

	PieceKindPawn
	PieceKindKnight
	PieceKindBishop
	PieceKindRook
	PieceKindQueen
	PieceKindKing

	PieceKindCount = 7 // includes "PieceNone"
)

type Color uint8

const (
	ColorBlack Color = iota
	ColorWhite

	ColorCount = 2
)

func NewPiece(color Color, kind PieceKind) Piece {
	return Piece(color<<3) | Piece(kind)
}

func (p Piece) String() string {
	if p.Color() == ColorWhite {
		return strings.ToUpper(p.Kind().String())
	}

	return p.Kind().String()
}

func (p Piece) Color() Color {
	return Color(p&_PieceColorMask) >> 3
}

func (p Piece) Kind() PieceKind {
	return PieceKind(p &^ _PieceColorMask)
}

func (p Piece) Is(kind PieceKind) bool {
	return p.Kind() == kind
}

func (pk PieceKind) String() string {
	switch pk {
	case PieceKindPawn:
		return "p"

	case PieceKindKnight:
		return "n"

	case PieceKindBishop:
		return "b"

	case PieceKindRook:
		return "r"

	case PieceKindQueen:
		return "q"

	case PieceKindKing:
		return "k"

	default:
		return UnknownNumeric(pk)
	}
}

func (c Color) String() string {
	switch c {
	case ColorBlack:
		return "b"

	case ColorWhite:
		return "w"

	default:
		return UnknownNumeric(c)
	}
}

func (c Color) Valid() bool {
	return c == ColorBlack || c == ColorWhite
}

func (c Color) Opponent() Color {
	return c ^ ColorWhite
}
