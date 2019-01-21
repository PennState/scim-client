package scim

//Group describes a SCIM user based on the RFC7643 specification
//https://tools.ietf.org/html/rfc7643#section-4.2
type Group struct {
	Resource
	DisplayName string      `json:"displayName"` //DisplayName is a human-readable name for the Group
	Members     []MemberRef `json:"members"`     //Members is a list of members in the group. While values MAY be added or removed, sub-attributes of members are "immutable".
}

//MemberRef provides a reference (URI) to a group member as well as a small
//amount of cargo data.
type MemberRef StringMultivalued
