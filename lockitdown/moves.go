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
