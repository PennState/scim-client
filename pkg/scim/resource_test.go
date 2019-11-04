package scim

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/PennState/additional-properties/pkg/ap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	ca.AdditionalProperties = make(map[string]json.RawMessage)
	ca.AdditionalProperties["urn:fake.extension"] = json.RawMessage(`{"name": "Fake Extension"}`)
	ca.AdditionalProperties["additionalPropertiesOne"] = json.RawMessage(`"additionalPropertiesOne"`)
	ca.AdditionalProperties["additionalPropertiesTwo"] = json.RawMessage(`"additionalPropertiesTwo"`)
	return ca
}

func TestAddExtension(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension
	require.Len(resource.AdditionalProperties, 3)

	err := resource.AddExtension(&fakeExtension)
	assert.NotNil(err)
	assert.Len(resource.AdditionalProperties, 3)

	err = resource.AddExtension(&worthlessExtension)
	assert.Nil(err)
	assert.Len(resource.AdditionalProperties, 4)
	value, exists := resource.AdditionalProperties["urn:worthless.extension"]
	assert.True(exists)
	assert.Equal(json.RawMessage("{}"), value)
}

func TestAddAndUpdateExtensionWithEmptyResource(t *testing.T) {
	var ca CommonAttributes
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension
	assert.NoError(t, ca.AddExtension(fakeExtension))
	assert.NoError(t, ca.AddExtension(worthlessExtension))
	assert.NoError(t, ca.UpdateExtension(&fakeExtension))
	assert.NoError(t, ca.UpdateExtension(&worthlessExtension))
}

func TestGetExtension(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()

	var fakeExtension fakeExtension
	err := resource.GetExtension(&fakeExtension)
	assert.Nil(err)
}

func TestGetExtensionWithEmptyResource(t *testing.T) {
	var ca CommonAttributes
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension
	assert.Error(t, ca.GetExtension(&fakeExtension))
	assert.Error(t, ca.GetExtension(&worthlessExtension))
}

func TestGetExtensionURNs(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()

	urns := resource.GetExtensionURNs()
	assert.Len(urns, 1)
	assert.Equal("urn:fake.extension", urns[0])

	resource.AdditionalProperties["urn:worthless.extension"] = json.RawMessage("{}")
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

func TestHasExtensionWithEmptyResource(t *testing.T) {
	var ca CommonAttributes
	var fakeExtension fakeExtension
	var worthlessExtension worthlessExtension
	assert.False(t, ca.HasExtension((fakeExtension)))
	assert.False(t, ca.HasExtension(worthlessExtension))
}

func TestRemoveExtension(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	require.Len(resource.AdditionalProperties, 3)

	resource.RemoveExtension(&fakeExtension)
	assert.Len(resource.AdditionalProperties, 2)
	_, exists := resource.AdditionalProperties["urn:fake.extension"]
	assert.False(exists)
}

func TestUpdateExtension(t *testing.T) {
	assert := assert.New(t)
	resource := getResourceWithAdditionalProperties()
	var fakeExtension fakeExtension
	fakeExtension.Name = "Updated Fake Extension"

	err := resource.UpdateExtension(&fakeExtension)
	assert.Nil(err)
	value, exists := resource.AdditionalProperties["urn:fake.extension"]
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

func TestResourceMarshaling(t *testing.T) {
	assert := assert.New(t)

	ca := getResourceWithAdditionalProperties()
	jap := ap.ConfigCompatibleWithStandardLibrary
	data, err := jap.Marshal(&ca)
	if err != nil {
		assert.Error(err)
	}

	var obj map[string]json.RawMessage
	err = jap.Unmarshal(data, &obj)
	if err != nil {
		assert.Error(err)
	}
	assert.Contains(obj, "id", "meta", "urn:fake.extension", "additionalPropertiesOne", "additionalPropertiesTwo")
	assert.JSONEq("{\"name\":\"Fake Extension\"}", string(obj["urn:fake.extension"]))
	assert.JSONEq("\"additionalPropertiesOne\"", string(obj["additionalPropertiesOne"]))
	assert.JSONEq("\"additionalPropertiesTwo\"", string(obj["additionalPropertiesTwo"]))
	meta := "{\"created\":\"0001-01-01T00:00:00Z\",\"lastModified\":\"0001-01-01T00:00:00Z\",\"location\":\"\",\"resourceType\":\"\",\"version\":\"\"}"
	assert.JSONEq(meta, string(obj["meta"]))
}

//
//
// Resource Unmarshaling tests
//
//

func TestResourceUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	const resourceJSON = `
	{
		"schemas": [
			"urn:ietf:params:scim:schemas:core:2.0:User"
		],
		"id": "2819c223-7f76-453a-919d-413861904646",
		"externalId": "43496746-7739-460b-bf99-3421f2970687",
		"meta": {
			"resourceType": "User",
			"created": "2010-01-23T04:56:22Z",
			"lastModified": "2011-05-13T04:42:34Z",
			"version": "W/\"3694e05e9dff590\"",
			"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
		}
	}
	`

	var ca CommonAttributes
	error := json.Unmarshal([]byte(resourceJSON), &ca)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(ca.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
	assert.Equal(ca.ExternalID, "43496746-7739-460b-bf99-3421f2970687")

	assert.Equal(ca.Meta.ResourceType, "User")
	assert.Equal(ca.Meta.Created, time.Date(2010, time.January, 23, 4, 56, 22, 0, time.UTC))
	assert.Equal(ca.Meta.LastModified, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC))
	assert.Equal(ca.Meta.Version, "W/\"3694e05e9dff590\"")
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
	err := json.Unmarshal([]byte(badResourceJSON), &ca)
	assert.NotNil(err)
}
