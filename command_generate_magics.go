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
	King       bool   `default:"true" negatable:"true" help:"Generate king moves."`
	Knight     bool   `default:"true" negatable:"true" help:"Generate king moves."`
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

	for _, dirs := range [][]Direction{DirectionsOrthogonal[:], DirectionsDiagonal[:]} {
		if !cmd.Orthogonal && !dirs[0].IsDiagonal() {
			continue
		} else if !cmd.Diagonal && dirs[0].IsDiagonal() {
			continue
		}

		dest := Ternary(dirs[0].IsDiagonal(), &magics.diagonal, &magics.orthogonal)
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

			dest[src] = entry

			slog.Debug("generated", "kind", kind, "square", src, "attempts", attempts, "duration", duration)

			total.attempts += attempts
			total.duration += duration
		}

		slog.Info("generated", "kind", kind, "attempts", total.attempts, "duration", total.duration)
	}

	if cmd.King {
		slog.Info("generating", "kind", "king")

		start := time.Now()

		for src := SquareFirst; src <= SquareLast; src++ {
			moves := Bitboard(0)

			for _, dir := range Directions {
				if dir.ToEdge(src) != 0 {
					dst := src + dir.Offset()

					moves.Set(dst.Bitboard())
				}
			}

			magics.king[src] = moves
		}

		slog.Info("generated", "kind", "king", "duration", time.Since(start))
	}

	if cmd.Knight {
		slog.Info("generating", "kind", "knight")

		start := time.Now()
		jumps := [...]Square{
			DirectionNorth.Offset()*2 + DirectionEast.Offset(),
			DirectionNorth.Offset()*2 + DirectionWest.Offset(),
			DirectionSouth.Offset()*2 + DirectionEast.Offset(),
			DirectionSouth.Offset()*2 + DirectionWest.Offset(),
			DirectionEast.Offset()*2 + DirectionNorth.Offset(),
			DirectionEast.Offset()*2 + DirectionSouth.Offset(),
			DirectionWest.Offset()*2 + DirectionNorth.Offset(),
			DirectionWest.Offset()*2 + DirectionSouth.Offset(),
		}

		for src := SquareFirst; src <= SquareLast; src++ {
			moves := Bitboard(0)

			for _, jump := range jumps {
				dst := src + jump

				if !dst.Valid() {
					continue
				}

				if Abs(src.File()-dst.File()) > 2 {
					continue
				}

				if Abs(src.Rank()-dst.Rank()) > 2 {
					continue
				}

				moves.Set(dst.Bitboard())
			}

			magics.knight[src] = moves
		}

		slog.Info("generated", "kind", "knight", "duration", time.Since(start))
	}

	return magics
}

// https://www.chessprogramming.org/Magic_Bitboards
// https://www.chessprogramming.org/Looking_for_Magics
// https://analog-hors.github.io/site/magic-bitboards/
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
