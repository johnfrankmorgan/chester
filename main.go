package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Log struct {
			Level  slog.Level `help:"Logging level." default:"info"`
			Source bool       `help:"Add source to logs." default:"false"`
			File   string     `help:"Write logs to the specified file." type:"path"`
		} `embed:"" prefix:"log-"`

		GenerateMagics CommandGenerateMagics `cmd:"" help:"Generate magic bitboards."`
		Perft          CommandPerft          `cmd:"" help:"Run Perft."`
		UCI            CommandUCI            `cmd:"" help:"Run UCI."`
	}

	ctx := kong.Parse(&cli)

	{
		log := os.Stderr

		if cli.Log.File != "" {
			err := error(nil)

			f, err := os.OpenFile(cli.Log.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			ctx.FatalIfErrorf(err)

			defer f.Close()

			log = f
		}

		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(log, &slog.HandlerOptions{
					AddSource: cli.Log.Source,
					Level:     cli.Log.Level,
				}),
			),
		)
	}

	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	}
}
