package scim

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	ji "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

//Named identifies the implementing code as including a SCIM URN
type Named interface {
	URN() string
}

//ServerDiscoveryResource identifies the implementing code as a SCIM
//service discovery resource (Schemas, ServiceProviderConfig and
//ResourceTypes).
type ServerDiscoveryResource interface {
	Named
	ResourceType() ResourceType
}

type resource interface {
	addAdditionalProperties(additionalProperties map[string]json.RawMessage)
}

//Resource identifies the implementing code as a SCIM resource.  Resources
//defined by the specfication are User and Group.
//https://tools.ietf.org/html/rfc7643#section-3
type Resource interface {
	resource
	ServerDiscoveryResource
}

//Extension is the SCIM method for adapting and extending Resources
//https://tools.ietf.org/html/rfc7643#section-3.3
type Extension Named

//CommonAttributes describes the base common attributes of all Scim Resources
//https://tools.ietf.org/html/rfc7643#section-3.1
type CommonAttributes struct {
	ID                   string               `json:"id"`
	ExternalID           string               `json:"externalId"`
	Meta                 Meta                 `json:"meta"`
	Schemas              []string             `json:"schemas"`
	AdditionalProperties AdditionalProperties `json:",inline"`
	additionalProperties AdditionalProperties
}

type AdditionalProperties map[string]json.RawMessage

//Meta is a complex attribute containing resource metadata.
type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created,string"`
	LastModified time.Time `json:"lastModified,string"`
	Location     string    `json:"location"`
	Version      string    `json:"version"`
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
	Value string `json:"value"` //The attribute's significant value, e.g., email address, phone	numbeca.
}

func (ca *CommonAttributes) addAdditionalProperties(additionalProperties map[string]json.RawMessage) {
	ca.additionalProperties = additionalProperties
	log.Debugf("Saved additional properties: %v", ca.additionalProperties)
}

//AddExtension adds a new SCIM extension to a SCIM resource.  This method is
//purposely designed to return an error if the provided extension's URN is
//already a key in the additionalProperties map.
func (ca *CommonAttributes) AddExtension(extension Extension) error {
	if ca.HasExtension(extension) {
		return errors.New("Extension to be added already exists in resource - use UpdateExtension() instead")
	}
	err := ca.putExtension(extension)
	return err
}

//GetExtension retrieves a SCIM Extension from the additionalProperties map
//by the Extension's URN.
func (ca *CommonAttributes) GetExtension(extension Extension) error {
	name := extension.URN()
	err := json.Unmarshal(ca.additionalProperties[name], extension)
	return err
}

//GetExtensionURNs returns a list of the keys in the additionalProperties
//map that start with "urn:".  Clearly this is not a perfect way to
//guarantee that the RawMessage stored in that key is an extension.
func (ca *CommonAttributes) GetExtensionURNs() []string {
	keys := make([]string, 0, len(ca.additionalProperties))
	for key := range ca.additionalProperties {
		log.Debugf("Incoming key: %s", key)
		if strings.HasPrefix(key, "urn:") {
			log.Debugf("Saved key: %s", key)
			keys = append(keys, key)
		}
	}
	return keys
}

//HasExtension indicates whether the URN included with the passed
//Extension is a key in the additionalProperties map.
func (ca *CommonAttributes) HasExtension(extension Extension) bool {
	urn := extension.URN()
	return ca.HasExtensionByURN(urn)
}

//HasExtensionByURN indicates whether the passed URN string is a key in
//the additionalProperties map.
func (ca *CommonAttributes) HasExtensionByURN(urn string) bool {
	_, exists := ca.additionalProperties[urn]
	return exists
}

func (ca *CommonAttributes) putExtension(extension Extension) error {
	urn := extension.URN()
	var err error
	var rawMessage json.RawMessage
	rawMessage, err = json.Marshal(extension)

	if err == nil {
		ca.additionalProperties[urn] = rawMessage
	}

	return nil
}

//RemoveExtension deletes the RawMessage with the URN included with the
//passed SCIM Extension from the additionalProperties map.
func (ca *CommonAttributes) RemoveExtension(extension Extension) {
	ca.RemoveExtensionByURN(extension.URN())
}

