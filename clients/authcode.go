package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/kofoworola/oauthcli/utils"
	"golang.org/x/oauth2"
)

type authCodeClient struct {
	client      *http.Client
	authURL     string
	tokenURL    string
	redirectURL string
}

func NewAuthCodeClient(authURL, tokenURL, redirectURL string, client *http.Client) *authCodeClient {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &authCodeClient{
		client:      client,
		authURL:     authURL,
		tokenURL:    tokenURL,
		redirectURL: redirectURL,
	}
}

func (a *authCodeClient) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*oauth2.Token, error) {
	a.checkTokenURL()
	// todo redirect uri
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {client_id},
		"state":         {"TEST STATE"},
		"scope":         {scopes},
	}
	// TODO support values having multiple values
	for key, val := range customParams {
		v[key] = []string{val}
	}
	redirectURL := fmt.Sprintf("%s?%s", strings.TrimRight(a.authURL, "?"), v.Encode())
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", redirectURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", redirectURL).Start()
	case "darwin":
		err = exec.Command("open", redirectURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("error opening redirect url, manually visit %s to authorize\n", redirectURL)
	}

	returnedURL, err := utils.ReadLine("Enter the link you were redirected to after authorization")
	// TODO handle error properly
	// TODO loop de loop
	if err != nil {
		fmt.Printf("error reading redirect URL: %v", err)
	}
	u, err := url.Parse(strings.TrimSpace(returnedURL))
	for err != nil {
		returnedURL, err = utils.ReadLine(fmt.Sprintf("Invalid URL (%v), please enter the URL you were redirected to:\n>", err))
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}
		u, err = url.Parse(returnedURL)
	}

	code := u.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("URL does not contain the `code` parameter")
	}
	tokenURL, err := url.Parse(a.tokenURL)
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
	resp, err := a.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error running request: %w", err)
	}

	var token oauth2.Token
	// TODO return json response
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &token, err
}

func (a *authCodeClient) Refresh(refreshToken string) (string, error) {
	return "", nil
}

// TODO move error handling out
func (a *authCodeClient) checkTokenURL() {
	for a.tokenURL == "" {
		lastSlash := strings.LastIndex(a.authURL, "/")
		defaultURL := a.authURL[:lastSlash+1] + "token"
		tokenURL, err := utils.ReadLine(fmt.Sprintf("Token URL(%s)", defaultURL))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if tokenURL == "" {
			a.tokenURL = defaultURL
		} else {
			a.tokenURL = tokenURL
		}
	}
}
