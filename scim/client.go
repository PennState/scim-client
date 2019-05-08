package scim

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"text/tabwriter"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//
//SCIM client configuration
//

const sEnvPrefix = "scim"

//ClientConfig ..
type ClientConfig struct {
	ServiceURL       string `split_words:"true" required:"true"` //ServiceURL is the base URI of the SCIM server's resources - see https://tools.ietf.org/html/rfc7644#section-1.3
	IgnoreRedirects  bool   `split_words:"true" default:"false"`
	DisableDiscovery bool   `split_words:"true" default:"false"`
	DisableEtag      bool   `split_words:"true" default:"false"`
}

//NewDefaultClientConfig ..
func NewDefaultClientConfig(serviceURL string) *ClientConfig {
	var sCfg ClientConfig
	sCfg.ServiceURL = serviceURL
	return &sCfg
}

//NewClientConfigFromEnv ..
func NewClientConfigFromEnv() (*ClientConfig, error) {
	var sCfg ClientConfig
	err := envconfig.Process(sEnvPrefix, &sCfg)
	return &sCfg, err
}

func clientConfigUsage(sCfg *ClientConfig) error {
	return envconfig.Usage(sEnvPrefix, sCfg)
}

var sEnvSpec = envSpec{
	prefix: sEnvPrefix,
	spec:   new(ClientConfig),
}

//
//OAuth2 configuration
//

const oEnvPrefix = "oauth"

//OAuthConfig holds the oauth parameters for API authentication
type OAuthConfig struct {
	TokenURL     string `split_words:"true" required:"true"`
	ClientID     string `split_words:"true" required:"true"`
	ClientSecret string `split_words:"true" required:"true"`
}

//NewOAuthConfigFromEnv searches the environment to get the required oauth parameters
func NewOAuthConfigFromEnv() (*OAuthConfig, error) {
	var oCfg OAuthConfig
	err := envconfig.Process(oEnvPrefix, &oCfg)
	return &oCfg, err
}

func oauthConfigUsage(oCfg *OAuthConfig) error {
	return envconfig.Usage(oEnvPrefix, oCfg)
}

var oEnvSpec = envSpec{
	prefix: oEnvPrefix,
	spec:   new(OAuthConfig),
}

//
//SCIM client
//

type client struct {
	sCfg    *ClientConfig
	hClient *http.Client
}

//Client allows request scim resources
type Client struct {
	client
}

//
//SCIM client constructors
//

//NewClient ..
func NewClient(sCfg *ClientConfig, hClient *http.Client) (*Client, error) {
	var sClient Client

	//Validate that the URL is at least formatted correctly
	_, err := url.Parse(sCfg.ServiceURL)
	if err != nil {
		return &sClient, err
	}

	sClient.hClient = hClient
	sClient.sCfg = sCfg
	return &sClient, err
}

