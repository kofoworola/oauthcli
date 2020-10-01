package clients

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var (
	authURL = `https://authserver.com/auth?response_type=code&client_id=client_id&redirect_uri=https://redirect_uri.com/callback&scope=create&state=rand`
)

func TestAuthCode_GenerateAccessToken(t *testing.T) {
	tests := []struct {
		description   string
		authURL       string
		responseCode  int
		response      string
		stdIN         string
		expectedError string
	}{
		{
			description:   "TestSuccess",
			authURL:       authURL,
			responseCode:  http.StatusOK,
			response:      successToken,
			stdIN:         "\nhttps://callback.com/callback?code=test_auth_code&state=state\n",
			expectedError: "",
		},
		{
			description:   "CodeNotPresent",
			authURL:       authURL,
			response:      "",
			stdIN:         "\nhttps://callback.com/callback?state=state\n",
			expectedError: "URL does not contain the `code` parameter",
		},
		{
			description:   "TestInvalidRequest",
			authURL:       authURL,
			responseCode:  http.StatusUnauthorized,
			response:      errResp,
			stdIN:         "\nhttps://callback.com/callback?code=test_auth_code&state=state\n",
			expectedError: "invalid_client - Invalid client_id passed",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test := test
			t.Parallel()
			client := &http.Client{
				Transport: roundtripFunction(func(*http.Request) *http.Response {
					return &http.Response{
						StatusCode: test.responseCode,
						Body:       ioutil.NopCloser(strings.NewReader(test.response)),
					}
				}),
			}
			util := &TestUtil{
				input:   test.stdIN,
				openURL: test.authURL,
			}

			c := NewAuthCode(
				"https://authserver.com/auth",
				"",
				"https://redirect_uri.com/callback",
				client,
				util,
			)

			tok, err := c.GenerateAccessToken("client_id", "client_secret", "create", nil)
			if err != nil && err.Error() != test.expectedError {
				t.Fatalf("expected error '%s' got error: %v", test.expectedError, err)
			}
			t.Logf("response gotten is %#v \n", tok)
		})
	}
}
