package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAttacks(t *testing.T) {
	t.Parallel()

	suite.Run(t, &AttacksTest{})
}

type AttacksTest struct {
	suite.Suite
}

func (t *AttacksTest) TestGeneration() {
	for _, test := range []struct {
		scenario string
		fn       func(*Attacks, *Board, Color)
		board    Board
		attacker Color
		expected Attacks
	}{
		{
			scenario: "king on edge",
			fn:       (*Attacks)._king,
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareD1.Bitboard() |
					SquareF1.Bitboard() |
					SquareD2.Bitboard() |
					SquareE2.Bitboard() |
					SquareF2.Bitboard(),
			},
		},

		{
			scenario: "king in center",
			fn:       (*Attacks)._king,
			board: SetupTestBoard([SquareCount]Piece{
				SquareD4: PieceBlackKing,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: SquareC3.Bitboard() |
					SquareD3.Bitboard() |
					SquareE3.Bitboard() |
					SquareC4.Bitboard() |
					SquareE4.Bitboard() |
					SquareC5.Bitboard() |
					SquareD5.Bitboard() |
					SquareE5.Bitboard(),
			},
		},

		{
			scenario: "white rook",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareB1: PieceWhiteRook,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: (BitboardFileB | BitboardRank1) &^ SquareB1.Bitboard(),
			},
		},

		{
			scenario: "black bishop",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareE3: PieceBlackBishop,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: (DirectionNorthWest.Mask(SquareE3) | DirectionNorthEast.Mask(SquareE3)) &^ SquareE3.Bitboard(),
			},
		},

		{
			scenario: "white queens",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareE4: PieceWhiteQueen,
				SquareE3: PieceWhiteQueen,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: (DirectionNorthWest.Mask(SquareE3)|DirectionNorthEast.Mask(SquareE3))&^SquareE3.Bitboard() |
					(DirectionNorthWest.Mask(SquareE4)|DirectionNorthEast.Mask(SquareE4))&^SquareE4.Bitboard() |
					BitboardFileE |
					BitboardRank3&^SquareE3.Bitboard() |
					BitboardRank4&^SquareE4.Bitboard(),
			},
		},

		{
			scenario: "rook check",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareA3: PieceWhiteRook,
				SquareA7: PieceBlackKing,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: BitboardRank3&^SquareA3.Bitboard() |
					SquareA1.Bitboard() |
					SquareA2.Bitboard() |
					SquareA4.Bitboard() |
					SquareA5.Bitboard() |
					SquareA6.Bitboard() |
					SquareA7.Bitboard() |
					SquareA8.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check: true,
					Rays: SquareA3.Bitboard() |
						SquareA4.Bitboard() |
						SquareA5.Bitboard() |
						SquareA6.Bitboard() |
						SquareA7.Bitboard(),
				},
			},
		},

		{
			scenario: "bishop check",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareG2: PieceBlackBishop,
				SquareF1: PieceWhiteKing,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: DirectionNorthWest.Mask(SquareG2)&^SquareG2.Bitboard() |
					SquareF1.Bitboard() |
					SquareH3.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check: true,
					Rays:  SquareG2.Bitboard() | SquareF1.Bitboard(),
				},
			},
		},

		{
			scenario: "bishop long diagonal",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareH8: PieceWhiteBishop,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: DirectionSouthWest.Mask(SquareH8) &^ SquareH8.Bitboard(),
			},
		},

		{
			scenario: "pinned piece",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareB1: PieceWhiteRook,
				SquareB5: PieceBlackBishop,
				SquareB8: PieceBlackKing,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: BitboardRank1&^SquareB1.Bitboard() |
					SquareB2.Bitboard() |
					SquareB3.Bitboard() |
					SquareB4.Bitboard() |
					SquareB5.Bitboard(),
				Pins: BitboardFileB,
			},
		},

		{
			scenario: "two pieces between attacker and king is not pinned",
			fn:       (*Attacks)._sliding,
			board: SetupTestBoard([SquareCount]Piece{
				SquareB1: PieceWhiteRook,
				SquareB5: PieceBlackBishop,
				SquareB7: PieceBlackPawn,
				SquareB8: PieceBlackKing,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: BitboardRank1&^SquareB1.Bitboard() |
					SquareB2.Bitboard() |
					SquareB3.Bitboard() |
					SquareB4.Bitboard() |
					SquareB5.Bitboard(),
			},
		},

		{
			scenario: "knight on edge",
			fn:       (*Attacks)._knight,
			board: SetupTestBoard([SquareCount]Piece{
				SquareB1: PieceWhiteKnight,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareA3.Bitboard() | SquareC3.Bitboard() | SquareD2.Bitboard(),
			},
		},

		{
			scenario: "knight in center",
			fn:       (*Attacks)._knight,
			board: SetupTestBoard([SquareCount]Piece{
				SquareF4: PieceWhiteKnight,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareH3.Bitboard() |
					SquareH5.Bitboard() |
					SquareE2.Bitboard() |
					SquareE6.Bitboard() |
					SquareG2.Bitboard() |
					SquareG6.Bitboard() |
					SquareD3.Bitboard() |
					SquareD5.Bitboard(),
			},
		},

		{
			scenario: "knight check",
			fn:       (*Attacks)._knight,
			board: SetupTestBoard([SquareCount]Piece{
				SquareF4: PieceWhiteKnight,
				SquareH3: PieceBlackKing,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareH3.Bitboard() |
					SquareH5.Bitboard() |
					SquareE2.Bitboard() |
					SquareE6.Bitboard() |
					SquareG2.Bitboard() |
					SquareG6.Bitboard() |
					SquareD3.Bitboard() |
					SquareD5.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check: true,
					Rays:  SquareF4.Bitboard() | SquareH3.Bitboard(),
				},
			},
		},

		{
			scenario: "white pawns",
			fn:       (*Attacks)._pawn,
			board: SetupTestBoard([SquareCount]Piece{
				SquareD4: PieceWhitePawn,
				SquareE4: PieceWhitePawn,
			}, nil),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareC5.Bitboard() |
					SquareD5.Bitboard() |
					SquareE5.Bitboard() |
					SquareF5.Bitboard(),
			},
		},

		{
			scenario: "black pawns",
			fn:       (*Attacks)._pawn,
			board: SetupTestBoard([SquareCount]Piece{
				SquareD3: PieceBlackPawn,
				SquareE4: PieceBlackPawn,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: SquareC2.Bitboard() |
					SquareE2.Bitboard() |
					SquareD3.Bitboard() |
					SquareF3.Bitboard(),
			},
		},

		{
			scenario: "pawn check",
			fn:       (*Attacks)._pawn,
			board: SetupTestBoard([SquareCount]Piece{
				SquareD3: PieceBlackPawn,
				SquareE2: PieceWhiteKing,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: SquareC2.Bitboard() |
					SquareE2.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check: true,
					Rays:  SquareD3.Bitboard() | SquareE2.Bitboard(),
				},
			},
		},

		{
			scenario: "double check",
			fn:       (*Attacks).Generate,
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareA3: PieceBlackRook,
				SquareA8: PieceBlackKing,
				SquareH8: PieceBlackBishop,
			}, nil),
			attacker: ColorBlack,
			expected: Attacks{
				All: DirectionSouthWest.Mask(SquareH8)&^SquareH8.Bitboard() |
					BitboardRank3&^SquareA3.Bitboard() |
					BitboardFileA&^SquareA3.Bitboard() |
					SquareB7.Bitboard() |
					SquareB8.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check:  true,
					Double: true,
					Rays: SquareA3.Bitboard() |
						SquareA2.Bitboard() |
						DirectionSouthWest.Mask(SquareH8),
				},
			},
		},

		{
			scenario: "multiple pieces",
			fn:       (*Attacks).Generate,
			board: func() Board {
				board, err := NewBoard("r1bqkbnr/p1p2ppp/p7/1B2pP1Q/3PP3/2P5/4K2P/RN5R b kq - 0 1")
				if err != nil {
					panic(err)
				}

				return board
			}(),
			attacker: ColorWhite,
			expected: Attacks{
				All: SquareA2.Bitboard() |
					SquareA3.Bitboard() |
					SquareA4.Bitboard() |
					SquareA5.Bitboard() |
					SquareA6.Bitboard() |
					SquareB1.Bitboard() |
					SquareB4.Bitboard() |
					SquareC1.Bitboard() |
					SquareC3.Bitboard() |
					SquareC4.Bitboard() |
					SquareC5.Bitboard() |
					SquareC6.Bitboard() |
					SquareD1.Bitboard() |
					SquareD2.Bitboard() |
					SquareD3.Bitboard() |
					SquareD4.Bitboard() |
					SquareD5.Bitboard() |
					SquareD7.Bitboard() |
					SquareE1.Bitboard() |
					SquareE2.Bitboard() |
					SquareE3.Bitboard() |
					SquareE5.Bitboard() |
					SquareE6.Bitboard() |
					SquareE8.Bitboard() |
					SquareF1.Bitboard() |
					SquareF2.Bitboard() |
					SquareF3.Bitboard() |
					SquareF5.Bitboard() |
					SquareF7.Bitboard() |
					SquareG1.Bitboard() |
					SquareG3.Bitboard() |
					SquareG4.Bitboard() |
					SquareG5.Bitboard() |
					SquareG6.Bitboard() |
					SquareH2.Bitboard() |
					SquareH3.Bitboard() |
					SquareH4.Bitboard() |
					SquareH6.Bitboard() |
					SquareH7.Bitboard(),
				Checks: struct {
					Check  bool
					Double bool
					Rays   Bitboard
				}{
					Check: true,
					Rays:  SquareB5.Bitboard() | SquareC6.Bitboard() | SquareD7.Bitboard() | SquareE8.Bitboard(),
				},
				Pins: SquareH5.Bitboard() | SquareG6.Bitboard() | SquareF7.Bitboard() | SquareE8.Bitboard(),
			},
		},
	} {
		t.Run(test.scenario, func() {
			attacks := Attacks{}

			test.fn(&attacks, &test.board, test.attacker)

			t.Assert().Equal(test.expected, attacks)
		})
	}
}

func (t *AttacksTest) TestIsAttacked() {
	attacks := Attacks{
		All: SquareA5.Bitboard() | SquareA3.Bitboard(),
	}

	t.Assert().True(attacks.IsAttacked(SquareA5))
	t.Assert().True(attacks.IsAttacked(SquareA3))
	t.Assert().False(attacks.IsAttacked(SquareB2))
}

func (t *AttacksTest) TestIsPinned() {
	attacks := Attacks{
		Pins: SquareA5.Bitboard() | SquareA3.Bitboard(),
	}

	t.Assert().True(attacks.IsPinned(SquareA5))
	t.Assert().True(attacks.IsPinned(SquareA3))
	t.Assert().False(attacks.IsPinned(SquareB2))
}
