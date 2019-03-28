package scim

import "sync"

type ResourceRegistry struct {
	resourceMap map[string]ResourceType
}

var instance *ResourceRegistry
var once sync.Once

func GetResourceRegistry() *ResourceRegistry {
	once.Do(func() {
		instance = &ResourceRegistry{
			resourceMap: make(map[string]ResourceType),
		}
		instance.Register(
			GroupResourceType,
			ResourceTypeResourceType,
			SchemaResourceType,
			ServiceProviderConfigResourceType,
			UserResourceType,
		)
	})
	return instance
}

func (rr ResourceRegistry) Lookup(name string) (ResourceType, bool) {
	rt, ok := rr.resourceMap[name]
	return rt, ok
}

func (rr ResourceRegistry) Register(rts ...ResourceType) {
	for _, rt := range rts {
		rr.resourceMap[rt.Name] = rt
	}
}
