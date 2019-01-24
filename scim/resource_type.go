package scim

const ResourceTypeURN = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"

type ResourceType interface {
	NewResourceType() resourceType
	NewResource() Resource
}

type resourceType struct {
	id               string            `json:"id" validate:"required"`
	schemas          []string          `json:"schemas" validate:"required"`
	name             string            `json:"name" validate:"required`
	description      string            `json:"description"`
	endpoint         string            `json:"endpoint"`
	schema           string            `json:"schema"`
	schemaExtensions []schemaExtension `json:"schemaExtensions"`
}

func (rt resourceType) ID() string                          { return rt.id }
func (rt resourceType) Schemas() []string                   { return rt.schemas }
func (rt resourceType) Name() string                        { return rt.name }
func (rt resourceType) Description() string                 { return rt.description }
func (rt resourceType) Endpoint() string                    { return rt.endpoint }
func (rt resourceType) Schema() string                      { return rt.schema }
func (rt resourceType) SchemaExtensions() []schemaExtension { return rt.schemaExtensions }

type schemaExtension struct {
	schema   string `json:"schema"`
	required bool   `json:"required"`
}

func NewSchemaExtension(schema string, required bool) schemaExtension {
	return schemaExtension{
		schema:   schema,
		required: required,
	}
}

func (se schemaExtension) Schema() string { return se.schema }
func (se schemaExtension) Required() bool { return se.required }

// var UserResourceType = NewResourceType{
// 	CommonAttributes: CommonAttributes{
// 		ID:      "User",
// 		Schemas: []string{"urn:"},
// 	},
// 	Name:        "User",
// 	Description: "",
// 	Endpoint:    "/Users",
// }
