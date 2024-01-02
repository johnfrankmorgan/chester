package piece

type Piece uint8

const (
	Empty Piece = 0

	_White = Piece(ColorWhite << 3)
	_Black = Piece(ColorBlack << 3)

	WhitePawn   Piece = _White | Piece(Pawn)
	WhiteKnight Piece = _White | Piece(Knight)
	WhiteBishop Piece = _White | Piece(Bishop)
	WhiteRook   Piece = _White | Piece(Rook)
	WhiteQueen  Piece = _White | Piece(Queen)
	WhiteKing   Piece = _White | Piece(King)

	BlackPawn   Piece = _Black | Piece(Pawn)
	BlackKnight Piece = _Black | Piece(Knight)
	BlackBishop Piece = _Black | Piece(Bishop)
	BlackRook   Piece = _Black | Piece(Rook)
	BlackQueen  Piece = _Black | Piece(Queen)
	BlackKing   Piece = _Black | Piece(King)
)

func New(color Color, kind Kind) Piece {
	return Piece(color<<3) | Piece(kind)
}

func (p Piece) Color() Color {
	return Color(p >> 3)
}

func (p Piece) Kind() Kind {
	return Kind(p &^ _White)
}
