package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommandIO(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CommandIOTest{})
}

type CommandIOTest struct {
	suite.Suite
}

type rcloser struct {
	io.Reader

	err error
}

func (r rcloser) Close() error {
	return r.err
}

type wcloser struct {
	io.Writer

	err error
}

func (w wcloser) Close() error {
	return w.err
}

func (t *CommandIOTest) TestClose() {
	err1 := errors.New("1")
	err2 := errors.New("2")
	err3 := errors.New("3")

	cmd := CommandIO{}
	cmd.SetIn(rcloser{nil, err1})
	cmd.SetOut(wcloser{nil, err2})
	cmd.SetErr(wcloser{nil, err3})

	err := cmd.Close()

	t.Assert().ErrorIs(err, err1)
	t.Assert().ErrorIs(err, err2)
	t.Assert().ErrorIs(err, err3)
}

func (t *CommandIOTest) TestIn() {
	cmd := CommandIO{}

	t.Assert().Equal(os.Stdin, cmd.In())
}

func (t *CommandIOTest) TestOut() {
	cmd := CommandIO{}

	t.Assert().Equal(os.Stdout, cmd.Out())
}

func (t *CommandIOTest) TestErr() {
	cmd := CommandIO{}

	t.Assert().Equal(os.Stderr, cmd.Err())
}
