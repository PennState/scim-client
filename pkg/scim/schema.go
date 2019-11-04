package scim

const SchemaURN = "urn:ietf:params:scim:schemas:core:2.0:Schema"

//https://tools.ietf.org/html/rfc7643#section-7
type Schema struct {
	CommonAttributes
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Attributes  []Attribute `json:"attributes"`
}

type Attribute struct {
	Name            string      `json:"name"`
	Type            Type        `json:"type"`
	SubAttributes   []Attribute `json:"subAttributes"`
	Multivalued     bool        `json:"multiValued"`
	Description     string      `json:"description"`
	Required        bool        `json:"required"`
	CanonicalValues []string    `json:"canonicalValues,omitempty"`
	CaseExact       bool        `json:"caseExact"`
	Mutability      Mutability  `json:"mutability"`
	Returned        Returned    `json:"returned"`
	Uniqueness      Uniqueness  `json:"uniqueness"`
	ReferenceTypes  []string    `json:"referenceTypes"`
}

type Type string

const (
	String    Type = "string"
	Boolean   Type = "boolean"
	Decimal   Type = "decimal"
	Integer   Type = "integer"
	DateTime  Type = "dateTime"
	Reference Type = "reference"
	Complex   Type = "complex"
)

type Mutability string

const (
	ReadOnly  Mutability = "readOnly"
	ReadWrite Mutability = "readWrite"
	Immutable Mutability = "immutable"
	WriteOnly Mutability = "writeOnly"
)

type Returned string

const (
	Always  Returned = "always"
	Never   Returned = "never"
	Default Returned = "default"
	Request Returned = "request"
)

type Uniqueness string

const (
	None   Uniqueness = "none"
	Server Uniqueness = "server"
	Global Uniqueness = "global"
)

var SchemaResourceType = ResourceType{
	CommonAttributes: CommonAttributes{
		Schemas: []string{
			ResourceTypeURN,
		},
		ID: "ResourceType",
	},
	Name:        "Schema",
	Endpoint:    "Schemas",
	Description: "SCIM Schema - See https://tools.ietf.org/html/rfc7643#section-7",
	Schema:      SchemaURN,
}

func (s Schema) URN() string {
	return SchemaURN
}

func (s Schema) ResourceType() ResourceType {
	return SchemaResourceType
}
