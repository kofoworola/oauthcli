package clients

import (
	"net/http"
	"net/url"
)

type OAuthClient interface {
	GenerateAccessToken(client_id, client_secret string) (string, error)
	Refresh(refreshToken string) (string, error)
}

type GenerateAccessTokenOption func(*http.Client, url.Values)
