package lockitdown

import "sort"

var cache map[int][]Placement = make(map[int][]Placement)

type Placement struct {
	position  Pair
	direction Pair
}

type ByCorner []Placement

func (a ByCorner) Len() int     