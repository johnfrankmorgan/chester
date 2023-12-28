package main

import (
	"log/slog"
)

type MoveGenerator struct {
	//
}

type MoveGeneratorOptions struct {
	CapturesOnly bool
}

func (mg MoveGenerator) Generate(board *Board, opts MoveGeneratorOptions) []Move {
	player := board.Player

	slog.Debug("generating moves", "player", player, "captures-only", opts.CapturesOnly)

	moves := mg._king(board, player, opts, make([]Move, 0, 256))

	if !board.Attacks.Checks.Double {
		moves = mg._sliding(board, player, opts, moves)
		moves = mg._knight(board, player, opts, moves)
		moves = mg._pawn(board, player, opts, moves)
	}

	return moves
}

func (mg MoveGenerator) _king(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	slog.Debug("generating king moves", "player", player, "captures-only", opts.CapturesOnly)

	src := board.Kings[player]

	for _, dir := range Directions {
		if dir.ToEdge(src) > 0 {
			dst := src + dir.Offset()

			if board.Bitboards.Colors[player].IsSet(dst.Bitboard()) {
				continue
			} else if !board.Attacks.IsAttacked(dst) {
				if board.Bitboards.Colors[player.Opponent()].IsSet(dst.Bitboard()) {
					moves = append(moves, NewMove(src, dst, MoveFlagsCapture))
				} else if !opts.CapturesOnly {
					moves = append(moves, NewMove(src, dst))
				}
			}
		}
	}

	// TODO: castling

	return moves
}

func (mg MoveGenerator) _sliding(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	slog.Debug("generating sliding moves", "player", player, "captures-only", opts.CapturesOnly)

	queens := board.Bitboards.Pieces[PieceQueen]
	rooks := board.Bitboards.Pieces[PieceRook]
	bishops := board.Bitboards.Pieces[PieceBishop]

	orthogonal := (queens | rooks) & board.Bitboards.Colors[player]
	diagonal := (queens | bishops) & board.Bitboards.Colors[player]

	for _, dir := range Directions {
		sliders := orthogonal

		if dir.IsDiagonal() {
			sliders = diagonal
		}

		for sliders > 0 {
			src := Square(sliders.PopLSB())

			if board.Attacks.Checks.Check && board.Attacks.IsPinned(src) {
				continue
			}

			legal := BitboardAll

			if board.Attacks.Checks.Check {
				legal &= board.Attacks.Checks.Rays
			} else if board.Attacks.IsPinned(src) {
				legal &= Precomputed.Masks.Alignment[src][board.Kings[player]]
			}

			for mul := Square(1); mul <= dir.ToEdge(src); mul++ {
				dst := src + dir.Offset()*mul

				if board.Bitboards.Colors[player].IsSet(dst.Bitboard()) {
					break
				}

				if legal.IsSet(dst.Bitboard()) {
					if board.Bitboards.Colors[player.Opponent()].IsSet(dst.Bitboard()) {
						moves = append(moves, NewMove(src, dst, MoveFlagsCapture))
						break
					} else if !opts.CapturesOnly {
						moves = append(moves, NewMove(src, dst))
					}
				}
			}
		}
	}

	return moves
}

func (mg MoveGenerator) _knight(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	slog.Debug("generating knight moves", "player", player, "captures-only", opts.CapturesOnly)

	knights := board.Bitboards.Pieces[PieceKnight] & board.Bitboards.Colors[player]

	for knights > 0 {
		src := Square(knights.PopLSB())

		if board.Attacks.IsPinned(src) {
			continue
		}

		legal := BitboardAll
		legal &^= board.Bitboards.Colors[player]

		if board.Attacks.Checks.Check {
			legal &= board.Attacks.Checks.Rays
		}

		for _, move := range Precomputed.Moves.Knight[src] {
			if legal.IsSet(move.To.Bitboard()) {
				if board.Bitboards.Colors[player.Opponent()].IsSet(move.To.Bitboard()) {
					moves = append(moves, NewMove(move.From, move.To, MoveFlagsCapture))
				} else if !opts.CapturesOnly {
					moves = append(moves, move)
				}
			}
		}
	}

	return moves
}

func (mg MoveGenerator) _pawn(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	slog.Debug("generating pawn moves", "player", player, "captures-only", opts.CapturesOnly)

	pawns := board.Bitboards.Pieces[PiecePawn] & board.Bitboards.Colors[player]

	for pawns > 0 {
		src := Square(pawns.PopLSB())

		if board.Attacks.Checks.Check && board.Attacks.IsPinned(src) {
			continue
		}

		legal := (BitboardAll &^ board.Bitboards.Colors[player])

		if board.Attacks.Checks.Check {
			legal &= board.Attacks.Checks.Rays
		} else if board.Attacks.IsPinned(src) {
			legal &= Precomputed.Masks.Alignment[src][board.Kings[player]]
		}

		if !opts.CapturesOnly {
			for _, move := range Precomputed.Moves.Pawn[player][src] {
				if board.Bitboards.Colors[player].IsSet(move.To.Bitboard()) {
					break
				} else if board.Bitboards.Colors[player.Opponent()].IsSet(move.To.Bitboard()) {
					break
				} else if legal.IsSet(move.To.Bitboard()) {
					moves = append(moves, move)
				}
			}
		}

		for _, dst := range Precomputed.Attacks.Pawn[player][src] {
			if board.Bitboards.Colors[player.Opponent()].IsSet(dst.Bitboard()) && legal.IsSet(dst.Bitboard()) {
				moves = append(moves, NewMove(src, dst, MoveFlagsCapture))
			} else if dst == board.EnPassant {
				ep := board.EnPassant + ternary(player == ColorWhite, DirectionSouth.Offset(), DirectionNorth.Offset())
				capture := true

				king := board.Kings[player]

				orthogonal := board.Bitboards.Pieces[PieceRook] | board.Bitboards.Pieces[PieceQueen]
				orthogonal &= board.Bitboards.Colors[player.Opponent()]

				for orthogonal > 0 {
					ortho := Square(orthogonal.PopLSB())

					for _, dir := range DirectionsOrthogonal {
						for mul := Square(1); mul <= dir.ToEdge(ortho); mul++ {
							attack := ortho + dir.Offset()*mul

							if attack == src || attack == ep {
								continue
							} else if king == attack {
								capture = false
								break
							} else if board.Bitboards.All.IsSet(attack.Bitboard()) {
								break
							}
						}
					}
				}

				if capture {
					moves = append(moves, NewMove(src, dst, MoveFlagsCapture, MoveFlagsCaptureEnPassant))
				}
			}
		}
	}

	return moves
}
