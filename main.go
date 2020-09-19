package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kofoworola/oauthcli/clients"
	"github.com/kofoworola/oauthcli/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	clientID     = kingpin.Flag("client-id", "The client ID to be used in the oAuth token request.").Required().String()
	clientSecret = kingpin.Flag("client-secret", "The client Secret to be used in the oAuth token request").Default("").String()
	scopes       = kingpin.Flag("scopes", "The requested scopes").Required().String()
	customParams = kingpin.Flag("extra", "Extra params to be sent along with the client_id during authentication request").StringMap()

	auth_code = kingpin.Command("authcode", "Perform an authorization_code grant flow")
	authURL   = auth_code.Arg("auth_url", "The URL that will be used for the authorization code generation").Required().String()
	tokenURL  = auth_code.Arg("token_url", "The URL to be used in token generation").Default("").String()
)

func main() {
	command := kingpin.Parse()

	clientSecretValue, err := validateClientSecret()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	var client clients.OAuthClient
	switch command {
	case "authcode":
		client = clients.NewAuthCodeClient(strings.TrimSpace(*authURL), strings.TrimSpace(*tokenURL), "", nil)
	}

	token, err := client.GenerateAccessToken(*clientID, clientSecretValue, *scopes, *customParams)
	if err != nil {
		utils.PrintError(err)
	}
	fmt.Printf("\n token is: %#v\n\n", token)
}

func validateClientSecret() (string, error) {
	val := *clientSecret
	var err error
	for val == "" {
		val, err = utils.ReadPassType("Enter client-secret")
		if err != nil {
			break
		}
	}
	return val, err
}
