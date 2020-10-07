package clients

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kofoworola/oauthcli/utils"
)

const (
	accTokenString = "access_token="
	errString      = "err="
)

type implicitClient struct {
	authURL     string
	redirectURL string
	util        utils.Util
}

func NewImplicitClient(authURL, redirectURL string, util utils.Util) *implicitClient {
	return &implicitClient{
		authURL:     authURL,
		redirectURL: redirectURL,
		util:        util,
	}
}

func (c *implicitClient) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error) {
	values := url.Values{
		"response_type": {"token"},
		"client_id":     {client_id},
		"scope":         {scopes},
		"state":         {"dummy_state"},
		"redirect_url":  {c.redirectURL},
	}
	for key, val := range customParams {
		values.Add(key, val)
	}

	u, err := url.Parse(c.authURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	authURL := fmt.Sprintf("%s?%s", strings.TrimRight(u.String(), "?"), values.Encode())
	c.util.OpenURL(authURL)

	redirected, err := c.util.ReadLine("Enter the link you were redirected to after authorization", false)
	if err != nil {
		return nil, fmt.Errorf("error reading redirect URL: %w", err)
	}

	redirectedURL, err := url.Parse(redirected)
	if err != nil {
		return nil, fmt.Errorf("error parsing redirected URL: %w", err)
	}

	fragment := redirectedURL.Fragment
	if fragment == "" {
		return nil, fmt.Errorf("redirected URL fragment should not be empty")
	}
	if !strings.Contains(fragment, accTokenString) {
		if strings.Contains(fragment, errString) {
			return nil, fmt.Errorf("%s", strings.TrimPrefix(fragment, errString))
		}
		return nil, fmt.Errorf("redirected URL does not have an access_token fragment: %s\n", redirectedURL.Fragment)
	}

	if strings.Contains(fragment, "&") {
		fragment = fragment[len(accTokenString):strings.Index(fragment, "&")]
	} else {
		fragment = strings.TrimPrefix(fragment, accTokenString)
	}

	token := Token{
		AccessToken: fragment,
	}
	expires := redirectedURL.Query().Get("expires_id")
	if expires != "" {
		expiresInt, err := strconv.Atoi(expires)
		if err != nil {
			fmt.Printf("WARNING: could not parse %s to expiry time", expires)
		} else {
			token.ExpiryTime = time.Now().Add(time.Second * time.Duration(expiresInt))
		}
	}

	return &token, nil
}
