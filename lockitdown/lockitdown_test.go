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
			Robot:     Pair{4, 0},
			Direction: Left,
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
		{&PlaceRobot{
			Robot:     Pair{0, -5},
			Direction: Pair{0, 1},
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
		{
			move: &AdvanceRobot{
				Robot: Pair{0, -5},
			},
			player: 1,
			err:    nil,
		},
	}

	for _, tt := range tests {
		m := NewMove(tt.move, tt.player)
		err := game.Move(m)
		if tt.err == nil && err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Errorf("expected error %s, none recieved", tt.err.Error())
		}
		if tt.err != nil && err != nil && tt.err.Error() != err.Error() {
			t.Errorf("expected error %s, recieved %s", tt.err.Error(), err.Error())
		}
	}

	assert.Equal(t, 1, game.Players[0].PlacedRobots, "wrong number of player 1 robots")
	assert.Equal(t, 3, game.Players[1].PlacedRobots, "wrong number of player 2 robots")
	assert.Equal(t, 3, game.Players[1].Points, "wrong number of player 2 points")
}

func TestGameOver(t *testing.T) {
	gameState := GameState{
		GameDef: TwoPlayerGameDef,
		Players: []*Player{
			{
				Points:       0,
				PlacedRobots: 6,
			},
			{
				Points:       9,
				PlacedRobots: 3,
			},
		},
		Robots: []Robot{
			{
				Position:      Pair{-4, 4},
				Direction:     NE,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        1,
			},
			{
				Position:      Pair{4, 0},
				Direction:     SW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        1,
			},
			{
				Position:      Pair{0, -4},
				Direction:     SE,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        1,
			},
			{
				Position:      Pair{4, -4},
				Direction:     SW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
			{
				Position:      Pair{5, -5},
				Direction:     SW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
			{
				Position:      Pair{0, 5},
				Direction:     NW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
		},
		PlayerTurn:       0,
		MovesThisTurn:    3,
		RequiresTieBreak: false,
		Winner:           -1,
	}

	testcases := []struct {
		move   Mover
		player PlayerPosition
		result error
	}{
		{
			&AdvanceRobot{
				Robot: Pair{4, -4},
			},
			0,
			nil,
		},
		{
			&AdvanceRobot{
				Robot: Pair{5, -5},
			}, 0,
			nil,
		},
		{
			&AdvanceRobot{
				Robot: Pair{0, 5},
			},
			0,
			nil,
		},
	}

	for _, tc := range testcases {
		m := NewMove(tc.move, tc.player)
		err := gameState.Move(m)
		assert.Nilf(t, err, "unexpected err for move %v", tc.move)
		assert.Equal(t, -1, gameState.Winner)
		if err != nil {
			fmt.Println(gameState.ToJson())
		}
	}

	err := gameState.Move(NewMove(&TurnRobot{
		Robot:     Pair{-4, 4},
		Direction: Right,
	}, 1))

	fmt.Println(gameState.ToJson())

	assert.EqualError(t, err, "winner is 2")
	assert.Equal(t, 1, gameState.Winner)
}

func TestPairCopy(t *testing.T) {
	p1 := Pair{
		Q: 32,
		R: 12,
	}

	p2 := p1.Copy()

	p1.Q = 2
	p1.R = 4

	p2.Q = 5
	p2.R = 7

	assert.Equal(t, p1, Pair{2, 4})
	assert.Equal(t, p2, Pair{5, 7})
}

func TestInBounds(t *tes