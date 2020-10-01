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

type authCode struct {
	client      *http.Client
	authURL     string
	tokenURL    string
	redirectURL string
	util        utils.Util
}

func NewAuthCode(authURL, tokenURL, redirectURL string, client *http.Client, util utils.Util) *authCode {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &authCode{
		client:      client,
		authURL:     authURL,
		tokenURL:    tokenURL,
		redirectURL: redirectURL,
		util:        util,
	}
}

func (c *authCode) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error) {
	c.checkTokenURL()

	// todo fix state
	// todo redirect uri
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {client_id},
		"state":         {"rand"},
		"scope":         {scopes},
		"redirect_uri":  {c.redirectURL},
	}
	// TODO support values having multiple values
	for key, val := range customParams {
		v[key] = []string{val}
	}
	redirectURL := fmt.Sprintf("%s?%s", strings.TrimRight(c.authURL, "?"), v.Encode())
	if err := c.util.OpenURL(redirectURL); err != nil {
		return nil, err
	}

	returnedURL, err := c.util.ReadLine("Enter the link you were redirected to after authorization", false)
	if err != nil {
		return nil, fmt.Errorf("error reading redirect URL: %w", err)
	}
	u, err := url.Parse(strings.TrimSpace(returnedURL))
	for err != nil {
		return nil, fmt.Errorf("%s is not a valid url: %w", returnedURL, err)
	}

	code := u.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("URL does not contain the `code` parameter")
	}
	tokenURL, err := url.Parse(c.tokenURL)
	if err != nil {
		return nil, fmt.Errorf("Invalid Token URL:%w", err)
	}

	accessTokenParams := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {client_id},
		"client_secret": {client_secret},
	}
	r := &http.Request{
		Method: "POST",
		URL:    tokenURL,
		Body:   ioutil.NopCloser(strings.NewReader(accessTokenParams.Encode())),
		Header: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error running request: %w", err)
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
	return &token, err
}

func (a *authCode) Refresh(refreshToken string) (string, error) {
	return "", nil
}

func (c *authCode) checkTokenURL() error {
	if c.tokenURL == "" {
		lastSlash := strings.LastIndex(c.authURL, "/")
		defaultURL := c.authURL[:lastSlash+1] + "token"
		tokenURL, err := c.util.ReadLine(fmt.Sprintf("Token URL(%s)", defaultURL), false)
		if err != nil {
			return err
		}
		if tokenURL == "" {
			c.tokenURL = defaultURL
		} else {
			c.tokenURL = tokenURL
		}
	}
	return nil
}
