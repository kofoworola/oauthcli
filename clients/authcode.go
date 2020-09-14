package clients

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

func (a *AuthCodeClient) GenerateAccessToken(client_id, client_secret, scopes string) (string, error) {
	// todo redirect uri
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {client_id},
		"state":         {"TEST STATE"},
		"scope":         {scopes},
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter the link you were redirected to after authorization:\n>")
	returnedURL, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	u, err := url.Parse(returnedURL)
	for err != nil {
		fmt.Printf("Invalid URL, please enter the URL you were redirected to:\n>")
		returnedURL, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}
		u, err = url.Parse(returnedURL)
	}

	code := u.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("URL does not contain the `code` parameter")
	}
	tokenURL, err := url.Parse(a.tokenURL)
	if err != nil {
		return "", fmt.Errorf("Invalid Token URL:%w", err)
	}

	accessTokenParams := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {client_id},
		"client_secret": {client_secret},
	}
	r := &http.Request{
		Method:   "POST",
		URL:      tokenURL,
		PostForm: accessTokenParams,
	}

	resp, err := a.client.Do(r)
	if err != nil {
		return "", fmt.Errorf("error running request: %w", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}
	return string(respBytes), err
}

func (a *AuthCodeClient) Refresh(refreshToken string) (string, error) {
	return "", nil
}
