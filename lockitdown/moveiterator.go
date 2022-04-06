package lockitdown

import "sync"

type (
	MoveIterator struct {
		game        *GameState
		currentMove *GameMove
		moveBuf     [3]GameMove
		moveIdx     int
		edgeIndex   int
		robotIndex  int
		botCache    map[Pair]*Robot
	}
)

var (
	advancePool = sync.Pool{
		New: func() any { return new(AdvanceRobot) },
	}

	turnPool = sync.Pool{
		New: func() any { return new(TurnRobot) },
	}

	placePool = sync.Pool{
		New: func() any { return new(PlaceRobot) },
	}
)

func NewMoveIterator(game *GameState) *MoveIterator {
	botCache := make(map[Pair]*Robot)
	for i := 0; i < len(game.Robots); i++ {
		robot := &game.Robots[i]
		botCache[robot.Posit