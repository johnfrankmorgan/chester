package main

import (
	"fmt"
	"strings"
)

type Move struct {
	From  Square
	To    Square
	Flags MoveFlags
}

type MoveFlags uint16

const (
	MoveFlagsEnPassant MoveFlags = 1 << iota
	MoveFlagsDoublePawnPush
	MoveFlagsCastleKingside
	MoveFlagsCastleQueenside
	MoveFlagsPromoteToQueen
	MoveFlagsPromoteToRook
	MoveFlagsPromoteToBishop
	MoveFlagsPromoteToKnight

	MoveFlagsCastle  = MoveFlagsCastleKingside | MoveFlagsCastleQueenside
	MoveFlagsPromote = MoveFlagsPromoteToQueen |
		MoveFlagsPromoteToRook |
		MoveFlagsPromoteToBishop |
		MoveFlagsPromoteToKnight
)

type MoveGenerator struct {
	//
}

type MoveGeneratorOptions struct {
	CapturesOnly bool
}

func NewMove(from, to Square, flags ...MoveFlags) Move {
	move := Move{
		From: from,
		To:   to,
	}

	for _, flag := range flags {
		move.Flags |= flag
	}

	return move
}

func (m Move) String() string {
	s := strings.Builder{}

	s.WriteString(m.From.String())
	s.WriteString(m.To.String())

	if m.Flags != 0 {
		fmt.Fprintf(&s, " (%s)", m.Flags)
	}

	return s.String()
}

func (m Move) Valid() bool {
	return m.From != 0 || m.To != 0
}

func (m Move) Promotion() PieceKind {
	if m.Flags.AnySet(MoveFlagsPromote) {
		if m.Flags.IsSet(MoveFlagsPromoteToQueen) {
			return PieceKindQueen
		}

		if m.Flags.IsSet(MoveFlagsPromoteToRook) {
			return PieceKindRook
		}

		if m.Flags.IsSet(MoveFlagsPromoteToBishop) {
			return PieceKindBishop
		}

		if m.Flags.IsSet(MoveFlagsPromoteToKnight) {
			return PieceKindKnight
		}
	}

	return PieceKindNone
}

func (mf MoveFlags) String() string {
	s := strings.Builder{}

	for i, c := range []byte{'e', '2', 'K', 'Q', 'q', 'r', 'b', 'n'} {
		if mf.IsSet(1 << i) {
			s.WriteByte(c)
		}
	}

	return s.String()
}

func (mf MoveFlags) IsSet(flags MoveFlags) bool {
	return mf&flags == flags
}

func (mf MoveFlags) AnySet(flags MoveFlags) bool {
	return mf&flags > 0
}

func (mg MoveGenerator) Generate(board *Board, opts MoveGeneratorOptions) []Move {
	player := board.Player

	moves := mg.king(board, player, opts, make([]Move, 0, 256))

	if !board.Attacks.Checks.Double {
		moves = mg.sliding(board, player, opts, moves)
		moves = mg.knight(board, player, opts, moves)
		moves = mg.pawn(board, player, opts, moves)
	}

	return moves
}

func (mg MoveGenerator) king(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	src := board.Kings[player]

	legal := Magic.KingMoves(src) & ^board.Attacks.All & ^board.Bitboards.Colors[player]

	// FIXME: do this without an if - eg:
	//     legal &= opts.LegalMask
	// board.Bitboards.Colors[player.Opponent()] when opts.CapturesOnly
	// BitboardAll otherwise
	if opts.CapturesOnly {
		legal &= board.Bitboards.Colors[player.Opponent()]
	}

	for legal != 0 {
		dst := Square(legal.PopLSB())

		moves = append(moves, NewMove(src, dst))
	}

	if !opts.CapturesOnly && !board.Attacks.Checks.Check {
		if board.Castling[player].Kingside {
			clear := Ternary(player == ColorWhite, Bitboard(0b01100000), 0b01100000<<56)
			check := clear

			if !board.Bitboards.All.AnySet(clear) && !board.Attacks.All.AnySet(check) {
				dst := Ternary(player == ColorWhite, SquareG1, SquareG8)

				moves = append(moves, NewMove(src, dst, MoveFlagsCastleKingside))
			}
		}

		if board.Castling[player].Queenside {
			clear := Ternary(player == ColorWhite, Bitboard(0b00001110), 0b00001110<<56)
			check := Ternary(player == ColorWhite, Bitboard(0b00001100), 0b00001100<<56)

			if !board.Bitboards.All.AnySet(clear) && !board.Attacks.All.AnySet(check) {
				dst := Ternary(player == ColorWhite, SquareC1, SquareC8)

				moves = append(moves, NewMove(src, dst, MoveFlagsCastleQueenside))
			}
		}
	}

	return moves
}

func (mg MoveGenerator) sliding(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	queens := board.Bitboards.Pieces[PieceKindQueen] & board.Bitboards.Colors[player]

	orthogonal := (board.Bitboards.Pieces[PieceKindRook] & board.Bitboards.Colors[player]) | queens
	diagonal := (board.Bitboards.Pieces[PieceKindBishop] & board.Bitboards.Colors[player]) | queens

	if board.Attacks.Checks.Check {
		// can't move pinned pieces when in check
		orthogonal &= ^board.Attacks.Pins
		diagonal &= ^board.Attacks.Pins
	}

	legal := ^board.Bitboards.Colors[player] & board.Attacks.Checks.Rays

	if opts.CapturesOnly {
		legal &= board.Bitboards.Colors[player.Opponent()]
	}

	for orthogonal != 0 {
		src := Square(orthogonal.PopLSB())

		legal := Magic.OrthogonalMoves(src, board.Bitboards.All) & legal

		if board.Attacks.Pins.IsSet(src.Bitboard()) {
			legal &= src.AlignMask(board.Kings[player])
		}

		for legal != 0 {
			dst := Square(legal.PopLSB())

			moves = append(moves, NewMove(src, dst))
		}
	}

	for diagonal != 0 {
		src := Square(diagonal.PopLSB())

		legal := Magic.DiagonalMoves(src, board.Bitboards.All) & legal

		if board.Attacks.Pins.IsSet(src.Bitboard()) {
			legal &= src.AlignMask(board.Kings[player])
		}

		for legal != 0 {
			dst := Square(legal.PopLSB())

			moves = append(moves, NewMove(src, dst))
		}
	}

	return moves
}

