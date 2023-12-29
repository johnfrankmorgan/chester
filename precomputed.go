package main

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"unsafe"
)

var Precomputed = struct {
	OpeningBook OpeningBook

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

	Evaluation struct {
		Squares PieceTable[ColorTable[struct{ Middlegame, Endgame [SquareCount]float64 }]]
	}

	Zobrist struct {
		Players   ColorTable[uint64]
		Pieces    ColorTable[PieceTable[[SquareCount]uint64]]
		Castling  ColorTable[struct{ Kingside, Queenside uint64 }]
		EnPassant [SquareCount]uint64
	}
}{
	Evaluation: struct {
		Squares PieceTable[ColorTable[struct{ Middlegame, Endgame [SquareCount]float64 }]]
	}{
		// https://www.chessprogramming.org/Simplified_Evaluation_Function
		Squares: PieceTable[ColorTable[struct{ Middlegame, Endgame [SquareCount]float64 }]]{
			PiecePawn: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						0, 0, 0, 0, 0, 0, 0, 0,
						50, 50, 50, 50, 50, 50, 50, 50,
						10, 10, 20, 30, 30, 20, 10, 10,
						5, 5, 10, 25, 25, 10, 5, 5,
						0, 0, 0, 20, 20, 0, 0, 0,
						5, -5, -10, 5, 5, -10, -5, 5,
						5, 10, 10, -20, -20, 10, 10, 5,
						0, 0, 0, 0, 0, 0, 0, 0,
					},

					Endgame: [SquareCount]float64{
						0, 0, 0, 0, 0, 0, 0, 0,
						80, 80, 80, 80, 80, 80, 80, 80,
						50, 50, 50, 50, 50, 50, 50, 50,
						30, 30, 30, 30, 30, 30, 30, 30,
						20, 20, 20, 20, 20, 20, 20, 20,
						10, 10, 10, 10, 10, 10, 10, 10,
						10, 10, 10, 10, 10, 10, 10, 10,
						0, 0, 0, 0, 0, 0, 0, 0,
					},
				},
			},

			PieceKnight: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						-50, -40, -30, -30, -30, -30, -40, -50,
						-40, -20, 0, 0, 0, 0, -20, -40,
						-30, 0, 10, 15, 15, 10, 0, -30,
						-30, 5, 15, 20, 20, 15, 5, -30,
						-30, 0, 15, 20, 20, 15, 0, -30,
						-30, 5, 10, 15, 15, 10, 5, -30,
						-40, -20, 0, 5, 5, 0, -20, -40,
						-50, -40, -30, -30, -30, -30, -40, -50,
					},
				},
			},

			PieceBishop: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						-20, -10, -10, -10, -10, -10, -10, -20,
						-10, 0, 0, 0, 0, 0, 0, -10,
						-10, 0, 5, 10, 10, 5, 0, -10,
						-10, 5, 5, 10, 10, 5, 5, -10,
						-10, 0, 10, 10, 10, 10, 0, -10,
						-10, 10, 10, 10, 10, 10, 10, -10,
						-10, 5, 0, 0, 0, 0, 5, -10,
						-20, -10, -10, -10, -10, -10, -10, -20,
					},
				},
			},

			PieceRook: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						0, 0, 0, 0, 0, 0, 0, 0,
						5, 10, 10, 10, 10, 10, 10, 5,
						-5, 0, 0, 0, 0, 0, 0, -5,
						-5, 0, 0, 0, 0, 0, 0, -5,
						-5, 0, 0, 0, 0, 0, 0, -5,
						-5, 0, 0, 0, 0, 0, 0, -5,
						-5, 0, 0, 0, 0, 0, 0, -5,
						0, 0, 0, 5, 5, 0, 0, 0,
					},
				},
			},

			PieceQueen: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						-20, -10, -10, -5, -5, -10, -10, -20,
						-10, 0, 0, 0, 0, 0, 0, -10,
						-10, 0, 5, 5, 5, 5, 0, -10,
						-5, 0, 5, 5, 5, 5, 0, -5,
						0, 0, 5, 5, 5, 5, 0, -5,
						-10, 5, 5, 5, 5, 5, 0, -10,
						-10, 0, 5, 0, 0, 0, 0, -10,
						-20, -10, -10, -5, -5, -10, -10, -20,
					},
				},
			},

			PieceKing: {
				ColorWhite: {
					Middlegame: [SquareCount]float64{
						-30, -40, -40, -50, -50, -40, -40, -30,
						-30, -40, -40, -50, -50, -40, -40, -30,
						-30, -40, -40, -50, -50, -40, -40, -30,
						-30, -40, -40, -50, -50, -40, -40, -30,
						-20, -30, -30, -40, -40, -30, -30, -20,
						-10, -20, -20, -20, -20, -20, -20, -10,
						20, 20, 0, 0, 0, 0, 20, 20,
						20, 30, 10, 0, 0, 10, 30, 20,
					},

					Endgame: [SquareCount]float64{
						-20, -10, -10, -10, -10, -10, -10, -20,
						-5, 0, 5, 5, 5, 5, 0, -5,
						-10, -5, 20, 30, 30, 20, -5, -10,
						-15, -10, 35, 45, 45, 35, -10, -15,
						-20, -15, 30, 40, 40, 30, -15, -20,
						-25, -20, 20, 25, 25, 20, -20, -25,
						-30, -25, 0, 0, 0, 0, -25, -30,
						-50, -30, -30, -30, -30, -30, -30, -50,
					},
				},
			},
		},
	},
}

