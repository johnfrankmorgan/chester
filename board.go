package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	Player    Color
	Squares   [SquareCount]Piece
	Bits      BoardBitboards
	Kings     [ColorCount]Square
	Castling  [ColorCount]BoardCastlingRights
	EnPassant Square
	Attacks   Attacks
	Moves     BoardMoves
	Zobrist   Zobrist
}

type BoardBitboards struct {
	All     Bitboard
	Players [ColorCount]Bitboard
	Pieces  [PieceTypeCount + 1]Bitboard
}

type BoardCastlingRights struct {
	Kingside  bool
	Queenside bool
}

type BoardMoves struct {
	Half int
	Full int
	Last Move
}

const BoardStartPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var ErrInvalidFEN = fmt.Errorf("invalid fen")

func BoardFromFEN(fen string) (Board, error) {
	b := Board{}

	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return Board{}, fmt.Errorf("%w: expected 6 fields, got %d", ErrInvalidFEN, len(fields))
	}

	pieces := fields[0]

	for rank := range RanksReversed() {
		if rank < RankLast {
			if pieces[0] != '/' {
				return Board{}, fmt.Errorf("%w: missing slash: %s", ErrInvalidFEN, pieces)
			}

			pieces = pieces[1:]
		}

		for file := FileFirst; file <= FileLast; file++ {
			ch := pieces[0]
			pieces = pieces[1:]

			if ch >= '1' && ch <= '8' {
				file += File(ch-'0') - 1
				continue
			}

			pc, ok := PieceFromString(string(ch))
			if !ok {
				return Board{}, fmt.Errorf("%w: invalid piece: %s", ErrInvalidFEN, string(ch))
			}

			sq := NewSquare(file, rank)

			b.Squares[sq] = pc

			b.Bits.Players[pc.Color()] = b.Bits.Players[pc.Color()].Occupy(sq)
			b.Bits.Pieces[pc.Type()] = b.Bits.Pieces[pc.Type()].Occupy(sq)

			if pc.Type() == King {
				b.Kings[pc.Color()] = sq
			}
		}
	}

	b.Bits.All = b.Bits.Players[White].Set(b.Bits.Players[Black])

	player, ok := ColorFromString(fields[1])
	if !ok {
		return Board{}, fmt.Errorf("%w: invalid player: %s", ErrInvalidFEN, fields[1])
	}

	b.Player = player

	if fields[2] != "-" {
		for _, c := range fields[2] {
			switch c {
			case 'k':
				b.Castling[Black].Kingside = true

			case 'q':
				b.Castling[Black].Queenside = true

			case 'K':
				b.Castling[White].Kingside = true

			case 'Q':
				b.Castling[White].Queenside = true

			default:
				return Board{}, fmt.Errorf("%w: invalid castling: %s", ErrInvalidFEN, fields[2])
			}
		}
	}

	if fields[3] != "-" {
		if b.EnPassant, ok = SquareFromString(fields[3]); !ok {
			return Board{}, fmt.Errorf("%w: invalid en passant target: %s", ErrInvalidFEN, fields[3])
		}
	}

	half, err := strconv.Atoi(fields[4])
	if err != nil {
		return Board{}, fmt.Errorf("%w: invalid half moves: %s", ErrInvalidFEN, fields[4])
	}

	b.Moves.Half = half

	full, err := strconv.Atoi(fields[5])
	if err != nil {
		return Board{}, fmt.Errorf("%w: invalid full moves: %s", ErrInvalidFEN, fields[5])
	}

	b.Moves.Full = full

	b.Attacks = GenerateAttacks(&b, b.Player.Opponent())
	b.Zobrist = CalculateZobrist(&b)

	return b, nil
}

func (b Board) String() string {
	s := strings.Builder{}

	s.WriteString("       a   b   c   d   e   f   g   h\n\n")
	s.WriteString("     +---+---+---+---+---+---+---+---+\n")

	for rank := range RanksReversed() {
		if rank < RankLast {
			s.WriteByte('\n')
		}

		s.WriteString("  ")
		s.WriteString(rank.String())
		s.WriteString(" ")

		for file := range Files() {
			s.WriteString(" | ")

			if piece := b.Squares[NewSquare(file, rank)]; piece != EmptySquare {
				s.WriteString(piece.String())
			} else {
				s.WriteByte(' ')
			}
		}

		s.WriteString(" |  ")
		s.WriteString(rank.String())
		s.WriteString("\n     +---+---+---+---+---+---+---+---+")
	}

	s.WriteString("\n\n       a   b   c   d   e   f   g   h\n\n")
	fmt.Fprintf(&s, "   Player: %s\n", b.Player)
	fmt.Fprintf(&s, "     Move: %s\n", b.Moves.Last)
	fmt.Fprintf(&s, "     Half: %d\n", b.Moves.Half)
	fmt.Fprintf(&s, "     Full: %d\n", b.Moves.Full)
	fmt.Fprintf(&s, "EnPassant: %s\n", b.EnPassant)
	fmt.Fprintf(&s, "  Zobrist: %d\n", b.Zobrist)

	return s.String()
}

