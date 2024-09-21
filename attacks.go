package main

import "math"

type Attacks struct {
	All       Bitboard
	Checks    int
	CheckRays Bitboard
	Pins      Bitboard
}

func GenerateAttacks(b *Board, attacker Color) Attacks {
	a := Attacks{}

	GenerateKingAttacks(b, &a, attacker)
	GenerateSlidingAttacks(b, &a, attacker)
	GenerateKnightAttacks(b, &a, attacker)
	GeneratePawnAttacks(b, &a, attacker)

	if a.Checks == 0 {
		a.CheckRays = math.MaxUint64
	}

	return a
}

func (a Attacks) IsAttacked(sq Square) bool {
	return a.All.IsOccupied(sq)
}

func (a Attacks) IsPinned(sq Square) bool {
	return a.Pins.IsOccupied(sq)
}

var KingAttacks = func() (ret [SquareCount]Bitboard) {
	for src := range Squares() {
		for dir := range Directions() {
			if SquaresToEdge(src, dir) > 0 {
				ret[src] = ret[src].Occupy(src + dir.Offset())
			}
		}
	}

	return ret
}()

func GenerateKingAttacks(b *Board, a *Attacks, attacker Color) {
	a.All = a.All.Set(KingAttacks[b.Kings[attacker]])
}

func GenerateSlidingAttacks(b *Board, a *Attacks, attacker Color) {
	orthogonal := b.Bits.Pieces[Rook].
		Set(b.Bits.Pieces[Queen]).
		And(b.Bits.Players[attacker])

	diagonal := b.Bits.Pieces[Bishop].
		Set(b.Bits.Pieces[Queen]).
		And(b.Bits.Players[attacker])

	blockers := b.Bits.All.Unoccupy(b.Kings[attacker.Opponent()])

	for src := range orthogonal.Occupied() {
		a.All = a.All.Set(MagicOrthogonalMoves(src, blockers))
	}

	for src := range diagonal.Occupied() {
		a.All = a.All.Set(MagicDiagonalMoves(src, blockers))
	}

	// find pins / checks
	src := b.Kings[attacker.Opponent()]

	for dir := range Directions() {
		// FIXME
		// sliders := orthogonal
		// if dir.IsDiagonal() {
		// sliders = diagonal
		// }
		// if !dir.Mask(src).AnySet(sliders) {
		//  continue
		// }

		pin := false
		ray := Bitboard(0)

		for mul := Square(1); mul <= SquaresToEdge(src, dir); mul++ {
			dst := src + dir.Offset()*mul
			ray = ray.Occupy(dst)

			blocker := b.Squares[dst]

			if blocker == EmptySquare {
				continue
			}

			if blocker.Color() == attacker {
				if (dir.IsDiagonal() && diagonal.IsOccupied(dst)) || (!dir.IsDiagonal() && orthogonal.IsOccupied(dst)) {
					if pin {
						a.Pins = a.Pins.Set(ray)
					} else {
						a.Checks++
						a.CheckRays = a.CheckRays.Set(ray)
					}

					break
				}
			}

			if !pin {
				// opponent's piece blocks the ray, could be pinned
				pin = true
			} else {
				// second piece blocking the ray, not a pin
				break
			}
		}

	}
}

var KnightAttacks = func() (ret [SquareCount]Bitboard) {
	for src := range Squares() {
		jumps := [...]Square{
			North.Offset()*2 + East.Offset(),
			North.Offset()*2 + West.Offset(),
			South.Offset()*2 + East.Offset(),
			South.Offset()*2 + West.Offset(),
			East.Offset()*2 + North.Offset(),
			East.Offset()*2 + South.Offset(),
			West.Offset()*2 + North.Offset(),
			West.Offset()*2 + South.Offset(),
		}

		for _, jump := range jumps {
			dst := src + jump

			if !dst.Valid() {
				continue
			}

			if abs(src.File()-dst.File()) > 2 || abs(src.Rank()-dst.Rank()) > 2 {
				continue
			}

			ret[src] = ret[src].Occupy(dst)
		}
	}

	return ret
}()

func GenerateKnightAttacks(b *Board, a *Attacks, attacker Color) {
	knights := b.Bits.Pieces[Knight] & b.Bits.Players[attacker]

	for src := range knights.Occupied() {
		atk := KnightAttacks[src]

		if atk.IsOccupied(b.Kings[attacker.Opponent()]) {
			a.Checks++
			a.CheckRays = a.CheckRays.Occupy(src)
		}

		a.All = a.All.Set(atk)
	}
}

var PawnAttacks = func() (ret [ColorCount][SquareCount]Bitboard) {
	for color := range Colors() {
		dirs := [2]Direction{NorthEast, NorthWest}
		if color == Black {
			dirs = [2]Direction{SouthWest, SouthEast}
		}

		for src := range Squares() {
			for _, dir := range dirs {
				dst := src + dir.Offset()

				if !dst.Valid() {
					continue
				}

				if abs(src.File()-dst.File()) == 1 {
					ret[color][src] = ret[color][src].Occupy(dst)
				}
			}
		}
	}

	return ret
}()

func GeneratePawnAttacks(b *Board, a *Attacks, attacker Color) {
	pawns := b.Bits.Pieces[Pawn].And(b.Bits.Players[attacker])

	for src := range pawns.Occupied() {
		atk := PawnAttacks[attacker][src]

		if atk.IsOccupied(b.Kings[attacker.Opponent()]) {
			a.Checks++
			a.CheckRays = a.CheckRays.Occupy(src)
		}

		a.All = a.All.Set(atk)
	}
}
