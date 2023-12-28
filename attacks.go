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

	*a = Attacks{}

	a._king(board, attacker)
	a._sliding(board, attacker)
	a._knight(board, attacker)
	a._pawn(board, attacker)
}

func (a *Attacks) IsAttacked(square Square) bool {
	return a.All.IsSet(square.Bitboard())
}

func (a *Attacks) IsPinned(square Square) bool {
	return a.Pins.IsSet(square.Bitboard())
}

func (a *Attacks) _king(board *Board, attacker Color) {
	slog.Debug("generating king attacks", "attacker", attacker)

	src := board.Kings[attacker]

	for _, dir := range Directions {
		if dir.ToEdge(src) > 0 {
			dst := src + dir.Offset()

			a.All.Set(dst.Bitboard())
		}
	}
}

func (a *Attacks) _sliding(board *Board, attacker Color) {
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

				ray.Set(dst.Bitboard())

				if board.Bitboards.Colors[attacker].IsSet(dst.Bitboard()) {
					break
				} else if board.Bitboards.Colors[attacker.Opponent()].IsSet(dst.Bitboard()) {
					if board.Kings[attacker.Opponent()] == dst {
						if pin {
							a.Pins.Set(ray)
						} else {
							if dir.ToEdge(dst) > 0 {
								// extend attacks through king by one square
								// so that we can't accidentally stay in check
								a.All.Set((dst + dir.Offset()).Bitboard())
							}

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
			}
		}
	}
}

func (a *Attacks) _knight(board *Board, attacker Color) {
	slog.Debug("generating knight attacks", "attacker", attacker)

	knights := board.Bitboards.Pieces[PieceKnight] & board.Bitboards.Colors[attacker]

	for knights > 0 {
		src := Square(knights.PopLSB())

		for _, dst := range Precomputed.Attacks.Knight[src] {
			a.All.Set(dst.Bitboard())

			if dst == board.Kings[attacker.Opponent()] {
				a.Checks.Double = a.Checks.Check
				a.Checks.Check = true
				a.Checks.Rays.Set(src.Bitboard() | dst.Bitboard())
			}
		}
	}
}

func (a *Attacks) _pawn(board *Board, attacker Color) {
	slog.Debug("generating pawn attacks", "attacker", attacker)

	pawns := board.Bitboards.Pieces[PiecePawn] & board.Bitboards.Colors[attacker]

	for pawns > 0 {
		src := Square(pawns.PopLSB())

		for _, dst := range Precomputed.Attacks.Pawn[attacker][src] {
			a.All.Set(dst.Bitboard())

			if dst == board.Kings[attacker.Opponent()] {
				a.Checks.Double = a.Checks.Check
				a.Checks.Check = true
				a.Checks.Rays.Set(src.Bitboard() | dst.Bitboard())
			}
		}
	}
}
