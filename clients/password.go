package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kofoworola/oauthcli/utils"
)

type passwordClients struct {
	authURL  string
	password string
	username string

	client *http.Client
	util   utils.Util
}

func NewPasswordClient(authURL, username, password string, client *http.Client, util utils.Util) *passwordClients {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &passwordClients{
		authURL:  authURL,
		password: password,
		username: username,
		client:   client,
		util:     util,
	}
}

func (c *passwordClients) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error) {
	if c.password == "" {
		var err error
		c.password, err = c.util.ReadLine("Enter password", true)
		if err != nil {
			return nil, fmt.Errorf("error reading password: %w", err)
		}
	}
	values := url.Values{
		"grant_type":    {"password"},
		"username":      {c.username},
		"password":      {c.password},
		"scope":         {scopes},
		"client_id":     {client_id},
		"client_secret": {client_secret},
	}
	for key, val := range customParams {
		values.Add(key, val)
	}

	u, err := url.Parse(c.authURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(strings.NewReader(values.Encode())),
		URL:    u,
		Header: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error running auth resuest: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %w", err)
		}
		return nil, fmt.Errorf("%s - %s", errResp.Error, errResp.Description)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	token.ExpiryTime = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return &token, nil
}
