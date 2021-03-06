package scim

import (
	"github.com/PennState/additional-properties/pkg/ap"
)

//Group describes a SCIM user based on the RFC7643 specification
//https://tools.ietf.org/html/rfc7643#section-4.2
type Group struct {
	CommonAttributes
	DisplayName string      `json:"displayName"`       //DisplayName is a human-readable name for the Group
	Members     []MemberRef `json:"members,omitempty"` //Members is a list of members in the group. While values MAY be added or removed, sub-attributes of members are "immutable".
}

//MemberRef provides a reference (URI) to a group member as well as a small
//amount of cargo data.
type MemberRef StringMultivalued

//GroupURN is the IANA registered SCIM name for the standardized SCIM
//Group.
const GroupURN = "urn:ietf:params:scim:schemas:core:2.0:Group"

//GroupResourceType provides the default structure which connects the User
//struct to its associated ResourceType.
var GroupResourceType = ResourceType{
	CommonAttributes: CommonAttributes{
		Schemas: []string{
			ResourceTypeURN,
		},
		ID: "Group",
	},
	Name:        "Group",
	Endpoint:    "/Groups",
	Description: "SCIM ResourceType - See https://tools.ietf.org/html/rfc7643#section-6",
	Schema:      ResourceTypeURN,
}

//URN returns the IANA registered SCIM name for the User data structure
//and, together with ResourceType() identifies this code as implementing
//the Resource interface.
func (g Group) URN() string {
	return GroupURN
}

//ResourceType returns the default structure describing the availability
//of the Group resource and, together with URN() identifies this code as
//implementing the Resource interface.
func (g Group) ResourceType() ResourceType {
	return GroupResourceType
}

//
// JSON marshaling and unmarshaling
//

// MarshalJSON implements https://golang.org/pkg/encoding/json/#Marshaler
func (g Group) MarshalJSON() ([]byte, error) {
	type Alias Group
	json := ap.ConfigCompatibleWithStandardLibrary
	return json.Marshal((Alias)(g))
}

// UnmarshalJSON implements https://golang.org/pkg/encoding/json/#Unmarshaler
func (g *Group) UnmarshalJSON(data []byte) error {
	type Alias Group
	json := ap.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, (*Alias)(g))
}
