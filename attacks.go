package main

import "log/slog"

type Attacks struct {
	All    Bitboard
	Checks struct {
		Check  bool
		Double bool
		Rays   Bitboard
	}
	Pins Bitboard
}

func (a *Attacks) Generate(board *Board, attacker Color) {
	slog.Debug("generating attacks", "attacker", attacker)

	a.King(board, attacker)
	a.Sliding(board, attacker)
	a.Knight(board, attacker)
	a.Pawn(board, attacker)
}

func (a *Attacks) King(board *Board, attacker Color) {
	slog.Debug("generating king attacks", "attacker", attacker)

	src := board.Kings[attacker]

	for _, dir := range Directions {
		if dir.ToEdge(src) > 0 {
			dst := src + dir.Offset()

			a.All.Set(dst.Bitboard())
		}
	}
}

func (a *Attacks) Sliding(board *Board, attacker Color) {
	slog.Debug("generating sliding attacks", "attacker", attacker)

	queens := board.Bitboards.Pieces[PieceQueen]
	rooks := board.Bitboards.Pieces[PieceRook]
	bishops := board.Bitboards.Pieces[PieceBishop]

	orthogonal := (queens | rooks) & board.Bitboards.Colors[attacker]
	diagonal := (queens | bishops) & board.Bitboards.Colors[attacker]

	for _, dir := range Directions {
		sliders := orthogonal

		if dir.IsDiagonal() {
			sliders = diagonal
		}

		for sliders > 0 {
			src := Square(sliders.PopLSB())
			ray := src.Bitboard()
			pin := false

			for mul := Square(1); mul <= dir.ToEdge(src); mul++ {
				dst := src + dir.Offset()*mul

				if !pin {
					a.All.Set(dst.Bitboard())
				}

				if board.Bitboards.Colors[attacker].IsSet(dst.Bitboard()) {
					break
				} else if board.Bitboards.Colors[attacker.Opponent()].IsSet(dst.Bitboard()) {
					if board.Kings[attacker.Opponent()] == dst {
						if pin {
							a.Pins.Set(ray)
						} else {
							a.Checks.Double = a.Checks.Check
							a.Checks.Check = true
							a.Checks.Rays.Set(ray)
						}

						break
					}

					if pin {
						break
					}

					pin = true
				}

				ray.Set(dst.Bitboard())
			}
		}
	}
}

func (a *Attacks) Knight(board *Board, attacker Color) {
	slog.Debug("generating knight attacks", "attacker", attacker)

	knights := board.Bitboards.Pieces[PieceKnight] & board.Bitboards.Colors[attacker]

	for knights > 0 {
		src := Square(knights.PopLSB())

		for _, dst := range Precomputed.Attacks.Knight[src] {
			a.All.Set(dst.Bitboard())

			if dst == board.Kings[attacker.Opponent()] {
				a.Checks.Double = a.Checks.Check
				a.Checks.Check = true
				a.Checks.Rays.Set(src.Bitboard())
			}
		}
	}
}

func (a *Attacks) Pawn(board *Board, attacker Color) {
	slog.Debug("generating pawn attacks", "attacker", attacker)

	pawns := board.Bitboards.Pieces[PiecePawn] & board.Bitboards.Colors[attacker]

	for pawns > 0 {
		src := Square(pawns.PopLSB())

		for _, dst := range Precomputed.Attacks.Pawn[attacker][src] {
			a.All.Set(dst.Bitboard())

			if dst == board.Kings[attacker.Opponent()] {
				a.Checks.Double = a.Checks.Check
				a.Checks.Check = true
				a.Checks.Rays.Set(src.Bitboard())
			}
		}
	}
}
