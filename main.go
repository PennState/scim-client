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
//- OAUTH_SERVICE_URL
//- OAUTH_CLIENT_ID
//- OAUTH_CLIENT_SECRET
//- SCIM_SERVICE_URL

import (
	"encoding/json"

	"github.com/PennState/golang_scimclient/scim"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.AddHook(filename.NewHook())

	sClient, err := scim.NewOAuthClientFromEnv()
	if err != nil {
		log.Error(err)
		return
	}

	// Returns a SCIM user by id

	var user scim.User
	err = sClient.RetrieveResource(&user, "9991533")
	if err != nil {
		log.Error(err)
	}
	log.Infof("User: %v", user)
	extensionURNs := user.GetExtensionURNs()
	log.Infof("User extensions: %v", extensionURNs)

	j, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
	}
	log.Info("User JSON: ", string(j))

	// Finding users who have a specific PSU id prefix

	log.Info("===== Finding users who have a specific PSU id prefix =====")
	sr := scim.SearchRequest{
		Filter: "externalId SW \"972806\"", // 9991533 has PSU Id 972806446
	}
	lr, err := sClient.SearchResource(scim.UserResourceType, sr)
	if err != nil {
		log.Error(err)
	}
	log.Info("ListResponse: ", lr)
	for ord, res := range lr.Resources {
		log.Infof("ListResponse resource %d: %v", ord, res)
	}

	// Search using the convenience methods

	log.Info("===== Search using the convenience method - SearchUserResourcesByUserName =====")
	lr, err = sClient.SearchUserResourcesByUserName("swm16")
	if err != nil {
		log.Error(err)
	}
	log.Info("ListResponse: ", lr)
	for ord, res := range lr.Resources {
		log.Infof("ListResponse resource %d: %v", ord, res)
	}

	log.Info("===== Search using the convenience method - SearchResourcesByExternalId =====")
	lr, err = sClient.SearchResourcesByExternalID(scim.UserResourceType, "972806446")
	if err != nil {
		log.Error(err)
	}
	log.Info("ListResponse: ", lr)
	for ord, res := range lr.Resources {
		log.Infof("ListResponse resource %d: %v", ord, res)
	}

	// Replace the user's resource

	log.Info("===== ReplaceResource =====")
	err = sClient.ReplaceResource(&user)
	if err != nil {
		log.Error(err)
	}
	//	user.Name.GivenName = "Stephen"
	log.Infof("User: %v", user)
	extensionURNs = user.GetExtensionURNs()
	log.Infof("User extensions: %v", extensionURNs)

	j, err = json.Marshal(user)
	if err != nil {
		log.Error(err)
	}
	log.Info("User JSON: ", string(j))

	// List the server's resource types

	var resourceTypes []scim.ResourceType
	resourceTypes, err = sClient.GetResourceTypes()
	if err != nil {
		log.Error(err)
	}
	log.Infof("ResourceTypes(s) %v", resourceTypes)

	// List the server's schemas

	var schemas []scim.Schema
	schemas, err = sClient.GetSchemas()
	if err != nil {
		log.Error(err)
	}
	log.Infof("Schema(s): %v", schemas)

	// List the server's service provider config

	var cfg scim.ServiceProviderConfig
	cfg, err = sClient.GetServerProviderConfig()
	if err != nil {
		log.Error(err)
	}
	log.Infof("ServiceProviderConfig: %v", cfg)
}
