package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommandUCI(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CommandUCITest{})
}

type CommandUCITest struct {
	suite.Suite

	stdin  *bytes.Buffer
	stdout *bytes.Buffer
	stderr *bytes.Buffer

	cmd CommandUCI
}

func (t *CommandUCITest) SetupTest() {
	t.stdin = bytes.NewBuffer(nil)
	t.stdout = bytes.NewBuffer(nil)
	t.stderr = bytes.NewBuffer(nil)

	t.cmd.SetIn(rcloser{t.stdin, nil})
	t.cmd.SetOut(wcloser{t.stdout, nil})
	t.cmd.SetErr(wcloser{t.stderr, nil})
}

func (t *CommandUCITest) TestRun() {
	t.stdin.WriteString("debug on\ninvalidcommand\nquit\n")

	err := t.cmd.Run()

	t.Assert().NoError(err)
	t.Assert().Equal("error: uci: 1:1: unexpected token \"invalidcommand\"\n", t.stdout.String())
}
