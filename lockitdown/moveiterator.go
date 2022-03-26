package lockitdown

import "sync"

type (
	MoveIterator struct {
		game        *G