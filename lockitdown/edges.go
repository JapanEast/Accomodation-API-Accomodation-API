package lockitdown

import "sort"

var cache map[int][]Placement = make(map[int][]Placement)

type Placement struct {
	position  Pair
	direction Pair
}

type ByCorner []Placement

func (a ByCorner) Len() int      { return len(a) }
func (a ByCorner) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCorner) Less(i, j int) bool {
	return min(a[i].position) 