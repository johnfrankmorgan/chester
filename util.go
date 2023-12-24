package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func istr[T constraints.Integer](value T) string {
	return fmt.Sprintf("%T(%d)", value, value)
}

func abs[T constraints.Signed](value T) T {
	if value < 0 {
		return -value
	}

	return value
}

func iabs[T constraints.Signed](value T) int {
	return int(abs(value))
}
