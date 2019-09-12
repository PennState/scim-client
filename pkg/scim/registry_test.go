package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryInstanceIsConstant(t *testing.T) {
	assert := assert.New(t)
	ref1 := GetResourceRegistry()
	ref2 := GetResourceRegistry()
	assert.Equal(ref1, ref2)
}

func TestRegistryInstanceContainsDefaultResources(t *testing.T) {
	assert := assert.New(t)
	registry := GetResourceRegistry()
	assert.Len(registry.resourceMap, 5)
	rt, ok := registry.Lookup("Group")
	assert.True(ok)
	assert.Equal(GroupResourceType, rt)
}

func TestResourceTypeCanBeAdded(t *testing.T) {
	myResourceType := ResourceType{
		CommonAttributes: CommonAttributes{
			Schemas: []string{
				"urn:fake.urn.goes.here",
			},
			ID: "Fake",
		},
		Name:        "Fake",
		Endpoint:    "/Fakes",
		Description: "SCIM ResourceType - See https://tools.ietf.org/html/rfc7643#section-6",
		Schema:      ResourceTypeURN,
	}

	assert := assert.New(t)
	registry := GetResourceRegistry()
	registry.Register(myResourceType)
	assert.Len(registry.resourceMap, 6)
	rt, ok := registry.Lookup("Fake")
	assert.True(ok)
	assert.Equal(myResourceType, rt)
}
