package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Log struct {
			Level  slog.Level `help:"Logging level." default:"info"`
			Source bool       `help:"Add source to logs." default:"false"`
		} `embed:"" prefix:"log-"`

		GenerateMagics CommandGenerateMagics `cmd:"" help:"Generate magic bitboards."`
	}

	ctx := kong.Parse(&cli)

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: cli.Log.Source,
				Level:     cli.Log.Level,
			}),
		),
	)

	ctx.FatalIfErrorf(ctx.Run())
}
