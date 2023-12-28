package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBoard(t *testing.T) {
	t.Parallel()

	suite.Run(t, &BoardTest{})
}

type BoardTest struct {
	suite.Suite
}

func SetupTestBoard(pieces [SquareCount]Piece, init func(*Board)) Board {
	b := Board{Pieces: pieces}

	for src, piece := range b.Pieces {
		src := Square(src)

		if piece.Is(PieceNone) {
			continue
		} else if piece.Is(PieceKing) {
			b.Kings[piece.Color()] = src
		}

		b.Bitboards.Colors[piece.Color()].Set(src.Bitboard())
		b.Bitboards.Pieces[piece.Kind()].Set(src.Bitboard())
	}

	b.Bitboards.All = b.Bitboards.Colors[ColorBlack] | b.Bitboards.Colors[ColorWhite]

	if init != nil {
		init(&b)
	}

	b.Attacks.Generate(&b, b.Player.Opponent())

	return b
}

func (t *BoardTest) TestNewBoard() {
	t.Run("standard start position can be parsed", func() {
		board, err := NewBoard(BoardStartPositionFEN)

		t.Assert().NoError(err)

		t.Assert().Equal(Board{
			Player: ColorWhite,
			Attacks: Attacks{
				All: (BitboardRank6 | BitboardRank7 | BitboardRank8) & ^(SquareA8.Bitboard() | SquareH8.Bitboard()),
			},
			Pieces: [SquareCount]Piece{
				SquareA8: PieceBlackRook,
				SquareB8: PieceBlackKnight,
				SquareC8: PieceBlackBishop,
				SquareD8: PieceBlackQueen,
				SquareE8: PieceBlackKing,
				SquareF8: PieceBlackBishop,
				SquareG8: PieceBlackKnight,
				SquareH8: PieceBlackRook,
				SquareA7: PieceBlackPawn,
				SquareB7: PieceBlackPawn,
				SquareC7: PieceBlackPawn,
				SquareD7: PieceBlackPawn,
				SquareE7: PieceBlackPawn,
				SquareF7: PieceBlackPawn,
				SquareG7: PieceBlackPawn,
				SquareH7: PieceBlackPawn,
				SquareA1: PieceWhiteRook,
				SquareB1: PieceWhiteKnight,
				SquareC1: PieceWhiteBishop,
				SquareD1: PieceWhiteQueen,
				SquareE1: PieceWhiteKing,
				SquareF1: PieceWhiteBishop,
				SquareG1: PieceWhiteKnight,
				SquareH1: PieceWhiteRook,
				SquareA2: PieceWhitePawn,
				SquareB2: PieceWhitePawn,
				SquareC2: PieceWhitePawn,
				SquareD2: PieceWhitePawn,
				SquareE2: PieceWhitePawn,
				SquareF2: PieceWhitePawn,
				SquareG2: PieceWhitePawn,
				SquareH2: PieceWhitePawn,
			},
			Kings: ColorTable[Square]{
				ColorWhite: SquareE1,
				ColorBlack: SquareE8,
			},
			Bitboards: struct {
				All    Bitboard
				Colors ColorTable[Bitboard]
				Pieces PieceTable[Bitboard]
			}{
				All: BitboardRank1 | BitboardRank2 | BitboardRank7 | BitboardRank8,
				Colors: ColorTable[Bitboard]{
					ColorWhite: BitboardRank1 | BitboardRank2,
					ColorBlack: BitboardRank7 | BitboardRank8,
				},
				Pieces: PieceTable[Bitboard]{
					PiecePawn:   BitboardRank2 | BitboardRank7,
					PieceKnight: SquareB1.Bitboard() | SquareG1.Bitboard() | SquareB8.Bitboard() | SquareG8.Bitboard(),
					PieceBishop: SquareC1.Bitboard() | SquareF1.Bitboard() | SquareC8.Bitboard() | SquareF8.Bitboard(),
					PieceRook:   SquareA1.Bitboard() | SquareH1.Bitboard() | SquareA8.Bitboard() | SquareH8.Bitboard(),
					PieceQueen:  SquareD1.Bitboard() | SquareD8.Bitboard(),
					PieceKing:   SquareE1.Bitboard() | SquareE8.Bitboard(),
				},
			},
			Castling: ColorTable[struct{ Kingside, Queenside bool }]{
				ColorWhite: {Kingside: true, Queenside: true},
				ColorBlack: {Kingside: true, Queenside: true},
			},
			EnPassant: 0,
			Moves: struct {
				Half int
				Full int
				Last Move
			}{
				Half: 0,
				Full: 1,
				Last: Move{},
			},
		}, board)

		t.Assert().Equal(BoardStartPositionFEN, board.FEN())
	})

	t.Run("invalid fens result in an error", func() {
		for _, test := range []struct {
			fen string
			err string
		}{
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1", "invalid number of fields: expected 6, got 5"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 1 extra", "invalid number of fields: expected 6, got 7"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - invalid 1", "invalid half moves value"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 invalid", "invalid full moves value"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR d KQkq - 1 1", "invalid player: d"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkqd - 1 1", "invalid castling rights: d"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b - a3d 1 1", "invalid en passant target: a3d"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b - a9 1 1", "invalid en passant target: a9"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP b - a3 1 1", "invalid number of ranks: expected 8, got 7"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNdQKBNR/ b - a3 1 1", "invalid number of ranks: expected 8, got 9"},
			{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNdQKBNR b - a3 1 1", "invalid piece: d"},
		} {
			t.Run(test.err, func() {
				_, err := NewBoard(test.fen)

				t.Assert().ErrorIs(err, ErrInvalidFEN)
				t.Assert().ErrorContains(err, test.err)
			})
		}
	})

	raw, err := os.ReadFile("testdata/fens.json")
	if err != nil {
		panic(err)
	}

	fens := []string(nil)

	if err := json.Unmarshal(raw, &fens); err != nil {
		panic(err)
	}

	for _, fen := range fens {
		t.Run(fen, func() {
			board, err := NewBoard(fen)

			t.Assert().NoError(err)
			t.Assert().Equal(fen, board.FEN())
		})
	}
}

