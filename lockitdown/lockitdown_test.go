package lockitdown

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TwoPlayerGameDef = GameDef{
	Players: 2,
	Board: Board{
		HexaBoard: BoardType{4},
	},
	RobotsPerPlayer: 6,
	WinCondition:    "Elimination",
	MovesPerTurn:    3,
}

func TestNewGame(t *testing.T) {
	game := NewGame(GameDef{
		Players: 2,
	})
	if game.PlayerTurn != 0 {
		t.Errorf("Wrong player turn")
	}
	if len(game.Players) != 2 {
		t.Errorf(("Wrong number of players"))
	}
	if len(game.Robots) != 0 {
		t.Error("Improperly initialized robots")
	}
}
func TestMoves(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)

	tests := []struct {
		move   Mover
		player PlayerPosition
		err    error
	}{
		{&PlaceRobot{
			Robot:     Pair{0, 5},
			Direction: Pair{0, -1},
		}, 0, nil},
		{&PlaceRobot{
			Robot:     Pair{5, 0},
			Direction: Pair{-1, 0},
		}, 1, nil},
		{&AdvanceRobot{
			Robot: Pair{0, 5},
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Left,
		}, 0, nil},
		{&PlaceRobot{
			Robot:     Pair{-5, 4},
			Direction: Pair{1, 0},
		}, 1, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&AdvanceRobot{
			Robot: Pair{-5, 4},
		}, 1, nil},
		{&AdvanceRobot{
			Robot: Pair{5, 0},
		}, 1, nil},
		{&TurnRobot{
			Robot:    