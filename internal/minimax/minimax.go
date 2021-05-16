package minimax

import (
	"math"
	"sync"
)

type (
	Node interface {
		Evaluate()
		Children([]Node) []Node
		Release()
		Move()
		Undo()
		ShouldMaximize() bool
		Score() int
		SetScore(int)
	}
)

var nodePool = sync.Pool{
	New: func() any {
		buf := make([]Node, 0, 