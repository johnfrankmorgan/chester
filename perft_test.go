package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerft(t *testing.T) {
	t.Parallel()

	raw := must(os.ReadFile("testdata/perft.json"))

	tests := []struct {
		FEN     string
		Results []PerftResult
	}{}

	check(json.Unmarshal(raw, &tests))

	for _, test := range tests {
		for depth, expected := range test.Results {
			if depth > 5 && testing.Short() {
				break
			}

			depth := depth
			expected := expected

			t.Run(fmt.Sprintf("%s (depth %d)", test.FEN, depth), func(t *testing.T) {
				game := must(NewGame(test.FEN))

				got := Perft(game, depth)

				assert.Equal(t, expected, got)
			})
		}
	}
}
