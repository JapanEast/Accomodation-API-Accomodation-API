// randobot.go makes random moves for the given lockitdown game.
// You need to run randobot with a username, gameId, and server address.
//
// $> randobot -username=randobot -gameId=00000000-...-0000 -server=https://boardbots.dev
//
// It is recommend to run setupbots before using randobot.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/rwsargent/boardbots-go/client"
	"github.com/rwsargent/boardbots-go/internal"
	"github.com/rwsargent/boardbots-go/lockitdown"
)

var edges = []struct{ pos, dir internal.Position }{
	{pos: internal.Position{Q: 0, R: -5}, dir: internal.Position{Q: 0, R: 1}},
	{pos: internal.Position{Q: 5, R: -5}, dir: internal.Position{Q: -1, R: 1}},
	{pos: internal.Position{Q: 5, R: 0}, dir: internal.Position{Q: -1, R: 0}}