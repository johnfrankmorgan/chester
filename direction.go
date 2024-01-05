package main

type Direction uint8

const (
	DirectionNorth = iota
	DirectionSouth
	DirectionEast
	DirectionWest
	DirectionNorthEast
	DirectionSouthWest
	DirectionNorthWest
	DirectionSouthEast

	DirectionCount = 8

	_DirectionDiagonalStart = DirectionNorthEast
)

var (
	Directions = [DirectionCount]Direction{
		DirectionNorth,
		DirectionSouth,
		DirectionEast,
		DirectionWest,
		DirectionNorthEast,
		DirectionSouthWest,
		DirectionNorthWest,
		DirectionSouthEast,
	}

	DirectionsOrthogonal = [...]Direction{
		DirectionNorth,
		DirectionSouth,
		DirectionEast,
		DirectionWest,
	}

	DirectionsDiagonal = [...]Direction{
		DirectionNorthEast,
		DirectionSouthWest,
		DirectionNorthWest,
		DirectionSouthEast,
	}

	_DirectionStrings = [DirectionCount]string{
		DirectionNorth:     "north",
		DirectionSouth:     "south",
		DirectionEast:      "east",
		DirectionWest:      "west",
		DirectionNorthEast: "north east",
		DirectionSouthWest: "south west",
		DirectionNorthWest: "north west",
		DirectionSouthEast: "south east",
	}

	_DirectionOffsets = [DirectionCount]Square{
		DirectionNorth:     FileCount,
		DirectionSouth:     -FileCount,
		DirectionEast:      1,
		DirectionWest:      -1,
		DirectionNorthEast: FileCount + 1,
		DirectionSouthWest: -FileCount - 1,
		DirectionNorthWest: FileCount - 1,
		DirectionSouthEast: -FileCount + 1,
	}

	_DirectionToEdge [SquareCount][DirectionCount]Square
)

func init() {
	for src := SquareFirst; src <= SquareLast; src++ {
		file := Square(src.File())
		rank := Square(src.Rank())

		north := RankCount - rank - 1
		south := rank
		east := FileCount - file - 1
		west := file

		_DirectionToEdge[src] = [DirectionCount]Square{
			DirectionNorth:     north,
			DirectionSouth:     south,
			DirectionEast:      east,
			DirectionWest:      west,
			DirectionNorthEast: min(north, east),
			DirectionSouthEast: min(south, east),
			DirectionNorthWest: min(north, west),
			DirectionSouthWest: min(south, west),
		}
	}
}

func (d Direction) String() string {
	if d.Valid() {
		return _DirectionStrings[d]
	}

	return UnknownNumeric(d)
}

func (d Direction) Valid() bool {
	return d < DirectionCount
}

func (d Direction) Offset() Square {
	return _DirectionOffsets[d]
}

func (d Direction) Opposite() Direction {
	if d.Offset() < 0 {
		return d - 1
	}

	return d + 1
}

func (dir Direction) ToEdge(src Square) Square {
	return _DirectionToEdge[src][dir]
}

func (d Direction) IsDiagonal() bool {
	return d >= _DirectionDiagonalStart
}
