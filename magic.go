package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"unsafe"

	_ "embed"
)

type MagicEntry struct {
	Magic uint64
	Mask  Bitboard
	Shift uint8
	Moves []Bitboard
}

func (e MagicEntry) Index(blockers Bitboard) uint64 {
	blockers &= e.Mask

	hash := uint64(blockers) * e.Magic

	return hash >> e.Shift
}

func (e MagicEntry) Get(blockers Bitboard) Bitboard {
	return e.Moves[e.Index(blockers)]
}

type Magics struct {
	orthogonal [SquareCount]MagicEntry
	diagonal   [SquareCount]MagicEntry
	king       [SquareCount]Bitboard
	knight     [SquareCount]Bitboard
	pawn       [ColorCount][SquareCount]Bitboard
}

func (m Magics) GobEncode() ([]byte, error) {
	e := *(*struct {
		Orthogonal, Diagonal [SquareCount]MagicEntry
		King, Knight         [SquareCount]Bitboard
		Pawn                 [ColorCount][SquareCount]Bitboard
	})(unsafe.Pointer(&m))

	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(e)

	return b.Bytes(), err
}

func (m *Magics) GobDecode(b []byte) error {
	e := struct {
		Orthogonal, Diagonal [SquareCount]MagicEntry
		King, Knight         [SquareCount]Bitboard
		Pawn                 [ColorCount][SquareCount]Bitboard
	}{}

	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&e)
	*m = *(*Magics)(unsafe.Pointer(&e))

	return err
}

var (
	Magic Magics

	//go:embed magic.gob
	_magraw []byte
)

func init() {
	PanicIfError(
		Magic.LoadDefault(),
	)
}

func (m *Magics) Load(r io.Reader) error {
	return gob.NewDecoder(r).Decode(m)
}

func (m *Magics) LoadDefault() error {
	return m.Load(bytes.NewReader(_magraw))
}

func (m *Magics) Orthogonal(src Square, blockers Bitboard) Bitboard {
	return m.orthogonal[src].Get(blockers)
}

func (m *Magics) Diagonal(src Square, blockers Bitboard) Bitboard {
	return m.diagonal[src].Get(blockers)
}

func (m *Magics) King(src Square) Bitboard {
	return m.king[src]
}

func (m *Magics) Knight(src Square) Bitboard {
	return m.knight[src]
}

func (m *Magics) PawnAttacks(color Color, src Square) Bitboard {
	return m.pawn[color][src]
}
