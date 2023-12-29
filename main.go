package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	cli := struct {
		Debug bool `help:"Enable debug logging."`

		Divide CommandDivide `cmd:""`
		UCI    CommandUCI    `cmd:""`
	}{}

	ctx := kong.Parse(&cli)

	if cli.Debug {
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: true,
					Level:     slog.LevelDebug,
				}),
			),
		)
	}

	ctx.FatalIfErrorf(ctx.Run())
}
