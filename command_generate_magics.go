package main

import (
	"encoding/gob"
	"math/rand"
	"os"
	"time"

	"log/slog"
)

type CommandGenerateMagics struct {
	Output     string `arg:"" type:"path" help:"Output path."`
	Orthogonal bool   `default:"true" negatable:"true" help:"Generate magics for orthogonal moves."`
	Diagonal   bool   `default:"true" negatable:"true" help:"Generate magics for diagonal moves."`
}

func (cmd CommandGenerateMagics) Run() error {
	f, err := os.Create(cmd.Output)
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(f).Encode(cmd.run()); err != nil {
		return err
	}

	return f.Close()
}

func (cmd CommandGenerateMagics) run() Magics {
	magics := Magics{}

	for mul, dirs := range [][]Direction{DirectionsOrthogonal[:], DirectionsDiagonal[:]} {
		if !cmd.Orthogonal && !dirs[0].IsDiagonal() {
			continue
		} else if !cmd.Diagonal && dirs[0].IsDiagonal() {
			continue
		}

		kind := Ternary(dirs[0].IsDiagonal(), "diagonal", "orthogonal")
		total := struct {
			attempts int
			duration time.Duration
		}{}

		slog.Info("generating", "kind", kind)

		for src := SquareFirst; src <= SquareLast; src++ {
			blockers := cmd.blockers(dirs, src)
			shift := uint8(blockers.SetCount())
			start := time.Now()

			entry, attempts := cmd.generate(dirs, src, shift)

			duration := time.Since(start)

			magics[src+Square(mul*SquareCount)] = entry

			slog.Debug("generated", "kind", kind, "square", src, "attempts", attempts, "duration", duration)

			total.attempts += attempts
			total.duration += duration
		}

		slog.Info("generated", "kind", kind, "attempts", total.attempts, "duration", total.duration)
	}

	return magics
}

func (cmd CommandGenerateMagics) generate(dirs []Direction, src Square, shift uint8) (MagicEntry, int) {
	for attempts := 1; ; attempts++ {
		entry := MagicEntry{
			Magic: rand.Uint64() & rand.Uint64() & rand.Uint64(),
			Mask:  cmd.blockers(dirs, src),
			Shift: 64 - shift,
		}

		if moves, ok := cmd.attempt(entry, dirs, src); ok {
			entry.Moves = moves

			return entry, attempts
		}
	}
}

func (cmd CommandGenerateMagics) attempt(entry MagicEntry, dirs []Direction, src Square) ([]Bitboard, bool) {
	shift := 64 - entry.Shift

	moves := make([]Bitboard, 1<<shift)

	for blockers := Bitboard(0); ; {
		if moves[entry.Index(blockers)] != 0 {
			return nil, false
		}

		moves[entry.Index(blockers)] = cmd.moves(dirs, src, blockers)

		if blockers = (blockers - 1) & entry.Mask; blockers == 0 {
			break
		}
	}

	return moves, true
}

func (cmd CommandGenerateMagics) moves(dirs []Direction, src Square, blockers Bitboard) Bitboard {
	moves := Bitboard(0)

	for _, dir := range dirs {
		for mul := Square(1); mul <= dir.ToEdge(src); mul++ {
			dst := src + dir.Offset()*mul

			moves.Set(dst.Bitboard())

			if blockers.IsSet(dst.Bitboard()) {
				break
			}
		}
	}

	return moves
}

func (cmd CommandGenerateMagics) blockers(dirs []Direction, src Square) Bitboard {
	blockers := Bitboard(0)

	for _, dir := range dirs {
		for mul := Square(1); mul < dir.ToEdge(src); mul++ {
			dst := src + dir.Offset()*mul

			blockers.Set(dst.Bitboard())
		}
	}

	return blockers
}
