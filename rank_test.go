package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRank(t *testing.T) {
	t.Parallel()

	suite.Run(t, &RankTest{})
}

type RankTest struct {
	suite.Suite
}

func (t *RankTest) TestString() {
	for _, test := range []struct {
		file     Rank
		expected string
	}{
		{Rank1, "1"},
		{Rank2, "2"},
		{Rank3, "3"},
		{Rank4, "4"},
		{Rank5, "5"},
		{Rank6, "6"},
		{Rank7, "7"},
		{Rank8, "8"},
		{10, "main.Rank(10)"},
		{-100, "main.Rank(-100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.file.String())
		})
	}
}

func (t *RankTest) TestValid() {
	for _, test := range []struct {
		file     Rank
		expected bool
	}{
		{Rank1, true},
		{Rank2, true},
		{Rank3, true},
		{Rank4, true},
		{Rank5, true},
		{Rank6, true},
		{Rank7, true},
		{Rank8, true},
		{10, false},
		{-100, false},
	} {
		t.Run(test.file.String(), func() {
			t.Assert().Equal(test.expected, test.file.Valid())
		})
	}
}
