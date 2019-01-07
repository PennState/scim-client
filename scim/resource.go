package scim

import (
	"encoding/json"
	"time"

	"github.com/PennState/golang_scimclient/schema"
)

//ScimResource describes the base common attributes of all Scim Resources
//https://tools.ietf.org/html/rfc7643#section-3.1
type ScimResource struct {
	Meta Meta `json:"meta"`

	ID string `json:"id"`

	ExternalID string `json:"externalId"`

	// TODO - Figure out JAXB equivalent of JsonAnyGetter and JsonAnySetter
	// (XmlElementAny?)
	//  private Map<String, ScimExtension> extensions = new HashMap<String, ScimExtension>();

	ScimExtensions []interface{}

	// private String baseUrn;

}

//Meta is a complex attribute containing resource metadata.
type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created,string"`
	LastModified time.Time `json:"lastModified,string"`

	Location string `json:"location"`
	Version  string `json:"version"`
}

//Multivalued attributes contain a list of elements using the JSON array format defined in Section 5 of [RFC7159].
//https://tools.ietf.org/html/rfc7643#section-2.4
type Multivalued struct {
	//Type is a label indicating the attribute's function; e.g., 'work' or 'home'.
	Type string `json:"type"`

	//Display is a  human readable name, primarily used for display purposes. READ-ONLY.
	Display string `json:"display"`

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.
	Primary bool `json:"primary"`

	//Reference is the reference URI of a target resource, if the attribute is a reference.
	Reference string `json:"$ref"`
}

//Provides a base structure for simple string multi-valued attributes.
type StringMultivalued struct {
	Multivalued

	//The attribute's significant value, e.g., email address, phone	number.
	Value string `json:"value"`
}

func Unmarshal(data []byte, resource *ScimResource, extensions ...schema.ScimExtension) error {
	var error error

	error = json.Unmarshal(data, resource)

	if error == nil && len(extensions) > 0 {
		//var attributes map[string]interface{}

	}

	return error
}