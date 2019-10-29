package examples

import (
	"github.com/PennState/additional-properties/pkg/ap"
	"github.com/PennState/scim-client/pkg/scim"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

//Organization represents some hierarchy of an arbitrary organization
//including (URI) an optional reference to a parent organization as well
//as to zero or more child organizations.  As with other SCIM references
//the Parent and Children references may be absolute or relative URIs.
type Organization struct {
	scim.CommonAttributes
	Name     string                  `json:"name"`      //Name is the organization's name - e.g. "Tour Promotion"
	Type     string                  `json:"type"`      //Type is the organization's type - e.g. "Department"
	Parent   OrganizationReference   `json:"$parent"`   //Parent is a URI reference to a parent organization
	Children []OrganizationReference `json:"$children"` //Children is a URI reference to zero or more child organizations
}

//OrganizationReference is a string containing an absolute or relative
//URL to another Organization.
type OrganizationReference string

const OrganizationURN = "urn:com:example:2.0:Organization"

//OrganizationResourceType provides the default structure which connects the
//Organization struct to its associated ResourceType.
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

//URN returns the IANA registered SCIM name for the Organization data structure
//and, together with ResourceType() identifies this code as implementing
//the Resource interface.
func (o Organization) URN() string {
	return OrganizationURN
}

//ResourceType returns the default structure describing the availability
//of the Organization resource and, together with URN() identifies this code as
//implementing the Resource interface.
func (o Organization) ResourceType() scim.ResourceType {
	return OrganizationResourceType
}

func (o Organization) MarshalJSON() ([]byte, error) {
	log.Info("Got to Organization's MarshalJSON()")
	type Alias Organization
	// alias := struct {
	// 	*Alias
	// }{
	// 	Alias: (*Alias)(u),
	// }

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	ap.RegisterAdditionalPropertiesExtension(json)
	return json.Marshal((Alias)(o))
}

func (o *Organization) UnmarshalJSON(data []byte) error {
	type Alias Organization
	// var a Alias
	// log.Info("Alias type: ", reflect.TypeOf(a).Name(), ", kind: ", reflect.TypeOf(a).Kind())
	// type Container struct {
	// 	Alias
	// }
	// alias := Container{
	// 	Alias: (Alias)(*u),
	// }
	// log.Info("alias type: ", reflect.TypeOf(alias).Name(), ", kind: ", reflect.TypeOf(alias).Kind())
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	ap.RegisterAdditionalPropertiesExtension(json)
	//return json.Unmarshal(data, &alias)
	return json.Unmarshal(data, (*Alias)(o))
}
