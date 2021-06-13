
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
		GameDef:          gameDef,
		Players:          players,
		Robots:           make([]Robot, 0),
		PlayerTurn:       0,
		MovesThisTurn:    gameDef.MovesPerTurn,
		RequiresTieBreak: false,
		Winner:           -1,
		saveStack:        make([]SaveState, 0),
	}
}

// Only intended for Unit pairs
func (p *Pair) Rotate(direction TurnDirection) {
	s := p.S()
	if direction == Right {
		p.Q = -p.R
		p.R = -s
	} else {
		p.R = -p.Q
		p.Q = -s
	}
}

func (game *GameState) Move(move *GameMove) error {
	if move.Player != PlayerPosition(game.PlayerTurn) {
		return fmt.Errorf("wrong player, expected %d, was %d", game.PlayerTurn, move.Player)
	}

	game.saveState()
	err := move.Move(game)

	if err != nil {
		return err
	}

	// Resolve move
	if err = game.resolveMove(); err != nil {
		return err
	}

	if game.MovesThisTurn == 0 {
		game.PlayerTurn = PlayerPosition((int(game.PlayerTurn) + 1) % len(game.Players))
		game.MovesThisTurn = 3
	}

	if over, winner := game.checkGameOver(); over {
		game.Winner = winner
		return fmt.Errorf("winner is %d", winner+1)
	}
	return nil
}

func (game *GameState) Undo(move *GameMove) error {
	save := game.saveStack[len(game.saveStack)-1]

	game.Robots = save.bots

	for i, player := range save.players {
		game.Players[i].PlacedRobots = player.PlacedRobots
		game.Players[i].Points = player.Points
	}

	game.PlayerTurn = save.player
	game.MovesThisTurn = save.movesThisTurn

	game.saveStack = game.saveStack[:len(game.saveStack)-1]
	return nil
}

func (game *GameState) RobotAt(hex Pair) *Robot {
	for i := 0; i < len(game.Robots); i++ {
		robot := &game.Robots[i]
		if robot.Position == hex {
			return robot
		}
	}
	return nil
}

func (game *GameState) resolveMove() error {
	for resolved := false; !resolved; {
		targeted := game.taretedRobots()

		if tiebreaks := game.checkForTieBreaks(targeted); len(tiebreaks) > 0 {
			game.RequiresTieBreak = true
			json, _ := game.ToJson()
			return TieBreak{
				Robots: tiebreaks,
				State:  json,
			}
		}

		resolved = game.updateLockedRobots(targeted)
	}
	return nil
}

func (game *GameState) updateLockedRobots(targeted map[*Robot][]*Robot) bool {
	resolved := true
	doomed := []int{}
	for i, _ := range game.Robots {
		robot := &game.Robots[i]
		attackers, found := targeted[robot]
		if !found || len(attackers) == 1 {
			if robot == game.activeBot {
				// The active bots state is controlled by the move, until
				// 'released'.
				continue
			}
			beam := robot.IsBeamEnabled
			lock := robot.IsLockedDown
			// Enable bot
			robot.IsLockedDown = false
			robot.IsBeamEnabled = !game.isCorridor(robot.Position)

			// State change, reevaluate
			if beam != robot.IsBeamEnabled || lock != robot.IsLockedDown {
				resolved = false
			}
		} else if len(attackers) == 3 {
			doomed = append(doomed, i)
			game.shutdownRobot(i, attackers)
			resolved = false
		} else if len(attackers) == 2 {
			robot.Disable()
		}
	}
	for _, doomedIdx := range doomed {
		game.Robots = append(game.Robots[:doomedIdx], game.Robots[doomedIdx+1:]...)
	}
	return resolved
}

// If any "doomed" robots (locked or shutdown) are also part of a lock or shut down,
// we need to break a tie.
func (game *GameState) checkForTieBreaks(targeted map[*Robot][]*Robot) []*Robot {
	tiebreaks := make([]*Robot, 0, 2)
	for doomed, attackers := range targeted {
		// TODO(rwsargent) update targeted to be a *Robot -> *Robot map.