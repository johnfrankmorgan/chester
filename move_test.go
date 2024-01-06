package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMove(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MoveTest{})
}

type MoveTest struct {
	suite.Suite
}

func (t *MoveTest) TestString() {
	for _, test := range []struct {
		move     Move
		expected string
	}{
		{NewMove(SquareA1, SquareD1), "a1d1"},
		{NewMove(SquareE2, SquareE4), "e2e4"},
		{NewMove(SquareE4, SquareD5), "e4d5"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.move.String())
		})
	}
}

func (t *MoveTest) TestPromotion() {
	for _, test := range []struct {
		move     Move
		expected PieceKind
	}{
		{NewMove(SquareA3, SquareA4), PieceKindNone},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToQueen), PieceKindQueen},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToRook), PieceKindRook},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToBishop), PieceKindBishop},
		{NewMove(SquareA3, SquareA4, MoveFlagsPromoteToKnight), PieceKindKnight},
	} {
		t.Run(test.move.String(), func() {
			t.Assert().Equal(test.expected, test.move.Promotion())
		})
	}
}

func (t *MoveTest) TestValid() {
	for _, test := range []struct {
		move     Move
		expected bool
	}{
		{Move{}, false},
		{NewMove(SquareA1, SquareA2), true},
	} {
		t.Run(test.move.String(), func() {
			t.Assert().Equal(test.expected, test.move.Valid())
		})
	}
}

func TestMoveFlags(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MoveFlagsTest{})
}

type MoveFlagsTest struct {
	suite.Suite
}

func (t *MoveFlagsTest) TestString() {
	for _, test := range []struct {
		mf       MoveFlags
		expected string
	}{
		{MoveFlagsEnPassant, "e"},
		{MoveFlagsDoublePawnPush, "2"},
		{MoveFlagsCastleKingside, "K"},
		{MoveFlagsCastleQueenside, "Q"},
		{MoveFlagsPromoteToQueen, "q"},
		{MoveFlagsPromoteToRook, "r"},
		{MoveFlagsPromoteToBishop, "b"},
		{MoveFlagsPromoteToKnight, "n"},
		{MoveFlagsCastle, "KQ"},
		{MoveFlagsPromote, "qrbn"},
	} {
		t.Run(test.expected, func() {
			t.Assert().Equal(test.expected, test.mf.String())
		})
	}
}

func (t *MoveFlagsTest) TestIsSet() {
	mf := MoveFlagsEnPassant | MoveFlagsPromoteToQueen

	t.Assert().True(mf.IsSet(MoveFlagsEnPassant))
	t.Assert().True(mf.IsSet(MoveFlagsPromoteToQueen))
	t.Assert().True(mf.IsSet(MoveFlagsEnPassant | MoveFlagsPromoteToQueen))
	t.Assert().False(mf.IsSet(MoveFlagsDoublePawnPush))
	t.Assert().False(mf.IsSet(MoveFlagsEnPassant | MoveFlagsDoublePawnPush))
}

func (t *MoveFlagsTest) TestAnySet() {
	mf := MoveFlagsCastleKingside

	t.Assert().True(mf.AnySet(MoveFlagsCastle))
}

func TestMoveGenerator(t *testing.T) {
	tests := []struct {
		FEN   string
		Depth int
		Nodes int
		Skip  bool
	}(nil)

	PanicIfError(
		json.Unmarshal(Must(os.ReadFile("testdata/perft.json")), &tests),
	)

	maxdepth := Ternary(testing.Short(), 3, 5)

	for _, test := range tests {
		test := test

		if test.Skip || test.Depth > maxdepth {
			continue
		}

		t.Run(fmt.Sprintf("%d %s", test.Depth, test.FEN), func(t *testing.T) {
			t.Parallel()

			game := Must(NewGame(test.FEN))

			cmd := CommandPerft{}
			cmd.SetOut(wcloser{io.Discard, nil})

			assert.Equal(t, test.Nodes, cmd.perft(game, test.Depth))
		})
	}
}