//NewClientFromEnv ..
func NewClientFromEnv(hClient *http.Client) (*Client, error) {
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

//NewDefaultClient ..
func NewDefaultClient(serviceURL string, hClient *http.Client) (*Client, error) {
	var sCfg ClientConfig
	sCfg.ServiceURL = serviceURL
	return NewClient(&sCfg, hClient)
}

//NewDefaultClientFromEnv ..
func NewDefaultClientFromEnv(httpClient *http.Client) {

}

//NewOAuthClient ..
func NewOAuthClient(sCfg *ClientConfig, oCfg *OAuthConfig) (*Client, error) {
	var ccc = clientcredentials.Config{
		TokenURL:     oCfg.TokenURL,
		ClientID:     oCfg.ClientID,
		ClientSecret: oCfg.ClientSecret,
	}
	hClient := ccc.Client(oauth2.NoContext)

	return NewClient(sCfg, hClient)
}

//NewOAuthClientFromEnv ..
func NewOAuthClientFromEnv() (*Client, error) {
	sCfg, err1 := NewClientConfigFromEnv()
	if err1 != nil {
		log.Error(err1)
	}

	oCfg, err2 := NewOAuthConfigFromEnv()
	if err2 != nil {
		log.Error(err2)
	}

	if err1 != nil || err2 != nil {
		usage(sEnvSpec, oEnvSpec)
		return nil, errors.New("Failed to configure the SCIM client - see preceding messages for cause")
	}

	return NewOAuthClient(sCfg, oCfg)
}

//
//Resource accessor/mutator methods
//

//RetrieveResource ..
func (c Client) RetrieveResource(res Resource, id string) error {
	path := c.sCfg.ServiceURL + res.ResourceType().Endpoint + "/" + id

	log.Debugf("Path: %s", path)
	resp, err := c.hClient.Get(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Body: %s", body)

	return Unmarshal(body, res)
}

//SearchResource ..
func (c Client) SearchResource(rt ResourceType, sr SearchRequest) (ListResponse, error) {
	path := c.sCfg.ServiceURL + rt.Endpoint + "/.search"
	return c.search(path, sr)
}

//SearchServer ..
func (c Client) SearchServer(sr SearchRequest) (ListResponse, error) {
	path := c.sCfg.ServiceURL + "/.search"
	return c.search(path, sr)
}

func (c Client) search(path string, sr SearchRequest) (ListResponse, error) {
	log.Debug("Path: ", path)
	var lr ListResponse

	// TODO: Remove this after SCIMple is fixed
	if sr.SortOrder == NotSpecified {
		sr.SortOrder = Ascending
	}

	//TODO if you ask for 0 records, you'll get what you ask for
	if sr.Count == 0 {
		sr.Count = 1000
	}

	srj, err := json.Marshal(sr)
	if err != nil {
		return lr, err
	}
	log.Debug("SearchRequest JSON: ", string(srj))

	resp, err := c.hClient.Post(path, "application/json", bytes.NewReader(srj))
	if err != nil {
		return lr, err
	}
	log.Debug("Search HTTP status code: ", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lr, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.Unmarshal(body, &er)
		if err != nil {
			log.Error("Couldn't unmarshal ErrorResponse") // TODO: Fix SCIMple to return the correct format
			return lr, err
		}
		log.Error("ErrorResponse: ", er)
		return lr, errors.New(resp.Status)
	}

	log.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &lr)
	if err != nil {
		return lr, err
	}
	return lr, nil
}

//func CreateResource(res *Resource)
//func ReplaceResource(res *Resource)
//func ModifyResource(res *Resource)
//func DeleteResource(res *Resource) error
//func Bulk

//
//Server Discovery
//

func (c Client) GetResourceTypes() ([]ResourceType, error) {
	var resourceTypes []ResourceType
	err := c.getServerDiscoveryResources(ResourceTypeResourceType, resourceTypes)
	return resourceTypes, err
}

func (c Client) GetSchemas() ([]Schema, error) {
	var schemas []Schema
	err := c.getServerDiscoveryResources(SchemaResourceType, &schemas)
	return schemas, err
}

func (c Client) GetServerProviderConfig() (ServiceProviderConfig, error) {
	var cfg ServiceProviderConfig
	err := c.getServerDiscoveryResource(&cfg)
	return cfg, err
}

func (c Client) getServerDiscoveryResources(typ ResourceType, res interface{}) error {
	return nil
}

func (c Client) getServerDiscoveryResource(r ServerDiscoveryResource) error {
	log.Debugf("Type: %v", reflect.TypeOf(r))
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
// Convenience methods

//RetrieveResourceByUserName is a helper method for retrieving resources by UserName
func (c Client) SearchUserResourcesByUserName(userName string) (ListResponse, error) {
	sr := SearchRequest{
		Filter: "userName EQ \"" + userName + "\"",
	}
	return c.SearchResource(UserResourceType, sr)
}

//RetrieveResourceByExternalID is a helper method for retrieving resources by ExternalId
func (c Client) SearchResourcesByExternalID(rt ResourceType, externalID string) (ListResponse, error) {
	sr := SearchRequest{
		Filter: "externalId EQ \"" + externalID + "\"",
	}
	return c.SearchResource(rt, sr)
}

//
//General HTTP client code
//

func getStringEntityBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("Body: %s", body)
	return body, err
}

//
//Error reporting for envconfig
//

const (
	multipleTableInstructionHeader = `The method call you've chosen causes this this library to be configured via the
environment. The following environment variables can be (or in the case of required
parameters must be) used:
`

	multipleTableFormatHeader = `
KEY	TYPE	DEFAULT	REQUIRED	DESCRIPTION
---	----	-------	--------	-----------
`

	multipleTableFormatTemplate = `{{range .}}{{usage_key .}}	{{usage_type .}}	{{usage_default .}}	{{usage_required .}}	{{usage_description .}}
{{end}}`
)

type envSpec struct {
	prefix string
	spec   interface{}
}

type empty struct {
}

func usage(specs ...envSpec) {
	buf := new(bytes.Buffer)
	buf.WriteString(multipleTableInstructionHeader)
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)
	envconfig.Usagef("", new(empty), tabs, multipleTableFormatHeader)
	for _, spec := range specs {
		envconfig.Usagef(spec.prefix, spec.spec, tabs, multipleTableFormatTemplate)
	}
	tabs.Flush()
	log.Info(buf)
}
