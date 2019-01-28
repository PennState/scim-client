package scim

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const resourceJSON = `
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

const resourceWithAdditionalPropertiesJSON = `
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
	Name string `json:"name"`
}

func (fe fakeExtension) URN() string {
	return "urn:fake.extension"
}

type worthlessExtension struct {
}

func (we worthlessExtension) URN() string {
	return "urn:worthless.extension"
}

//
//
// Extension management tests
//
//

func getResourceWithAdditionalProperties() CommonAttributes {
	var ca CommonAttributes

	ca.ID = "2819c223-7f76-453a-919d-413861904646"
	ca.additionalProperties = make(map[string]json.RawMessage)
	ca.additionalProperties["urn:fake.extension"] = json.RawMessage(`{"name": "Fake Extension"}`)
	ca.additionalProperties["additionalPropertiesOne"] = json.RawMessage(`"additionalPropertiesOne"`)
	ca.additionalProperties["additionalPropertiesTwo"] = json.RawMessage(`"additionalPropertiesTwo"`)
	return ca
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
	fakeExtension.Name = "Updated Fake Extension"

	err := resource.UpdateExtension(&fakeExtension)
	assert.Nil(err)
	value, exists := resource.additionalProperties["urn:fake.extension"]
	assert.True(exists)
	assert.Equal(json.RawMessage(`{"name":"Updated Fake Extension"}`), value)

	var worthlessExtension worthlessExtension
	err = resource.UpdateExtension(worthlessExtension)
	assert.NotNil(err)
}

//
//
// Resource Marshaling tests
//
//

func ResourceMarshaling(t *testing.T) {
	assert := assert.New(t)

	var ca = CommonAttributes{ID: "2819c223-7f76-453a-919d-413861904646"}

	b, err := json.Marshal(ca)

	if err != nil {

	}

	expected := json.RawMessage(resourceJSON)
	actual := json.RawMessage(b)
	assert.Equal(actual, expected)
}

//
//
// Resource Unmarshaling tests
//
//

func TestResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	var ca CommonAttributes
	error := Unmarshal([]byte(resourceJSON), &ca)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(ca.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
	assert.Equal(ca.ExternalID, "43496746-7739-460b-bf99-3421f2970687")

	assert.Equal(ca.Meta.ResourceType, "User")
	assert.Equal(ca.Meta.Created, time.Date(2010, time.January, 23, 4, 56, 22, 0, time.UTC))
	assert.Equal(ca.Meta.LastModified, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC))
	assert.Equal(ca.Meta.Version, "W/3694e05e9dff590")
	assert.Equal(ca.Meta.Location, "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646")
}

func TestBadResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	const badResourceJSON = `
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

	var ca CommonAttributes
	err := Unmarshal([]byte(badResourceJSON), &ca)
	assert.NotNil(err)
}
