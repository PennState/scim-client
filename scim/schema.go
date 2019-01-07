package scim

//https://tools.ietf.org/html/rfc7643#section-7
type Schema struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Name string `json:"name"`
	Type Type `json:"type"`
	SubAttributes []Attribute `json:"subAttributes"`
	Multivalued bool `json:"multiValued"`
	Description string `json:"description"`
	Required bool `json:"required"`
	CanonicalValues []string `json:"canonicalValues"`
	CaseExact bool `json:"caseExact"`
	Mutability Mutability `json:"mutability"`
	Returned Returned `json:"returned"`
	Uniqueness Uniqueness `json:"uniqueness"`
	ReferenceTypes []string `json:"referenceTypes"`
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
