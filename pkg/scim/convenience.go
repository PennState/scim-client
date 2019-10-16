package scim

import (
	"context"
)

const (
	resourcesByExternalId = "externalId EQ \"%s\""
	userByUserName        = "userName EQ \"%s\""
)

// QueryResourcesByExternalID is a helper method for retrieving resources
// from a ResourceType by ExternalID
func (c Client) QueryResourceTypeByExternalID(ctx context.Context, rt ResourceType, externalID string) (ListResponse, error) {
	return c.QueryResourceType(ctx, rt, NewSearchRequestFromFormat(resourcesByExternalId, externalID))
}

// QueryServerByExternalID is a helper method for retrieving any resources
// from the server by ExternalID
func (c Client) QueryServerByExternalID(ctx context.Context, externalID string) (ListResponse, error) {
	return c.QueryServer(ctx, NewSearchRequestFromFormat(resourcesByExternalId, externalID))
}

// QueryUserResourcesByUserName is a helper method for retrieving User
// resources by UserName
func (c Client) QueryUserResourceTypeByUserName(ctx context.Context, userName string) (ListResponse, error) {
	return c.QueryResourceType(ctx, UserResourceType, NewSearchRequestFromFormat(userByUserName, userName))
}
