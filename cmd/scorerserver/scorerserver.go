package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/rwsargent/boardbots-go/lockitdown"
)

type (
	ScoreRequest struct {
		GameType  string                    `json:"gameType"`
		GameState lockitdown.TransportState `json:"state"`
		Strategy  string                    `json:"strategy"`
		Player    int                       `json:"player"`
	}

	ScoreResponse struct {
		Score int `json:"score"`
	}
)

func main() {

	port := flag.String("port", ":8888", "server port")
	flag.Parse()

	http.HandleFunc("/api/score", score)

	fmt.Printf("Now listening on port %s\n", *port)
	http.ListenAndServe(*port, nil)

	fmt.Println("done")
}

func score(w http.ResponseWriter, req *http.Request) {
	var scoreReqest ScoreRequest
	err := json.NewDecoder(req.Body).Decode(&scoreReqest)
	defer req.Body.Close()

	fmt.Printf("Request:\n%+v\n", scoreReqest)

	if err != nil {
		fmt.Printf("error reading body, %v\n", err)
	