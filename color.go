package main

import "iter"

type Color uint8

const (
	Black Color = iota
	White

	ColorCount = 2
)

func ColorFromString(s string) (Color, bool) {
	switch s {
	case "b":
		return Black, true

	case "w":
		return White, true

	default:
		return 0, false
	}
}

func (c Color) String() string {
	switch c {
	case White:
		return "w"

	case Black:
		return "b"

	default:
		return repr(c)
	}
}

func (c Color) Opponent() Color {
	return c ^ 1
}

func Colors() iter.Seq[Color] {
	return func(yield func(Color) bool) {
		for color := Color(0); color < ColorCount; color++ {
			if !yield(color) {
				break
			}
		}
	}
}
