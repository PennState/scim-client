package client

import (
	"net/http"
	"net/url"

	"github.com/PennState/golang_scimclient/scim"
)

type Client struct {
	server *url.URL
	client *http.Client
}

func NewClient(svr string, client *http.Client) (*Client, error) {
	_, err := url.Parse(svr)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (c Client) Get(resourceType string, id string) *scim.Resource {
	return nil
}