func init() {
	slog.Debug("initializing knight moves")
	{
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
	}

	slog.Debug("initializing pawn moves")
	{
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
	}

	slog.Debug("initializing knight attacks")
	{
		for _, moves := range Precomputed.Moves.Knight {
			for _, move := range moves {
				Precomputed.Attacks.Knight[move.From] = append(Precomputed.Attacks.Knight[move.From], move.To)
			}
		}
	}

	slog.Debug("initializing pawn attacks")
	{
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
	}

	slog.Debug("initializing direction masks")
	{
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
	}

	slog.Debug("initializing alignment masks")
	{
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

	slog.Debug("initializing evaluation square tables")
	{
		for _, p := range [...]PieceKind{PieceKnight, PieceBishop, PieceRook, PieceQueen} {
			Precomputed.Evaluation.Squares[p][ColorWhite].Endgame = Precomputed.Evaluation.Squares[p][ColorWhite].Middlegame
		}

		for p := range Precomputed.Evaluation.Squares {
			p := PieceKind(p)

			if p == PieceNone {
				continue
			}

			for rank := RankFirst; rank <= RankLast; rank++ {
				for file := FileFirst; file <= FileLast; file++ {
					src := NewSquare(file, rank)
					dst := NewSquare(file, RankLast-rank)

					Precomputed.Evaluation.Squares[p][ColorBlack].Middlegame[dst] = Precomputed.Evaluation.Squares[p][ColorWhite].Middlegame[src]
					Precomputed.Evaluation.Squares[p][ColorBlack].Endgame[dst] = Precomputed.Evaluation.Squares[p][ColorWhite].Endgame[src]
				}
			}
		}
	}

	slog.Debug("loading opening book")
	{
		check(json.Unmarshal(_openings, &Precomputed.OpeningBook))
	}

	slog.Debug("initializing zobrist values")
	{
		size := unsafe.Sizeof(Precomputed.Zobrist)
		if size%unsafe.Sizeof(uint64(0)) != 0 {
			panic("unsafe.Sizeof(Precomputed.Zobrist)%unsafe.Sizeof(uint64(0)) != 0")
		}

		ptr := unsafe.Pointer(&Precomputed.Zobrist)

		for offset := uintptr(0); offset < size; offset += unsafe.Sizeof(uint64(0)) {
			*(*uint64)(unsafe.Pointer(uintptr(ptr) + offset)) = rand.Uint64()
		}
	}
}
