package client

import (
	"fmt"
	"testing"

	"github.com/rwsargent/boardbots-go/lockitdown"
)

func TestMoves(t *testing.T) {

	bbClient, err := NewBoardBotClient[lockitdown.TransportState](Credentials{
		Username: "tester",
	}, "http://localhost:8080")

	if err != nil {
		fmt.Printf("failed to start cl