package examples

import (
	"github.com/PennState/additional-properties/pkg/ap"
	"github.com/PennState/scim-client/pkg/scim"
	log "github.com/sirupsen/logrus"
)

// Organization represents some hierarchy of an arbitrary organization
// including (URI) an optional reference to a parent organization as well
// as to zero or more child organizations.  As with other SCIM references
// the Parent and Children references may be absolute or relative URIs.
type Organization struct {
	scim.CommonAttributes
	Name     string                  `json:"name,omitempty"`      //Name is the organization's name - e.g. "Tour Promotion"
	Type     string                  `json:"type,omitempty"`      //Type is the organization's type - e.g. "Department"
	Parent   OrganizationReference   `json:"$parent,omitempty"`   //Parent is a URI reference to a parent organization
	Children []OrganizationReference `json:"$children,omitempty"` //Children is a URI reference to zero or more child organizations
}

// OrganizationReference is a string containing an absolute or relative
// URL to another Organization.
type OrganizationReference string

const OrganizationURN = "urn:com:example:2.0:Organization"

// OrganizationResourceType provides the default structure which connects the
// Organization struct to its associated ResourceType.
var OrganizationResourceType = scim.ResourceType{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			scim.ResourceTypeURN,
		},
		ID: "Organization",
	},
	Name:        "Organization",
	Endpoint:    "/Organizations",
	Description: "SCIM ResourceType - See https://tools.ietf.org/html/rfc7643#section-6",
	Schema:      scim.ResourceTypeURN,
}

// URN returns the IANA registered SCIM name for the Organization data structure
// and, together with ResourceType() identifies this code as implementing
// the Resource interface.
func (o Organization) URN() string {
	return OrganizationURN
}

// ResourceType returns the default structure describing the availability
// of the Organization resource and, together with URN() identifies this code as
// implementing the Resource interface.
func (o Organization) ResourceType() scim.ResourceType {
	return OrganizationResourceType
}

//
// JSON marshaling and unmarshaling
//

// MarshalJSON implements https://golang.org/pkg/encoding/json/#Marshaler
func (o Organization) MarshalJSON() ([]byte, error) {
	type Alias Organization
	json := ap.ConfigCompatibleWithStandardLibrary
	return json.Marshal((Alias)(o))
}

// UnmarshalJSON implements https://golang.org/pkg/encoding/json/#Unmarshaler
func (o *Organization) UnmarshalJSON(data []byte) error {
	type Alias Organization
	json := ap.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, (*Alias)(o))
}
