package main

import "log/slog"

const (
	EvaluationMax = 100000
	EvaluationMin = -EvaluationMax
)

type Evaluator struct {
	Name string
	Func func(*Board) int
}

var _Evaluators = []Evaluator{
	{"material", func(board *Board) int {
		const (
			endq = 45
			endr = 20
			endb = 10
			endn = 10

			endstart = endq + endr*2 + endb*2 + endn*2
		)

		endgame := struct {
			white, black float64
		}{}

		for piece, bits := range board.Bitboards.Pieces {
			piece := PieceKind(piece)

			if piece == PieceNone {
				continue
			}

			white := bits & board.Bitboards.Colors[ColorWhite]
			black := bits & board.Bitboards.Colors[ColorBlack]

			switch piece {
			case PieceQueen:
				endgame.white += float64(white.SetCount() * endq)
				endgame.black += float64(black.SetCount() * endq)

			case PieceRook:
				endgame.white += float64(white.SetCount() * endr)
				endgame.black += float64(black.SetCount() * endr)

			case PieceBishop:
				endgame.white += float64(white.SetCount() * endb)
				endgame.black += float64(black.SetCount() * endb)

			case PieceKnight:
				endgame.white += float64(white.SetCount() * endn)
				endgame.black += float64(black.SetCount() * endn)
			}
		}

		endgame.white = 1 - min(1, endgame.white/endstart)
		endgame.black = 1 - min(1, endgame.black/endstart)

		eval := 0

		for piece, bits := range board.Bitboards.Pieces {
			piece := PieceKind(piece)

			if piece == PieceNone {
				continue
			}

			white := bits & board.Bitboards.Colors[ColorWhite]
			black := bits & board.Bitboards.Colors[ColorBlack]

			eval += white.SetCount() * piece.Value()
			eval -= black.SetCount() * piece.Value()

			for white != 0 {
				src := Square(white.PopLSB())

				eval += int(Precomputed.Evaluation.Squares[piece][ColorWhite].Middlegame[src] * (1 - endgame.white))
				eval += int(Precomputed.Evaluation.Squares[piece][ColorWhite].Endgame[src] * (endgame.white))
			}

			for black != 0 {
				src := Square(black.PopLSB())

				eval -= int(Precomputed.Evaluation.Squares[piece][ColorBlack].Middlegame[src] * (1 - endgame.black))
				eval -= int(Precomputed.Evaluation.Squares[piece][ColorBlack].Endgame[src] * (endgame.black))
			}
		}

		return eval
	}},

	{"castling", func(board *Board) int {
		eval := 0

		if board.Castled[ColorWhite] {
			eval += 25
		}

		if board.Castled[ColorBlack] {
			eval -= 25
		}

		return eval
	}},
}

func Evaluate(board *Board) int {
	eval := 0

	for _, evaluator := range _Evaluators {
		slog.Debug("running evaluator", "name", evaluator.Name)

		eval += evaluator.Func(board)
	}

	if board.Player == ColorBlack {
		return -eval
	}

	return eval
}
