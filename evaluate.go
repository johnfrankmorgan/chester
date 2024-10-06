package main

type Eval int

const (
	EvalInf  Eval = 10000000
	EvalMate Eval = 1000000
)

func Evaluate(b *Board) Eval {
	const (
		mid = 0
		end = 1
	)

	evals := [ColorCount][2]Eval{}

	for src, piece := range b.Squares {
		if piece != EmptySquare {
			color := piece.Color()
			ptype := piece.Type()
			value := EvaluatePiece(ptype)

			evals[color][mid] += value + EvaluatePiecePositionMiddlegame(Square(src), color, ptype)
			evals[color][end] += value + EvaluatePiecePositionEndgame(Square(src), color, ptype)
		}
	}

	phase := Phase(b)

	eval := [ColorCount]Eval{
		Black: (evals[Black][mid]*(256-phase) + evals[Black][end]*phase) / 256,
		White: (evals[White][mid]*(256-phase) + evals[White][end]*phase) / 256,
	}

	return eval[b.Player] - eval[b.Player.Opponent()]
}

func Phase(b *Board) Eval {
	const (
		pawn   = 0
		knight = 1
		bishop = 1
		rook   = 2
		queen  = 4

		total = pawn*16 + knight*4 + bishop*4 + rook*4 + queen*2
	)

	phase := total

	phase -= pawn * b.Bits.Pieces[Pawn].OnesCount()
	phase -= knight * b.Bits.Pieces[Knight].OnesCount()
	phase -= bishop * b.Bits.Pieces[Bishop].OnesCount()
	phase -= rook * b.Bits.Pieces[Rook].OnesCount()
	phase -= queen * b.Bits.Pieces[Queen].OnesCount()

	return Eval((phase*256)+(total/2)) / total
}

func (e Eval) MateIn() (int, bool) {
	e = abs(e)

	if e < EvalMate-256 {
		return 0, false
	}

	return int(EvalMate - e), true
}

func EvaluatePiece(ptype PieceType) Eval {
	return [PieceTypeCount + 1]Eval{
		Pawn:   100,
		Knight: 320,
		Bishop: 330,
		Rook:   500,
		Queen:  900,
		King:   20000,
	}[ptype]
}

var EvaluatePiecePositionMiddlegame = func() func(Square, Color, PieceType) Eval {
	lookup := [ColorCount][PieceTypeCount + 1][SquareCount]Eval{
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

	for ptype, board := range lookup[Black] {
		for src, value := range board {
			src := Square(src)

			file := src.File()
			rank := RankLast - src.Rank()

			lookup[White][ptype][NewSquare(file, rank)] = value
		}
	}

	return func(src Square, color Color, ptype PieceType) Eval {
		return lookup[color][ptype][src]
	}
}()

var EvaluatePiecePositionEndgame = func() func(Square, Color, PieceType) Eval {
	lookup := [ColorCount][PieceTypeCount + 1][SquareCount]Eval{
		Black: {
			Pawn: {
				0, 0, 0, 0, 0, 0, 0, 0,
				80, 80, 80, 80, 80, 80, 80, 80,
				50, 50, 50, 50, 50, 50, 50, 50,
				30, 30, 30, 30, 30, 30, 30, 30,
				20, 20, 20, 20, 20, 20, 20, 20,
				10, 10, 10, 10, 10, 10, 10, 10,
				10, 10, 10, 10, 10, 10, 10, 10,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			King: {
				-50, -40, -30, -20, -20, -30, -40, -50,
				-30, -20, -10, 0, 0, -10, -20, -30,
				-30, -10, 20, 30, 30, 20, -10, -30,
				-30, -10, 30, 40, 40, 30, -10, -30,
				-30, -10, 30, 40, 40, 30, -10, -30,
				-30, -10, 20, 30, 30, 20, -10, -30,
				-30, -30, 0, 0, 0, 0, -30, -30,
				-50, -30, -30, -30, -30, -30, -30, -50,
			},
		},
	}

	for ptype, board := range lookup[Black] {
		for src, value := range board {
			src := Square(src)

			file := src.File()
			rank := RankLast - src.Rank()

			lookup[White][ptype][NewSquare(file, rank)] = value
		}
	}

	return func(src Square, color Color, ptype PieceType) Eval {
		if ptype == King || ptype == Pawn {
			return lookup[color][ptype][src]
		}

		return EvaluatePiecePositionMiddlegame(src, color, ptype)
	}
}()
