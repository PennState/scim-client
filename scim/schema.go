package scim

type Schema struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Attributes []Attribute `json:"attributes"`
}

//https://tools.ietf.org/html/rfc7643#section-2.4
type Attribute struct {
	Name string
	Type Type
	Subattributes []Attribute
	Multivalued bool
	Description string
	Required bool
	CanonicalValues []string
	CaseExact bool
	Mutability Mutability
	Returned Returned
	Uniqueness Uniqueness
	ReferenceTypes []string
}

type Type = string
const (
	String Type = "string"
	Boolean Type = "boolean"
	Decimal Type = "decimal"
	Integer Type = "integer"
	DateTime Type = "dateTime"
	Reference Type = "reference"
	Complex Type = "complex"
)

type Mutability = string
const (
	ReadOnly Mutability = "readOnly"
	ReadWrite Mutability = "readWrite"
	Immutable Mutability = "immutable"
	WriteOnly Mutability = "writeOnly"
)

type Returned = string
const (
	Always Returned = "always"
	Never Returned = "never"
	Default Returned = "default"
	Request Returned = "request"
)

type Uniqueness = string
const (
	None Uniqueness = "none"
	Server Uniqueness = "server"
	Global Uniqueness = "global"
)
