package clients

import (
	"fmt"
	"testing"
)

func TestImplicitClient_GenerateAccessToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description   string
		redirectedURL string
		expectedToken string
		expectedError string
	}{
		{
			description:   "TestSuccessful",
			redirectedURL: "https://callback.com.com/callback#access_token=test_token&token_type=Bearer&expires_in=600&state=dummy_state",
			expectedToken: "test_token",
			expectedError: "",
		},
		{
			description:   "TestSuccessFullNoExtra",
			redirectedURL: "https://callback.com.com/callback#access_token=test_token",
			expectedToken: "test_token",
			expectedError: "",
		},
		{
			description:   "TestError",
			redirectedURL: "https://callback.com.com/callback#err=access_denied",
			expectedToken: "",
			expectedError: "access_denied",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test := test
			t.Parallel()
			util := &TestUtil{
				openURL: `https://authserver.com/auth?response_type=token&client_id=client_id&redirect_uri=https://callback.com/callback&scope=create&state=rand`,
				input:   fmt.Sprintf("%s\n", test.redirectedURL),
			}

			client := NewImplicitClient("https://authserver.com/auth", "https://callback.com.comi", util)
			token, err := client.GenerateAccessToken("client_id", "client_secret", "scope", nil)
			if err != nil {
				if err.Error() != test.expectedError {
					t.Fatalf("expected error '%s' got '%s'", test.expectedError, err.Error())
				}
				// end test if err != nil and error was expected
				return
			}
			if token.AccessToken != test.expectedToken {
				t.Fatalf("expected token '%s', got '%s'", test.expectedToken, token.AccessToken)
			}
		})
	}
}
