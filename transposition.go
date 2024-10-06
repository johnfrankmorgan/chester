package main

import (
	"log/slog"
	"unsafe"
)

type Transposition struct {
	Key   Zobrist
	Eval  Eval
	Bound Bound
	Best  Move
	Depth int
}

type Bound byte

const (
	BoundBeta Bound = iota
	BoundAlpha
	BoundExact
)

type TranspositionTable struct {
	entries []Transposition
}

const TranspositionTableSize = 128 * 1024 * 1024

func NewTranspositionTable() *TranspositionTable {
	size := TranspositionTableSize / unsafe.Sizeof(Transposition{})

	slog.Debug("initializing transposition table", "size", size)

	return &TranspositionTable{
		entries: make([]Transposition, size),
	}
}

func (tt *TranspositionTable) index(key Zobrist) Zobrist {
	return key % Zobrist(len(tt.entries))
}

func (tt *TranspositionTable) Get(key Zobrist) (Transposition, bool) {
	if entry := tt.entries[tt.index(key)]; entry.Key == key {
		return entry, true
	}

	return Transposition{}, false
}

func (tt *TranspositionTable) Store(entry Transposition) {
	tt.entries[tt.index(entry.Key)] = entry
}