//RemoveExtensionByURN deletes the RawMessage with the key matching
//the passed URN from the additionalProperties map.
func (ca *CommonAttributes) RemoveExtensionByURN(urn string) {
	delete(ca.additionalProperties, urn)
}

//UpdateExtension changes an existing SCIM extension already stored in a SCIM
//resource.  This method is purposely designed to return an error if the
//provided extension's URN is not a key in the additionalProperties map.
func (ca *CommonAttributes) UpdateExtension(extension Extension) error {
	if !ca.HasExtension(extension) {
		return errors.New("Extension to be updated does not exist in resource - use AddExtension() instead")
	}
	err := ca.putExtension(extension)
	return err
}

type internalJSONResource CommonAttributes

//Unmarshal attempts to decode the JSON provided in the passed data parameter
//into the Resource provided by the resource parameteca.  Any JSON properties
//(using the JSON schema vernacular) that are not included in the resource's
//fields are added to the additionalProperties map as RawMessages.  These
//additional parameters may included SCIM extensions which can be manipulated
//by methods of the Resource as well as properties that are simply cargo data.
//In both cases, the additionalProperties are maintained so that a client
//will (by default) return all the parameters that were originally provided.
//
//TODO: Convert to UnmarshalJSON interface implementation
func Unmarshal(data []byte, resource resource) error {
	log.Info("Got here too: ", reflect.TypeOf(resource).Elem())
	err := json.Unmarshal(data, resource)
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
	log.Debugf("Discovered additional properties: %v", additionalProperties)
	resource.addAdditionalProperties(additionalProperties)

	return err
}

// func (r resource) blah(data []byte) error {
// 	return nil
// }

// func (r resource) UnmarshalJSON(j []byte) error {
// 	return nil
// }

func (r *CommonAttributes) UnmarshalJSON1(j []byte) error {
	log.Info("Got here")
	err := ji.Unmarshal(j, r)
	return err
	// log.Info("Got here!")
	// log.Info("JSON: ", string(j))
	// log.Info("UnmarshalJSON type: ", reflect.TypeOf(*r))
	// t := reflect.TypeOf(r).Elem()
	// log.Info("Kind: ", t.Kind())
	// if t.Kind() != reflect.Struct {
	// 	return errors.New("blah")
	// }

	// var additionalProperties map[string]json.RawMessage
	// err := json.Unmarshal(j, &additionalProperties)
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	// for i := 0; i < t.NumField(); i++ {
	// 	f := t.Field(i)
	// 	n := jsonName(f)
	// 	log.Infof("     Name: %s, Type: %s, Kind: %s", n, f.Type, f.Type.Kind())
	// }

	// return nil
}

func (r CommonAttributes) MarshalJSON() ([]byte, error) {
	u, err := json.Marshal(internalJSONResource(r))
	if err != nil {
		return nil, err
	}

	var um map[string]json.RawMessage
	err = json.Unmarshal(u, &um)
	if err != nil {
		return nil, err
	}

	for k, v := range r.additionalProperties {
		um[k] = v
	}

	return json.Marshal(um)
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

//
//
// TODO - The code below this line is a temporary fix for a bug in SCIMple
//        that renders a LocalDateTime to JSON without the trailing Z (or
//        presumbably TZ).  Remove this code when SCIMple is fixed.
//
//

//UnmarshalJSON turns the meta string values on the wire into the correct types
func (m *Meta) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		key := strings.ToLower(k)
		if key == "resourcetype" {
			m.ResourceType = v
		} else if key == "created" {
			value := fixTimeZone(v)
			t, err := time.Parse(time.RFC3339, value)

			if err != nil {
				log.Warnf("Unable to parse \"created\" %s", v)
				return err
			}
			m.Created = t
		} else if key == "lastmodified" {
			value := fixTimeZone(v)
			t, err := time.Parse(time.RFC3339, value)

			if err != nil {
				log.Warnf("Unable to parse \"lastModified\" %s", v)
				return err
			}
			m.LastModified = t
		} else if key == "location" {
			m.Location = v
		} else if key == "version" {
			m.Version = v
		}
	}
	return nil
}

func fixTimeZone(in string) string {
	if strings.HasSuffix(in, "Z") {
		return in
	}

	return in + "Z"
}
