package main

import (
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
}

func (rcloser) Close() error {
	return nil
}

type wcloser struct {
	io.Writer
}

func (wcloser) Close() error {
	return nil
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
