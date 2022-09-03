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
			Player:        0,
		},
		{
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{0, 4},
			Direction:     NE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{2, 3},
			Direction:     W,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
	}

	move := NewMove(&AdvanceRobot{
		Robot: Pair{2, 3},
	}, 0)
	err := state.Move(move)

	assert.Nil(t, err)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.True(t, state.RobotAt(Pair{4, 0}).IsBeamEnabled)

	err = state.Undo(move)
	assert.Nil(t, err)
	assert.True(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsBeamEnabled)
}

func TestAdvanceRemovesBot(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = []Robot{
		{
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{-4, 0},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{-1, 5},
			Direction:     NE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
	}

	state.PlayerTurn = 1
	move := NewMove(&AdvanceRobot{
		Robot: Pair{-1, 5},
	}, 1)

	err := state.Move(move)
	assert.Nil(t, err)

	assert.Nil(t, state.RobotAt(Pair{4, 0}))
	assert.Equal(t, 3, state.Players[1].Points)
	assert.Len(t, state.Robots, 3)

	state.Undo(move)

	assert.NotNil(t, state.RobotAt(Pair{4, 0}))
	assert.Equal(t, 0, state.Players[1].Points)
	assert.Len(t, state.Robots, 4)
}

func TestTurnLockUnlock(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)

	state.Robots = []Robot{
		{
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{-4, 0},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{0, -4},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
	}
	state.PlayerTurn = 1
	state.MovesThisTurn = 3

	move1 := NewMove(&TurnRobot{
		Robot:     Pair{-4, 0},
		Direction: Left,
	}, 1)

	err := state.Move(move1)
	assert.Nil(t, err)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.False(t, state.RobotAt(Pair{0, -4}).IsLockedDown)

	move2 := NewMove(&TurnRobot{
		Robot:     Pair{4, -4},
		Direction: Right,
	}, 1)
	err = state.Move(move2)
	assert.Nil(t, err)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.False(t, state.RobotAt(Pair{0, -4}).IsLockedDown)

	move3 := NewMove(&TurnRobot{
		Robot:     Pair{4, -4},
		Direction: Right,
	}, 1)
	err = state.Move(move3)
	assert.Nil(t, err)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.True(t, state.RobotAt(Pair{0, -4}).IsLockedDown)

	// REVERSE!

	err = state.Undo(move3)
	assert.Nil(t, err)
	assert.False(t, state.RobotAt(Pair{4, 0}).IsLockedDown)
	assert.False(t, state.RobotAt(Pair{0, -4}).IsLockedDown)

	err = state.Undo(move2)
	a