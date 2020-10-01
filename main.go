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

	auth_code   = kingpin.Command("authcode", "Perform an authorization_code grant flow")
	authCodeURL = auth_code.Arg("auth_url", "The URL that will be used for the authorization code generation").Required().String()
	tokenURL    = auth_code.Arg("token_url", "The URL to be used in token generation").Default("").String()

	client_credentials = kingpin.Command("client_credentials", "Perform client_credential grant flow")
	clienCredsURL      = client_credentials.Arg("auth_url", "URL for generating client credential based access token").Required().String()
)

func main() {
	command := kingpin.Parse()
	util := &utils.MainUtil{}

	clientSecretValue, err := validateClientSecret(util)
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	var client clients.OAuthClient
	switch command {
	case "authcode":
		client = clients.NewAuthCode(strings.TrimSpace(*authCodeURL), strings.TrimSpace(*tokenURL), "", nil, util)
	case "client_credentials":
		client = clients.NewClientCredentialClient(strings.TrimSpace(*clienCredsURL), nil)
	}

	token, err := client.GenerateAccessToken(*clientID, clientSecretValue, *scopes, *customParams)
	if err != nil {
		utils.PrintError(err)
	}
	fmt.Printf("\n token is: %#v\n\n", token)
}

func validateClientSecret(util *utils.MainUtil) (string, error) {
	val := *clientSecret
	var err error
	for val == "" {
		val, err = util.ReadLine("Enter client-secret", true)
		if err != nil {
			break
		}
	}
	return val, err
}
