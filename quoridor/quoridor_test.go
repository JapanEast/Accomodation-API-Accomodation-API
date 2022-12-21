package quoridor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestIds = []uuid.UUID{
	uuid.MustParse("98ae983e-3f04-42ab-928a-c399d6d18375"),
	uuid.MustParse("5341acab-6e28-4d28-8530-8716e0c3dd2e"),
	uuid.MustParse("790bcc3f-6e72-4a0e-a6ea-bc806aa8aa03"),
	uuid.MustParse("6c8420b5-e7f5-4328-ae29-4dbdf7537612"),
	uuid.MustParse("f3282245-f546-4c71-92ca-5bada1f9c037"),
	uuid.MustParse("9cae2aa0-d21a-48ab-a877-4b78942259e4"),
	uuid.MustParse("0ad943b2-6ea9-45ad-9098-f67714652fcd"),
	uuid.MustParse("93ded37f-57d3-4b43-8933-1164e086a881"),
	uuid.MustParse("5b399bd3-aa3e-4754-bb51-175b30b77400"),
	uuid.MustParse("f7ea9019-033b-41e7-a671-26231952cd8c"),
}

func Test_NewGame(t *testing.T) {
	var testCases = []struct {
		id             uuid.UUID
		gameName       string
		expectGame     bool
		expectedErrMsg string
	}{
		{
			TestIds[0], "game 1", true, "",
		},
		{
			uuid.Nil, "game name", false, "unable to create game, need valid id",
		},
		{
			TestIds[0], "", false, "unable to create game, need non-empty name",
		},
	}
	for _, tc := range testCases {
		game, err := NewGame(tc.id, tc.gameName)
		if tc.expectGame {
			if game == nil {
				t.Fail()
			}
			assert.NotNil(t, game.Board)
			assert.Empty(t, game.Board)
			assert.NotNil(t, game.Players)
			assert.Empty(t, game.Players)
			assert.Equal(t, tc.gameName, game.Name)
			assert.Equal(t, PlayerOne, game.CurrentTurn)
			assert.Equal(t, tc.id, game.Id)
			assert.Equal(t, PlayerPosition(-1), game.Winner)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedErrMsg, err.Error())
		}
	}
}

func Test_AddPlayer(t *testing.T) {
	var testCases = []struct {
		name                 string
		expectedBarrier      int
		expectedPawnLocation Position
	}{
		{"playerOne", 10, Position{X: 8, Y: 16}},
		{"playerTwo", 10, Position{X: 8, Y: 0}},
		{"playerThree", 5, Position{X: 0, Y: 8}},
		{"playerFour", 5, Position{X: 16, Y: 8}},
	}
	game, _ := NewGame(TestIds[0], "AddPlayerGame")
	for idx, tc := range testCases {
		playerPosition, err := game.AddPlayer(TestIds[idx], tc.name)
		assert.Nil(t, err)
		assert.Equal(t, PlayerPosition(idx), playerPosition)
		assert.Len(t, game.Players, idx+1)
		for _, p := range game.Players {
			assert.Equal(t, tc.expectedBarrier, p.Barriers)
		}

		addedPlayer := game.Players[PlayerPosition(idx)]
		assert.Equal(t, tc.name, addedPlayer.PlayerName)
		assert.Equal(t, TestIds[idx], addedPlayer.PlayerId)
		assert.Equal