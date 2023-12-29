package main

import (
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/alecthomas/kong"
)

func main() {
	cli := struct {
		Debug   bool   `help:"Enable debug logging."`
		Profile string `help:"Write profile information here."`

		Divide           CommandDivide           `cmd:""`
		UCI              CommandUCI              `cmd:"" default:"true"`
		GenerateOpenings CommandGenerateOpenings `cmd:""`
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

	profile := (*os.File)(nil)

	if cli.Profile != "" {
		err := error(nil)

		profile, err = os.Create(cli.Profile)

		ctx.FatalIfErrorf(err)
		ctx.FatalIfErrorf(pprof.StartCPUProfile(profile))
	}

	err := ctx.Run()

	if profile != nil {
		pprof.StopCPUProfile()
		ctx.FatalIfErrorf(profile.Close())
	}

	ctx.FatalIfErrorf(err)
}
