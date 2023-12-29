package main

import "log/slog"

type Transpositions struct {
	_entries []Transposition
	_enabled bool
}

type Transposition struct {
	Zobrist    uint64
	Type       TranspositionType
	Depth      int
	Evaluation int
	Move       Move
}

type TranspositionType int

const (
	TranspositionExact TranspositionType = iota
	TranspositionLower
	TranspositionUpper
)

func NewTranspositions(size int, enabled bool) *Transpositions {
	return &Transpositions{
		_entries: make([]Transposition, size),
		_enabled: enabled,
	}
}

func (t *Transpositions) Store(board *Board, depth, eval int, typ TranspositionType, move Move) {
	if !t._enabled {
		return
	}

	index := board.Zobrist() % uint64(len(t._entries))

	if zobrist := t._entries[index].Zobrist; zobrist != 0 && zobrist != board.Zobrist() {
		slog.Debug("overwriting transposition", "old", zobrist, "new", board.Zobrist())
	}

	t._entries[board.Zobrist()%uint64(len(t._entries))] = Transposition{
		Zobrist:    board.Zobrist(),
		Type:       typ,
		Depth:      depth,
		Evaluation: eval,
		Move:       move,
	}
}

func (t *Transpositions) Get(board *Board, depth, alpha, beta int) (Transposition, bool) {
	if !t._enabled {
		return Transposition{}, false
	}

	entry := t._entries[board.Zobrist()%uint64(len(t._entries))]

	if entry.Zobrist == board.Zobrist() && entry.Depth >= depth {
		if entry.Type == TranspositionExact {
			return entry, true
		} else if entry.Type == TranspositionUpper && entry.Evaluation <= alpha {
			return entry, true
		} else if entry.Type == TranspositionLower && entry.Evaluation >= beta {
			return entry, true
		}
	}

	return Transposition{}, false
}

func (t *Transpositions) Move(board *Board) (Move, bool) {
	if !t._enabled {
		return Move{}, false
	}

	entry := t._entries[board.Zobrist()%uint64(len(t._entries))]

	return entry.Move, !entry.Move.IsZero()
}
