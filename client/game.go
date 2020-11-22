package client

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwsargent/boardbots-go/internal"
)

// game.go is responsible for the game specific APIs in boardbots.dev

type (
	Game[S any] struct {
		Id        uu