func (t *BoardTest) TestMakeMove() {
	for _, test := range []struct {
		scenario string
		board    Board
		move     Move
		assert   func(Board)
	}{
		{
			scenario: "making a move updates player",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE2: PieceWhitePawn,
			}, nil),
			move: NewMove(SquareE2, SquareE3),
			assert: func(b Board) {
				t.Assert().Equal(ColorWhite, b.Player)
			},
		},

		{
			scenario: "pawn moves clear half move counter",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE2: PieceWhitePawn,
			}, func(b *Board) { b.Moves.Half = 10 }),
			move: NewMove(SquareE2, SquareE3),
			assert: func(b Board) {
				t.Assert().Equal(1, b.Moves.Half)
			},
		},

		{
			scenario: "captures clear half move counter",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE2: PieceWhitePawn,
			}, func(b *Board) { b.Moves.Half = 10 }),
			move: NewMove(SquareE2, SquareD3, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().Equal(1, b.Moves.Half)
			},
		},

		{
			scenario: "black moving increments full move counter",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE2: PieceWhitePawn,
			}, func(b *Board) { b.Moves.Full = 10 }),
			move: NewMove(SquareE2, SquareD3, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().Equal(11, b.Moves.Full)
			},
		},

		{
			scenario: "king moves updates king position and castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[b.Player].Kingside = true
				b.Castling[b.Player].Queenside = true
			}),
			move: NewMove(SquareE1, SquareD1),
			assert: func(b Board) {
				t.Assert().Equal(SquareD1, b.Kings[ColorWhite])
				t.Assert().False(b.Castling[ColorWhite].Kingside)
				t.Assert().False(b.Castling[ColorWhite].Queenside)
			},
		},

		{
			scenario: "castling white kingside moves rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
				SquareH1: PieceWhiteRook,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareE1, SquareG1, MoveFlagsCastleKingside),
			assert: func(b Board) {
				t.Assert().Equal(PieceEmpty, b.Pieces[SquareH1])
				t.Assert().Equal(PieceWhiteRook, b.Pieces[SquareF1])
				t.Assert().Equal(SquareF1.Bitboard()|SquareG1.Bitboard(), b.Bitboards.All)
				t.Assert().Equal(SquareF1.Bitboard()|SquareG1.Bitboard(), b.Bitboards.Colors[ColorWhite])
				t.Assert().Equal(SquareF1.Bitboard(), b.Bitboards.Pieces[PieceRook])
			},
		},

		{
			scenario: "castling white queenside moves rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE1: PieceWhiteKing,
				SquareA1: PieceWhiteRook,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareE1, SquareC1, MoveFlagsCastleQueenside),
			assert: func(b Board) {
				t.Assert().Equal(PieceEmpty, b.Pieces[SquareA1])
				t.Assert().Equal(PieceWhiteRook, b.Pieces[SquareD1])
				t.Assert().Equal(SquareC1.Bitboard()|SquareD1.Bitboard(), b.Bitboards.All)
				t.Assert().Equal(SquareC1.Bitboard()|SquareD1.Bitboard(), b.Bitboards.Colors[ColorWhite])
				t.Assert().Equal(SquareD1.Bitboard(), b.Bitboards.Pieces[PieceRook])
			},
		},

		{
			scenario: "castling black kingside moves rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE8: PieceBlackKing,
				SquareH8: PieceBlackRook,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareE8, SquareG8, MoveFlagsCastleKingside),
			assert: func(b Board) {
				t.Assert().Equal(PieceEmpty, b.Pieces[SquareH8])
				t.Assert().Equal(PieceBlackRook, b.Pieces[SquareF8])
				t.Assert().Equal(SquareF8.Bitboard()|SquareG8.Bitboard(), b.Bitboards.All)
				t.Assert().Equal(SquareF8.Bitboard()|SquareG8.Bitboard(), b.Bitboards.Colors[ColorBlack])
				t.Assert().Equal(SquareF8.Bitboard(), b.Bitboards.Pieces[PieceRook])
			},
		},

		{
			scenario: "castling black queenside moves rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareE8: PieceBlackKing,
				SquareA8: PieceBlackRook,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareE8, SquareC8, MoveFlagsCastleQueenside),
			assert: func(b Board) {
				t.Assert().Equal(PieceEmpty, b.Pieces[SquareA8])
				t.Assert().Equal(PieceBlackRook, b.Pieces[SquareD8])
				t.Assert().Equal(SquareC8.Bitboard()|SquareD8.Bitboard(), b.Bitboards.All)
				t.Assert().Equal(SquareC8.Bitboard()|SquareD8.Bitboard(), b.Bitboards.Colors[ColorBlack])
				t.Assert().Equal(SquareD8.Bitboard(), b.Bitboards.Pieces[PieceRook])
			},
		},

		{
			scenario: "white promotion to queen",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB7: PieceWhitePawn,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareB7, SquareB8, MoveFlagsPromoteToQueen),
			assert: func(b Board) {
				t.Assert().Equal(PieceWhiteQueen, b.Pieces[SquareB8])
				t.Assert().Equal(SquareB8.Bitboard(), b.Bitboards.Pieces[PieceQueen])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "white promotion to rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB7: PieceWhitePawn,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareB7, SquareB8, MoveFlagsPromoteToRook),
			assert: func(b Board) {
				t.Assert().Equal(PieceWhiteRook, b.Pieces[SquareB8])
				t.Assert().Equal(SquareB8.Bitboard(), b.Bitboards.Pieces[PieceRook])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "white promotion to bishop",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB7: PieceWhitePawn,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareB7, SquareB8, MoveFlagsPromoteToBishop),
			assert: func(b Board) {
				t.Assert().Equal(PieceWhiteBishop, b.Pieces[SquareB8])
				t.Assert().Equal(SquareB8.Bitboard(), b.Bitboards.Pieces[PieceBishop])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "white promotion to knight",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB7: PieceWhitePawn,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareB7, SquareB8, MoveFlagsPromoteToKnight),
			assert: func(b Board) {
				t.Assert().Equal(PieceWhiteKnight, b.Pieces[SquareB8])
				t.Assert().Equal(SquareB8.Bitboard(), b.Bitboards.Pieces[PieceKnight])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "black promotion to queen",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB2: PieceBlackPawn,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareB2, SquareB1, MoveFlagsPromoteToQueen),
			assert: func(b Board) {
				t.Assert().Equal(PieceBlackQueen, b.Pieces[SquareB1])
				t.Assert().Equal(SquareB1.Bitboard(), b.Bitboards.Pieces[PieceQueen])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "black promotion to rook",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB2: PieceBlackPawn,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareB2, SquareB1, MoveFlagsPromoteToRook),
			assert: func(b Board) {
				t.Assert().Equal(PieceBlackRook, b.Pieces[SquareB1])
				t.Assert().Equal(SquareB1.Bitboard(), b.Bitboards.Pieces[PieceRook])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "black promotion to bishop",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB2: PieceBlackPawn,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareB2, SquareB1, MoveFlagsPromoteToBishop),
			assert: func(b Board) {
				t.Assert().Equal(PieceBlackBishop, b.Pieces[SquareB1])
				t.Assert().Equal(SquareB1.Bitboard(), b.Bitboards.Pieces[PieceBishop])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "black promotion to knight",
			board: SetupTestBoard([SquareCount]Piece{
				SquareB2: PieceBlackPawn,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareB2, SquareB1, MoveFlagsPromoteToKnight),
			assert: func(b Board) {
				t.Assert().Equal(PieceBlackKnight, b.Pieces[SquareB1])
				t.Assert().Equal(SquareB1.Bitboard(), b.Bitboards.Pieces[PieceKnight])
				t.Assert().EqualValues(0, b.Bitboards.Pieces[PiecePawn])
			},
		},

		{
			scenario: "white double pawn push updates en passant target",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA2: PieceWhitePawn,
			}, func(b *Board) { b.Player = ColorWhite }),
			move: NewMove(SquareA2, SquareA4, MoveFlagsDoublePawnPush),
			assert: func(b Board) {
				t.Assert().Equal(SquareA3, b.EnPassant)
			},
		},

		{
			scenario: "black double pawn push updates en passant target",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD7: PieceBlackPawn,
			}, func(b *Board) { b.Player = ColorBlack }),
			move: NewMove(SquareD7, SquareD5, MoveFlagsDoublePawnPush),
			assert: func(b Board) {
				t.Assert().Equal(SquareD6, b.EnPassant)
			},
		},

		{
			scenario: "moving white kingside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareH1: PieceWhiteRook,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorWhite].Kingside = true
				b.Castling[ColorWhite].Queenside = true
			}),
			move: NewMove(SquareH1, SquareA1),
			assert: func(b Board) {
				t.Assert().False(b.Castling[ColorWhite].Kingside)
				t.Assert().True(b.Castling[ColorWhite].Queenside)
			},
		},

		{
			scenario: "moving black kingside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareH8: PieceBlackRook,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Kingside = true
				b.Castling[ColorBlack].Queenside = true
			}),
			move: NewMove(SquareH8, SquareA8),
			assert: func(b Board) {
				t.Assert().False(b.Castling[ColorBlack].Kingside)
				t.Assert().True(b.Castling[ColorBlack].Queenside)
			},
		},

		{
			scenario: "moving white queenside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA1: PieceWhiteRook,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorWhite].Kingside = true
				b.Castling[ColorWhite].Queenside = true
			}),
			move: NewMove(SquareA1, SquareA2),
			assert: func(b Board) {
				t.Assert().True(b.Castling[ColorWhite].Kingside)
				t.Assert().False(b.Castling[ColorWhite].Queenside)
			},
		},

		{
			scenario: "moving black queenside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackRook,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorBlack].Kingside = true
				b.Castling[ColorBlack].Queenside = true
			}),
			move: NewMove(SquareA8, SquareA2),
			assert: func(b Board) {
				t.Assert().True(b.Castling[ColorBlack].Kingside)
				t.Assert().False(b.Castling[ColorBlack].Queenside)
			},
		},

		{
			scenario: "capturing white kingside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareH8: PieceBlackRook,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorWhite].Kingside = true
				b.Castling[ColorWhite].Queenside = true
			}),
			move: NewMove(SquareH8, SquareH1, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().False(b.Castling[ColorWhite].Kingside)
				t.Assert().True(b.Castling[ColorWhite].Queenside)
			},
		},

		{
			scenario: "capturing black kingside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareH1: PieceWhiteRook,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorBlack].Kingside = true
				b.Castling[ColorBlack].Queenside = true
			}),
			move: NewMove(SquareH1, SquareH8, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().False(b.Castling[ColorBlack].Kingside)
				t.Assert().True(b.Castling[ColorBlack].Queenside)
			},
		},

		{
			scenario: "capturing white queenside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareA8: PieceBlackRook,
			}, func(b *Board) {
				b.Player = ColorBlack
				b.Castling[ColorWhite].Kingside = true
				b.Castling[ColorWhite].Queenside = true
			}),
			move: NewMove(SquareA8, SquareA1, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().True(b.Castling[ColorWhite].Kingside)
				t.Assert().False(b.Castling[ColorWhite].Queenside)
			},
		},

		{
			scenario: "capturing black queenside rook updates castling rights",
			board: SetupTestBoard([SquareCount]Piece{
				SquareH1: PieceWhiteRook,
			}, func(b *Board) {
				b.Player = ColorWhite
				b.Castling[ColorBlack].Kingside = true
				b.Castling[ColorBlack].Queenside = true
			}),
			move: NewMove(SquareA1, SquareA8, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().True(b.Castling[ColorBlack].Kingside)
				t.Assert().False(b.Castling[ColorBlack].Queenside)
			},
		},

		{
			scenario: "bitboards are updated",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD7: PieceBlackKing,
				SquareH2: PieceBlackRook,
				SquareE1: PieceWhiteKing,
				SquareH1: PieceWhiteBishop,
			}, nil),
			move: NewMove(SquareH2, SquareH1, MoveFlagsCapture),
			assert: func(b Board) {
				expected := struct {
					All    Bitboard
					Colors ColorTable[Bitboard]
					Pieces PieceTable[Bitboard]
				}{
					All: SquareD7.Bitboard() | SquareH1.Bitboard() | SquareE1.Bitboard(),
					Colors: ColorTable[Bitboard]{
						ColorBlack: SquareD7.Bitboard() | SquareH1.Bitboard(),
						ColorWhite: SquareE1.Bitboard(),
					},
					Pieces: PieceTable[Bitboard]{
						PieceKing: SquareD7.Bitboard() | SquareE1.Bitboard(),
						PieceRook: SquareH1.Bitboard(),
					},
				}

				t.Assert().Equal(expected, b.Bitboards)
			},
		},

		{
			scenario: "attacks are updated",
			board: SetupTestBoard([SquareCount]Piece{
				SquareD7: PieceBlackKing,
				SquareH2: PieceBlackRook,
				SquareE1: PieceWhiteKing,
				SquareH1: PieceWhiteBishop,
			}, nil),
			move: NewMove(SquareH2, SquareH1, MoveFlagsCapture),
			assert: func(b Board) {
				t.Assert().Equal(Attacks{
					All: BitboardFileH&^SquareH1.Bitboard() |
						SquareG1.Bitboard() | SquareF1.Bitboard() | SquareE1.Bitboard() | SquareD1.Bitboard() |
						SquareC8.Bitboard() | SquareD8.Bitboard() | SquareE8.Bitboard() |
						SquareC7.Bitboard() | SquareE7.Bitboard() |
						SquareC6.Bitboard() | SquareD6.Bitboard() | SquareE6.Bitboard(),
					Checks: struct {
						Check  bool
						Double bool
						Rays   Bitboard
					}{
						Check: true,
						Rays:  SquareH1.Bitboard() | SquareG1.Bitboard() | SquareF1.Bitboard() | SquareE1.Bitboard(),
					},
				}, b.Attacks)
			},
		},
	} {
		t.Run(test.scenario, func() {
			board := test.board.MakeMove(test.move)

			test.assert(board)

			t.Assert().Equal(test.move, board.Moves.Last)
		})
	}
}

func (t *BoardTest) TestGenerateMoves() {
	board, err := NewBoard(BoardStartPositionFEN)
	if err != nil {
		panic(err)
	}

	t.Assert().Len(board.GenerateMoves(MoveGeneratorOptions{}), 20)
}
