package scim

import (
	"bytes"
	"context"
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

// clientConfig ..
// ServiceURL is the base URI of the SCIM server's resources - see https://tools.ietf.org/html/rfc7644#section-1.3
type clientCfg struct {
	ServiceURL       string `split_words:"true" required:"true"`
	IgnoreRedirects  bool   `split_words:"true" default:"false"`
	DisableDiscovery bool   `split_words:"true" default:"false"`
	DisableEtag      bool   `split_words:"true" default:"false"`
}

//
//SCIM client options
//

type ClientOpt func(*clientCfg)

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
		return nil, errors.New("provided service URL is not valid")
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
func (c Client) RetrieveResource(ctx context.Context, res Resource, id string) error {
	path := c.cfg.ServiceURL + res.ResourceType().Endpoint + "/" + id

	log.Debugf("Path: %s", path)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
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
			return err
		}
		return er
	}

	return json.Unmarshal(body, res)
}

//SearchResource ..
func (c Client) QueryResourceType(ctx context.Context, rt ResourceType, sr SearchRequest) (ListResponse, error) {
	path := c.cfg.ServiceURL + rt.Endpoint + "/.search"
	return c.query(ctx, path, sr)
}

//SearchServer ..
func (c Client) QueryServer(ctx context.Context, sr SearchRequest) (ListResponse, error) {
	path := c.cfg.ServiceURL + "/.search"
	return c.query(ctx, path, sr)
}

func (c Client) query(ctx context.Context, path string, sr SearchRequest) (ListResponse, error) {
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

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(srj))
	if err != nil {
		return lr, err
	}
	req.Header.Set("Content-Type", "application/scim+json")
	resp, err := c.http.Do(req)
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
			return lr, err
		}
		return lr, er
	}

	log.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &lr)
	if err != nil {
		return lr, err
	}
	return lr, nil
}

//CreateResource ..
func (c Client) CreateResource(ctx context.Context, res Resource) error {
	log.Trace("(c Client) ReplaceResource(res)")
	rj, err := json.Marshal(res)
	if err != nil {
		return err
	}
	log.Debug("Marshaled resource: ", string(rj))

	path := c.cfg.ServiceURL + res.ResourceType().Endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(rj))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/scim+json")

	return c.resourceOrError(res, req)
}

//ReplaceResource ..
func (c Client) ReplaceResource(ctx context.Context, res Resource) error {
	log.Trace("(c Client) ReplaceResource(res)")
	rj, err := json.Marshal(res)
	if err != nil {
		return err
	}
	log.Debug("Marshaled resource: ", string(rj))

	path := c.cfg.ServiceURL + res.ResourceType().Endpoint + "/" + res.getID()
	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(rj))
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
	log.Debug("Replace HTTP status code: ", resp.StatusCode)

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
			return err
		}
		return errors.New(resp.Status)
	}

	return json.Unmarshal(body, res)
}

//func ModifyResource(res *Resource)
//func DeleteResource(res *Resource) error
//func Bulk

//
//Server Discovery
//

func (c Client) GetResourceTypes(ctx context.Context) ([]ResourceType, error) {
	var resourceTypes []ResourceType
	err := c.getServerDiscoveryResources(ctx, ResourceTypeResourceType, resourceTypes)
	return resourceTypes, err
}

func (c Client) GetSchemas(ctx context.Context) ([]Schema, error) {
	var schemas []Schema
	err := c.getServerDiscoveryResources(ctx, SchemaResourceType, &schemas)
	return schemas, err
}

func (c Client) GetServerProviderConfig(ctx context.Context) (ServiceProviderConfig, error) {
	var cfg ServiceProviderConfig
	err := c.getServerDiscoveryResource(ctx, &cfg)
	return cfg, err
}

func (c Client) getServerDiscoveryResources(ctx context.Context, typ ResourceType, res interface{}) error {
	return nil
}

func (c Client) getServerDiscoveryResource(ctx context.Context, r ServerDiscoveryResource) error {
	log.Debugf("Type: %v", reflect.TypeOf(r))
	path := c.cfg.ServiceURL + r.ResourceType().Endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/scim+json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	body, err := c.body(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, r)
}

//
//General HTTP client code
//

func (c Client) body(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("Body: %s", body)
	return body, err
}

func (c Client) error(resp *http.Response) error {
	body, err := c.body(resp)
	if err != nil {
		return err
	}

	var er ErrorResponse
	err = json.Unmarshal(body, &er)
	if err != nil {
		return err
	}
	return er
}

func (c Client) etag(res Resource, req *http.Request) {
	if !c.cfg.DisableEtag {
		req.Header.Set("If-Match", res.getMeta().Version)
	}
}

func (c Client) mime(req *http.Request) {
	req.Header.Set("Accept", "application/scim+json")
	req.Header.Set("Content-Type", "application/scim+json")
}

func (c Client) resourceOrError(res Resource, req *http.Request) error {
	c.mime(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return c.error(resp)
	}
	return c.resource(resp, res)
}

func (c Client) resource(resp *http.Response, res Resource) error {
	body, err := c.body(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, res)
}
