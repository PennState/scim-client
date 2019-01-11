package scim

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//https://tools.ietf.org/html/rfc7643#section-3
type resource interface {
	addAdditionalProperties(additionalProperties map[string]json.RawMessage)
}

//Resource describes the base common attributes of all Scim Resources
//https://tools.ietf.org/html/rfc7643#section-3.1
type Resource struct {
	ID                   string   `json:"id"`
	ExternalID           string   `json:"externalId"`
	Meta                 Meta     `json:"meta"`
	Schemas              []string `json:"schemas"`
	additionalProperties map[string]json.RawMessage
}

//Meta is a complex attribute containing resource metadata.
type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created,string"`
	LastModified time.Time `json:"lastModified,string"`
	Location     string    `json:"location"`
	Version      string    `json:"version"`
}

//Extension is the SCIM method for adapting and extending Resources
//https://tools.ietf.org/html/rfc7643#section-3.3
type Extension interface {
	GetUrn() string
}

//Multivalued attributes contain a list of elements using the JSON array format defined in Section 5 of [RFC7159].
//https://tools.ietf.org/html/rfc7643#section-2.4
type Multivalued struct {
	Type      string `json:"type"`    //Type is a label indicating the attribute's function; e.g., 'work' or 'home'.
	Display   string `json:"display"` //Display is a  human readable name, primarily used for display purposes. READ-ONLY.
	Primary   bool   `json:"primary"` //Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.
	Reference string `json:"$ref"`    //Reference is the reference URI of a target resource, if the attribute is a reference.
}

//StringMultivalued provides a base structure for simple string multi-valued attributes.
type StringMultivalued struct {
	Multivalued
	Value string `json:"value"` //The attribute's significant value, e.g., email address, phone	number.
}

func (cr *Resource) addAdditionalProperties(additionalProperties map[string]json.RawMessage) {
	cr.additionalProperties = additionalProperties
}

//AddExtension adds a new SCIM extension to a SCIM resource.  This method is
//purposely designed to return an error if the provided extension's URN is
//already a key in the additionalProperties map.
func (cr *Resource) AddExtension(extension Extension) error {
	if cr.HasExtension(extension) {
		return errors.New("Extension to be added already exists in resource - use UpdateExtension() instead")
	}
	err := cr.putExtension(extension)
	return err
}

//GetExtension retrieves a SCIM Extension from the additionalProperties map
//by the Extension's URN.
func (cr *Resource) GetExtension(extension Extension) error {
	name := extension.GetUrn()
	err := json.Unmarshal(cr.additionalProperties[name], extension)
	return err
}

//GetExtensionURNs returns a list of the keys in the additionalProperties
//map that start with "urn:".  Clearly this is not a perfect way to
//guarantee that the RawMessage stored in that key is an extension.
func (cr *Resource) GetExtensionURNs() []string {
	keys := make([]string, 0, len(cr.additionalProperties))
	for key := range cr.additionalProperties {
		if strings.HasPrefix(key, "urn:") {
			keys = append(keys, key)
		}
	}
	return keys
}

//HasExtension indicates whether the URN included with the passed
//Extension is a key in the additionalProperties map.
func (cr *Resource) HasExtension(extension Extension) bool {
	urn := extension.GetUrn()
	return cr.HasExtensionByURN(urn)
}

//HasExtensionByURN indicates whether the passed URN string is a key in
//the additionalProperties map.
func (cr *Resource) HasExtensionByURN(urn string) bool {
	_, exists := cr.additionalProperties[urn]
	return exists
}

func (cr *Resource) putExtension(extension Extension) error {
	urn := extension.GetUrn()
	var err error
	var rawMessage json.RawMessage
	rawMessage, err = json.Marshal(extension)

	if err == nil {
		cr.additionalProperties[urn] = rawMessage
	}

	return nil
}

//RemoveExtension deletes the RawMessage with the URN included with the
//passed SCIM Extension from the additionalProperties map.
func (cr *Resource) RemoveExtension(extension Extension) {
	cr.RemoveExtensionByURN(extension.GetUrn())
}

//RemoveExtensionByURN deletes the RawMessage with the key matching
//the passed URN from the additionalProperties map.
func (cr *Resource) RemoveExtensionByURN(urn string) {
	delete(cr.additionalProperties, urn)
}

//UpdateExtension changes an existing SCIM extension already stored in a SCIM
//resource.  This method is purposely designed to return an error if the
//provided extension's URN is not a key in the additionalProperties map.
func (cr *Resource) UpdateExtension(extension Extension) error {
	if !cr.HasExtension(extension) {
		return errors.New("Extension to be updated does not exist in resource - use AddExtension() instead")
	}
	err := cr.putExtension(extension)
	return err
}

//Unmarshal attempts to decode the JSON provided in the passed data parameter
//into the Resource provided by the resource parameter.  Any JSON properties
//(using the JSON schema vernacular) that are not included in the resource's
//fields are added to the additionalProperties map as RawMessages.  These
//additional parameters may included SCIM extensions which can be manipulated
//by methods of the Resource as well as properties that are simply cargo data.
//In both cases, the additionalProperties are maintained so that a client
//will (by default) return all the parameters that were originally provided.
func Unmarshal(data []byte, resource resource) error {
	var err error
	err = json.Unmarshal(data, resource)

	if err != nil {
		log.Error(err)
		return err
	}

	var additionalProperties map[string]json.RawMessage
	err = json.Unmarshal(data, &additionalProperties)
	if err != nil {
		log.Error(err)
		return err
	}

	t := reflect.TypeOf(resource).Elem()
	removeKnownProperties(additionalProperties, t)
	resource.addAdditionalProperties(additionalProperties)

	return err
}

func jsonName(sf reflect.StructField) string {
	t := sf.Tag.Get("json")
	log.Debugf("Tag: %s", t)

	if t != "" {
		if idx := strings.Index(t, ","); idx != -1 {
			return t[:idx]
		}
		return t
	}

	return sf.Name
}

func removeKnownProperties(additionalProperties map[string]json.RawMessage, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		n := jsonName(f)
		if strings.HasSuffix(f.Type.Name(), n) && f.Type.Kind() == reflect.Struct {
			log.Debugf("Recursing into: %s", n)
			removeKnownProperties(additionalProperties, f.Type)
		} else {
			log.Debugf("Name: %s, Type: %s, Kind: %s", n, f.Type, f.Type.Kind())
			delete(additionalProperties, n)
		}
	}
}
