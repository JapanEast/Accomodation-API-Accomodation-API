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
		GameState lockitdown.TransportSt