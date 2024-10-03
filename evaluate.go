package main

type Eval int

const (
	EvalInf  Eval = 10000000
	EvalMate Eval = 1000000
)

func Evaluate(b *Board) Eval {
	eval := [ColorCount]Eval{}

	for src, piece := range b.Squares {
		if piece != EmptySquare {
			color := piece.Color()

			eval[color] += EvaluatePiece(Square(src), color, piece.Type())
		}
	}

	return eval[b.Player] - eval[b.Player.Opponent()]
}

func (e Eval) MateIn() (int, bool) {
	e = abs(e)

	if e < EvalMate {
		return 0, false
	}

	return int(EvalMate - e), true
}

var EvaluatePiece = func() func(Square, Color, PieceType) Eval {
	ptypes := [PieceTypeCount + 1]Eval{
		Pawn:   100,
		Knight: 320,
		Bishop: 330,
		Rook:   500,
		Queen:  900,
		King:   20000,
	}

	psquares := [ColorCount][PieceTypeCount + 1][SquareCount]Eval{
		Black: {
			Pawn: {
				0, 0, 0, 0, 0, 0, 0, 0,
				50, 50, 50, 50, 50, 50, 50, 50,
				10, 10, 20, 30, 30, 20, 10, 10,
				5, 5, 10, 25, 25, 10, 5, 5,
				0, 0, 0, 20, 20, 0, 0, 0,
				5, -5, -10, 0, 0, -10, -5, 5,
				5, 10, 10, -20, -20, 10, 10, 5,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			Knight: {
				-50, -40, -30, -30, -30, -30, -40, -50,
				-40, -20, 0, 0, 0, 0, -20, -40,
				-30, 0, 10, 15, 15, 10, 0, -30,
				-30, 5, 15, 20, 20, 15, 5, -30,
				-30, 0, 15, 20, 20, 15, 0, -30,
				-30, 5, 10, 15, 15, 10, 5, -30,
				-40, -20, 0, 5, 5, 0, -20, -40,
				-50, -40, -30, -30, -30, -30, -40, -50,
			},
			Bishop: {
				-20, -10, -10, -10, -10, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 10, 10, 5, 0, -10,
				-10, 5, 5, 10, 10, 5, 5, -10,
				-10, 0, 10, 10, 10, 10, 0, -10,
				-10, 10, 10, 10, 10, 10, 10, -10,
				-10, 5, 0, 0, 0, 0, 5, -10,
				-20, -10, -10, -10, -10, -10, -10, -20,
			},
			Rook: {
				0, 0, 0, 0, 0, 0, 0, 0,
				5, 10, 10, 10, 10, 10, 10, 5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				0, 0, 0, 5, 5, 0, 0, 0,
			},
			Queen: {
				-20, -10, -10, -5, -5, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 5, 5, 5, 0, -10,
				-5, 0, 5, 5, 5, 5, 0, -5,
				0, 0, 5, 5, 5, 5, 0, -5,
				-10, 5, 5, 5, 5, 5, 0, -10,
				-10, 0, 5, 0, 0, 0, 0, -10,
				-20, -10, -10, -5, -5, -10, -10, -20,
			},
			King: {
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-20, -30, -30, -40, -40, -30, -30, -20,
				-10, -20, -20, -20, -20, -20, -20, -10,
				20, 20, 0, 0, 0, 0, 20, 20,
				20, 30, 10, 0, 0, 10, 30, 20,
			},
		},
	}

	for ptype, board := range psquares[Black] {
		for src, value := range board {
			src := Square(src)

			file := src.File()
			rank := RankLast - src.Rank()

			psquares[White][ptype][NewSquare(file, rank)] = value
		}
	}

	// TODO: endgame tables

	return func(src Square, color Color, ptype PieceType) Eval {
		return ptypes[ptype] + psquares[color][ptype][src]
	}
}()
