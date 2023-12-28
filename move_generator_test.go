package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMoveGenerator(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MoveGeneratorTest{})
}

type MoveGeneratorTest struct {
	suite.Suite
}

func (t *MoveGeneratorTest) TestGenerate() {
	for _, test := range []struct {
		scenario string
		board    Board
		opts     MoveGeneratorOptions
		expected []Move
	}{
		{
			scenario: "king moves",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE3: PieceWhiteKing,
				SquareF3: PieceWhitePawn,
				SquareD2: PieceBlackRook,
				SquareC5: PieceBlackBishop,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareE3, SquareE4),
				NewMove(SquareE3, SquareF4),
				NewMove(SquareE3, SquareD2, MoveFlagsCapture),
			},
		},

		{
			scenario: "capturing king moves",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE3: PieceWhiteKing,
				SquareF3: PieceWhitePawn,
				SquareD2: PieceBlackRook,
				SquareC5: PieceBlackBishop,
			}, func(b *Board) { b.Player = ColorWhite }),
			opts: MoveGeneratorOptions{CapturesOnly: true},
			expected: []Move{
				NewMove(SquareE3, SquareD2, MoveFlagsCapture),
			},
		},

		{
			scenario: "white king can castle kingside",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
				SquareH1: PieceWhiteRook,
				SquareA8: PieceBlackKing,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorWhite].Kingside = true
			}),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareE2),
				NewMove(SquareE1, SquareD2),
				NewMove(SquareE1, SquareF2),
				NewMove(SquareH1, SquareG1),
				NewMove(SquareH1, SquareF1),
				NewMove(SquareH1, SquareH2),
				NewMove(SquareH1, SquareH3),
				NewMove(SquareH1, SquareH4),
				NewMove(SquareH1, SquareH5),
				NewMove(SquareH1, SquareH6),
				NewMove(SquareH1, SquareH7),
				NewMove(SquareH1, SquareH8),
				NewMove(SquareE1, SquareG1, MoveFlagsCastleKingside),
			},
		},

		{
			scenario: "white king can castle queenside",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteRook,
				SquareE1: PieceWhiteKing,
				SquareH8: PieceBlackKing,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorWhite].Queenside = true
			}),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareE2),
				NewMove(SquareE1, SquareD2),
				NewMove(SquareE1, SquareF2),
				NewMove(SquareA1, SquareB1),
				NewMove(SquareA1, SquareC1),
				NewMove(SquareA1, SquareD1),
				NewMove(SquareA1, SquareA2),
				NewMove(SquareA1, SquareA3),
				NewMove(SquareA1, SquareA4),
				NewMove(SquareA1, SquareA5),
				NewMove(SquareA1, SquareA6),
				NewMove(SquareA1, SquareA7),
				NewMove(SquareA1, SquareA8),
				NewMove(SquareE1, SquareC1, MoveFlagsCastleQueenside),
			},
		},

		{
			scenario: "black king can castle kingside",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareE8: PieceBlackKing,
				SquareH8: PieceBlackRook,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Kingside = true
			}),
			expected: []Move{
				NewMove(SquareE8, SquareD8),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareE7),
				NewMove(SquareE8, SquareD7),
				NewMove(SquareE8, SquareF7),
				NewMove(SquareH8, SquareG8),
				NewMove(SquareH8, SquareF8),
				NewMove(SquareH8, SquareH7),
				NewMove(SquareH8, SquareH6),
				NewMove(SquareH8, SquareH5),
				NewMove(SquareH8, SquareH4),
				NewMove(SquareH8, SquareH3),
				NewMove(SquareH8, SquareH2),
				NewMove(SquareH8, SquareH1),
				NewMove(SquareE8, SquareG8, MoveFlagsCastleKingside),
			},
		},

		{
			scenario: "black king can castle queenside",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackRook,
				SquareE8: PieceBlackKing,
				SquareH1: PieceWhiteKing,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Queenside = true
			}),
			expected: []Move{
				NewMove(SquareE8, SquareD8),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareE7),
				NewMove(SquareE8, SquareD7),
				NewMove(SquareE8, SquareF7),
				NewMove(SquareA8, SquareB8),
				NewMove(SquareA8, SquareC8),
				NewMove(SquareA8, SquareD8),
				NewMove(SquareA8, SquareA7),
				NewMove(SquareA8, SquareA6),
				NewMove(SquareA8, SquareA5),
				NewMove(SquareA8, SquareA4),
				NewMove(SquareA8, SquareA3),
				NewMove(SquareA8, SquareA2),
				NewMove(SquareA8, SquareA1),
				NewMove(SquareE8, SquareC8, MoveFlagsCastleQueenside),
			},
		},

		{
			scenario: "can't castle out of check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteRook,
				SquareH1: PieceWhiteKing,
				SquareA8: PieceBlackRook,
				SquareE8: PieceBlackKing,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Queenside = true
			}),
			expected: []Move{
				NewMove(SquareE8, SquareD8),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareD7),
				NewMove(SquareE8, SquareF7),
			},
		},

		{
			scenario: "can't castle through check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD1: PieceWhiteRook,
				SquareH1: PieceWhiteKing,
				SquareA8: PieceBlackRook,
				SquareE8: PieceBlackKing,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Queenside = true
			}),
			expected: []Move{
				NewMove(SquareE8, SquareE7),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareF7),
				NewMove(SquareA8, SquareB8),
				NewMove(SquareA8, SquareC8),
				NewMove(SquareA8, SquareD8),
				NewMove(SquareA8, SquareA7),
				NewMove(SquareA8, SquareA6),
				NewMove(SquareA8, SquareA5),
				NewMove(SquareA8, SquareA4),
				NewMove(SquareA8, SquareA3),
				NewMove(SquareA8, SquareA2),
				NewMove(SquareA8, SquareA1),
			},
		},

		{
			scenario: "can't castle into check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareC1: PieceWhiteRook,
				SquareH1: PieceWhiteKing,
				SquareA8: PieceBlackRook,
				SquareE8: PieceBlackKing,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Queenside = true
			}),
			expected: []Move{
				NewMove(SquareE8, SquareD7),
				NewMove(SquareE8, SquareE7),
				NewMove(SquareE8, SquareD8),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareF7),
				NewMove(SquareA8, SquareB8),
				NewMove(SquareA8, SquareC8),
				NewMove(SquareA8, SquareD8),
				NewMove(SquareA8, SquareA7),
				NewMove(SquareA8, SquareA6),
				NewMove(SquareA8, SquareA5),
				NewMove(SquareA8, SquareA4),
				NewMove(SquareA8, SquareA3),
				NewMove(SquareA8, SquareA2),
				NewMove(SquareA8, SquareA1),
			},
		},

		{
			scenario: "only king moves generated in double check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareG4: PieceWhiteKnight,
				SquareF5: PieceWhitePawn,
				SquareD3: PieceBlackKing,
				SquareA8: PieceBlackRook,
				SquareH8: PieceBlackBishop,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareA1, SquareB1),
			},
		},

		{
			scenario: "basic sliding moves",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackRook,
				SquareB7: PieceBlackBishop,
				SquareE8: PieceBlackKing,
				SquareG8: PieceBlackRook,
				SquareF7: PieceBlackQueen,
				SquareG7: PieceBlackBishop,
				SquareF2: PieceWhiteRook,
				SquareB1: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorBlack }),
			expected: []Move{
				// king
				NewMove(SquareE8, SquareD8),
				NewMove(SquareE8, SquareF8),
				NewMove(SquareE8, SquareD7),
				NewMove(SquareE8, SquareE7),
				// rooks
				NewMove(SquareA8, SquareB8),
				NewMove(SquareA8, SquareC8),
				NewMove(SquareA8, SquareD8),
				NewMove(SquareA8, SquareA7),
				NewMove(SquareA8, SquareA6),
				NewMove(SquareA8, SquareA5),
				NewMove(SquareA8, SquareA4),
				NewMove(SquareA8, SquareA3),
				NewMove(SquareA8, SquareA2),
				NewMove(SquareA8, SquareA1),
				NewMove(SquareG8, SquareF8),
				NewMove(SquareG8, SquareH8),
				// bishops
				NewMove(SquareB7, SquareA6),
				NewMove(SquareB7, SquareC8),
				NewMove(SquareB7, SquareC6),
				NewMove(SquareB7, SquareD5),
				NewMove(SquareB7, SquareE4),
				NewMove(SquareB7, SquareF3),
				NewMove(SquareB7, SquareG2),
				NewMove(SquareB7, SquareH1),
				NewMove(SquareG7, SquareF8),
				NewMove(SquareG7, SquareH8),
				NewMove(SquareG7, SquareH6),
				NewMove(SquareG7, SquareF6),
				NewMove(SquareG7, SquareE5),
				NewMove(SquareG7, SquareD4),
				NewMove(SquareG7, SquareC3),
				NewMove(SquareG7, SquareB2),
				NewMove(SquareG7, SquareA1),
				// queen
				NewMove(SquareF7, SquareF8),
				NewMove(SquareF7, SquareE7),
				NewMove(SquareF7, SquareD7),
				NewMove(SquareF7, SquareC7),
				NewMove(SquareF7, SquareF6),
				NewMove(SquareF7, SquareF5),
				NewMove(SquareF7, SquareF4),
				NewMove(SquareF7, SquareF3),
				NewMove(SquareF7, SquareF2, MoveFlagsCapture),
				NewMove(SquareF7, SquareG6),
				NewMove(SquareF7, SquareH5),
				NewMove(SquareF7, SquareE6),
				NewMove(SquareF7, SquareD5),
				NewMove(SquareF7, SquareC4),
				NewMove(SquareF7, SquareB3),
				NewMove(SquareF7, SquareA2),
			},
		},

		{
			scenario: "pinned sliding pieces can't move when in check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareB2: PieceWhiteRook,
				SquareA8: PieceBlackRook,
				SquareH8: PieceBlackKing,
				SquareG7: PieceBlackBishop,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareA1, SquareB1),
			},
		},

		{
			scenario: "pinned sliding pieces can move along pin mask",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareA2: PieceWhiteRook,
				SquareA8: PieceBlackRook,
				SquareH8: PieceBlackKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareA1, SquareB1),
				NewMove(SquareA1, SquareB2),
				NewMove(SquareA2, SquareA3),
				NewMove(SquareA2, SquareA4),
				NewMove(SquareA2, SquareA5),
				NewMove(SquareA2, SquareA6),
				NewMove(SquareA2, SquareA7),
				NewMove(SquareA2, SquareA8, MoveFlagsCapture),
			},
		},

		{
			scenario: "sliding pieces can block checks",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareB2: PieceWhiteRook,
				SquareA8: PieceBlackRook,
				SquareH8: PieceBlackKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareA1, SquareB1),
				NewMove(SquareB2, SquareA2),
			},
		},

		{
			scenario: "knights can block checks",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD3: PieceBlackKing,
				SquareE3: PieceBlackKnight,
				SquareD8: PieceWhiteRook,
				SquareH8: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorBlack }),
			expected: []Move{
				NewMove(SquareD3, SquareC2),
				NewMove(SquareD3, SquareC3),
				NewMove(SquareD3, SquareC4),
				NewMove(SquareD3, SquareE2),
				NewMove(SquareD3, SquareE4),
				NewMove(SquareE3, SquareD5),
			},
		},

		{
			scenario: "knights can capture checking pieces",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD3: PieceBlackKing,
				SquareE3: PieceBlackKnight,
				SquareD5: PieceWhiteRook,
				SquareH8: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorBlack }),
			expected: []Move{
				NewMove(SquareD3, SquareC2),
				NewMove(SquareD3, SquareC3),
				NewMove(SquareD3, SquareC4),
				NewMove(SquareD3, SquareE2),
				NewMove(SquareD3, SquareE4),
				NewMove(SquareE3, SquareD5, MoveFlagsCapture),
			},
		},

		{
			scenario: "pinned knights can't move",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD3: PieceBlackKing,
				SquareD5: PieceBlackKnight,
				SquareD8: PieceWhiteRook,
				SquareH8: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorBlack }),
			expected: []Move{
				NewMove(SquareD3, SquareC2),
				NewMove(SquareD3, SquareD2),
				NewMove(SquareD3, SquareE2),
				NewMove(SquareD3, SquareC3),
				NewMove(SquareD3, SquareE3),
				NewMove(SquareD3, SquareC4),
				NewMove(SquareD3, SquareD4),
				NewMove(SquareD3, SquareE4),
			},
		},

		{
			scenario: "capturing knight moves",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE4: PieceWhiteKnight,
				SquareF6: PieceBlackRook,
			}, func(b *Board) { b.Player = ColorWhite }),
			opts: MoveGeneratorOptions{CapturesOnly: true},
			expected: []Move{
				NewMove(SquareE4, SquareF6, MoveFlagsCapture),
			},
		},

		{
			scenario: "basic pawn moves",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackKing,
				SquareE7: PieceWhitePawn,
				SquareD2: PieceWhitePawn,
				SquareC2: PieceWhitePawn,
				SquareC3: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareC3, SquareB2),
				NewMove(SquareC3, SquareB3),
				NewMove(SquareC3, SquareD3),
				NewMove(SquareC3, SquareB4),
				NewMove(SquareC3, SquareC4),
				NewMove(SquareC3, SquareD4),
				NewMove(SquareE7, SquareE8, MoveFlagsPromoteToQueen),
				NewMove(SquareE7, SquareE8, MoveFlagsPromoteToRook),
				NewMove(SquareE7, SquareE8, MoveFlagsPromoteToBishop),
				NewMove(SquareE7, SquareE8, MoveFlagsPromoteToKnight),
				NewMove(SquareD2, SquareD3),
				NewMove(SquareD2, SquareD4, MoveFlagsDoublePawnPush),
			},
		},

		{
			scenario: "pinned pawns can't move when in check",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackKing,
				SquareE8: PieceBlackRook,
				SquareE2: PieceWhitePawn,
				SquareE1: PieceWhiteKing,
				SquareA1: PieceBlackRook,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareE1, SquareD2),
				NewMove(SquareE1, SquareF2),
			},
		},

		{
			scenario: "pinned pawns can capture pinning piece",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackKing,
				SquareC3: PieceBlackBishop,
				SquareD2: PieceWhitePawn,
				SquareE1: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareE2),
				NewMove(SquareE1, SquareF2),
				NewMove(SquareD2, SquareC3, MoveFlagsCapture),
			},
		},

		{
			scenario: "pawns can block checks by capturing",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackKing,
				SquareE8: PieceBlackRook,
				SquareE3: PieceBlackRook,
				SquareF2: PieceWhitePawn,
				SquareE1: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareD2),
				NewMove(SquareF2, SquareE3, MoveFlagsCapture),
			},
		},

		{
			scenario: "pawns can't move onto squares already occupied by an enemy",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackKing,
				SquareB5: PieceBlackPawn,
				SquareB4: PieceWhitePawn,
				SquareE1: PieceWhiteKing,
			}, func(b *Board) { b.Player = ColorWhite }),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareD2),
				NewMove(SquareE1, SquareE2),
				NewMove(SquareE1, SquareF2),
			},
		},

		{
			scenario: "pawns can capture en passant",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE8: PieceBlackKing,
				SquareD5: PieceBlackPawn,
				SquareC5: PieceWhitePawn,
				SquareE1: PieceWhiteKing,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.EnPassant = SquareD6
			}),
			expected: []Move{
				NewMove(SquareE1, SquareD1),
				NewMove(SquareE1, SquareF1),
				NewMove(SquareE1, SquareD2),
				NewMove(SquareE1, SquareE2),
				NewMove(SquareE1, SquareF2),
				NewMove(SquareC5, SquareC6),
				NewMove(SquareC5, SquareD6, MoveFlagsCapture, MoveFlagsCaptureEnPassant),
			},
		},

		{
			scenario: "en passant can't reveal attacks on king (rank)",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE8: PieceBlackKing,
				SquareD5: PieceBlackPawn,
				SquareH5: PieceBlackRook,
				SquareC5: PieceWhitePawn,
				SquareA5: PieceWhiteKing,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.EnPassant = SquareD6
			}),
			expected: []Move{
				NewMove(SquareA5, SquareA4),
				NewMove(SquareA5, SquareA6),
				NewMove(SquareA5, SquareB4),
				NewMove(SquareA5, SquareB5),
				NewMove(SquareA5, SquareB6),
				NewMove(SquareC5, SquareC6),
			},
		},

		{
			scenario: "en passant can't reveal attacks on king (file)",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareC8: PieceBlackKing,
				SquareC2: PieceWhiteRook,
				SquareB4: PieceWhitePawn,
				SquareC4: PieceBlackPawn,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.EnPassant = SquareB3
			}),
			expected: []Move{
				NewMove(SquareC8, SquareB8),
				NewMove(SquareC8, SquareD8),
				NewMove(SquareC8, SquareB7),
				NewMove(SquareC8, SquareC7),
				NewMove(SquareC8, SquareD7),
				NewMove(SquareC4, SquareC3),
			},
		},

		{
			scenario: "pawn captures to rank 1 or 8 result in promotion",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteKing,
				SquareC1: PieceBlackKing,
				SquareB1: PieceWhiteRook,
				SquareC2: PieceBlackPawn,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.EnPassant = SquareB3
			}),
			opts: MoveGeneratorOptions{CapturesOnly: true},
			expected: []Move{
				NewMove(SquareC2, SquareB1, MoveFlagsCapture, MoveFlagsPromoteToQueen),
				NewMove(SquareC2, SquareB1, MoveFlagsCapture, MoveFlagsPromoteToRook),
				NewMove(SquareC2, SquareB1, MoveFlagsCapture, MoveFlagsPromoteToBishop),
				NewMove(SquareC2, SquareB1, MoveFlagsCapture, MoveFlagsPromoteToKnight),
			},
		},
	} {
		t.Run(test.scenario, func() {
			moves := MoveGenerator{}.Generate(&test.board, test.opts)

			expected := make([]string, len(test.expected))
			got := make([]string, len(moves))

			for i, move := range test.expected {
				expected[i] = move.String()
			}

			for i, move := range moves {
				got[i] = move.String()
			}

			t.Assert().ElementsMatch(expected, got)
		})
	}
}
