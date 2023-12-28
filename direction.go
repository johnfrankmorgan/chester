package main

import "log/slog"

type Direction uint8

const (
	DirectionNorth Direction = iota
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
	_DirectionNames = [DirectionCount]string{
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
)

func init() {
	slog.Debug("initializing squares to edge")

	for src := 0; src < SquareCount; src++ {
		src := Square(src)

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

func (dir Direction) String() string {
	return _DirectionNames[dir]
}

func (dir Direction) Offset() Square {
	return _DirectionOffsets[dir]
}

func (dir Direction) Opposite() Direction {
	if dir.Offset() < 0 {
		return dir - 1
	}

	return dir + 1
}

func (dir Direction) ToEdge(src Square) Square {
	return _DirectionToEdge[src][dir]
}

func (dir Direction) Mask(src Square) Bitboard {
	return Precomputed.Masks.Direction[src][dir]
}

func (dir Direction) IsOrthogonal() bool {
	return dir < _DirectionDiagonalStart
}

func (dir Direction) IsDiagonal() bool {
	return dir >= _DirectionDiagonalStart
}
