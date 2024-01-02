package magic

import (
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/johnfrankmorgan/chester/internal/chester/bb"
	"github.com/johnfrankmorgan/chester/internal/chester/square"
	"github.com/johnfrankmorgan/chester/internal/chester/util"
)

func Generate() Entries {
	entries := Entries{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	for _, slider := range []slider{{square.DirectionsOrthogonal}, {square.DirectionsDiagonal}} {
		slider := slider

		go func() {
			defer wg.Done()

			start := time.Now()
			attempts := 0

			for src := square.First; src <= square.Last; src++ {
				blockers := slider.blockers(src)
				shift := uint8(blockers.SetCount())

				entry, tries := generate(slider, src, shift)

				if slider.dirs[0].IsDiagonal() {
					entries.Diagonal[src] = entry
				} else {
					entries.Orthogonal[src] = entry
				}

				attempts += tries
			}

			slog.Info("generated magics",
				"kind", util.Ternary(slider.dirs[0].IsDiagonal(), "diagonal", "orthogonal"),
				"attempts", attempts,
				"duration", time.Since(start))
		}()
	}

	wg.Wait()

	return entries
}

// https://www.chessprogramming.org/Magic_Bitboards
// https://www.chessprogramming.org/Looking_for_Magics
// https://analog-hors.github.io/site/magic-bitboards/
func generate(slider slider, src square.Square, shift uint8) (Entry, int) {
	for attempts := 1; ; attempts++ {
		entry := Entry{
			Magic: rand.Uint64() & rand.Uint64() & rand.Uint64(),
			Mask:  slider.blockers(src),
			Shift: 64 - shift,
		}

		if moves, ok := try(entry, slider, src); ok {
			entry.Moves = moves

			return entry, attempts
		}
	}
}

func try(entry Entry, slider slider, src square.Square) ([]bb.Bitboard, bool) {
	shift := 64 - entry.Shift

	moves := make([]bb.Bitboard, 1<<shift)

	for blockers := bb.Bitboard(0); ; {
		if moves[entry.index(blockers)] != 0 {
			return nil, false
		}

		moves[entry.index(blockers)] = slider.moves(src, blockers)

		if blockers = (blockers - 1) & entry.Mask; blockers == 0 {
			break
		}
	}

	return moves, true
}

type slider struct {
	dirs [4]square.Direction
}

func (s slider) moves(src square.Square, blockers bb.Bitboard) bb.Bitboard {
	moves := bb.Bitboard(0)

	for _, dir := range s.dirs {
		for mul := square.Square(1); mul <= dir.ToEdge(src); mul++ {
			dst := src + dir.Offset()*mul

			if blockers.IsSet(dst.Bitboard()) {
				break
			}

			moves.Set(dst.Bitboard())
		}
	}

	return moves
}

func (s slider) blockers(src square.Square) bb.Bitboard {
	blockers := bb.Bitboard(0)

	for _, dir := range s.dirs {
		for mul := square.Square(1); mul < dir.ToEdge(src); mul++ {
			dst := src + dir.Offset()*mul

			blockers.Set(dst.Bitboard())
		}
	}

	return blockers
}
