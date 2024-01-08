package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	Player    Color
	Attacks   Attacks
	Pieces    [SquareCount]Piece
	Kings     [ColorCount]Square
	Bitboards struct {
		All    Bitboard
		Colors [ColorCount]Bitboard
		Pieces [PieceKindCount]Bitboard
	}
	Castling  [ColorCount]struct{ Kingside, Queenside bool }
	EnPassant Square
	Moves     struct {
		Half int
		Full int
		Last Move
	}
}

type Attacks struct {
	All    Bitboard
	Checks struct {
		Check  bool
		Double bool
		Rays   Bitboard
	}
	Pins Bitboard
}

const BoardStartPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var ErrInvalidFEN = errors.New("invalid fen")

func NewBoard(fen string) (Board, error) {
	fields := strings.Fields(fen)

	if len(fields) != 6 {
		return Board{}, fmt.Errorf("%w: invalid number of fields: expected 6, got %d", ErrInvalidFEN, len(fields))
	}

	pieces := fields[0]
	player := fields[1]
	castling := fields[2]
	enpassant := fields[3]

	hmoves, err := strconv.Atoi(fields[4])
	if err != nil {
		return Board{}, fmt.Errorf("%w: invalid half moves value: %w", ErrInvalidFEN, err)
	}

	fmoves, err := strconv.Atoi(fields[5])
	if err != nil {
		return Board{}, fmt.Errorf("%w: invalid full moves value: %w", ErrInvalidFEN, err)
	}

	b := Board{}

	b.Moves.Half = hmoves
	b.Moves.Full = fmoves

	switch player {
	case "b":
		b.Player = ColorBlack

	case "w":
		b.Player = ColorWhite

	default:
		return Board{}, fmt.Errorf("%w: invalid player: %s", ErrInvalidFEN, player)
	}

	if castling != "-" {
		for _, ch := range castling {
			switch ch {
			case 'q':
				b.Castling[ColorBlack].Queenside = true

			case 'Q':
				b.Castling[ColorWhite].Queenside = true

			case 'k':
				b.Castling[ColorBlack].Kingside = true

			case 'K':
				b.Castling[ColorWhite].Kingside = true

			default:
				return Board{}, fmt.Errorf("%w: invalid castling rights: %c", ErrInvalidFEN, ch)
			}
		}
	}

	if enpassant != "-" {
		if len(enpassant) != 2 {
			return Board{}, fmt.Errorf("%w: invalid en passant target: %s", ErrInvalidFEN, enpassant)
		}

		file := File(enpassant[0] - 'a')
		rank := Rank(enpassant[1] - '1')

		if !file.Valid() || !rank.Valid() {
			return Board{}, fmt.Errorf("%w: invalid en passant target: %s", ErrInvalidFEN, enpassant)
		}

		b.EnPassant = NewSquare(file, rank)
	}

	ranks := strings.Split(pieces, "/")

	if len(ranks) != RankCount {
		return Board{}, fmt.Errorf("%w: invalid number of ranks: expected %d, got %d", ErrInvalidFEN, RankCount, len(ranks))
	}

	for rank, pieces := range ranks {
		rank := RankLast - Rank(rank)

		for file := FileFirst; file <= FileLast; file++ {
			ch := pieces[0]
			pieces = pieces[1:]

			if ch > '0' && ch < '9' {
				file += File(ch-'0') - 1
				continue
			}

			piece := PieceNone

			switch ch {
			case 'k':
				piece = PieceBlackKing

			case 'q':
				piece = PieceBlackQueen

			case 'r':
				piece = PieceBlackRook

			case 'b':
				piece = PieceBlackBishop

			case 'n':
				piece = PieceBlackKnight

			case 'p':
				piece = PieceBlackPawn

			case 'K':
				piece = PieceWhiteKing

			case 'Q':
				piece = PieceWhiteQueen

			case 'R':
				piece = PieceWhiteRook

			case 'B':
				piece = PieceWhiteBishop

			case 'N':
				piece = PieceWhiteKnight

			case 'P':
				piece = PieceWhitePawn

			default:
				return Board{}, fmt.Errorf("%w: invalid piece: %c", ErrInvalidFEN, ch)
			}

			square := NewSquare(file, rank)

			b.Pieces[square] = piece
			b.Bitboards.All.Set(square.Bitboard())
			b.Bitboards.Colors[piece.Color()].Set(square.Bitboard())
			b.Bitboards.Pieces[piece.Kind()].Set(square.Bitboard())

			if piece.Is(PieceKindKing) {
				b.Kings[piece.Color()] = square
			}
		}
	}

	b.Attacks.Generate(&b, b.Player.Opponent())

	return b, nil
}

