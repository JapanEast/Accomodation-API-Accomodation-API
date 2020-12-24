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
	server := flag.String("server", "http://loc