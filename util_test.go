package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUtil(t *testing.T) {
	t.Parallel()

	suite.Run(t, &UtilTest{})
}

type UtilTest struct {
	suite.Suite
}

func (t *UtilTest) TestIstr() {
	t.Assert().Equal("int(100)", istr(100))
}

func (t *UtilTest) TestAbs() {
	for _, test := range []struct {
		value    int
		expected int
	}{
		{100, 100},
		{0, 0},
		{-100, 100},
		{-80, 80},
		{123, 123},
	} {
		t.Run(strconv.Itoa(test.value), func() {
			t.Assert().Equal(test.expected, iabs(test.value))
		})
	}
}
