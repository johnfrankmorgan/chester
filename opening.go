package main

import (
	_ "embed"
	"math/rand"
	"strings"
	"unsafe"
)

type OpeningMove struct {
	From      Square
	To        Square
	Promotion int8
	Priority  int8
	Games     int32
	Won       int32
	Lost      int32
	Plies     int32

	next    int32
	sibling int32
}

//go:embed embed/openings/Perfect2023.abk
var _OpeningBook []byte

const OpeningBookStartPos = 900

func RandomOpeningMove(moves ...string) *OpeningMove {
	available := AvaiableOpeningMoves(moves...)
	if len(available) == 0 {
		return nil
	}

	return available[rand.Intn(len(available))]
}

func AvaiableOpeningMoves(moves ...string) []*OpeningMove {
	move := FindOpeningMove(moves...)
	if move == nil {
		return nil
	}

	available := []*OpeningMove(nil)

	for s := move; s != nil; s = s.Sibling() {
		available = append(available, s)
	}

	return available
}

func OpeningMoveAt(index uintptr) *OpeningMove {
	return (*OpeningMove)(unsafe.Pointer(&_OpeningBook[index*unsafe.Sizeof(OpeningMove{})]))
}

func FindOpeningMove(moves ...string) *OpeningMove {
	opening := OpeningMoveAt(OpeningBookStartPos)

	for _, move := range moves {
		found := false

		for o := opening; o != nil; o = o.Sibling() {
			if move == o.String() {
				opening = o.Next()
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return opening
}

func (m *OpeningMove) String() string {
	s := strings.Builder{}

	s.WriteString(m.From.String())
	s.WriteString(m.To.String())

	switch abs(m.Promotion) {
	case 1:
		s.WriteByte('r')

	case 2:
		s.WriteByte('n')

	case 3:
		s.WriteByte('b')

	case 4:
		s.WriteByte('q')
	}

	return s.String()
}

func (m *OpeningMove) Next() *OpeningMove {
	if m.next <= 0 {
		return nil
	}

	return OpeningMoveAt(uintptr(m.next))
}

func (m *OpeningMove) Sibling() *OpeningMove {
	if m.sibling <= 0 {
		return nil
	}

	return OpeningMoveAt(uintptr(m.sibling))
}
