package main

import (
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		),
	)
}

func main() {
	//
}
