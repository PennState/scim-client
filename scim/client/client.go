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

func (c Client) ForResourceType(resourceType scim.ResourceType) (Client, error) {
	return c, nil
}

func (c Client) Get(path string, resource scim.Resource) error {
	resp, err := c.client.Get(c.server.String() + path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Body: %s", body)

	return scim.Unmarshal(body, resource)
}

func (c Client) GetResourceTypes() ([]scim.ResourceType, error) {
	var resourceTypes []scim.ResourceType

	resp, err := c.client.Get(c.server.String() + "/ResourceTypes")
	if err != nil {
		return resourceTypes, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Infof("Body: %s", body)

	err = json.Unmarshal(body, &resourceTypes)
	return resourceTypes, nil
}

func (c Client) GetSchemas() ([]scim.Schema, error) {
	var schemas []scim.Schema

	resp, err := c.client.Get(c.server.String() + "/Schemas")
	if err != nil {
		return schemas, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return schemas, err
	}
	log.Infof("Body: %s", body)

	//return schemas, scim.Unmarshal(body, &schemas)
	return schemas, nil
}

func (c Client) GetServerProviderConfig() (scim.ServiceProviderConfig, error) {
	var cfg scim.ServiceProviderConfig
	err := c.getServiceDiscoveryEndpoint(&cfg)
	return cfg, err

}

func (c Client) getServiceDiscoveryEndpoint(r scim.ServiceDiscoveryResource) error {
	log.Infof("Type: %v", reflect.TypeOf(r))
	resp, err := c.client.Get(c.server.String() + r.ServiceDiscoveryResourceType().Endpoint())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Body: %s", body)

	return json.Unmarshal(body, r)
}

// func (c Client) Get()
