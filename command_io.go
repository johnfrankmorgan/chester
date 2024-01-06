package main

import (
	"errors"
	"io"
	"os"
)

type CommandIO struct {
	stdin  io.ReadCloser
	stdout io.WriteCloser
	stderr io.WriteCloser
}

func (cmd *CommandIO) init() {
	if cmd.stdin == nil {
		cmd.SetIn(os.Stdin)
	}

	if cmd.stdout == nil {
		cmd.SetOut(os.Stdout)
	}

	if cmd.stderr == nil {
		cmd.SetErr(os.Stderr)
	}
}

func (cmd *CommandIO) Close() error {
	cmd.init()

	err := error(nil)

	err = errors.Join(err, cmd.stdin.Close())
	err = errors.Join(err, cmd.stdout.Close())
	err = errors.Join(err, cmd.stderr.Close())

	return err
}

func (cmd *CommandIO) In() io.Reader {
	cmd.init()

	return cmd.stdin
}

func (cmd *CommandIO) SetIn(r io.ReadCloser) {
	cmd.stdin = r
}

func (cmd *CommandIO) Out() io.Writer {
	cmd.init()

	return cmd.stdout
}

func (cmd *CommandIO) SetOut(w io.WriteCloser) {
	cmd.stdout = w
}

func (cmd *CommandIO) Err() io.Writer {
	cmd.init()

	return cmd.stderr
}

func (cmd *CommandIO) SetErr(w io.WriteCloser) {
	cmd.stderr = w
}
