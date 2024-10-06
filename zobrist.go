package main

import (
	"math/rand"
	"unsafe"
)

type Zobrist uint64

var Zobrists struct {
	Players   [ColorCount]Zobrist
	Castling  [16]Zobrist
	Pieces    [ColorCount][PieceTypeCount + 1][SquareCount]Zobrist
	EnPassant [SquareCount]Zobrist
}

func init() {
	rand := rand.New(rand.NewSource(0xdeadbeef)) // use constant seed for reproducibility

	for off := uintptr(0); off < unsafe.Sizeof(Zobrists); off += unsafe.Sizeof(uint64(0)) {
		*(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&Zobrists)) + off)) = rand.Uint64()
	}
}

func CastlingZobristIndex(b *Board) int {
	index := 0

	if b.Castling[White].Kingside {
		index |= 1
	}

	if b.Castling[White].Queenside {
		index |= 2
	}

	if b.Castling[Black].Kingside {
		index |= 4
	}

	if b.Castling[Black].Queenside {
		index |= 8
	}

	return index
}

func CalculateZobrist(b *Board) Zobrist {
	zobrist := Zobrist(0)

	zobrist ^= Zobrists.Players[b.Player]
	zobrist ^= Zobrists.Castling[CastlingZobristIndex(b)]

	for src := range Squares() {
		piece := b.Squares[src]
		if piece == EmptySquare {
			continue
		}

		zobrist ^= Zobrists.Pieces[piece.Color()][piece.Type()][src]
	}

	zobrist ^= Zobrists.EnPassant[b.EnPassant]

	return zobrist
}
