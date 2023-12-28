package main

import "log/slog"

var Precomputed struct {
	Moves struct {
		Pawn   ColorTable[[SquareCount][]Move]
		Knight [SquareCount][]Move
	}

	Attacks struct {
		Pawn   ColorTable[[SquareCount][]Square]
		Knight [SquareCount][]Square
	}

	Masks struct {
		Direction [SquareCount][DirectionCount]Bitboard
		Alignment [SquareCount][SquareCount]Bitboard
	}
}

func init() {
	slog.Debug("initializing knight moves")

	for src := range Precomputed.Moves.Knight {
		src := Square(src)

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

		for _, jump := range jumps {
			if dst := src + jump; dst.Valid() && max(iabs(src.File()-dst.File()), iabs(src.Rank()-dst.Rank())) == 2 {
				Precomputed.Moves.Knight[src] = append(Precomputed.Moves.Knight[src], NewMove(src, dst))
			}
		}
	}

	slog.Debug("initializing pawn moves")

	for color, squares := range Precomputed.Moves.Pawn {
		color := Color(color)

		dir := DirectionNorth
		if color == ColorBlack {
			dir = DirectionSouth
		}

		for src := range squares {
			src := Square(src)
			dst := src + dir.Offset()

			if !dst.Valid() {
				continue
			}

			if dst.Rank() == Rank1 || dst.Rank() == Rank8 {
				Precomputed.Moves.Pawn[color][src] = append(
					Precomputed.Moves.Pawn[color][src],
					NewMove(src, dst, MoveFlagsPromoteToQueen),
					NewMove(src, dst, MoveFlagsPromoteToRook),
					NewMove(src, dst, MoveFlagsPromoteToBishop),
					NewMove(src, dst, MoveFlagsPromoteToKnight),
				)
			} else {
				Precomputed.Moves.Pawn[color][src] = append(Precomputed.Moves.Pawn[color][src], NewMove(src, dst))
			}

			if (color == ColorWhite && src.Rank() == Rank2) || (color == ColorBlack && src.Rank() == Rank7) {
				Precomputed.Moves.Pawn[color][src] = append(Precomputed.Moves.Pawn[color][src], NewMove(src, dst+dir.Offset(), MoveFlagsDoublePawnPush))
			}
		}
	}

	slog.Debug("initializing knight attacks")

	for _, moves := range Precomputed.Moves.Knight {
		for _, move := range moves {
			Precomputed.Attacks.Knight[move.From] = append(Precomputed.Attacks.Knight[move.From], move.To)
		}
	}

	slog.Debug("initializing pawn attacks")

	for color, squares := range Precomputed.Attacks.Pawn {
		color := Color(color)

		for src := range squares {
			src := Square(src)

			dirs := [...]Direction{
				DirectionNorthEast,
				DirectionNorthWest,
			}

			if color == ColorBlack {
				dirs = [...]Direction{
					DirectionSouthEast,
					DirectionSouthWest,
				}
			}

			for _, dir := range dirs {
				if dst := src + dir.Offset(); dst.Valid() && abs(src.File()-dst.File()) == 1 {
					Precomputed.Attacks.Pawn[color][src] = append(Precomputed.Attacks.Pawn[color][src], dst)
				}
			}
		}
	}

	slog.Debug("initializing direction masks")

	for src := 0; src < SquareCount; src++ {
		src := Square(src)

		Precomputed.Masks.Direction[src] = [DirectionCount]Bitboard{
			DirectionNorth: BitboardFiles[src.File()],
			DirectionSouth: BitboardFiles[src.File()],
			DirectionEast:  BitboardRanks[src.Rank()],
			DirectionWest:  BitboardRanks[src.Rank()],
		}

		for _, dir := range DirectionsDiagonal {
			ray := src.Bitboard()

			for _, off := range [...]Square{dir.Offset(), dir.Opposite().Offset()} {
				src := src

				for {
					dst := src + off

					if !dst.Valid() || abs(src.File()-dst.File()) != 1 || abs(src.Rank()-dst.Rank()) != 1 {
						break
					}

					ray.Set(dst.Bitboard())

					src = dst
				}
			}

			Precomputed.Masks.Direction[src][dir] = ray
		}
	}

	slog.Debug("initializing alignment masks")

	for src, dirs := range Precomputed.Masks.Direction {
		for _, mask := range dirs {
			alignment := mask

			for mask > 0 {
				dst := mask.PopLSB()

				Precomputed.Masks.Alignment[src][dst] = alignment
			}
		}
	}
}
