package scim

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//
//SCIM client configuration
//

const sEnvPrefix = "scim"

type ClientConfig struct {
	ServiceURL       string `split_words:"true" required:"true"` //ServiceURL is the base URI of the SCIM server's resources - see https://tools.ietf.org/html/rfc7644#section-1.3
	IgnoreRedirects  bool   `split_words:"true" default:"false"`
	DisableDiscovery bool   `split_words:"true" default:"false"`
	DisableEtag      bool   `split_words:"true" default:"false"`
}

func NewDefaultClientConfig(serviceUrl string) *ClientConfig {
	var sCfg ClientConfig
	sCfg.ServiceURL = serviceUrl
	return &sCfg
}

func NewClientConfigFromEnv() (*ClientConfig, error) {
	var sCfg ClientConfig
	err := envconfig.Process(sEnvPrefix, &sCfg)
	return &sCfg, err
}

func clientConfigUsage(sCfg *ClientConfig) error {
	return envconfig.Usage(sEnvPrefix, sCfg)
}

//
//OAuth2 configuration
//

const oEnvPrefix = "oauth"

type OAuthConfig struct {
	ServiceURL   string `split_words:"true" required:"true"`
	ClientID     string `split_words:"true" required:"true"`
	ClientSecret string `split_words:"true" required:"true"`
}

func NewOAuthConfigFromEnv() (*OAuthConfig, error) {
	var oCfg OAuthConfig
	err := envconfig.Process(oEnvPrefix, &oCfg)
	return &oCfg, err
}

func oauthConfigUsage(oCfg *OAuthConfig) error {
	return envconfig.Usage(oEnvPrefix, oCfg)
}

//
//SCIM client
//

type client struct {
	sCfg    *ClientConfig
	hClient *http.Client
}

//
//SCIM client constructors
//

func NewClient(sCfg *ClientConfig, hClient *http.Client) (*client, error) {
	var sClient client

	//Validate that the URL is at least formatted correctly
	_, err := url.Parse(sCfg.ServiceURL)
	if err != nil {
		return &sClient, err
	}

	sClient.hClient = hClient
	sClient.sCfg = sCfg
	return &sClient, err
}

func NewClientFromEnv(hClient *http.Client) (*client, error) {
	sCfg, err := NewClientConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewClient(sCfg, hClient)
}

// func NewClient(cfg *ClientConfig, http *http.Client)
// func NewUnauthenticatedClient()
// func NewUnauthenticatedClientFromEnv()
// func NewBasicAuthClient()
// func NewBasicAuthClientFromEnv()

func NewDefaultClient(serviceUrl string, hClient *http.Client) (*client, error) {
	var sCfg ClientConfig
	sCfg.ServiceURL = serviceUrl
	return NewClient(&sCfg, hClient)
}

func NewDefaultClientFromEnv(httpClient *http.Client) {

}

func NewOAuthClient(sCfg *ClientConfig, oCfg *OAuthConfig) (*client, error) {
	var ccc = clientcredentials.Config{
		TokenURL:     oCfg.ServiceURL,
		ClientID:     oCfg.ClientID,
		ClientSecret: oCfg.ClientSecret,
	}
	hClient := ccc.Client(oauth2.NoContext)

	return NewClient(sCfg, hClient)
}

func NewOAuthClientFromEnv() (*client, error) {
	sCfg, err1 := NewClientConfigFromEnv()
	oCfg, err2 := NewOAuthConfigFromEnv()
	if err1 != nil || err2 != nil {
		log.Error(clientConfigUsage(sCfg), err1)
		log.Error(oauthConfigUsage(oCfg), err2)
		return nil, err1
	}

	return NewOAuthClient(sCfg, oCfg)
}

//
//Resource accessor/mutator methods
//

func (c client) RetrieveResource(res Resource, id string) error {
	path := c.sCfg.ServiceURL + res.ResourceType().Endpoint + "/" + id
	log.Infof("Path: %s", path)
	resp, err := c.hClient.Get(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Body: %s", body)

	return Unmarshal(body, res)
}

//func CreateResource(res *Resource)
//func RetrieveResource(res *Resource, id string)
//func QueryResource(rt ResourceType, sr SearchRequest)
//func ReplaceResource(res *Resource)
//func ModifyResource(res *Resource)
//func DeleteResource(res *Resource) error
//func Search(lr *ListResponse, sr SearchRequest)
//func Bulk

//
//Server Discovery
//

func (c client) GetResourceTypes() ([]ResourceType, error) {
	var resourceTypes []ResourceType
	err := c.getServerDiscoveryResources(ResourceTypeResourceType, resourceTypes)
	return resourceTypes, err
}

func (c client) GetSchemas() ([]Schema, error) {
	var schemas []Schema
	err := c.getServerDiscoveryResources(SchemaResourceType, &schemas)
	return schemas, err
}

func (c client) GetServerProviderConfig() (ServiceProviderConfig, error) {
	var cfg ServiceProviderConfig
	err := c.getServerDiscoveryResource(&cfg)
	return cfg, err
}

func (c client) getServerDiscoveryResources(typ ResourceType, res interface{}) error {
	return nil
}

func (c client) getServerDiscoveryResource(r ServerDiscoveryResource) error {
	log.Infof("Type: %v", reflect.TypeOf(r))
	resp, err := c.hClient.Get(c.sCfg.ServiceURL + r.ResourceType().Endpoint)
	if err != nil {
		return err
	}

	body, err := getStringEntityBody(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, r)
}

//
//General HTTP client code
//

func getStringEntityBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Infof("Body: %s", body)
	return body, err
}
