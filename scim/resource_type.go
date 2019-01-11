package scim

type ResourceType struct {
	Resource

	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Endpoint         string            `json:"endpoint"`
	SchemaExtensions []SchemaExtension `json:"schemaExtensions"`
}

type SchemaExtension struct {
	Schema   string `json:"schema"`
	Required bool   `json:"required"`
}
