package client

import (
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func GetOauth2AuthenticatedHTTPClient(tokenURL string, clientID string, clientSecret string) *http.Client {
	var ccc clientcredentials.Config
	ccc.TokenURL = tokenURL
	ccc.ClientID = clientID
	ccc.ClientSecret = clientSecret

	return ccc.Client(oauth2.NoContext)
}

func GetOauth2AuthenticatedHTTPClientFromEnvironment() *http.Client {
	return GetOauth2AuthenticatedHTTPClient(
		os.Getenv(`OAUTH_SERVICE_URL`),
		os.Getenv(`OAUTH_CLIENT_ID`),
		os.Getenv(`OAUTH_CLIENT_SECRET`),
	)
}
