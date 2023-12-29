package main

import _ "embed"

//go:embed openings.json
var _openings []byte

type OpeningBook struct {
	Depth int
	ECOs  []string
	Names []string
	Moves map[string][]OpeningMove
}

type OpeningMove struct {
	ECO  int
	Name int
	Move Move
}
