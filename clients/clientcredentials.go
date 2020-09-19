package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type clientCredentialClient struct {
	client  *http.Client
	authURL string
}

func NewClientCredentialClient(authURL string, client *http.Client) (*clientCredentialClient, error) {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	return &clientCredentialClient{
		client:  client,
		authURL: authURL,
	}, nil
}

func (c *clientCredentialClient) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*oauth2.Token, error) {
	values := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {client_id},
		"client_secret": {client_secret},
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

	bod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	fmt.Printf("resp is %s", string(bod))
	return nil, nil
}
