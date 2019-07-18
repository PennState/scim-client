package scim

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

//
//SCIM client configuration
//

const envPrefix = "scim"

//ClientConfig ..
type clientCfg struct {
	ServiceURL       string `split_words:"true" required:"true"` //ServiceURL is the base URI of the SCIM server's resources - see https://tools.ietf.org/html/rfc7644#section-1.3
	IgnoreRedirects  bool   `split_words:"true" default:"false"`
	DisableDiscovery bool   `split_words:"true" default:"false"`
	DisableEtag      bool   `split_words:"true" default:"false"`
}

//
//SCIM client options
//

type ClientOpt func(*clientCfg)

// func ServiceUrl(serviceUrl string) ClientOpt {
// 	return func(cfg *clientCfg) {
// 		cfg.ServiceURL = serviceUrl
// 	}
// }

func IgnoreRedirects(ignoreRedirects bool) ClientOpt {
	return func(cfg *clientCfg) {
		cfg.IgnoreRedirects = ignoreRedirects
	}
}

func DisableDiscovery(disableDiscovery bool) ClientOpt {
	return func(cfg *clientCfg) {
		cfg.DisableDiscovery = disableDiscovery
	}
}

func DisableEtag(disableEtag bool) ClientOpt {
	return func(cfg *clientCfg) {
		cfg.DisableEtag = disableEtag
	}
}

//
//SCIM client
//

type client struct {
	cfg  *clientCfg
	http *http.Client
}

//Client allows request scim resources
type Client struct {
	*client
}

//
//SCIM client constructors
//

func NewClient(http *http.Client, url string, opts ...ClientOpt) (*Client, error) {
	var cfg clientCfg
	cfg.ServiceURL = url
	for _, opt := range opts {
		opt(&cfg)
	}
	return newClient(http, &cfg)
}

func NewClientFromEnv(http *http.Client) (*Client, error) {
	var cfg clientCfg
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return nil, err
	}
	return newClient(http, &cfg)
}

func newClient(http *http.Client, cfg *clientCfg) (*Client, error) {

	//Validate that the URL exists and is formatted correctly
	if cfg.ServiceURL == "" {
		return nil, errors.New("ServiceURL is a required configuration parameter")
	}
	_, err := url.Parse(cfg.ServiceURL)
	if err != nil {
		return nil, errors.New("Provided service URL is not valid")
	}

	return &Client{
		client: &client{
			http: http,
			cfg:  cfg,
		},
	}, nil
}

//
//Resource accessor/mutator methods
//

//RetrieveResource ..
func (c Client) RetrieveResource(res Resource, id string) error {
	path := c.cfg.ServiceURL + res.ResourceType().Endpoint + "/" + id

	log.Debugf("Path: %s", path)
	resp, err := c.http.Get(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Body: %s", body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.Unmarshal(body, &er)
		if err != nil {
			log.Error("Couldn't unmarshal ErrorResponse") // TODO: Fix SCIMple to return the correct format
			return err
		}
		return er
	}

	return Unmarshal(body, res)
}

//SearchResource ..
func (c Client) SearchResource(rt ResourceType, sr SearchRequest) (ListResponse, error) {
	path := c.cfg.ServiceURL + rt.Endpoint + "/.search"
	return c.search(path, sr)
}

//SearchServer ..
func (c Client) SearchServer(sr SearchRequest) (ListResponse, error) {
	path := c.cfg.ServiceURL + "/.search"
	return c.search(path, sr)
}

func (c Client) search(path string, sr SearchRequest) (ListResponse, error) {
	log.Debug("Path: ", path)
	var lr ListResponse

	// TODO: Remove this after SCIMple is fixed
	if sr.SortOrder == NotSpecified {
		sr.SortOrder = Ascending
	}

	srj, err := json.Marshal(sr)
	if err != nil {
		return lr, err
	}
	log.Debug("SearchRequest JSON: ", string(srj))

	resp, err := c.http.Post(path, "application/json", bytes.NewReader(srj))
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

//CreateResource ..
func (c Client) CreateResource(res Resource) error {
	log.Info("(c Client) ReplaceResource(res)")
	rj, err := json.Marshal(res)
	if err != nil {
		return err
	}
	log.Info("Marshaled resource: ", string(rj))

	path := c.cfg.ServiceURL + res.ResourceType().Endpoint
	resp, err := c.http.Post(path, "application/scim+json", bytes.NewReader(rj))
	if err != nil {
		return err
	}
	log.Info(resp)
	return nil
}

//ReplaceResource ..
func (c Client) ReplaceResource(res Resource) error {
	log.Info("(c Client) ReplaceResource(res)")
	rj, err := json.Marshal(res)
	if err != nil {
		return err
	}
	log.Info("Marshaled resource: ", string(rj))

	path := c.cfg.ServiceURL + res.ResourceType().Endpoint + "/" + res.getID()
	req, err := http.NewRequest("PUT", path, bytes.NewReader(rj))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/scim+json")
	if !c.cfg.DisableEtag {
		req.Header.Set("If-Match", res.getMeta().Version)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	log.Info("Replace HTTP status code: ", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Body: %s", body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.Unmarshal(body, &er)
		if err != nil {
			log.Error("Couldn't unmarshal ErrorResponse") // TODO: Fix SCIMple to return the correct format
			return err
		}
		log.Error("ErrorResponse: ", er)
		return errors.New(resp.Status)
	}

	return Unmarshal(body, res)
}

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
	resp, err := c.http.Get(c.cfg.ServiceURL + r.ResourceType().Endpoint)
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
//Convenience methods
//

//SearchUserResourcesByUserName is a helper method for retrieving resources by UserName
func (c Client) SearchUserResourcesByUserName(userName string) (ListResponse, error) {
	sr := SearchRequest{
		Filter: "userName EQ \"" + userName + "\"",
	}
	return c.SearchResource(UserResourceType, sr)
}

//SearchResourcesByExternalID is a helper method for retrieving resources by ExternalId
func (c Client) SearchResourcesByExternalID(rt ResourceType, externalID string) (ListResponse, error) {
	sr := SearchRequest{
		Filter: "externalId EQ \"" + externalID + "\"",
	}
	return c.SearchResource(rt, sr)
}

//
//SCIM client code
//

//type ScimError

//
//General HTTP client code
//

func getStringEntityBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("Body: %s", body)
	return body, err
}
