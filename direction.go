package main

import "iter"

type Direction int8

const (
	North Direction = iota
	South
	East
	West
	NorthEast
	SouthWest
	NorthWest
	SouthEast

	DirectionCount = 8

	_DiagonalStart = NorthEast
)

func (d Direction) String() string {
	return [...]string{
		North:     "north",
		South:     "south",
		East:      "east",
		West:      "west",
		NorthEast: "north-east",
		SouthWest: "south-west",
		NorthWest: "north-west",
		SouthEast: "south-east",
	}[d]
}

func (d Direction) Offset() Square {
	return [...]Square{
		North:     FileCount,
		South:     -FileCount,
		East:      1,
		West:      -1,
		NorthEast: FileCount + 1,
		SouthWest: -FileCount - 1,
		NorthWest: FileCount - 1,
		SouthEast: -FileCount + 1,
	}[d]
}

func (d Direction) Opposite() Direction {
	return [...]Direction{
		North:     South,
		South:     North,
		East:      West,
		West:      East,
		NorthEast: SouthWest,
		SouthWest: NorthEast,
		NorthWest: SouthEast,
		SouthEast: NorthWest,
	}[d]
}

func (d Direction) IsDiagonal() bool {
	return d >= _DiagonalStart
}

func Directions() iter.Seq[Direction] {
	return func(yield func(Direction) bool) {
		for dir := North; dir < DirectionCount; dir++ {
			if !yield(dir) {
				break
			}
		}
	}
}

func Orthogonals() iter.Seq[Direction] {
	return func(yield func(Direction) bool) {
		for dir := North; dir < _DiagonalStart; dir += 1 {
			if !yield(dir) {
				break
			}
		}
	}
}

func Diagonals() iter.Seq[Direction] {
	return func(yield func(Direction) bool) {
		for dir := _DiagonalStart; dir < DirectionCount; dir += 1 {
			if !yield(dir) {
				break
			}
		}
	}
}
