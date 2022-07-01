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
