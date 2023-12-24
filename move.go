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

func (m Move) String() string {
	s := strings.Builder{}

	s.WriteString(m.From.String())
	s.WriteString(m.To.String())

	if m.Flags != 0 {
		fmt.Fprintf(&s, " (%s)", m.Flags)
	}

	return s.String()
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
