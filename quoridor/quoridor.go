package quoridor

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	// Refers to the position in the Game's Player slice.
	PlayerPosition int

	// Identifies the type of the piece - either a 'p' for Pawn, or 'b' for Barrier.
	TypeId rune

	// The coordinates to a cell on the game board.
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	// Wrapper type around the map from position to Piece. The board represents pieces that are on it. If a position
	// isn't in the map, that position doesn't have a piece.
	Board map[Position]Piece

	//
	Player struct {
		// The number of remaining barriers left to the player
		Barriers int

		// A copy of this players pawn in this game's board.
		Pawn Piece

		PlayerId uuid.UUID

		PlayerName string
	}

	// The full representation of a Quoridor game.
	Game struct {
		// The game board.
		Board              Board
		Players            map[PlayerPosition]*Player
		Id                 uuid.UUID
		CurrentTurn        PlayerPosition
		StartDate, EndDate time.Time
		Winner             PlayerPosition
		Name               string
	}

	Piece struct {
		Position Position
		Owner    PlayerPosition
		Type     TypeId
	}

	Move struct {
		Player PlayerPosition
		Delta  []Position
	}
)

// An enumeration of all possible player positions.
const (
	PlayerO