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
		user