package client

import (
	"io/ioutil"
	"net/http"
	"net/url"

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

func (c Client) ForEndpoint(endPoint string) (Client, error) {
	var rt scim.ResourceType
	rt.Endpoint = endPoint
	return c.ForResourceType(rt)
}

func (c Client) ForResourceType(resourceType scim.ResourceType) (Client, error) {
	return c, nil
}

func (c Client) Get(path string, user *scim.User) error {
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

	return scim.Unmarshal(body, user)
}
