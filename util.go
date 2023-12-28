package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](value T, err error) T {
	check(err)

	return value
}

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

func sign[T constraints.Signed](value T) T {
	if value == 0 {
		return 0
	}

	if value < 0 {
		return -1
	}

	return 1
}

func isign[T constraints.Signed](value T) int {
	return int(sign(value))
}
