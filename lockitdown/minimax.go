package lockitdown

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/rwsargent/boardbots-go/internal/minimax"
)

type (
	Evaluator func(*GameState, PlayerPosition) int

	MinimaxNode struct {
		GameState    *GameState
		GameMove     GameMove
		Searcher     PlayerPosition
		Evaluator    Evaluator
		MinimaxValue int
	}
)

var nodePool = sync.Pool{
	New: func() any {
		return &MinimaxNode{}
	},
}

func (n *MinimaxNode) Evaluate() {
	n.MinimaxValue = n.Evaluator(n.GameState, n.Searcher)
}

func (n *MinimaxNode) Children(nodeBuffer []minimax.Node) []minimax.Node {
	moveBuffer := moveBufferPool.Get().(*[]*GameMove)
	defer moveBufferPool.Put(moveBuffer)

	nextMoves := n.GameState.PossibleMoves([]GameMove{})

	for _, nextMove := range nextMoves {
		node := nodePool.Get().(*MinimaxNode)
		node.GameState = n.GameState
		node.GameMove = nextMove
		node.Searcher = n.Searcher
		node.Evaluator = n.Evaluator
		nodeBuffer = append(nodeBuffer, node)
	}
	return nodeBuffer
}

func (n *MinimaxNode) ShouldMaximize() bool {
	return n.Searcher == n.GameState.PlayerTurn
}

func (n *MinimaxNode) Move() {
	err := n.GameState.Move(&n.GameMove)
	if err != nil {
		json, _ := n.GameState.ToJson()
		panic(fmt.Errorf("%s.\n\n%s", err, json))
	}
}

func (n *MinimaxNode) Undo() {
	n.GameState.Undo(&n.GameMove)
}

func (n *MinimaxNode) Score() int {
	return n.MinimaxValue
}

func (n *MinimaxNode) SetScore(score int) {
	n.MinimaxValue = score
}

func (n *MinimaxNode) Release() {
	movePool.Put(&n.GameMove)
	nodePool.Put(n)
}

func MinimaxWithIterator(node *MinimaxNode, depth int) MinimaxNode {
	it := NewMoveIt