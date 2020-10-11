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
		Authenticate()