package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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
		{MoveFlagsCapture, "c"},
		{MoveFlagsCaptureEnPassant, "e"},
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
	mf := MoveFlagsCapture | MoveFlagsPromoteToQueen

	t.Assert().True(mf.IsSet(MoveFlagsCapture))
	t.Assert().True(mf.IsSet(MoveFlagsPromoteToQueen))
	t.Assert().True(mf.IsSet(MoveFlagsCapture | MoveFlagsPromoteToQueen))
	t.Assert().False(mf.IsSet(MoveFlagsDoublePawnPush))
	t.Assert().False(mf.IsSet(MoveFlagsCapture | MoveFlagsDoublePawnPush))
}

func (t *MoveFlagsTest) TestAnySet() {
	mf := MoveFlagsCastleKingside

	t.Assert().True(mf.AnySet(MoveFlagsCastle))
}
