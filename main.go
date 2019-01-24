package main

//Build/run this program with the following four commands:
//
//	go clean
//	go build ./...
//	go build -o clienttest .
//	./clienttest
//
//Required env variables are:
//
//	OAUTH_SERVICE_URL
//	OAUTH_CLIENT_ID
//	OAUTH_CLIENT_SECRET

import (
	"github.com/PennState/golang_scimclient/scim"
	"github.com/PennState/golang_scimclient/scim/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	httpClient := client.GetOauth2AuthenticatedHTTPClientFromEnvironment()
	scimClient, err := client.New("https://dev.apps.psu.edu/cpr/resources", httpClient)

	var user scim.User
	//TODO: err = scimClient.Get("/Users/9991533", &user)
	if err != nil {
		log.Error(err)
	}

	log.Infof("User: %v", user)

	var resourceTypes []scim.ResourceType
	resourceTypes, err = scimClient.GetResourceTypes()
	if err != nil {
		log.Error(err)
	}

	log.Infof("ResourceTypes(s) %v", resourceTypes)

	var schemas []scim.Schema
	schemas, err = scimClient.GetSchemas()
	if err != nil {
		log.Error(err)
	}

	log.Infof("Schema(s): %v", schemas)

	var cfg scim.ServiceProviderConfig
	cfg, err = scimClient.GetServerProviderConfig()
	if err != nil {
		log.Error(err)
	}

	log.Infof("ServiceProviderConfig: %v", cfg)
}

func search() {
	searchRequest := scim.NewSearchRequest()
	searchRequest.SortOrder = scim.Ascending
	searchRequest.SortOrder = "something else"
	searchRequest.Attributes = []string{"one", "two"}
}
