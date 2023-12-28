package main

import (
	"io"
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

func (t *UtilTest) TestCheck() {
	t.Assert().NotPanics(func() {
		check(nil)
	})

	t.Assert().PanicsWithValue(io.EOF, func() {
		check(io.EOF)
	})
}

func (t *UtilTest) TestMust() {
	t.Assert().Equal(1, must(1, nil))

	t.Assert().PanicsWithValue(io.EOF, func() {
		must(1, io.EOF)
	})
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

func (t *UtilTest) TestSign() {
	for _, test := range []struct {
		value    int
		expected int
	}{
		{100, 1},
		{0, 0},
		{-100, -1},
		{-80, -1},
		{123, 1},
	} {
		t.Run(strconv.Itoa(test.value), func() {
			t.Assert().Equal(test.expected, isign(test.value))
		})
	}
}

func (t *UtilTest) TestTernary() {
	for _, test := range []struct {
		condition bool
		tvalue    int
		fvalue    int
		expected  int
	}{
		{true, 1, 100, 1},
		{false, 1, 100, 100},
	} {
		t.Run(strconv.FormatBool(test.condition), func() {
			t.Assert().Equal(test.expected, ternary(test.condition, test.tvalue, test.fvalue))
		})
	}
}
