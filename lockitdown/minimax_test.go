
package lockitdown

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rwsargent/boardbots-go/internal/minimax"
	"github.com/stretchr/testify/assert"
)

func TestDepthOfOne(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)

	root := MinimaxNode{
		GameState: game,
		GameMove:  GameMove{},
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	score := minimax.Minimax(&root, 1)
	move, _ := score.(*MinimaxNode)
	fmt.Printf("%T: %+v\n", move.GameMove, move.GameMove)
}
