package scim

import (
	"time"
)

//URI is a string that contains a URI according to:
//
//* RFC3986 (See: http://www.rfc-editor.org/info/rfc3986)
type URI string

//URN is a string that contains a URN according to:
//
//* RFC2141 (See: http://www.rfc-editor.org/info/rfc2119)
//* RFC3553 (See: http://www.rfc-editor.org/info/rfc3553)
type URN string

type commonAttributes struct {
	iD      string `json:"id" validation:"required"`
	schemas []URN  `json:"schemas" validation:"required"`
}

//ID returns the read-only valud of the SCIM resource's id field.
func (car commonAttributes) ID() string {
	return car.iD
}

func (ca commonAttributes) Schemas() []URN {
	return ca.schemas
}

//CommonAttributes is
//See: https://tools.ietf.org/html/rfc7643#section-3.1
type CommonAttributes1 struct {
	commonAttributes
	ExternalId string `json:"externalId"`
	meta       meta   `json:"meta"`
}

type meta struct {
	resourceType string    `json:"resourceType" validation:"required"`
	created      time.Time `json:"created" validation:"required"`
	lastModified time.Time `json:"lastModified" validation:"required"`
	version      string    `json:"version" validation:"required"`
	location     string    `json:"url" validation:"required"`
}

func (m meta) ResourceType() string {
	return m.resourceType
}
