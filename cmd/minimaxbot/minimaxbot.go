package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/rwsargent/boardbots-go/client"
	"github.com/rwsargent/boardbots-go/internal"
	"github.com/rwsargent/boardbots-go/lockitdown"
)

func main() {
	server := flag.String("server", "http://localhost:8080", "Host of the boardbots server to play on.")
	username := flag.String("username", "", "Username")
	gameId := flag.String("gameId", "", "Game ID")

	flag.Parse()

	if *gameId == "" || *username == "" {
		fmt.Println("Require a game ID and username")
	}

	bbClient, err := client.NewBoardBotClient[lockitdown.TransportState](client.Credentials{
		Username: *username,
	}, *serve