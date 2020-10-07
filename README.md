# OAuthCLI
## Features
[![Build Status](https://travis-ci.org/kofoworola/oauthcli.svg?branch=master)](https://travis-ci.org/kofoworola/oauthcli)


`oauthcli` is a command line tool for generating access tokens for diferent oauth grant types via the command line. It allows you to carry out OAuth authentication requests on the following grant_types:
- [Authorization Code](https://oauth.net/2/grant-types/authorization-code/)
- [Client Credentials](https://oauth.net/2/grant-types/client-credentials/)
- [Implicit (Legacy)](https://oauth.net/2/grant-types/implicit/)
- [Password (Legacy)](https://oauth.net/2/grant-types/password/)

## Installation
### From Source
Currently, since it is written in go, you can only install it with
```
go install github.com/kofoworola/oauthcli
```

This installs it to the `bin` folder of your `$GOPATH` env variable, which should be a part of your `$PATH` variable by default. If it isn't then be sure to add it and restart your terminal.

## Usage
### Commands
Run `oauthcli --help` to see a list of commands and their parameters. Make sure to run `oauthcli help <command>` for more info on a specific command

```
usage: oauthcli --client-id=CLIENT-ID --scopes=SCOPES [<flags>] <command> [<args> ...]

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --client-id=CLIENT-ID  The client ID to be used in getting the access_token.
  --client-secret=""     The client Secret to be used in the getting the access_token
  --scopes=SCOPES        The requested scopes
  --extra=EXTRA ...      Extra params to be sent along with the client_id during authentication request

Commands:
  help [<command>...]
    Show help.

  authcode [<flags>] <auth_url> [<token_url>]
    Perform an authorization_code grant flow; https://oauth.net/2/grant-types/authorization-code/

  client_credentials <auth_url>
    Perform client_credential grant flow; https://oauth.net/2/grant-types/client-credentials/

  implicit [<flags>] <auth_url>
    Perform (Legacy) implicit grant type; https://oauth.net/2/grant-types/implicit/

  password --username=USERNAME [<flags>] <auth_url>
    Perform (Legacy) password grant type; https://oauth.net/2/grant-types/password/
```

### Dealing with callbacks
Since grant flows like `auth_code` and `implicit` involves the a redirect to a callback URL, you would need to enter the URL redirected to during the callback to the CLI, at which point, an access_token can be exctracted.

![redirect link](https://imgur.com/a/C1f8kCH)
