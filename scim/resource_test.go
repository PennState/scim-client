package scim

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

const commonResourceWithAdditionalPropertiesJSON = `
{
	"id": "2819c223-7f76-453a-919d-413861904646",
	"externalId": "43496746-7739-460b-bf99-3421f2970687",
	"meta": {
		"resourceType": "User",
		"created": "2010-01-23T04:56:22Z",
		"lastModified": "2011-05-13T04:42:34Z",
		"version": "W/3694e05e9dff590",
		"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
	},
	"urn:fake.extension": {
		"urn": "Fake Extension"
	},
	"additionalPropertyOne": "additionalPropertyOne",
	"additionalPropertyTwo": "additionalPropertyTwo"
}`

type fakeExtension struct {
	URN string `json:"urn"`
}

func (fe fakeExtension) GetURN() string {
	return "urn:fake.extension"
}

type worthlessExtension struct {
}

func (we worthlessExtension) GetURN() string {
	return "urn:worthless.extension"
}

//
//
// Extension management tests
//
//

func getResourceWithAdditionalProperties() Resource {
	var resource Resource

	resource.ID = "2819c223-7f76-453a-919d-413861904646"
	resource.additionalProperties = make(map[string]json.RawMessage)
	resource.additionalProperties["urn:fake.extension"] = json.RawMessage(`{"urn": "Fake Extension"}`)
	resource.additionalProperties["additionalPropertiesOne"] = json.RawMessage(`"additionalPropertiesOne"`)
	resource.additionalProperties["additionalPropertiesTwo"] = json.RawMessage(`"additionalPropertiesTwo"`)
	return resource
}

func TestAddExtension(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension
	require.Len(resource.additionalProperties, 3)

	err := resource.AddExtension(&fakeExtension)
	assert.NotNil(err)
	assert.Len(resource.additionalProperties, 3)

	err = resource.AddExtension(&worthlessExtension)
	assert.Nil(err)
	assert.Len(resource.additionalProperties, 4)
	value, exists := resource.additionalProperties["urn:worthless.extension"]
	assert.True(exists)
	assert.Equal(json.RawMessage("{}"), value)
}

func TestGetExtension(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()

	var fakeExtension fakeExtension
	err := resource.GetExtension(&fakeExtension)
	assert.Nil(err)
}

func TestGetExtensionURNs(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()

	urns := resource.GetExtensionURNs()
	assert.Len(urns, 1)
	assert.Equal("urn:fake.extension", urns[0])

	resource.additionalProperties["urn:worthless.extension"] = json.RawMessage("{}")
	urns = resource.GetExtensionURNs()
	assert.Len(urns, 2)
	sort.Strings(urns)
	assert.Equal("urn:worthless.extension", urns[1])
}

func TestHasExtension(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension

	assert.True(resource.HasExtension(fakeExtension))
	assert.False(resource.HasExtension(worthlessExtension))
}

func TestRemoveExtension(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	require.Len(resource.additionalProperties, 3)

	resource.RemoveExtension(&fakeExtension)
	assert.Len(resource.additionalProperties, 2)
	_, exists := resource.additionalProperties["urn:fake.extension"]
	assert.False(exists)
}

func TestUpdateExtension(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	fakeExtension.URN = "Updated Fake Extension"

	err := resource.UpdateExtension(&fakeExtension)
	assert.Nil(err)
	value, exists := resource.additionalProperties["urn:fake.extension"]
	assert.True(exists)
	assert.Equal(json.RawMessage(`{"urn":"Updated Fake Extension"}`), value)

	var worthlessExtension worthlessExtension
	err = resource.UpdateExtension(worthlessExtension)
	assert.NotNil(err)
}

//
//
// Resource Marshaling tests
//
//

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

//
//
// Resource Unmarshaling tests
//
//

func TestCommonResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	var commonResource Resource
	error := Unmarshal([]byte(commonResourceJSON), &commonResource)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(commonResource.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
	assert.Equal(commonResource.ExternalID, "43496746-7739-460b-bf99-3421f2970687")

	assert.Equal(commonResource.Meta.ResourceType, "User")
	assert.Equal(commonResource.Meta.Created, time.Date(2010, time.January, 23, 4, 56, 22, 0, time.UTC))
	assert.Equal(commonResource.Meta.LastModified, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC))
	assert.Equal(commonResource.Meta.Version, "W/3694e05e9dff590")
	assert.Equal(commonResource.Meta.Location, "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646")
}

func TestBadCommonResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	const badCommonResourceJSON = `
	{
		"id": [
			"2819c223-7f76-453a-919d-413861904646"
		],
		"externalId": {
			"id": "43496746-7739-460b-bf99-3421f2970687"
		},
		"meta": {
			"resourceType": "User",
			"created": "2010-01-23T04:56:22Z",
			"lastModified": "2011-05-13T04:42:34Z",
			"version": "W/3694e05e9dff590",
			"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
		}
	}`

	var badCommonResource Resource
	error := Unmarshal([]byte(badCommonResourceJSON), &badCommonResource)
	assert.NotNil(error)
}
