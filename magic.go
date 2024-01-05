package main

import (
	"bytes"
	"encoding/gob"
	"io"

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

type Magics [SquareCount * 2]MagicEntry

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
	return m[src].Get(blockers)
}

func (m *Magics) Diagonal(src Square, blockers Bitboard) Bitboard {
	return m[src+SquareCount].Get(blockers)
}
