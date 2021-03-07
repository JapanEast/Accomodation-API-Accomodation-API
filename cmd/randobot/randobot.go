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
	{pos: internal.Position{Q: 5, R: 0}, dir: internal.Position{Q: -1, R: 0}},
	{pos: internal.Position{Q: 0, R: 5}, dir: internal.Position{Q: 0, R: -1}},
	{pos: internal.Position{Q: -5, R: 5}, dir: internal.Position{Q: 1, R: -1}},
	{pos: internal.Position{Q: -5, R: 0}, dir: internal.Position{Q: 1, R: 0}},
}

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
	}, *server)

	if err != nil {
		fmt.Printf("failed to start client, %s\n", err.Error())
		return
	}

	err = bbClient.Authenticate()
	if err != nil {
		fmt.Printf("could not authenticate %s", err.Error())
	}

	tGame, err := bbClient.Game(*gameId)
	if err != nil {
		panic(err)
	}
	playerPosition := getPlayerPosition(tGame, bbClient.Credentials.Username)

	game := lockitdown.StateFromTransport(&tGame.State)

	rand.Seed(63)

	for game.Winner < 0 {
		if playerPosi