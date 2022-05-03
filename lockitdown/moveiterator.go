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
		botCache[robot.Position] = &game.Robots[i]
	}
	it := &MoveIterator{
		game:        game,
		currentMove: new(GameMove),
		moveBuf:     [3]GameMove{},
		edgeIndex:   0,
		moveIdx:     -1,
		robotIndex:  0,
		botCache:    botCache,
	}
	return it
}

func (it *MoveIterator) Get() *GameMove {
	return it.currentMove
}

func (it *MoveIterator) Next() bool {
	it.findNext()
	return it.currentMove != nil
}

func (it *MoveIterator) findNext() {
	// Check to see if we have any buffered moves
	// already calculated
	for it.moveIdx >= 0 && it.moveIdx < len(it.moveBuf) {
		next := &it.moveBuf[it.moveIdx]
		it.currentMove = next

		it.moveIdx++
		if it.currentMove.Mover != nil {
			return
		}
	}

	// Check all owned bots. robotIndex refers to the next bot
	// to buffer moves for.

	ringSize := it.game.GameDef.Board.HexaBoard.ArenaRadiu