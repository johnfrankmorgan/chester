package util

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func UnknownNumeric[T constraints.Integer](value T) string {
	return fmt.Sprintf("%T(%d)", value, value)
}

func Ternary[T any](condition bool, tvalue T, fvalue T) T {
	if condition {
		return tvalue
	}

	return fvalue
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](value T, err error) T {
	Check(err)

	return value
}
