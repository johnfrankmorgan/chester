package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownNumeric(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "int(-100)", UnknownNumeric(-100))
}
