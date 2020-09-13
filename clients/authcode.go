package clients

import (
	"net/http"
	"net/url"
	"time"
)

type AuthCodeClient struct {
	client      *http.Client
	authURL     string
	tokenURL    string
	redirectURL string
}

func NewAuthCodeClient(authURL, tokenURL, redirectURL string, client *http.Client) *AuthCodeClient {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &AuthCodeClient{
		client:      client,
		authURL:     authURL,
		tokenURL:    tokenURL,
		redirectURL: redirectURL,
	}
}

func (a *AuthCodeClient) AccessToken(client_id, client_secret string) (string, error) {
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {client_id},
		"state":         {"TEST STATE"},
		"scope":         {"TEST SCOPE"},
	}

	return v.Encode(), nil
}
