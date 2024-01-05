package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func UnknownNumeric[T constraints.Integer](value T) string {
	return fmt.Sprintf("%T(%d)", value, value)
}
