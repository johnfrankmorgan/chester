package main

import (
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownNumeric(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "int(-100)", UnknownNumeric(-100))
}

func TestTernary(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		condition bool
		tvalue    int
		fvalue    int
		expected  int
	}{
		{true, 1, 0, 1},
		{false, 1, 0, 0},
	} {
		test := test

		t.Run(strconv.FormatBool(test.condition), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, Ternary(test.condition, test.tvalue, test.fvalue))
		})
	}
}

func TestPanicIfError(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		PanicIfError(nil)
	})

	assert.PanicsWithValue(t, io.EOF, func() {
		PanicIfError(io.EOF)
	})
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 100, Must(100, nil))

	assert.PanicsWithValue(t, io.EOF, func() {
		Must(100, io.EOF)
	})
}

func TestAbs(t *testing.T) {
	for _, test := range []struct {
		value    int
		expected int
	}{
		{100, 100},
		{-100, 100},
	} {
		test := test

		t.Run(strconv.Itoa(test.value), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, Abs(test.value))
		})
	}
}
