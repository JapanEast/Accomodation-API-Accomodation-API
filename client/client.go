// An implementaiton of the public API for boardbots.dev.
// BoardBotClient wraps an HTTPClient
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

type (
	Client interface {
		Authenticate() error
		MakeMove() error
		StartGame() error
		JoinLobby() error
		CreateLobby() error
	}

	Credentials struct {
		Username string
	}
	BoardBotClient[State any] struct {
		Credentials Credentials
		httpClient  *http.Client
		domain      string
		user        User
	}

	User struct {
		Id   json.Number `json:"id"`
		Name string      `json:"name"`
	}
)

// NewBoardBotClient builds a BoardBotClient with the users credentials, and server address.
// Parameterized with the type of Game (via its state). You'll need to instantiate
// multip