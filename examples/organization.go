package examples

import "github.com/PennState/golang_scimclient/scim"

//Organization represents some hierarchy of an arbitrary organization
//including (URI) an optional reference to a parent organization as well
//as to zero or more child organizations.  As with other SCIM references
//the Parent and Children references may be absolute or relative URIs.
type Organization struct {
	scim.Resource
	Name     string                `json:"name"`      //Name is the organization's name - e.g. "Tour Promotion"
	Type     string                `json:"type"`      //Type is the organization's type - e.g. "Department"
	Parent   OrganizationReference `json:"$parent"`   //Parent is a URI reference to a parent organization
	Children OrganizationReference `json:"$children"` //Children is a URI reference to zero or more child organizations
}

//OrganizationReference is a string containing an absolute or relative
//URL to another Organization.
type OrganizationReference string