func (b Board) MakeMove(move Move) Board {
	piece := b.Squares[move.From]
	ptype := piece.Type()
	color := piece.Color()

	if b.EnPassant != 0 {
		b.Zobrist ^= Zobrists.EnPassant[b.EnPassant.File()]
		b.EnPassant = 0
	}

	b.Zobrist ^= Zobrists.Players[b.Player]
	b.Zobrist ^= Zobrists.Pieces[color][ptype][move.From]
	b.Zobrist ^= Zobrists.Castling[CastlingZobristIndex(&b)]

	if ptype == Pawn || move.Flags&MoveFlagCapture != 0 {
		b.Moves.Half = 1
	} else {
		b.Moves.Half++
	}

	switch ptype {
	case King:
		b.Kings[color] = move.To
		b.Castling[color] = BoardCastlingRights{}

		if move.Flags&MoveFlagCastleAny != 0 {
			rook := Move{}

			if move.Flags&MoveFlagCastleKingside != 0 {
				if color == White {
					rook = Move{From: SquareH1, To: SquareF1}
				} else {
					rook = Move{From: SquareH8, To: SquareF8}
				}
			} else if move.Flags&MoveFlagCastleQueenside != 0 {
				if color == White {
					rook = Move{From: SquareA1, To: SquareD1}
				} else {
					rook = Move{From: SquareA8, To: SquareD8}
				}
			}

			b.Squares[rook.To] = b.Squares[rook.From]
			b.Squares[rook.From] = EmptySquare

			b.Bits.Players[color] = b.Bits.Players[color].Unoccupy(rook.From)
			b.Bits.Players[color] = b.Bits.Players[color].Occupy(rook.To)

			b.Bits.Pieces[Rook] = b.Bits.Pieces[Rook].Unoccupy(rook.From)
			b.Bits.Pieces[Rook] = b.Bits.Pieces[Rook].Occupy(rook.To)

			b.Zobrist ^= Zobrists.Pieces[color][Rook][rook.From]
			b.Zobrist ^= Zobrists.Pieces[color][Rook][rook.To]
		}

	case Rook:
		if move.From == SquareA1 || move.From == SquareA8 {
			b.Castling[color].Queenside = false
		} else if move.From == SquareH1 || move.From == SquareH8 {
			b.Castling[color].Kingside = false
		}

	case Pawn:
		if promotion, ok := move.Promotion(); ok {
			b.Bits.Pieces[Pawn] = b.Bits.Pieces[Pawn].Unoccupy(move.From)

			piece = NewPiece(color, promotion)
			ptype = promotion
		} else if move.Flags&MoveFlagDoublePawnPush != 0 {
			if color == White {
				b.EnPassant = move.From + North.Offset()
			} else {
				b.EnPassant = move.From + South.Offset()
			}

			b.Zobrist ^= Zobrists.EnPassant[b.EnPassant]
		} else if move.Flags&MoveFlagCaptureEnPassant != 0 {
			target := move.To

			if color == White {
				target += South.Offset()
			} else {
				target += North.Offset()
			}

			b.Squares[target] = EmptySquare
			b.Bits.All = b.Bits.All.Unoccupy(target)
			b.Bits.Players[color.Opponent()] = b.Bits.Players[color.Opponent()].Unoccupy(target)
			b.Bits.Pieces[Pawn] = b.Bits.Pieces[Pawn].Unoccupy(target)
		}
	}

	if move.Flags&MoveFlagCapture != 0 {
		switch move.To {
		case SquareA1:
			b.Castling[White].Queenside = false

		case SquareH1:
			b.Castling[White].Kingside = false

		case SquareA8:
			b.Castling[Black].Queenside = false

		case SquareH8:
			b.Castling[Black].Kingside = false
		}

		b.Zobrist ^= Zobrists.Pieces[color.Opponent()][b.Squares[move.To].Type()][move.To]
	}

	b.Squares[move.From] = EmptySquare
	b.Squares[move.To] = piece

	b.Bits.Players[color] = b.Bits.Players[color].Unoccupy(move.From)
	b.Bits.Players[color] = b.Bits.Players[color].Occupy(move.To)
	b.Bits.Players[color.Opponent()] = b.Bits.Players[color.Opponent()].Unoccupy(move.To)

	for p := range b.Bits.Pieces {
		b.Bits.Pieces[p] = b.Bits.Pieces[p].Unoccupy(move.To)
	}

	b.Bits.Pieces[ptype] = b.Bits.Pieces[ptype].Unoccupy(move.From)
	b.Bits.Pieces[ptype] = b.Bits.Pieces[ptype].Occupy(move.To)
	b.Bits.All = b.Bits.Players[Black].Set(b.Bits.Players[White])

	if color == Black {
		b.Moves.Full++
	}

	b.Moves.Last = move

	b.Attacks = GenerateAttacks(&b, b.Player)
	b.Player = b.Player.Opponent()

	b.Zobrist ^= Zobrists.Pieces[color][ptype][move.To]
	b.Zobrist ^= Zobrists.Castling[CastlingZobristIndex(&b)]
	b.Zobrist ^= Zobrists.Players[b.Player]

	return b
}
