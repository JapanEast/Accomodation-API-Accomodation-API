package lockitdown

import "sort"

var cache map[int][]Placement = make(map[int][]Placeme