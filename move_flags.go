package main

import "strings"

type MoveFlags uint16

const (
	MoveFlagsCapture MoveFlags = 1 << iota
	MoveFlagsCaptureEnPassant
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

func (mf MoveFlags) String() string {
	s := strings.Builder{}

	for i, c := range []byte{'c', 'e', '2', 'K', 'Q', 'q', 'r', 'b', 'n'} {
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
