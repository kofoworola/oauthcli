package clients

import (
	"net/http"
	"net/url"
	"time"
)

// Token is a custom defined token struct for json decoding
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`

	ExpiryTime time.Time
}

type errorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
	ErrorURI    string `json:"error_uri"`
}

type OAuthClient interface {
	GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error)
	//	Refresh(refreshToken string) (string, error)
}

type GenerateAccessTokenOption func(*http.Client, url.Values)
