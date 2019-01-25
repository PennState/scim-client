package scim

//ResourceTypeURN is the SCIM URN registered to identify the ResourceType
//data structure.
//https://tools.ietf.org/html/rfc7643#section-6
const ResourceTypeURN = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"

//ResourceType specifies the metadata about a SCIM resource type.
//https://tools.ietf.org/html/rfc7643#section-6
type ResourceType struct {
	CommonAttributes
	Name             string            `json:"name" validate:"required`
	Description      string            `json:"description"`
	Endpoint         string            `json:"endpoint"`
	Schema           string            `json:"schema"`
	SchemaExtensions []SchemaExtension `json:"schemaExtensions"`
}

//SchemaExtension is a list of URIs that specify the allowed
//SCIM extensions for a ResourceType
type SchemaExtension struct {
	Schema   string `json:"schema"`
	Required bool   `json:"required"`
}

//ResourceTypeResourceType is the constant data that defines the metadata
//for the ResourceType server discovery endpoint
var ResourceTypeResourceType = ResourceType{
	CommonAttributes: CommonAttributes{
		Schemas: []string{
			ResourceTypeURN,
		},
		ID: "ResourceType",
	},
	Name:        "ResourceType",
	Endpoint:    "/ResourceTypes",
	Description: "SCIM ResourceType - See https://tools.ietf.org/html/rfc7643#section-6",
	Schema:      ResourceTypeURN,
}

func (rt ResourceType) URN() string {
	return ResourceTypeURN
}

//ServerDiscoveryResourceType allows a ResourceType to retrieve the related
//metadata for its SCIM type.
func (rt ResourceType) ResourceType() ResourceType {
	return ResourceTypeResourceType
}
