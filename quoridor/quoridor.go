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

// Initialize with default values, and supplied game Id and Name.
//
// The game is initialized with an empty board and player slice. Defaults current turn to PlayerOne and winner to
// -1. Everything else defaults to their zero value.
func NewGame(id uuid.UUID, name string) (*Game, error) {
	if id == uuid.Nil {
		return nil, errors.New("unable to create game, need valid id")
	}
	if name == "" {
		return nil, errors.New("unable to create game, need non-empty name")
	}
	return &Game{
		Board:       make(map[Position]Piece),
		Players:     make(map[PlayerPosition]*Player),
		CurrentTurn: PlayerOne,
		Id:          id,
		Name:        name,
		Winner:      -1,
	}, nil
}

// Adds a new player to the player map at the next possible player position. Will also update the barrier count when
// the player count goes from two to three.
// Players can only be added if the game has not yet started, and they don't already exist in the game.
func (game *Game) AddPlayer(id uuid.UUID, name string) (PlayerPosition, error) {
	if !game.StartDate.IsZero() {
		return -1, errors.New(fmt.Sprintf("cannot add player %s, game has already started", name))
	}
	// Make sure the player isn't already a part of this game. The same player cannot play against themselves.
	for _, player := range game.Players {
		if player.PlayerId == id {
			return 0, errors.New(fmt.Sprintf("player with id %s alreayd in this game", id.String()))
		}
	}
	barriersForPlayer := 10
	if len(game.Players) >= 2 {
		barriersForPlayer = 5
	}
	// For each possible player
	for playerNumber := PlayerOne; playerNumber <= PlayerFour; playerNumber++ {
		p, present := game.Players[playerNumber]
		if present {
			// Make sure they have the correct number of barriers
			p.Barriers = barriersForPlayer
			game.Players[playerNumber] = p
		} else {
			playerPawn := Piece{
				Position: startingPositions[playerNumber],
				Owner:    playerNumber,
				Type:     Pawn,
			}
			// Create a new player with barrier count, starting position, etc.
			game.Players[playerNumber] = &Player{
				Barriers:   barriersForPlayer,
				PlayerId:   id,
				PlayerName: name,
				Pawn:       playerPawn,
			}
			// Add pawn to board
			game.Board[playerPawn.Position] = playerPawn
			return playerNumber, nil
		}
	}
	return -1, errors.New("no open player positions in this game")
}

// Starts a game by setting the StartDate to the current instant of time. Returns an error if there aren't enough
// players, or the game has already started.
func (game *Game) StartGame() error {
	if !(len(game.Players) == 2 || len(game.Players) == 4) {
		return errors.New(fmt.Sprintf("can't start game, wrong number of players (%d)", len(game.Players)))
	}
	if !game.StartDate.IsZero() {
		return errors.New(fmt.Sprintf("game already started"))
	}
	game.StartDate = time.Now()
	return nil
}

// Moves a pawn to the given new position for the give player. Returns an error if the move is invalid.
//
// The move is invalid if it's an invalid pawn location, the wrong player's turn, or the game is over.
func (game *Game) MovePawn(newPosition Position, player PlayerPosition) error {
	pawn := &game.Players[player].Pawn
	if !isValidPawnLocation(newPosition) {
		return errors.New("invalid Pawn Location")
	}
	if game.CurrentTurn != player {
		return errors.New(fmt.Sprintf("wrong turn, current turn is for Player: %d", game.CurrentTurn))
	}
	if moveError := isValidPawnMove(newPosition, pawn.Position, game.Board); moveError 