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
	log "github.com/sirupsen/logrus"
)

func main() {
	sClient, err := scim.NewOAuthClientFromEnv()
	if err != nil {
		return
	}

	var user scim.User
	err = sClient.RetrieveResource(&user, "9991533")
	if err != nil {
		log.Error(err)
	}

	log.Infof("User: %v", user)
	extensionURNs := user.GetExtensionURNs()
	log.Infof("User extensions: %v", extensionURNs)

	var resourceTypes []scim.ResourceType
	resourceTypes, err = sClient.GetResourceTypes()
	if err != nil {
		log.Error(err)
	}

	log.Infof("ResourceTypes(s) %v", resourceTypes)

	var schemas []scim.Schema
	schemas, err = sClient.GetSchemas()
	if err != nil {
		log.Error(err)
	}

	log.Infof("Schema(s): %v", schemas)

	var cfg scim.ServiceProviderConfig
	cfg, err = sClient.GetServerProviderConfig()
	if err != nil {
		log.Error(err)
	}

	log.Infof("ServiceProviderConfig: %v", cfg)
}

func search() {
	var searchRequest scim.SearchRequest
	searchRequest.SortOrder = scim.Ascending
	searchRequest.SortOrder = "something else"
	searchRequest.Attributes = []string{"one", "two"}
}
