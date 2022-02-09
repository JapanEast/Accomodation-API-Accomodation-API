package lockitdown

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/rwsargent/boardbots-go/internal/minimax"
)

type (
	Evaluator func(*GameState, PlayerPosition) 