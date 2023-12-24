package main

import "strings"

type Piece uint8

const (
	_PieceBlack = Piece(ColorBlack << 3)
	_PieceWhite = Piece(ColorWhite << 3)

	_PieceColorMask = _PieceWhite

	PieceEmpty Piece = 0

	PieceBlackPawn   = _PieceBlack | Piece(PiecePawn)
	PieceBlackKnight = _PieceBlack | Piece(PieceKnight)
	PieceBlackBishop = _PieceBlack | Piece(PieceBishop)
	PieceBlackRook   = _PieceBlack | Piece(PieceRook)
	PieceBlackQueen  = _PieceBlack | Piece(PieceQueen)
	PieceBlackKing   = _PieceBlack | Piece(PieceKing)

	PieceWhitePawn   = _PieceWhite | Piece(PiecePawn)
	PieceWhiteKnight = _PieceWhite | Piece(PieceKnight)
	PieceWhiteBishop = _PieceWhite | Piece(PieceBishop)
	PieceWhiteRook   = _PieceWhite | Piece(PieceRook)
	PieceWhiteQueen  = _PieceWhite | Piece(PieceQueen)
	PieceWhiteKing   = _PieceWhite | Piece(PieceKing)
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
