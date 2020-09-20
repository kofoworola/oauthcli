package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type clientCredentials struct {
	client  *http.Client
	authURL string
}

func NewClientCredentialClient(authURL string, client *http.Client) *clientCredentials {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	return &clientCredentials{
		client:  client,
		authURL: authURL,
	}
}

func (c *clientCredentials) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error) {
	values := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"scope":         {scopes},
	}
	for key, val := range customParams {
		values.Add(key, val)
	}

	// parse url
	u, err := url.Parse(c.authURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    u,
		Body:   ioutil.NopCloser(strings.NewReader(values.Encode())),
		Header: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error running http request: %w", err)
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
	token.ExpiryTime = time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	return &token, nil
}
