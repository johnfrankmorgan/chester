package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"slices"

	_ "embed"
)

type MagicEntry struct {
	Moves []Bitboard

	Magic uint64
	Mask  uint64
	Shift uint8
}

func (e MagicEntry) Index(blockers Bitboard) uint64 {
	return ((uint64(blockers) & e.Mask) * e.Magic) >> e.Shift
}

func (e MagicEntry) Get(blockers Bitboard) Bitboard {
	return e.Moves[e.Index(blockers)]
}

var (
	//go:embed embed/magics/orthogonal.json
	_MagicOrthogonalRaw []byte

	//go:embed embed/magics/diagonal.json
	_MagicDiagonalRaw []byte

	_MagicOrthogonal = func() (ret [SquareCount]MagicEntry) {
		if err := json.Unmarshal(_MagicOrthogonalRaw, &ret); err != nil {
			panic(err)
		}

		return ret
	}()

	_MagicDiagonal = func() (ret [SquareCount]MagicEntry) {
		if err := json.Unmarshal(_MagicDiagonalRaw, &ret); err != nil {
			panic(err)
		}

		return ret
	}()
)

func MagicOrthogonalMoves(src Square, blockers Bitboard) Bitboard {
	return _MagicOrthogonal[src].Get(blockers)
}

func MagicDiagonalMoves(src Square, blockers Bitboard) Bitboard {
	return _MagicDiagonal[src].Get(blockers)
}

type MagicGen struct {
	Output     string `help:"Output directory path." default:"embed/magics"`
	Orthogonal bool   `help:"Generate orthogonal magics." default:"true" negatable:""`
	Diagonal   bool   `help:"Generate diagonal magics." default:"true" negatable:""`

	orthogonal [SquareCount]MagicEntry
	diagonal   [SquareCount]MagicEntry
}

func (m *MagicGen) Run(ctx context.Context) error {
	if m.Orthogonal {
		m.run(ctx, slices.Collect(Orthogonals()), &m.orthogonal)
	}

	if m.Diagonal {
		m.run(ctx, slices.Collect(Diagonals()), &m.diagonal)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return m.write()
}

func (m *MagicGen) write() error {
	if err := os.MkdirAll(m.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	entries := make(map[string][SquareCount]MagicEntry)

	if m.Orthogonal {
		entries["orthogonal.json"] = m.orthogonal
	}

	if m.Diagonal {
		entries["diagonal.json"] = m.diagonal
	}

	for name, entry := range entries {
		path := filepath.Join(m.Output, name)

		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to open file for writing: %w", err)
		}

		if err := json.NewEncoder(file).Encode(entry); err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}

	return nil
}

func (m *MagicGen) run(ctx context.Context, dirs []Direction, dest *[SquareCount]MagicEntry) {
	for src := range Squares() {
		if ctx.Err() != nil {
			return
		}

		blockers := m.blockers(dirs, src)

		shift := uint8(blockers.OnesCount())

		dest[src] = m.generate(ctx, dirs, src, shift)
	}
}

func (m *MagicGen) generate(ctx context.Context, dirs []Direction, src Square, shift uint8) MagicEntry {
	for attempt := 0; ; attempt++ {
		if ctx.Err() != nil {
			return MagicEntry{}
		}

		if attempt > 0 && attempt%10000 == 0 {
			slog.Debug("still trying", "attempt", attempt)
		}

		entry := MagicEntry{
			Magic: rand.Uint64() & rand.Uint64() & rand.Uint64(),
			Mask:  uint64(m.blockers(dirs, src)),
			Shift: 64 - shift,
		}

		if moves, ok := m.attempt(entry, dirs, src); ok {
			slog.Info("success", "attempt", attempt, "src", src, "magic", entry.Magic)

			entry.Moves = moves
			return entry
		}
	}
}

func (m *MagicGen) attempt(entry MagicEntry, dirs []Direction, src Square) ([]Bitboard, bool) {
	moves := make([]Bitboard, 1<<(64-entry.Shift))

	for blockers := Bitboard(0); ; {
		if moves[entry.Index(blockers)] != 0 {
			return nil, false
		}

		moves[entry.Index(blockers)] = m.moves(dirs, src, blockers)

		if blockers = (blockers - 1) & Bitboard(entry.Mask); blockers == 0 {
			break
		}
	}

	return moves, true
}

func (m *MagicGen) moves(dirs []Direction, src Square, blockers Bitboard) Bitboard {
	var moves Bitboard

	for _, dir := range dirs {
		for mul := Square(1); mul <= SquaresToEdge(src, dir); mul++ {
			dst := src + dir.Offset()*mul

			moves = moves.Occupy(dst)

			if blockers.IsOccupied(dst) {
				break
			}
		}
	}

	return moves
}

func (m *MagicGen) blockers(dirs []Direction, src Square) Bitboard {
	var blockers Bitboard

	for _, dir := range dirs {
		for mul := Square(1); mul < SquaresToEdge(src, dir); mul++ {
			blockers = blockers.Occupy(src + dir.Offset()*mul)
		}
	}

	return blockers
}
