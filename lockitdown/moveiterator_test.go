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
			Dire