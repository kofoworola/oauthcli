package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type roundtripFunction func(*http.Request) *http.Response

func (r roundtripFunction) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

var successToken = `{
  "access_token":"dummy_token",
  "token_type":"bearer",
  "expires_in":3600,
  "refresh_token":"dummy_refresh_token",
  "scope":"create"
}`

var errResp = `
{
  "error": "invalid_client",
  "error_description": "Invalid client_id passed",
  "error_uri": "See the full API docs at https://authorization-server.com/docs/access_token"
}
`

func TestClientCredentials_GenerateAccessToken(t *testing.T) {
	t.Parallel()
	tests := []*struct {
		description string
		respCode    int
		resp        string
		eval        func(*Token, error) error
	}{
		{
			description: "TestSuccess",
			respCode:    http.StatusOK,
			resp:        successToken,
			eval: func(t *Token, err error) error {
				if t.AccessToken != "dummy_token" {
					return fmt.Errorf("wanted dummy_token as token, got %s", t.AccessToken)
				}
				return nil
			},
		},
		{
			description: "TestInvalidRequest",
			respCode:    http.StatusUnauthorized,
			resp:        errResp,
			eval: func(t *Token, err error) error {
				expectedError := "invalid_client - Invalid client_id passed"
				if err == nil || err.Error() != expectedError {
					return fmt.Errorf("invalid error response: %v", err)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			c := &http.Client{
				Transport: roundtripFunction(func(req *http.Request) *http.Response {
					req.ParseForm()
					if grant := req.PostForm.Get("grant_type"); grant != "client_credentials" {
						t.Fatalf("grant_type should be client_credentials instead got %s", grant)
					}
					switch "" {
					case req.PostFormValue("scope"):
						t.Fatalf("scope should not be empty")
					case req.PostFormValue("client_id"):
						t.Fatal("client_id should not be empty")
					case req.PostFormValue("client_secret"):
						t.Fatal("client_secret should not be empty")
					}
					return &http.Response{
						StatusCode: test.respCode,
						Body:       ioutil.NopCloser(strings.NewReader(test.resp)),
					}
				}),
			}
			client := NewClientCredentialClient("http://dummyurl.com", c)
			token, err := client.GenerateAccessToken("client_id", "client_secret", "create delete", nil)
			if err := test.eval(token, err); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
