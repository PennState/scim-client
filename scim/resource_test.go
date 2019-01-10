package scim

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const commonResourceJSON = `
{
	"id": "2819c223-7f76-453a-919d-413861904646",
	"externalId": "43496746-7739-460b-bf99-3421f2970687",
	"meta": {
		"resourceType": "User",
		"created": "2010-01-23T04:56:22Z",
		"lastModified": "2011-05-13T04:42:34Z",
		"version": "W/3694e05e9dff590",
		"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
	}
}`

func CommonResourceMarshaling(t *testing.T) {
	assert := assert.New(t)

	var commonResource = Resource{ID: "2819c223-7f76-453a-919d-413861904646"}

	b, err := json.Marshal(commonResource)

	if err != nil {

	}

	expected := json.RawMessage(commonResourceJSON)
	actual := json.RawMessage(b)
	assert.Equal(actual, expected)
}

func TestCommonResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	var commonResource Resource
	error := json.Unmarshal([]byte(commonResourceJSON), &commonResource)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(commonResource.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
	assert.Equal(commonResource.ExternalID, "43496746-7739-460b-bf99-3421f2970687")

	assert.Equal(commonResource.Meta.ResourceType, "User")
	assert.Equal(commonResource.Meta.Created, time.Date(2010, time.January, 23, 4, 56, 22, 0, time.UTC))
	assert.Equal(commonResource.Meta.LastModified, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC))
	assert.Equal(commonResource.Meta.Version, "W/3694e05e9dff590")
	assert.Equal(commonResource.Meta.Location, "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646")
}
