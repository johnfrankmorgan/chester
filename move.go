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

func NewUCIMove(board *Board, uci string) (Move, error) {
	if len(uci) < 4 || len(uci) > 5 {
		return Move{}, fmt.Errorf("%w: invalid move: %s", ErrUCI, uci)
	}

	from := struct {
		file File
		rank Rank
	}{
		file: File(uci[0]) - 'a',
		rank: Rank(uci[1] - '1'),
	}

	to := struct {
		file File
		rank Rank
	}{
		file: File(uci[2]) - 'a',
		rank: Rank(uci[3] - '1'),
	}

	if !from.file.Valid() {
		return Move{}, fmt.Errorf("%w: invalid source file: %c", ErrUCI, uci[0])
	}

	if !from.rank.Valid() {
		return Move{}, fmt.Errorf("%w: invalid source rank: %c", ErrUCI, uci[1])
	}

	if !to.file.Valid() {
		return Move{}, fmt.Errorf("%w: invalid destination file: %c", ErrUCI, uci[2])
	}

	if !to.rank.Valid() {
		return Move{}, fmt.Errorf("%w: invalid destination rank: %c", ErrUCI, uci[3])
	}

	move := Move{
		From: NewSquare(from.file, from.rank),
		To:   NewSquare(to.file, to.rank),
	}

	if len(uci) == 5 {
		switch uci[4] {
		case 'q':
			move.Flags |= MoveFlagsPromoteToQueen

		case 'r':
			move.Flags |= MoveFlagsPromoteToRook

		case 'b':
			move.Flags |= MoveFlagsPromoteToBishop

		case 'n':
			move.Flags |= MoveFlagsPromoteToKnight

		default:
			return Move{}, fmt.Errorf("%w: invalid promotion: %c", ErrUCI, uci[4])
		}
	}

	piece := board.Pieces[move.From]

	if piece.Is(PiecePawn) {
		if abs(to.rank-from.rank) == 2 {
			move.Flags |= MoveFlagsDoublePawnPush
		} else if to.file != from.file && board.Pieces[move.To].Is(PieceNone) {
			move.Flags |= MoveFlagsCapture | MoveFlagsCaptureEnPassant
		}
	} else if piece.Is(PieceKing) && from.file == FileE {
		if to.file == FileG {
			move.Flags |= MoveFlagsCastleKingside
		} else if to.file == FileC {
			move.Flags |= MoveFlagsCastleQueenside
		}
	}

	if !board.Pieces[move.To].Is(PieceNone) {
		move.Flags |= MoveFlagsCapture
	}

	return move, nil
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

func (m Move) UCI() string {
	suffix := ""

	if m.Flags.AnySet(MoveFlagsPromote) {
		if m.Flags.IsSet(MoveFlagsPromoteToQueen) {
			suffix = "q"
		} else if m.Flags.IsSet(MoveFlagsPromoteToRook) {
			suffix = "r"
		} else if m.Flags.IsSet(MoveFlagsPromoteToBishop) {
			suffix = "b"
		} else if m.Flags.IsSet(MoveFlagsPromoteToKnight) {
			suffix = "n"
		}
	}

	return m.From.String() + m.To.String() + suffix
}

func (m Move) Promotion() PieceKind {
	if m.Flags.AnySet(MoveFlagsPromote) {
		if m.Flags.IsSet(MoveFlagsPromoteToQueen) {
			return PieceQueen
		}

		if m.Flags.IsSet(MoveFlagsPromoteToRook) {
			return PieceRook
		}

		if m.Flags.IsSet(MoveFlagsPromoteToBishop) {
			return PieceBishop
		}

		if m.Flags.IsSet(MoveFlagsPromoteToKnight) {
			return PieceKnight
		}
	}

	return PieceNone
}
