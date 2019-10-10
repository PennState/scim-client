package main

import (
	"net/http"

	"github.com/PennState/scim-client/pkg/scim"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//
//OAuth2 configuration
//

const oEnvPrefix = "oauth"

//OAuthConfig holds the oauth parameters for API authentication
type OAuthCfg struct {
	TokenURL     string `split_words:"true" required:"true"`
	ClientID     string `split_words:"true" required:"true"`
	ClientSecret string `split_words:"true" required:"true"`
}

//NewOAuthConfigFromEnv searches the environment to get the required oauth parameters
func NewOAuthConfigFromEnv() (*OAuthCfg, error) {
	var oCfg OAuthCfg
	err := envconfig.Process(oEnvPrefix, &oCfg)
	return &oCfg, err
}

//NewOAuthClient ..
func NewOAuthClient(oCfg *OAuthCfg, url string, opts ...scim.ClientOpt) (*scim.Client, error) {
	http := newOAuthClient(oCfg)
	return scim.NewClient(http, url, opts...)
}

//NewOAuthClientFromEnv ..
func NewOAuthClientFromEnv() (*scim.Client, error) {
	cfg, err := NewOAuthConfigFromEnv()
	if err != nil {
		return nil, err
	}

	http := newOAuthClient(cfg)
	return scim.NewClientFromEnv(http)
}

func newOAuthClient(oCfg *OAuthCfg) *http.Client {
	var ccc = clientcredentials.Config{
		TokenURL:     oCfg.TokenURL,
		ClientID:     oCfg.ClientID,
		ClientSecret: oCfg.ClientSecret,
	}
	return ccc.Client(oauth2.NoContext)

}
