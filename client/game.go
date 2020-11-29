package client

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwsargent/boardbots-go/internal"
)

// game.go is responsible for the game specific APIs in boardbots.dev

type (
	Game[S any] struct {
		Id        uuid.UUID   `json:"id"`
		LobbyId   uuid.UUID   `json:"lobbyId"`
		Players   []User      `json:"players"`
		GameType  string      `json:"gameType"`
		State     S           `json:"state"`
		Status    string      `json:"status"`
		NumMoves  int         `json:"numMoves"`
		Starte