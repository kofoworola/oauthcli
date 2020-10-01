package clients

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type implicitClient struct {
	authURL string

	client *http.Client
}

func NewImplicitClient(authURL string, client *http.Client) *implicitClient {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	return &implicitClient{
		client:  client,
		authURL: authURL,
	}
}

func (c *implicitClient) GenerateAccessToken(client_id, client_secret, scopes string, customParams map[string]string) (*Token, error) {
	values := url.Values{
		"response_type": {"token"},
		"client_id":     {client_id},
		"scope":         {scopes},
		"state":         {"dummy_state"},
	}
	for key, val := range customParams {
		values.Add(key, val)
	}

	_, err := url.Parse(c.authURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	//	authURL := fmt.Sprintf("%s?%s", strings.TrimRight(u.String(), "?"), values.Encode())
	//	utils.OpenURL(authURL)
	//
	//	returnedURL, err := utils.ReadLine("Enter the link you were redirected to after authorization", os.Stdin)
	//	if err != nil {
	//		return nil, fmt.Errorf("error reading redirect URL: %w", err)
	//	}

	//	u, err = url.Parse(strings.TrimSpace(returnedURL))
	//	for err != nil {
	//		returnedURL, err = utils.ReadLine(fmt.Sprintf("Invalid URL (%v), please enter the URL you were redirected to:\n>", err), os.Stdin)
	//		if err != nil {
	//			return nil, fmt.Errorf("error reading input: %w", err)
	//		}
	//		u, err = url.Parse(returnedURL)
	//	}

	return nil, nil
}
