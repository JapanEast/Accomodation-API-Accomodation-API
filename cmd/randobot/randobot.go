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

	"github.com/rwsargent/board