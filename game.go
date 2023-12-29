package main

type Game struct {
	_boards []Board
}

func NewGame(fen string) (*Game, error) {
	board, err := NewBoard(fen)
	if err != nil {
		return nil, err
	}

	return &Game{
		_boards: []Board{board},
	}, nil
}

func (g *Game) Board() *Board {
	return &g._boards[len(g._boards)-1]
}

func (g *Game) Boards() []Board {
	return g._boards
}

func (g *Game) MakeMove(move Move) {
	g._boards = append(g._boards, g.Board().MakeMove(move))
}

func (g *Game) UnmakeMove() {
	g._boards = g._boards[:len(g._boards)-1]
}