func (mg MoveGenerator) knight(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	knights := board.Bitboards.Pieces[PieceKindKnight] & board.Bitboards.Colors[player] & ^board.Attacks.Pins

	legal := ^board.Bitboards.Colors[player] & board.Attacks.Checks.Rays

	if opts.CapturesOnly {
		legal &= board.Bitboards.Colors[player.Opponent()]
	}

	for knights != 0 {
		src := Square(knights.PopLSB())

		legal := Magic.KnightMoves(src) & legal

		for legal != 0 {
			dst := Square(legal.PopLSB())

			moves = append(moves, NewMove(src, dst))
		}
	}

	return moves
}

func (mg MoveGenerator) pawn(board *Board, player Color, opts MoveGeneratorOptions, moves []Move) []Move {
	pawns := board.Bitboards.Pieces[PieceKindPawn] & board.Bitboards.Colors[player]

	if board.Attacks.Checks.Check {
		// can't move pinned pieces when in check
		pawns &= ^board.Attacks.Pins
	}

	dir := Ternary(player == ColorWhite, DirectionNorth, DirectionSouth)

	for pawns > 0 {
		src := Square(pawns.PopLSB())

		if !opts.CapturesOnly {
			dst := src + dir.Offset()

			// Yikes! FIXME
			if !board.Bitboards.All.IsSet(dst.Bitboard()) {
				if !board.Attacks.Pins.IsSet(src.Bitboard()) || src.AlignMask(board.Kings[player]).IsSet(dst.Bitboard()) {
					if board.Attacks.Checks.Rays.IsSet(dst.Bitboard()) {
						if dst.Rank() == Rank1 || dst.Rank() == Rank8 {
							moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToQueen))
							moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToRook))
							moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToBishop))
							moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToKnight))
						} else {
							moves = append(moves, NewMove(src, dst))
						}
					}

					if (player == ColorWhite && src.Rank() == Rank2) || (player == ColorBlack && src.Rank() == Rank7) {
						dst += dir.Offset()

						if !board.Bitboards.All.IsSet(dst.Bitboard()) && board.Attacks.Checks.Rays.IsSet(dst.Bitboard()) {
							moves = append(moves, NewMove(src, dst, MoveFlagsDoublePawnPush))
						}
					}
				}
			}
		}

		attacks := Magic.PawnAttacks(player, src)

		if board.Attacks.Checks.Check && board.EnPassant != 0 && attacks.IsSet(board.EnPassant.Bitboard()) {
			// FIXME: this is repeated below
			// check if taking en passant will reveal a check
			ep := board.EnPassant + Ternary(player == ColorWhite, DirectionSouth, DirectionNorth).Offset()

			orthogonal := board.Bitboards.Pieces[PieceKindRook]
			orthogonal |= board.Bitboards.Pieces[PieceKindQueen]
			orthogonal &= board.Bitboards.Colors[player.Opponent()]

			blockers := board.Bitboards.All & ^ep.Bitboard() & ^src.Bitboard()

			if !Magic.OrthogonalMoves(board.Kings[player], blockers).AnySet(orthogonal) {
				moves = append(moves, NewMove(src, board.EnPassant, MoveFlagsEnPassant))
			}
		}

		attacks &= board.Attacks.Checks.Rays

		if board.EnPassant != 0 {
			attacks &= (board.Bitboards.Colors[player.Opponent()] | board.EnPassant.Bitboard())
		} else {
			attacks &= board.Bitboards.Colors[player.Opponent()]
		}

		if board.Attacks.Pins.IsSet(src.Bitboard()) {
			attacks &= src.AlignMask(board.Kings[player])
		}

		for attacks != 0 {
			dst := Square(attacks.PopLSB())

			if dst.Rank() == Rank1 || dst.Rank() == Rank8 {
				moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToQueen))
				moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToRook))
				moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToBishop))
				moves = append(moves, NewMove(src, dst, MoveFlagsPromoteToKnight))
			} else if dst == board.EnPassant {
				// check if taking en passant will reveal a check
				ep := board.EnPassant + Ternary(player == ColorWhite, DirectionSouth, DirectionNorth).Offset()

				orthogonal := board.Bitboards.Pieces[PieceKindRook]
				orthogonal |= board.Bitboards.Pieces[PieceKindQueen]
				orthogonal &= board.Bitboards.Colors[player.Opponent()]

				blockers := board.Bitboards.All & ^ep.Bitboard() & ^src.Bitboard()

				if !Magic.OrthogonalMoves(board.Kings[player], blockers).AnySet(orthogonal) {
					moves = append(moves, NewMove(src, dst, MoveFlagsEnPassant))
				}
			} else {
				moves = append(moves, NewMove(src, dst))
			}
		}
	}

	return moves
}
