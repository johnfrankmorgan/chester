package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestColor(t *testing.T) {
	t.Parallel()

	suite.Run(t, &ColorTest{})
}

type ColorTest struct {
	suite.Suite
}

func (t *ColorTest) TestString() {
	for _, test := range []struct {
		color    Color
		expected string
	}{
		{ColorWhite, "w"},
		{ColorBlack, "b"},
		{10, "main.Color(10)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.color.String())
		})
	}
}

func (t *ColorTest) TestValid() {
	for _, test := range []struct {
		color    Color
		expected bool
	}{
		{ColorWhite, true},
		{ColorBlack, true},
		{10, false},
	} {
		t.Run(test.color.String(), func() {
			t.Assert().Equal(test.expected, test.color.Valid())
		})
	}
}

func (t *ColorTest) TestOpponent() {
	for _, test := range []struct {
		color    Color
		expected Color
	}{
		{ColorWhite, ColorBlack},
		{ColorBlack, ColorWhite},
	} {
		t.Run(test.color.String(), func() {
			t.Assert().Equal(test.expected, test.color.Opponent())
		})
	}
}
