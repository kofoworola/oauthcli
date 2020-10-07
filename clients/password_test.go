package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestPasswordClient_GenerateAccessToken(t *testing.T) {
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
			pass := strconv.Itoa(int(time.Now().Unix()))
			c := &http.Client{
				Transport: roundtripFunction(func(req *http.Request) *http.Response {
					req.ParseForm()
					if grant := req.PostForm.Get("grant_type"); grant != "password" {
						t.Fatalf("grant_type should be password instead got %s", grant)
					}
					switch "" {
					case req.PostFormValue("scope"):
						t.Fatalf("scope should not be empty")
					case req.PostFormValue("client_id"):
						t.Fatal("client_id should not be empty")
					case req.PostFormValue("client_secret"):
						t.Fatal("client_secret should not be empty")
					case req.PostFormValue("password"):
						t.Fatalf("password should not be empty")
					}

					if p := req.PostFormValue("password"); p != pass {
						t.Fatalf("expected password %s got %s", pass, p)
					}
					return &http.Response{
						StatusCode: test.respCode,
						Body:       ioutil.NopCloser(strings.NewReader(test.resp)),
					}
				}),
			}
			u := &TestUtil{
				password: pass,
			}
			client := NewPasswordClient("http://dummyurl.com", "username", "", c, u)
			token, err := client.GenerateAccessToken("client_id", "client_secret", "create delete", nil)
			if err := test.eval(token, err); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