func (b *Board) FEN() string {
	s := strings.Builder{}

	for rank := RankLast; rank >= RankFirst; rank-- {
		empty := 0

		for file := FileFirst; file <= FileLast; file++ {
			piece := b.Pieces[NewSquare(file, rank)]

			if piece.Is(PieceKindNone) {
				empty++
			} else if empty > 0 {
				s.WriteByte('0' + byte(empty))
				empty = 0
			}

			switch piece {
			case PieceBlackKing:
				s.WriteByte('k')

			case PieceBlackQueen:
				s.WriteByte('q')

			case PieceBlackRook:
				s.WriteByte('r')

			case PieceBlackBishop:
				s.WriteByte('b')

			case PieceBlackKnight:
				s.WriteByte('n')

			case PieceBlackPawn:
				s.WriteByte('p')

			case PieceWhiteKing:
				s.WriteByte('K')

			case PieceWhiteQueen:
				s.WriteByte('Q')

			case PieceWhiteRook:
				s.WriteByte('R')

			case PieceWhiteBishop:
				s.WriteByte('B')

			case PieceWhiteKnight:
				s.WriteByte('N')

			case PieceWhitePawn:
				s.WriteByte('P')
			}
		}

		if empty > 0 {
			s.WriteByte('0' + byte(empty))
		}

		if rank != RankFirst {
			s.WriteByte('/')
		}
	}

	s.WriteByte(' ')

	switch b.Player {
	case ColorBlack:
		s.WriteByte('b')

	case ColorWhite:
		s.WriteByte('w')
	}

	s.WriteByte(' ')

	if !b.Castling[ColorWhite].Kingside && !b.Castling[ColorWhite].Queenside && !b.Castling[ColorBlack].Kingside && !b.Castling[ColorBlack].Queenside {
		s.WriteByte('-')
	} else {
		if b.Castling[ColorWhite].Kingside {
			s.WriteByte('K')
		}

		if b.Castling[ColorWhite].Queenside {
			s.WriteByte('Q')
		}

		if b.Castling[ColorBlack].Kingside {
			s.WriteByte('k')
		}

		if b.Castling[ColorBlack].Queenside {
			s.WriteByte('q')
		}
	}

	s.WriteByte(' ')

	if b.EnPassant == 0 {
		s.WriteByte('-')
	} else {
		s.WriteString(b.EnPassant.String())
	}

	s.WriteByte(' ')

	fmt.Fprintf(&s, "%d %d", b.Moves.Half, b.Moves.Full)

	return s.String()
}

func (b *Board) GenerateMoves(opts MoveGeneratorOptions) []Move {
	return MoveGenerator{}.Generate(b, opts)
}

func (b Board) MakeMove(move Move) Board {
	b.EnPassant = 0
	b.Moves.Last = move

	piece := b.Pieces[move.From]
	color := b.Player

	capture := !b.Pieces[move.To].Is(PieceKindNone)

	if capture || piece.Is(PieceKindPawn) {
		b.Moves.Half = 0
	}

	b.Moves.Half++

	if piece.Is(PieceKindKing) {
		b.Kings[color] = move.To

		b.Castling[color].Kingside = false
		b.Castling[color].Queenside = false

		if move.Flags.AnySet(MoveFlagsCastle) {
			rook := struct{ src, dst Square }{}

			if move.Flags.IsSet(MoveFlagsCastleKingside) {
				if color == ColorWhite {
					rook.src = SquareH1
					rook.dst = SquareF1
				} else {
					rook.src = SquareH8
					rook.dst = SquareF8
				}
			} else if move.Flags.IsSet(MoveFlagsCastleQueenside) {
				if color == ColorWhite {
					rook.src = SquareA1
					rook.dst = SquareD1
				} else {
					rook.src = SquareA8
					rook.dst = SquareD8
				}
			}

			b.Pieces[rook.dst] = b.Pieces[rook.src]
			b.Pieces[rook.src] = PieceNone

			b.Bitboards.Colors[color].Clear(rook.src.Bitboard())
			b.Bitboards.Colors[color].Set(rook.dst.Bitboard())

			b.Bitboards.Pieces[PieceKindRook].Clear(rook.src.Bitboard())
			b.Bitboards.Pieces[PieceKindRook].Set(rook.dst.Bitboard())
		}
	} else if piece.Is(PieceKindPawn) {
		if promotion := move.Promotion(); promotion != PieceKindNone {
			b.Bitboards.Pieces[PieceKindPawn].Clear(move.From.Bitboard())

			piece = NewPiece(color, promotion)
		} else if move.Flags.IsSet(MoveFlagsDoublePawnPush) {
			if color == ColorWhite {
				b.EnPassant = move.From + DirectionNorth.Offset()
			} else {
				b.EnPassant = move.From + DirectionSouth.Offset()
			}
		} else if move.Flags.IsSet(MoveFlagsEnPassant) {
			captured := move.To + DirectionSouth.Offset()

			if color == ColorBlack {
				captured = move.To + DirectionNorth.Offset()
			}

			b.Pieces[captured] = PieceNone
			b.Bitboards.All.Clear(captured.Bitboard())
			b.Bitboards.Colors[color.Opponent()].Clear(captured.Bitboard())
			b.Bitboards.Pieces[PieceKindPawn].Clear(captured.Bitboard())
		}
	} else if piece.Is(PieceKindRook) {
		if move.From == SquareA1 || move.From == SquareA8 {
			b.Castling[color].Queenside = false
		} else if move.From == SquareH1 || move.From == SquareH8 {
			b.Castling[color].Kingside = false
		}
	}

	if capture {
		switch move.To {
		case SquareA1:
			b.Castling[ColorWhite].Queenside = false

		case SquareH1:
			b.Castling[ColorWhite].Kingside = false

		case SquareA8:
			b.Castling[ColorBlack].Queenside = false

		case SquareH8:
			b.Castling[ColorBlack].Kingside = false
		}
	}

	b.Pieces[move.To] = piece
	b.Pieces[move.From] = PieceNone

	b.Bitboards.Colors[color].Clear(move.From.Bitboard())
	b.Bitboards.Colors[color].Set(move.To.Bitboard())
	b.Bitboards.Colors[color.Opponent()].Clear(move.To.Bitboard())

	for p := range b.Bitboards.Pieces {
		b.Bitboards.Pieces[p].Clear(move.To.Bitboard())
	}

	b.Bitboards.Pieces[piece.Kind()].Clear(move.From.Bitboard())
	b.Bitboards.Pieces[piece.Kind()].Set(move.To.Bitboard())

	b.Bitboards.All = b.Bitboards.Colors[ColorBlack] | b.Bitboards.Colors[ColorWhite]

	if color == ColorBlack {
		b.Moves.Full++
	}

	b.Attacks.Generate(&b, b.Player)
	b.Player = b.Player.Opponent()

	return b
}

