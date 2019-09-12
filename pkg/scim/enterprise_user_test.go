package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// https://tools.ietf.org/html/rfc7643#section-8.2
func TestEnterpriseUserUnmarshaling(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	const enterpriseUserJSON = `
	{
		"schemas": [
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
		],
		"id": "2819c223-7f76-453a-919d-413861904646",
		"externalId": "701984",
		"userName": "bjensen@example.com",
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": {
			"employeeNumber": "701984",
			"costCenter": "4130",
			"organization": "Universal Studios",
			"division": "Theme Park",
			"department": "Tour Operations",
			"manager": {
			  "value": "26118915-6090-4610-87e4-49d8ca9f808d",
			  "$ref": "../Users/26118915-6090-4610-87e4-49d8ca9f808d",
			  "displayName": "John Smith"
			}
		},
		"meta": {
			"resourceType": "User",
			"created": "2010-01-23T04:56:22Z",
			"lastModified": "2011-05-13T04:42:34Z",
			"version": "W/a330bc54f0671c9",
			"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
		}
	}`

	var user User
	err := Unmarshal([]byte(enterpriseUserJSON), &user)

	require.Nil(err, "Error unmarshaling the User resource - %q", err)
	assert.Len(user.Schemas, 2)
	assert.Equal("urn:ietf:params:scim:schemas:core:2.0:User", user.Schemas[0])
	assert.Equal("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User", user.Schemas[1])
	assert.Len(user.AdditionalProperties, 1, "Error! There should be exactly one SCIM extension")

	var enterpriseUser EnterpriseUser
	err = user.GetExtension(&enterpriseUser)

	require.Nil(err, "Error unmarshaling the ExterpriseUser extension - %q", err)
	assert.Equal("701984", enterpriseUser.EmployeeNumber)
	assert.Equal("4130", enterpriseUser.CostCenter)
	assert.Equal("Universal Studios", enterpriseUser.Organization)
	assert.Equal("Theme Park", enterpriseUser.Division)
	assert.Equal("Tour Operations", enterpriseUser.Department)
	assert.Equal("26118915-6090-4610-87e4-49d8ca9f808d", enterpriseUser.Manager.Value)
	assert.Equal("../Users/26118915-6090-4610-87e4-49d8ca9f808d", enterpriseUser.Manager.Reference)
	assert.Equal("John Smith", enterpriseUser.Manager.DisplayName)
}
