package magic

import (
	"bytes"
	"encoding/gob"

	"github.com/johnfrankmorgan/chester/internal/chester/bb"
	"github.com/johnfrankmorgan/chester/internal/chester/square"
	"github.com/johnfrankmorgan/chester/internal/chester/util"

	_ "embed"
)

func Orthogonal(src square.Square, blockers bb.Bitboard) bb.Bitboard {
	return entries.Orthogonal[src].Get(blockers)
}

func Diagonal(src square.Square, blockers bb.Bitboard) bb.Bitboard {
	return entries.Diagonal[src].Get(blockers)
}

var (
	//go:embed entries.gob
	_entries []byte

	entries Entries
)

func init() {
	util.Check(gob.NewDecoder(bytes.NewReader(_entries)).Decode(&entries))
}

type Entry struct {
	Magic uint64
	Mask  bb.Bitboard
	Shift uint8
	Moves []bb.Bitboard
}

type Entries struct {
	Orthogonal [square.Count]Entry
	Diagonal   [square.Count]Entry
}

func (e Entry) index(blockers bb.Bitboard) uint64 {
	blockers &= e.Mask

	hash := uint64(blockers) * e.Magic

	return hash >> e.Shift
}

func (e Entry) Get(blockers bb.Bitboard) bb.Bitboard {
	return e.Moves[e.index(blockers)]
}
