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
	if moveError := isValidPawnMove(newPosition, pawn.Position, game.Board); moveError != nil {
		return moveError
	}
	if game.IsOver() {
		return errors.New("invalid move, game is already over")
	}
	delete(game.Board, pawn.Position)
	pawn.Position = newPosition
	game.Board[pawn.Position] = *pawn
	checkGameOver(game)
	game.nextTurn()
	return nil
}

// GetValidMoveByDirection returns all possible valid positions a pawn can land in a given direction.
// Returns nil if there is a barrier present. If there is a pawn present on the destination square, check to see if
// a barrier is behind that pawn. If so, return possible diagonal positions.
func (board Board) getValidMoveByDirection(pawn, direction Position) []Position {
	// check if there is a barrier in direction
	cursor := Position{Y: pawn.Y + direction.Y, X: pawn.X + direction.X}
	if _, barrierPresent := board[cursor]; barrierPresent {
		return nil
	}

	// Advance pawn again - this is the square should land in.
	cursor.Y += direction.Y
	cursor.X += direction.X

	// check for pawn
	validPositions := make([]Position, 0, 2)
	if _, pawnPresent := board[cursor]; pawnPresent {
		// check for possible jumps
		if _, barrierBeyondPawn := board[Position{Y: cursor.Y + direction.Y, X: cursor.X + direction.X}]; barrierBeyondPawn {
			// look at diagonals instead
			validPositions = append(validPositions, getDiagonalPositions(direction, cursor, board)...)
		} else { // no barrier, final check for a pawn.
			jumpPos := Position{Y: cursor.Y + 2*direction.Y, X: cursor.X + (2 * direction.X)}
			_, finalPawn := board[jumpPos]
			if !finalPawn && isOnBoard(jumpPos) {
				validPositions = append(validPositions, jumpPos)
			}
		}
	} else if isOnBoard(cursor) {
		validPositions = append(validPositions, cursor)
	}
	return validPositions
}

// getDiagonalPositions will return the two positions to the left and right in a given direction.
// For example, the cursor is at position (4, 4) and a vector pointing to the top of the board (0, -1). The diagonal
// positions are [(2, 2), (6, 2)]
func getDiagonalPositions(vector Position, cursor Position, board Board) []Position {
	validPositions := make([]Position, 0, 2)
	leftVector := Position{Y: -1 * vector.X, X: -1 * vector.Y}
	leftTurn := getValidTurnPosition(leftVector, cursor, board)
	if leftTurn.Y != -1 {
		validPositions = append(validPositions, leftTurn)
	}
	rightVector := Position{Y: vector.X, X: vector.Y}
	rightTurn := getValidTurnPosition(rightVector, cursor, board)
	if rightTurn.Y != -1 {
		validPositions = append(validPositions, rightTurn)
	}
	return validPositions
}

// getValidTurnPosition returns a single valid pawn position for a given turn vector. Checks to make sure there are no
// barriers or pawns in the way, and the resulting position is on the board.
func getValidTurnPosition(vector Position, cursor Position, board Board) Position {
	turnCursor := Position{Y: cursor.Y + vector.Y, X: cursor.X + vector.X}
	_, turnBarrier := board[turnCursor]
	if !turnBarrier {
		turnCursor.Y += vector.Y
		turnCursor.X += vector.X

		if _, turnPawn := board[turnCursor]; !turnPawn && isOnBoard(turnCursor) {
			return turnCursor
		}
	}
	return Position{Y: -1, X: -1}
}

// Bounds check for a Position relative to the board.
func isOnBoard(position Position) bool {
	return !(position.Y < 0 || position.Y >= BoardSize || position.X < 0 || position.X >= BoardSize)
}

// Places a barrier on the board at the specified position, for the specified player. Returns an error if the barrier
// cannot be placed.
//
// This function will advance the player's turn and decrement the player's barrier count.
func (game *Game) PlaceBarrier(position Position, player PlayerPosition) error {
	if game.CurrentTurn != player {
		return errors.New(fmt.Sprintf("wrong turn, current turn is for Player: %d", game.CurrentTurn))
	}
	if invalidPosition(position) {
		return errors.New("invalid location for a barrier")
	}
	if playerHasNoMoreBarriers(game.Players[player]) {
		return errors.New("the player has no more barriers to play")
	}
	barrierPositions := createBarrierPositions(position)
	if barriersAreInTheWay(barrierPositions, game.Board) {
		return errors.New("the new barrier intersects with another")
	}
	if barrierPreventsWin(barrierPositions, game) {
		return errors.New("the barrier prevents a players victory")
	}
	if game.IsOver() {
		return errors.New("invalid move, game is already over")
	}
	game.Players[player].Barriers--
	for _, pos := range barrierPositions {
		game.Board[pos] = Piece{Position: pos, Owner: player, Type: Barrier}
	}
	game.nextTurn()
	return nil
}

// You can never place a pawn or barrier at a double-odd position (the intersections of the gutters), or on the very
// last row and column of the board.
func invalidPosition(position Position) bool {
	return position.Y&0x1 == position.X&0x1 || // both col and row are even or odd
		// can't be on the last valid row/
		!(position.Y < BoardSize-1 &&
			(position.X < BoardSize-1))
}

func playerHasNoMoreBarriers(player *Player) bool {
	return player.Barriers <= 0
}

// Check if a placed barrier will block a player's path to the goal. Uses a simple A-* algorithm.
func barrierPreventsWin(positions [3]Position, game *Game) bool {
	for _, position := range positions {
		game.Board[position] = Piece{Position: position, Owner: PlayerOne}
	}
	//remove those temporary barriers no matter what
	defer func() {
		for _, position := range positions {
			delete(game.Board, position)
		}
	}()

	for playerPosition, player := range game.Players {
		path := game.FindPath(pla