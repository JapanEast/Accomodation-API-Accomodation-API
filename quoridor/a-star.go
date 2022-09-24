
package quoridor

import (
	"container/heap"
	"math"
)

/**
 * A struct for the priority queue, holds the Position and priority of the Node in the Board graph
 */
type PQNode struct {
	position Position
	prev     *PQNode

	distance, priority int
}

func absInt(val int) int {
	y := val >> 31
	return (val ^ y) - y
}
