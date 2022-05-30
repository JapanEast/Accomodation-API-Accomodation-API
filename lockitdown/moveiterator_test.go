package lockitdown

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorThirdPly(t *testing.T) {

	gameState := NewGame(TwoPlayerGameDef)

	gameState.Robots = []Robot{
		{
			Position:      Pair{-5, 0},
			Direction:     E,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
		{
			Position:      Pair{5, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
		{
			Position:      Pair{-5, 5},
			Direction:     NE,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        1,
		},
		{
			Position:      Pair{5, -5},
			Direction:     SW,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        1,
		},
	}
	gameState.PlayerTurn = PlayerPosition(0)

	it := NewMoveIterator(gameState)

	for i := 0; i < 6; i++ {
		assert.True(t, it.Next())
		m := it.Get()
		assert.N