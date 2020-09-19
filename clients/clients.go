package clients

import (
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type OAuthClient interface {
	GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*oauth2.Token, error)
	Refresh(refreshToken string) (string, error)
}

type GenerateAccessTokenOption func(*http.Client, url.Values)
