package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/repr"
)

func init() {
	level := slog.Level(0)

	if env, ok := os.LookupEnv("CHESTER_LOG_LEVEL"); ok {
		if err := level.UnmarshalText([]byte(env)); err != nil {
			fmt.Fprintf(os.Stderr, "invalid log level %q, defaulting to debug\n", env)
			level = slog.LevelDebug
		}
	}

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     level,
			}),
		),
	)
}

func main() {
	game := must(NewGame(BoardStartPositionFEN))

	result := Perft(game, 5)

	repr.Println(result)
}
