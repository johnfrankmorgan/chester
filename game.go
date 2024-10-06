package main

import "iter"

type Game struct {
	boards []Board
	moves  []string
}

func GameFromFEN(fen string) (*Game, error) {
	board, err := BoardFromFEN(fen)
	if err != nil {
		return nil, err
	}

	return &Game{
		boards: []Board{
			board,
		},
	}, nil
}

func (g *Game) Board() *Board {
	return &g.boards[len(g.boards)-1]
}

func (g *Game) Boards() iter.Seq[*Board] {
	return func(yield func(*Board) bool) {
		for i := 0; i < len(g.boards); i++ {
			if !yield(&g.boards[i]) {
				break
			}
		}
	}
}

func (g *Game) Moves() []string {
	return g.moves
}

func (g *Game) MakeMove(move Move) {
	g.boards = append(g.boards, g.Board().MakeMove(move))
}

func (g *Game) MakeUCIMove(uci string) bool {
	if len(uci) < 4 || len(uci) > 5 {
		return false
	}

	from, ok := SquareFromString(uci[:2])
	if !ok {
		return false
	}

	to, ok := SquareFromString(uci[2:4])
	if !ok {
		return false
	}

	move := NewMove(from, to)

	if len(uci) == 5 {
		switch uci[4] {
		case 'q':
			move.Flags |= MoveFlagPromoteToQueen

		case 'r':
			move.Flags |= MoveFlagPromoteToRook

		case 'b':
			move.Flags |= MoveFlagPromoteToBishop

		case 'n':
			move.Flags |= MoveFlagPromoteToKnight

		default:
			return false
		}
	}

	b := g.Board()
	p := b.Squares[move.From]

	if b.EnPassant != 0 && p.Type() == Pawn && move.To == b.EnPassant {
		move.Flags |= MoveFlagCapture
		move.Flags |= MoveFlagCaptureEnPassant
	} else if p.Type() == Pawn && abs(move.From.Rank()-move.To.Rank()) == 2 {
		move.Flags |= MoveFlagDoublePawnPush
	}

	if b.Squares[move.To] != EmptySquare {
		move.Flags |= MoveFlagCapture
	}

	if p.Type() == King && abs(move.From.File()-move.To.File()) > 1 {
		if move.To.File() == FileC {
			move.Flags |= MoveFlagCastleQueenside
		} else {
			move.Flags |= MoveFlagCastleKingside
		}
	}

	g.MakeMove(move)

	g.moves = append(g.moves, uci)

	return true
}

func (g *Game) UnmakeMove() {
	g.boards = g.boards[:len(g.boards)-1]
}
