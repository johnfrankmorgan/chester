package move

import "strings"

type Flags uint16

const (
	FlagsCapture Flags = 1 << iota
	FlagsCaptureEnPassant
	FlagsDoublePawnPush
	FlagsCastleKingside
	FlagsCastleQueenside
	FlagsPromoteToQueen
	FlagsPromoteToRook
	FlagsPromoteToBishop
	FlagsPromoteToKnight

	FlagsCastle  = FlagsCastleKingside | FlagsCastleQueenside
	FlagsPromote = FlagsPromoteToQueen |
		FlagsPromoteToRook |
		FlagsPromoteToBishop |
		FlagsPromoteToKnight
)

func (f Flags) String() string {
	s := strings.Builder{}

	for i, ch := range "ce2KQqrbn" {
		if f.IsSet(1 << i) {
			s.WriteRune(ch)
		}
	}

	return s.String()
}

func (f Flags) IsSet(flags Flags) bool {
	return f&flags == flags
}

func (f Flags) AnySet(flags Flags) bool {
	return f&flags > 0
}
