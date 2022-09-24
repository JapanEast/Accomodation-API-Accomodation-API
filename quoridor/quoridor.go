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
	Typ