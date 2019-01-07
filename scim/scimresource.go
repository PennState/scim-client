package scim

import (
	"github.com/PennState/golang_scimclient/schema"
)
//ScimResource describes the base shared elements of all Scim Resources
type ScimResource struct {
	Meta schema.Meta `json:"meta"`

	ID string `json:"id"`

	ExternalID string `json:"externalId"`

	// TODO - Figure out JAXB equivalent of JsonAnyGetter and JsonAnySetter
	// (XmlElementAny?)
	//  private Map<String, ScimExtension> extensions = new HashMap<String, ScimExtension>();

	ScimExtensions []interface{}

	// private String baseUrn;

}
