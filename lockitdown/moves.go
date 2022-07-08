package lockitdown

import (
	"errors"
	"fmt"
)

type (
	Mover interface {
		Move(*GameState, PlayerPosition) error
		ToTransport() BoardbotsMove
	}

	GameMove struct {
		Player PlayerPosition
		Mover
	}

	BoardbotsMove struct {
		Position Pair `json:"pos"`
		Action   any  `json:"action"`
	}

	TurnRobot struct {
		Robot     Pair
		Direction TurnDirection
	}

	InnerTurnRobotT struct {
		Side string `json:"side"`
	}
	TurnRobotT struct {
		Turn InnerTurnRobotT `json:"Turn"`
	}

	PlaceRobot struct {
		Robot, Direction Pair
	}

	InnerPlaceRobotT struct {
		Dir Pair `json:"dir"`
	}

	PlaceRobotT struct {
		PlaceRobot InnerPlaceRobotT
	}

	AdvanceRobot struct {
		Robot Pair
	}
)

func NewMove(m Mover, p PlayerPosition) *GameMove {
	return &GameMove{
		Player: p,
		Mover:  m,
	}
}
func (m *GameMove) Move(state *GameState) error {
	return m.Mover.Move(state, m.Player)
}

func (m *AdvanceRobot) Move(game *GameState, player PlayerPosition) error {
	robot := game.RobotAt(m.Robot)
	if robot == nil {
		return fmt.Errorf("no robot at location %v", m.Robot)
	}
	if robot.IsLockedDown {
		return errors.New("cannot advance, robot is locked down")
	}
	if robot.Player != player {
		return fmt.Errorf("cannot move %s, it belongs to Player %d", m.Robot.String(), robot.Player)
	}
	advanceSpot := robot.Position.Copy()
	advanceSpot.Plus(robot.Direction)
	if block := game.RobotAt(advanceSpot); block != nil {
		// Undo move
		return errors.New("cannot advance, another bot in the way")
	}
	robot.Position.Plus(robot.Direction)

	game.MovesThisTurn -= 1

	// // Evaluate state before turning on beam
	game.resolveMove()

	robot.IsBeamEnabled = !game.isCorridor(robot.Position) && !robot.IsLockedDown
	return nil
}

func (m AdvanceRobot) ToTransport() BoardbotsMove {
	return BoardbotsMove{
		Position: m.Robot,
		Action:   "Advance",
	}

}

func (m AdvanceRobot) String() string {
	return fmt.Sprintf("Move %s", m.Ro