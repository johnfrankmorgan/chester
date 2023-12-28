package main

import (
	"errors"
	"io"
	"os"
)

type CommandIO struct {
	_in  io.ReadCloser
	_out io.WriteCloser
	_err io.WriteCloser
}

func (cmd *CommandIO) _init() {
	if cmd._in == nil {
		cmd.SetIn(os.Stdin)
	}

	if cmd._out == nil {
		cmd.SetOut(os.Stdout)
	}

	if cmd._err == nil {
		cmd.SetErr(os.Stderr)
	}
}

func (cmd *CommandIO) Close() error {
	cmd._init()

	err := error(nil)

	err = errors.Join(err, cmd._in.Close())
	err = errors.Join(err, cmd._out.Close())
	err = errors.Join(err, cmd._err.Close())

	return err
}

func (cmd *CommandIO) In() io.Reader {
	cmd._init()

	return cmd._in
}

func (cmd *CommandIO) SetIn(r io.ReadCloser) {
	cmd._in = r
}

func (cmd *CommandIO) Out() io.Writer {
	cmd._init()

	return cmd._out
}

func (cmd *CommandIO) SetOut(w io.WriteCloser) {
	cmd._out = w
}

func (cmd *CommandIO) Err() io.Writer {
	cmd._init()

	return cmd._err
}

func (cmd *CommandIO) SetErr(w io.WriteCloser) {
	cmd._err = w
}
