package main

import (
	"fmt"

	"github.com/johnfrankmorgan/chester/internal/chester/magic"
	"github.com/johnfrankmorgan/chester/internal/chester/square"
)

func main() {
	moves := magic.Diagonal(square.G2, square.A2.Bitboard()|square.A3.Bitboard()|square.B7.Bitboard())
	fmt.Println(moves)
}
