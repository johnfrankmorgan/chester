package main

import (
	"strings"
)

type Move struct {
	From  Square
	To    Square
	Flags MoveFlags
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

type MoveFlags uint16

const (
	MoveFlagCapture MoveFlags = 1 << iota
	MoveFlagCaptureEnPassant
	MoveFlagDoublePawnPush
	MoveFlagCastleKingside
	MoveFlagCastleQueenside
	MoveFlagPromoteToQueen
	MoveFlagPromoteToRook
	MoveFlagPromoteToBishop
	MoveFlagPromoteToKnight

	MoveFlagCastleAny  = MoveFlagCastleKingside | MoveFlagCastleQueenside
	MoveFlagPromoteAny = MoveFlagPromoteToQueen |
		MoveFlagPromoteToRook |
		MoveFlagPromoteToBishop |
		MoveFlagPromoteToKnight
)

func (m Move) String() string {
	s := strings.Builder{}

	s.WriteString(m.From.String())
	s.WriteString(m.To.String())

	if p, ok := m.Promotion(); ok {
		s.WriteString(p.String())
	}

	return s.String()
}

func (m Move) IsZero() bool {
	return m == Move{}
}

func (m Move) Promotion() (PieceType, bool) {
	if m.Flags&MoveFlagPromoteAny != 0 {
		if m.Flags&MoveFlagPromoteToQueen != 0 {
			return Queen, true
		} else if m.Flags&MoveFlagPromoteToRook != 0 {
			return Rook, true
		} else if m.Flags&MoveFlagPromoteToBishop != 0 {
			return Bishop, true
		} else if m.Flags&MoveFlagPromoteToKnight != 0 {
			return Knight, true
		}
	}
	return 0, false
}

type MoveGenerationOptions struct {
	CapturesOnly bool
}

func GenerateMoves(b *Board, opts MoveGenerationOptions) []Move {
	moves := GenerateKingMoves(b, make([]Move, 0, 256), opts)

	if b.Attacks.Checks < 2 {
		moves = GenerateSlidingMoves(b, moves, opts)
		moves = GenerateKnightMoves(b, moves, opts)
		moves = GeneratePawnMoves(b, moves, opts)
	}

	return moves
}

var KingMoves = KingAttacks

func GenerateKingMoves(b *Board, moves []Move, opts MoveGenerationOptions) []Move {
	src := b.Kings[b.Player]

	legal := KingMoves[src].
		Clear(b.Bits.Players[b.Player]).
		Clear(b.Attacks.All)

	for dst := range legal.Occupied() {
		if b.Bits.Players[b.Player.Opponent()].IsOccupied(dst) {
			moves = append(moves, NewMove(src, dst, MoveFlagCapture))
		} else if !opts.CapturesOnly {
			moves = append(moves, NewMove(src, dst))
		}
	}

	if opts.CapturesOnly {
		return moves
	}

	if b.Attacks.Checks == 0 && b.Castling[b.Player].Kingside {
		if b.Player == White {
			mask := Bitboard(0).Occupy(SquareF1).Occupy(SquareG1)

			if !b.Bits.All.AnySet(mask) && !b.Attacks.All.AnySet(mask) {
				moves = append(moves, NewMove(SquareE1, SquareG1, MoveFlagCastleKingside))
			}
		} else {
			mask := Bitboard(0).Occupy(SquareF8).Occupy(SquareG8)

			if !b.Bits.All.AnySet(mask) && !b.Attacks.All.AnySet(mask) {
				moves = append(moves, NewMove(SquareE8, SquareG8, MoveFlagCastleKingside))
			}
		}
	}

	if b.Attacks.Checks == 0 && b.Castling[b.Player].Queenside {
		if b.Player == White {
			blockers := Bitboard(0).Occupy(SquareD1).Occupy(SquareC1).Occupy(SquareB1)
			attacks := Bitboard(0).Occupy(SquareD1).Occupy(SquareC1)

			if !b.Bits.All.AnySet(blockers) && !b.Attacks.All.AnySet(attacks) {
				moves = append(moves, NewMove(SquareE1, SquareC1, MoveFlagCastleQueenside))
			}
		} else {
			blockers := Bitboard(0).Occupy(SquareD8).Occupy(SquareC8).Occupy(SquareB8)
			attacks := Bitboard(0).Occupy(SquareD8).Occupy(SquareC8)

			if !b.Bits.All.AnySet(blockers) && !b.Attacks.All.AnySet(attacks) {
				moves = append(moves, NewMove(SquareE8, SquareC8, MoveFlagCastleQueenside))
			}
		}
	}

	return moves
}

func GenerateSlidingMoves(b *Board, moves []Move, opts MoveGenerationOptions) []Move {
	orthogonal := b.Bits.Pieces[Rook].
		Set(b.Bits.Pieces[Queen]).
		And(b.Bits.Players[b.Player])

	diagonal := b.Bits.Pieces[Bishop].
		Set(b.Bits.Pieces[Queen]).
		And(b.Bits.Players[b.Player])

	if b.Attacks.Checks > 0 {
		// can't move pinned pieces when in check
		orthogonal = orthogonal.Clear(b.Attacks.Pins)
		diagonal = diagonal.Clear(b.Attacks.Pins)
	}

	legal := b.Attacks.CheckRays.Clear(b.Bits.Players[b.Player])

	if opts.CapturesOnly {
		legal = legal.And(b.Bits.Players[b.Player.Opponent()])
	}

	for src := range orthogonal.Occupied() {
		legal := MagicOrthogonalMoves(src, b.Bits.All).And(legal)

		if b.Attacks.Pins.IsOccupied(src) {
			legal = legal.And(BitboardAlignedAlong(src, b.Kings[b.Player]))
		}

		for dst := range legal.Occupied() {
			if b.Bits.Players[b.Player.Opponent()].IsOccupied(dst) {
				moves = append(moves, NewMove(src, dst, MoveFlagCapture))
			} else {
				moves = append(moves, NewMove(src, dst))
			}
		}
	}

	for src := range diagonal.Occupied() {
		legal := MagicDiagonalMoves(src, b.Bits.All).And(legal)

		if b.Attacks.Pins.IsOccupied(src) {
			legal = legal.And(BitboardAlignedAlong(src, b.Kings[b.Player]))
		}

		for dst := range legal.Occupied() {
			if b.Bits.Players[b.Player.Opponent()].IsOccupied(dst) {
				moves = append(moves, NewMove(src, dst, MoveFlagCapture))
			} else {
				moves = append(moves, NewMove(src, dst))
			}
		}
	}

	return moves
}

var KnightMoves = KnightAttacks

func GenerateKnightMoves(b *Board, moves []Move, opts MoveGenerationOptions) []Move {
	knights := b.Bits.Pieces[Knight].
		And(b.Bits.Players[b.Player]).
		Clear(b.Attacks.Pins)

	legal := b.Attacks.CheckRays.
		Clear(b.Bits.Players[b.Player])

	if opts.CapturesOnly {
		legal = legal.And(b.Bits.Players[b.Player.Opponent()])
	}

	for src := range knights.Occupied() {
		jumps := KnightMoves[src].And(legal)

		for dst := range jumps.Occupied() {
			if b.Bits.Players[b.Player.Opponent()].IsOccupied(dst) {
				moves = append(moves, NewMove(src, dst, MoveFlagCapture))
			} else {
				moves = append(moves, NewMove(src, dst))
			}
		}
	}

	return moves
}

func GeneratePawnMoves(b *Board, moves []Move, opts MoveGenerationOptions) []Move {
	dir := North

	if b.Player == Black {
		dir = South
	}

	pawns := b.Bits.Pieces[Pawn].And(b.Bits.Players[b.Player])

	for src := range pawns.Occupied() {
		if b.Attacks.Checks > 0 && b.Attacks.IsPinned(src) {
			continue
		}

		if !opts.CapturesOnly {
			dst := src + dir.Offset()

			if b.Bits.All.IsOccupied(dst) {
				goto captures
			}

			if b.Attacks.IsPinned(src) && !BitboardAlignedAlong(src, b.Kings[b.Player]).IsOccupied(dst) {
				goto captures
			}

			if b.Attacks.CheckRays.IsOccupied(dst) {
				if dst.Rank() == Rank1 || dst.Rank() == Rank8 {
					moves = append(
						moves,
						NewMove(src, dst, MoveFlagPromoteToQueen),
						NewMove(src, dst, MoveFlagPromoteToRook),
						NewMove(src, dst, MoveFlagPromoteToBishop),
						NewMove(src, dst, MoveFlagPromoteToKnight),
					)
				} else {
					moves = append(moves, NewMove(src, dst))
				}
			}

			dst += dir.Offset()

			if dst.Valid() && b.Attacks.CheckRays.IsOccupied(dst) {
				if b.Player == White && src.Rank() == Rank2 && !b.Bits.All.IsOccupied(dst) {
					moves = append(moves, NewMove(src, dst, MoveFlagDoublePawnPush))
				} else if b.Player == Black && src.Rank() == Rank7 && !b.Bits.All.IsOccupied(dst) {
					moves = append(moves, NewMove(src, dst, MoveFlagDoublePawnPush))
				}
			}
		}

	captures:
		enemies := b.Bits.Players[b.Player.Opponent()]

		if b.EnPassant != 0 {
			enemies = enemies.Occupy(b.EnPassant)
		}

		attacks := PawnAttacks[b.Player][src].
			And(enemies).
			And(b.Attacks.CheckRays)

		if b.Attacks.IsPinned(src) {
			attacks = attacks.And(BitboardAlignedAlong(src, b.Kings[b.Player]))
		}

		for dst := range attacks.Occupied() {
			if b.EnPassant != 0 && dst == b.EnPassant {
				moves = append(moves, NewMove(src, dst, MoveFlagCapture, MoveFlagCaptureEnPassant))
			} else if dst.Rank() == Rank1 || dst.Rank() == Rank8 {
				moves = append(
					moves,
					NewMove(src, dst, MoveFlagCapture|MoveFlagPromoteToQueen),
					NewMove(src, dst, MoveFlagCapture|MoveFlagPromoteToRook),
					NewMove(src, dst, MoveFlagCapture|MoveFlagPromoteToBishop),
					NewMove(src, dst, MoveFlagCapture|MoveFlagPromoteToKnight),
				)
			} else {
				moves = append(moves, NewMove(src, dst, MoveFlagCapture))
			}
		}
	}

	return moves
}
