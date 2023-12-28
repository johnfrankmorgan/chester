package main

type PerftResult struct {
	Moves             int
	Captures          int
	EnPassantCaptures int
	Checks            int
	DoubleChecks      int
	Castles           int
	Promotions        int
}

func Perft(game *Game, depth int) PerftResult {
	if depth == 0 {
		return PerftResult{
			Moves:             1,
			Captures:          ternary(game.Board().Moves.Last.Flags.IsSet(MoveFlagsCapture), 1, 0),
			EnPassantCaptures: ternary(game.Board().Moves.Last.Flags.IsSet(MoveFlagsCaptureEnPassant), 1, 0),
			Checks:            ternary(game.Board().Attacks.Checks.Check, 1, 0),
			DoubleChecks:      ternary(game.Board().Attacks.Checks.Double, 1, 0),
			Castles:           ternary(game.Board().Moves.Last.Flags.AnySet(MoveFlagsCastle), 1, 0),
			Promotions:        ternary(game.Board().Moves.Last.Flags.AnySet(MoveFlagsPromote), 1, 0),
		}
	}

	result := PerftResult{}

	for _, move := range game.Board().GenerateMoves(MoveGeneratorOptions{}) {
		game.MakeMove(move)
		result.Add(Perft(game, depth-1))
		game.UnmakeMove()
	}

	return result
}

func (pr *PerftResult) Add(other PerftResult) {
	pr.Moves += other.Moves
	pr.Captures += other.Captures
	pr.EnPassantCaptures += other.EnPassantCaptures
	pr.Checks += other.Checks
	pr.DoubleChecks += other.DoubleChecks
	pr.Castles += other.Castles
	pr.Promotions += other.Promotions
}
