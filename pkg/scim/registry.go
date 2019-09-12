package scim

import "sync"

//ResourceRegistry contains a map which will be treated as a singleton.
type ResourceRegistry struct {
	resourceMap map[string]ResourceType
}

var instance *ResourceRegistry
var once sync.Once

//GetResourceRegistry returns creates (if necessary) and returns the
//singleton instance of the ResourceRegistry.  The five ResourceTypes
//defined in the SCIM specifications are included in the registry by
//default.
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

//Lookup returns the ResourceType associated with the name parameter
//as well as a boolean that indicates whether the registry contained
//the requested ResourceType.
func (rr ResourceRegistry) Lookup(name string) (ResourceType, bool) {
	rt, ok := rr.resourceMap[name]
	return rt, ok
}

//Register allows one or more custom ResourceTypes to be added to the
//registry.
func (rr ResourceRegistry) Register(rts ...ResourceType) {
	for _, rt := range rts {
		rr.resourceMap[rt.Name] = rt
	}
}
