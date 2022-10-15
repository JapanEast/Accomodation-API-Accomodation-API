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
	PlayerOne PlayerPosition = iota
	PlayerTwo
	PlayerThree
	PlayerFour
)

const (
	Pawn    TypeId = 'p'
	Barrier TypeId = 'b'
)
const BoardSize int = 17

var (
	// Represents the row or column a Players pawn has to be to win the game. A value of -1 in X or Y means any value on
	// that row or column is part of a winning position.
	//
	// For example, PlayerOne can win when their pawn reaches the 'top' row. If the pawn reaches {Y: 0, X:2..16}
	// PlayerOne wins.
	winningPositions = map[PlayerPosition]Position{
		PlayerOne:   {Y: 0, X: -1},
		PlayerTwo:   {X: -1, Y: 16},
		PlayerThree: {Y: -1, X: 16},
		PlayerFour:  {Y: -1, X: 0},
	}

	startingPositions = map[PlayerPosition]Position{
		PlayerOne:   {X: 8, Y: 16}, // Bottom
		PlayerTwo:   {X: 8, Y: 0},  // Top
		PlayerThree: {X: 0, Y: 8},  // Left
		PlayerFour:  {X: 16, Y: 8}, // Right
	}

	directions = []Position{
		{X: 1},
		{Y: 1},
		{X: -1},
		{Y: -1},
	}
)

// Initialize with default values, a