package examples

import (
	"testing"

	"github.com/PennState/additional-properties/pkg/ap"
	"github.com/PennState/scim-client/pkg/scim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pantryJSON = `
{
	"schemas": [
		"urn:ietf:params:scim:schemas:core:2.0:User",
		"urn:com:example:2.0:Pantry"
	],
	"id": "2819c223-7f76-453a-919d-413861904646",
	"externalId": "701984",
	"userName": "bjensen@example.com",
	"urn:com:example:2.0:Pantry": {
		"building": "Technology Support Building",
		"office": "202BB",
		"balance": 4.25
	},
	"meta": {
		"resourceType": "User",
		"created": "2010-01-23T04:56:22Z",
		"lastModified": "2011-05-13T04:42:34Z",
		"version": "W/a330bc54f0671c9",
		"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
	}
}`

func TestPantryExtension(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var user scim.User
	ap := ap.ConfigCompatibleWithStandardLibrary
	err := ap.Unmarshal([]byte(pantryJSON), &user)
	require.Nil(err)

	var pantry Pantry
	require.True(user.HasExtension(&pantry))
	err = user.GetExtension(&pantry)
	assert.Nil(err)
	assert.Equal("Technology Support Building", pantry.Building)
	assert.Equal("202BB", pantry.Office)
	assert.Equal(4.25, pantry.Balance)
}