func (a *Attacks) Generate(board *Board, attacker Color) {
	*a = Attacks{}

	a.king(board, attacker)
	a.sliding(board, attacker)
	a.knight(board, attacker)
	a.pawn(board, attacker)

	if !board.Attacks.Checks.Check {
		board.Attacks.Checks.Rays = BitboardAll
	}
}

func (a *Attacks) king(board *Board, attacker Color) {
	a.All.Set(Magic.KingMoves(board.Kings[attacker]))
}

func (a *Attacks) sliding(board *Board, attacker Color) {
	queens := board.Bitboards.Pieces[PieceKindQueen] & board.Bitboards.Colors[attacker]

	orthogonal := (board.Bitboards.Pieces[PieceKindRook] & board.Bitboards.Colors[attacker]) | queens
	diagonal := (board.Bitboards.Pieces[PieceKindBishop] & board.Bitboards.Colors[attacker]) | queens

	// don't include opponent's king to prevent them from accidentally
	// staying in check
	blockers := board.Bitboards.All & ^board.Kings[attacker.Opponent()].Bitboard()

	for orthogonal := orthogonal; orthogonal != 0; {
		src := Square(orthogonal.PopLSB())

		a.All.Set(Magic.OrthogonalMoves(src, blockers))
	}

	for diagonal := diagonal; diagonal != 0; {
		src := Square(diagonal.PopLSB())

		a.All.Set(Magic.DiagonalMoves(src, blockers))
	}

	// find pins / checks
	src := board.Kings[attacker.Opponent()]

	for _, dir := range Directions {
		sliders := Bitboard(0)

		if dir.IsDiagonal() {
			sliders = diagonal
		} else {
			sliders = orthogonal
		}

		if !dir.Mask(src).AnySet(sliders) {
			continue
		}

		pin := false
		ray := Bitboard(0)

		for mul := Square(1); mul <= dir.ToEdge(src); mul++ {
			dst := src + dir.Offset()*mul
			ray.Set(dst.Bitboard())

			blocker := board.Pieces[dst]

			if blocker.Is(PieceKindNone) {
				continue
			}

			if blocker.Color() == attacker {
				if (dir.IsDiagonal() && diagonal.IsSet(dst.Bitboard())) || (!dir.IsDiagonal() && orthogonal.IsSet(dst.Bitboard())) {
					if pin {
						a.Pins.Set(ray)
					} else {
						a.Checks.Double = a.Checks.Check
						a.Checks.Check = true
						a.Checks.Rays.Set(ray)
					}

					break
				}
			}

			if !pin {
				// opponent's piece blocks the ray, could be pinned
				pin = true
			} else {
				// second piece blocking the ray, not a pin
				break
			}
		}
	}
}

func (a *Attacks) knight(board *Board, attacker Color) {
	knights := board.Bitboards.Pieces[PieceKindKnight] & board.Bitboards.Colors[attacker]

	for knights != 0 {
		src := Square(knights.PopLSB())

		attacks := Magic.KnightMoves(src)

		if attacks.IsSet(board.Kings[attacker.Opponent()].Bitboard()) {
			a.Checks.Double = a.Checks.Check
			a.Checks.Check = true
			a.Checks.Rays.Set(src.Bitboard())
		}

		a.All.Set(attacks)
	}
}

func (a *Attacks) pawn(board *Board, attacker Color) {
	pawns := board.Bitboards.Pieces[PieceKindPawn] & board.Bitboards.Colors[attacker]

	for pawns != 0 {
		src := Square(pawns.PopLSB())

		attacks := Magic.PawnAttacks(attacker, src)

		if attacks.IsSet(board.Kings[attacker.Opponent()].Bitboard()) {
			a.Checks.Double = a.Checks.Check
			a.Checks.Check = true
			a.Checks.Rays.Set(src.Bitboard())
		}

		a.All.Set(attacks)
	}
}
