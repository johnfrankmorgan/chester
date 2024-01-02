package piece

import (
	"github.com/johnfrankmorgan/chester/internal/chester/util"
)

type Color uint8

const (
	ColorBlack Color = iota
	ColorWhite

	ColorCount = 2
)

func (c Color) String() string {
	switch c {
	case ColorBlack:
		return "black"

	case ColorWhite:
		return "white"
	}

	return util.UnknownNumeric(c)
}

func (c Color) Valid() bool {
	return c == ColorBlack || c == ColorWhite
}

func (c Color) Opponent() Color {
	return c &^ ColorWhite
}
