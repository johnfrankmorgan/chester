package move

import (
	"fmt"
	"strings"

	"github.com/johnfrankmorgan/chester/internal/chester/piece"
	"github.com/johnfrankmorgan/chester/internal/chester/square"
)

type Move struct {
	From  square.Square
	To    square.Square
	Flags Flags
}

func New(from, to square.Square, flags ...Flags) Move {
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
	return m.From.Valid() || m.To.Valid()
}

func (m Move) Promotion() piece.Kind {
	if m.Flags.AnySet(FlagsPromote) {
		switch {
		case m.Flags.IsSet(FlagsPromoteToQueen):
			return piece.Queen

		case m.Flags.IsSet(FlagsPromoteToRook):
			return piece.Rook

		case m.Flags.IsSet(FlagsPromoteToBishop):
			return piece.Bishop

		case m.Flags.IsSet(FlagsPromoteToKnight):
			return piece.Knight
		}
	}

	return piece.None
}
