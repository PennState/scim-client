package client

import "net/url"

type Authenticator interface {
	getAuthenticationHeader()
}

type OAuth2Authenticator struct {
	Server       url.URL
	ClientID     string
	ClientSecret string
}

func NewOAuth2Authenticator(server string, clientID string, clientSecret string) (*OAuth2Authenticator, error) {
	_, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	var oauth2Authenticator OAuth2Authenticator
	//oauth2Authenticator.Server
	oauth2Authenticator.ClientID = clientID
	oauth2Authenticator.ClientSecret = clientSecret

	return &oauth2Authenticator, nil
}
