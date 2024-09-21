package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var cli struct {
		UCI       *UCI      `cmd:"" default:"" help:"Run UCI engine"`
		GenMagics *MagicGen `cmd:"" help:"Generate magic bitboards"`
		Log       struct {
			Level slog.Level `help:"Set the log level" enum:"DEBUG,INFO,WARN,ERROR" default:"DEBUG"`
		} `embed:"" prefix:"log-"`
	}

	kctx := kong.Parse(&cli, kong.BindTo(ctx, (*context.Context)(nil)))

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cli.Log.Level,
	})))

	return kctx.Run()
}
