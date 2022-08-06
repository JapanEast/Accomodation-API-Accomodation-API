package lockitdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurn(t *testing.T) {

	testcases := []struct {
		direction TurnDirection
		expected  Pair
	}{
		{Left,
			Pair{0, -1}},
		{Left,
			Pair{-1, 0}},
		{Left,
			Pair{-1, 1}},
		{Left,
			Pair{0, 1}},

		// Roll it back!

		{Right,
			Pair{-1, 1}},
		{Right,
			Pair{-1, 0}},
		{Right,
			Pair{0, -1}},
	}

	direction := Pair{1, -1}
	for _, tc := range testcases {
		direction.Rotate(tc.direction)
		assert.EqualValues(t, tc.expected, direction, "Wrong turn!")
	}
}

func TestAdvance(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = []Robot{
		{
			Position:      Pair{2, 3},
			Direction:     NW,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
	}
	move := NewMove(&AdvanceRobot{
		Robot: Pair{2, 3},
	}, 0)
	err := state.Move(move)
	assert.Nil(t, err)
	assert.Equal(t, 2, state.MovesThisTurn)

	err = state.Undo(move)

	assert.Nil(t, err)
	assert.Equal(t, 3, state.MovesThisTurn)
	bot := state.RobotAt(Pair{2, 3})
	assert.True(t, bot != nil)
	assert.Equal(t, Pair{2, 3}, bot.Position)
}

func TestAdvanceBlocksLockdown(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = []Robot{
		{
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0