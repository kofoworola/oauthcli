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
	clientID     = kingpin.Flag("client-id", "The client ID to be used in getting the access_token.").Required().String()
	clientSecret = kingpin.Flag("client-secret", "The client Secret to be used in the getting the access_token").Default("").String()
	scopes       = kingpin.Flag("scopes", "The requested scopes").Required().String()
	customParams = kingpin.Flag("extra", "Extra params to be sent along with the client_id during authentication request").StringMap()

	authCode                = kingpin.Command("authcode", "Perform an authorization_code grant flow; https://oauth.net/2/grant-types/authorization-code/")
	authCodeURL             = authCode.Arg("auth_url", "The URL that will be used for the authorization code generation").Required().String()
	tokenURL                = authCode.Arg("token_url", "The URL to be used in token generation").Default("").String()
	authCodeRedirectURLFLag = authCode.Flag("redirect_url", "Redirect URL to be passed with the request").Default("").String()

	client_credentials = kingpin.Command("client_credentials", "Perform client_credential grant flow; https://oauth.net/2/grant-types/client-credentials/")
	clienCredsURL      = client_credentials.Arg("auth_url", "URL for generating client credential based access token").Required().String()

	implicitGrant           = kingpin.Command("implicit", "Perform (Legacy) implicit grant type; https://oauth.net/2/grant-types/implicit/")
	implicitURL             = implicitGrant.Arg("auth_url", "URL for generating implicit access_token").Required().String()
	implicitRedirectURLFlag = implicitGrant.Flag("redirect_url", "Redirect URL to be passed with the request").Default("").String()

	passwordGrant = kingpin.Command("password", "Perform (Legacy) password grant type; https://oauth.net/2/grant-types/password/")
	passwordURL   = passwordGrant.Arg("auth_url", "URL for generatining password flow access_token").Required().String()
	usernameFlag  = passwordGrant.Flag("username", "Username for authorization").Required().String()
	passwordFlag  = passwordGrant.Flag("password", "User password").Default("").String()
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
		client = clients.NewAuthCode(strings.TrimSpace(*authCodeURL), strings.TrimSpace(*tokenURL), *authCodeRedirectURLFLag, nil, util)
	case "client_credentials":
		client = clients.NewClientCredentialClient(strings.TrimSpace(*clienCredsURL), nil)
	case "implicit":
		client = clients.NewImplicitClient(*implicitURL, *implicitRedirectURLFlag, util)
	case "password":
		client = clients.NewPasswordClient(*passwordURL, *usernameFlag, *passwordFlag, nil, util)
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
