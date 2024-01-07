package main

import (
	"errors"

	"log/slog"
)

type CommandUCI struct {
	CommandIO

	uci *UCI
}

func (cmd CommandUCI) Run() error {
	cmd.uci = NewUCI(cmd.In(), cmd.Out())

	for {
		if err := cmd.uci.Run(); err != nil {
			if errors.Is(err, ErrUCIQuit) {
				return nil
			}

			slog.Error("uci error", "error", err)

			cmd.uci.Respond("error: %s", err)
		}
	}
}
