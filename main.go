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
	err = scimClient.Get("/Users/9991533", &user)
	if err != nil {
		log.Error(err)
	}

	log.Infof("User: %s", user)
}
