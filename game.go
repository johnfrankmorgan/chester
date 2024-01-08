package main

type Game struct {
	boards []Board
}

func NewGame(fen string) (*Game, error) {
	board, err := NewBoard(fen)
	if err != nil {
		return nil, err
	}

	return &Game{
		boards: []Board{board},
	}, nil
}

func (g *Game) Board() *Board {
	return &g.boards[len(g.boards)-1]
}

func (g *Game) MakeMove(move Move) {
	g.boards = append(g.boards, g.Board().MakeMove(move))
}

func (g *Game) UnmakeMove() {
	g.boards = g.boards[:len(g.boards)-1]
}
