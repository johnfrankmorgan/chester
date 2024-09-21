package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerft(t *testing.T) {
	data, err := os.ReadFile("testdata/perft.json")
	require.NoError(t, err)

	tests := map[string][]uint64(nil)
	require.NoError(t, json.Unmarshal(data, &tests))

	for fen, test := range tests {
		for depth, expected := range test {
      if testing.Short() && depth > 3 {
        break
      }

			t.Run(fen, func(t *testing.T) {
				game, err := GameFromFEN(fen)
				require.NoError(t, err)

				got := Perft(game, depth, nil)

				assert.Equal(t, expected, got)
			})
		}
	}
}
