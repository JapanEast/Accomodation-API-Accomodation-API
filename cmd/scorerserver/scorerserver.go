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
		Player    int                       `json:"pla