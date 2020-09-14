package main

import (
	"fmt"
	"strings"

	"github.com/kofoworola/oauthcli/clients"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	clientID     = kingpin.Flag("client_id", "The client ID to be used in the oAuth token request.").Required().String()
	clientSecret = kingpin.Flag("client_secret", "The client Secret to be used in the oAuth token request").Required().String()
	scopes       = kingpin.Flag("scopes", "The requested scopes").Required().String()

	auth_code = kingpin.Command("authorization_code", "Perform an authorization_code grant flow")
	authURL   = auth_code.Arg("auth_url", "The URL that will be used for the authorization code generation").Required().String()
	tokenURL  = auth_code.Arg("token_url", "The URL to be used in token generation").Required().String()
)

func main() {
	var client clients.OAuthClient
	switch kingpin.Parse() {
	case "authorization_code":
		client = clients.NewAuthCodeClient(strings.TrimSpace(*authURL), strings.TrimSpace(*tokenURL), "", nil)
	}

	token, _ := client.GenerateAccessToken(*clientID, *clientSecret, *scopes)
	fmt.Printf("\n token is: %s\n\n", token)
}
