package client

import (
	"fmt"
	"testing"

	"github.com/rwsargent/boardbots-go/lockitdown"
)

func TestMoves(t *testing.T) {

	bbClient, err := NewBoardBotClient[lockitdown.TransportSta