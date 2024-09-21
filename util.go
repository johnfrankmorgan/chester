package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func repr(v any) string {
	return fmt.Sprintf("%T(%#v)", v, v)
}

func abs[T constraints.Signed](n T) T {
	if n < 0 {
		return -n
	}

	return n
}
