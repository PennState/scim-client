package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/PennState/golang_scimclient/scim"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	server       *url.URL
	client       *http.Client
	resourceType scim.ResourceType
}

func New(scimServerUrl string, httpClient *http.Client) (Client, error) {
	var scimClient Client
	var err error

	scimClient.server, err = url.Parse(scimServerUrl)
	if err != nil {
		return scimClient, err
	}

	scimClient.client = httpClient

	return scimClient, err
}

func (c Client) RetrieveResource(res scim.Resource, id string) error {
	path := c.server.String() + res.ResourceType().Endpoint + "/" + id
	log.Infof("Path: %s", path)
	resp, err := c.client.Get(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Body: %s", body)

	return scim.Unmarshal(body, res)
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

func (c Client) GetResourceTypes() ([]scim.ResourceType, error) {
	var resourceTypes []scim.ResourceType
	err := c.getServerDiscoveryResources(scim.ResourceTypeResourceType, resourceTypes)
	return resourceTypes, err
}

func (c Client) GetSchemas() ([]scim.Schema, error) {
	var schemas []scim.Schema
	err := c.getServerDiscoveryResources(scim.SchemaResourceType, &schemas)
	return schemas, err
}

func (c Client) GetServerProviderConfig() (scim.ServiceProviderConfig, error) {
	var cfg scim.ServiceProviderConfig
	err := c.getServerDiscoveryResource(&cfg)
	return cfg, err
}

func (c Client) getServerDiscoveryResources(typ scim.ResourceType, res interface{}) error {
	return nil
}

func (c Client) getServerDiscoveryResource(r scim.ServerDiscoveryResource) error {
	log.Infof("Type: %v", reflect.TypeOf(r))
	resp, err := c.client.Get(c.server.String() + r.ResourceType().Endpoint)
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
