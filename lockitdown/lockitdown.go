
// An implmentation of boardbots.dev LockItDown game. Users can use package
// to make and undo moves on an internal state of the board, and then apply
// those moves to a game hosted on boardbots.dev.
package lockitdown

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
)

type (
	Pair struct {
		Q int `json:"q"`
		R int `json:"r"`
	}
	Player struct {
		Points       int `json:"points"`
		PlacedRobots int `json:"placedRobots"`
	}

	PlayerPosition int

	BoardType struct {
		ArenaRadius int `json:"arenaRadius"`
	}

	Board struct {
		HexaBoard BoardType
	}

	GameDef struct {
		Board           Board  `json:"board"`
		Players         int    `json:"numOfPlayers"`
		MovesPerTurn    int    `json:"movesPerTurn"`
		RobotsPerPlayer int    `json:"robotsPerPlayer"`
		WinCondition    string `json:"winCondition"`
	}

	Robot struct {
		Position                    Pair
		Direction                   Pair
		IsBeamEnabled, IsLockedDown bool
		Player                      PlayerPosition
	}

	GameState struct {
		GameDef          GameDef
		Players          []*Player
		Robots           []Robot
		PlayerTurn       PlayerPosition
		MovesThisTurn    int
		RequiresTieBreak bool
		Winner           int
		activeBot        *Robot
		saveStack        []SaveState
	}

	TurnDirection int

	TieBreak struct {
		Robots []*Robot
		State  string
	}
)

const (
	Left TurnDirection = iota
	Right
)

var (
	NW Pair = Pair{0, -1}
	NE Pair = Pair{1, -1}
	E  Pair = Pair{1, 0}
	SE Pair = Pair{0, 1}
	SW Pair = Pair{-1, 1}
	W  Pair = Pair{-1, 0}

	Cardinals = []Pair{E, SE, SW, W, NW, NE}

	moveBufferPool = sync.Pool{
		New: func() any {
			s := make([]*GameMove, 0, 128)
			return &s
		},
	}

	movePool = sync.Pool{
		New: func() any {
			return new(GameMove)
		},
	}
)

func (p *Pair) Plus(that Pair) {
	p.Q += that.Q
	p.R += that.R
}

func (p *Pair) Minus(that Pair) {
	p.Q -= that.Q
	p.R -= that.R
}

func (p Pair) String() string {
	return fmt.Sprintf("{%d, %d}", p.Q, p.R)
}

func (p Pair) S() int {
	return -p.Q - p.R
}

func (p Pair) Copy() Pair {
	return p
}

func (p Pair) Dist() int {
	return (intAbs(p.Q) + intAbs(p.R) + intAbs(p.S())) / 2
}

func (r *Robot) Disable() {
	r.IsBeamEnabled = false
	r.IsLockedDown = true
}

func (r *Robot) Enable() {
	r.IsBeamEnabled = true
	r.IsLockedDown = false
}

func NewGame(gameDef GameDef) *GameState {
	players := make([]*Player, gameDef.Players)
	for i := 0; i < len(players); i++ {
		players[i] = &Player{}
	}
	return &GameState{