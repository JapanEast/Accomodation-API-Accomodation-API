
package client

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// lobby.go is responsible for the lobby specific APIs in boardbots.dev

type (
	CreateLobbyReq struct {
		GameType       string `json:"gameType"`
		InitialPlayers []int  `json:"initialPlayers"`
	}

	CreateLobbyResponse struct {
		Id        uuid.UUID   `json:"id"`
		Host      User        `json:"host"`
		GameType  string      `json:"gameType"`
		Status    string      `json:"status"`
		Players   []User      `json:"players"`
		CreatedAt json.Number `json:"createdAt"`
	}

	JoinLobbyReq struct {
	}

	JoinLobbyResponse struct {
	}

	Empty struct{}

	StartGameResp struct {
	}
)

func (c *BoardBotClient[S]) CreateLobby() (CreateLobbyResponse, error) {
	lobbyReq := CreateLobbyReq{
		GameType:       "lockitdown",
		InitialPlayers: []int{},
	}

	return Post[CreateLobbyReq, CreateLobbyResponse](c, "/api/lobby/create", lobbyReq)
}

func (c *BoardBotClient[S]) JoinLobby(lobbyId string) error {