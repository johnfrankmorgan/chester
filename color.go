package main

type Color uint8

const (
	ColorBlack Color = iota
	ColorWhite

	ColorCount = 2
)

func (c Color) String() string {
	switch c {
	case ColorBlack:
		return "b"

	case ColorWhite:
		return "w"

	default:
		return istr(c)
	}
}

func (c Color) Valid() bool {
	return c == ColorBlack || c == ColorWhite
}

func (c Color) Opponent() Color {
	return c ^ ColorWhite
}
