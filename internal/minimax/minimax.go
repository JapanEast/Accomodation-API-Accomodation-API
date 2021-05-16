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
		Un