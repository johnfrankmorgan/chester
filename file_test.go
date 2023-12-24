package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestFile(t *testing.T) {
	t.Parallel()

	suite.Run(t, &FileTest{})
}

type FileTest struct {
	suite.Suite
}

func (t *FileTest) TestString() {
	for _, test := range []struct {
		file     File
		expected string
	}{
		{FileA, "a"},
		{FileB, "b"},
		{FileC, "c"},
		{FileD, "d"},
		{FileE, "e"},
		{FileF, "f"},
		{FileG, "g"},
		{FileH, "h"},
		{10, "main.File(10)"},
		{-100, "main.File(-100)"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.file.String())
		})
	}
}

func (t *FileTest) TestValid() {
	for _, test := range []struct {
		file     File
		expected bool
	}{
		{FileA, true},
		{FileB, true},
		{FileC, true},
		{FileD, true},
		{FileE, true},
		{FileF, true},
		{FileG, true},
		{FileH, true},
		{10, false},
		{-100, false},
	} {
		t.Run(test.file.String(), func() {
			t.Assert().Equal(test.expected, test.file.Valid())
		})
	}
}
