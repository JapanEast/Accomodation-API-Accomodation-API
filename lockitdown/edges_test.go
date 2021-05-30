package lockitdown

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	assert.Len(t, edges(1), 18)
	assert.Len(t, edges(2), 42)
}

func TestCardinalDirections(t *testing.T) {
	edges := edges(1)
	directions := make(map[Pair][]Pair)
	for _, edge := range edges {
		if dir, found := directions[edge.position]; found {
			directions[edge.position] = append(dir, edge.direction)