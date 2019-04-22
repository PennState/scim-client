package scim

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

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

type JSONResource CommonAttributes

func Marshal(resource resource) ([]byte, error) {
	return json.Marshal(resource)
}

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
	var additionalProperties map[string]json.RawMessage
	err := json.Unmarshal(data, &additionalProperties)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("----- Additional properties -----")
	for k, v := range additionalProperties {
		log.Info("Key: ", k, ", Value: ", string(v))
	}

	sv := reflect.ValueOf(resource).Elem()
	log.Info("Struct value: ", sv)
	log.Info("Struct value field count: ", sv.NumField())
	st := reflect.TypeOf(resource).Elem()
	log.Info("Struct type: ", st)
	log.Info("Struct type field count: ", st.NumField())
	setReflectedStruct(sv, st, additionalProperties)
	// for i := 0; i < sv.NumField(); i++ {
	// 	fv := sv.Field(i)
	// 	log.Infof("Field value -- Type: %s, Kind: %s, Value: %v", fv.Type(), fv.Kind(), fv)
	// 	ft := st.Field(i)
	// 	log.Infof("Field type -- Name: %s, Type: %s, Kind: %s, Field: %v", ft.Name, ft.Type, ft.Type.Kind(), ft)
	// 	n := jsonName(ft)
	// 	log.Info("Effective name: ", n)
	// 	rm, ok := additionalProperties[n]
	// 	if !ok {
	// 		log.Warn("Field not found: ", n)
	// 		continue
	// 	}
	// 	log.Info("Incoming: ", string(rm))
	// 	setReflectedField(fv, ft, additionalProperties)

	// if ok && v.Kind() == reflect.Struct {
	// 	log.Infof("Recursing into: %s", n)
	// 	removeKnownProperties(additionalProperties, v)
	// }
	// if ok {
	// 	value := reflect.New(f.Type)
	// 	_ = json.Unmarshal(rm, &value)
	// 	v.Set(value)
	// }
	// delete(additionalProperties, n)
	// 	log.Info("Altered struct value: ", sv)
	// }

	// removeKnownProperties(additionalProperties, s)
	resource.addAdditionalProperties(additionalProperties)

	return err
}

func setReflectedStruct(sv reflect.Value, st reflect.Type, ap map[string]json.RawMessage) {
	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		log.Infof("Field value -- Type: %s, Kind: %s, Value: %v", fv.Type(), fv.Kind(), fv)
		ft := st.Field(i)
		log.Infof("Field type -- Name: %s, Type: %s, Kind: %s, Field: %v", ft.Name, ft.Type, ft.Type.Kind(), ft)
		// n := jsonName(ft)
		// log.Info("Effective name: ", n)
		// rm, ok := ap[n]
		// if !ok {
		// 	log.Warn("Field not found: ", n)
		// 	continue
		// }
		// log.Info("Incoming: ", string(rm))
		setReflectedField(fv, ft, ap)

		// if ok && v.Kind() == reflect.Struct {
		// 	log.Infof("Recursing into: %s", n)
		// 	removeKnownProperties(additionalProperties, v)
		// }
		// if ok {
		// 	value := reflect.New(f.Type)
		// 	_ = json.Unmarshal(rm, &value)
		// 	v.Set(value)
		// }
		log.Info("Altered struct value: ", sv)
	}
}

func setReflectedField(fv reflect.Value, ft reflect.StructField, ap map[string]json.RawMessage) {
	//Embedded structs
	if ft.Type.Kind() == reflect.Struct && strings.HasSuffix(ft.Type.String(), ft.Name) {
		//TODO: recursion goes here
		log.Info("Descending into embedded struct: ", ft.Name)
		setReflectedStruct(fv, ft.Type, ap)
		return
	}

	n := jsonName(ft)
	log.Info("Effective name: ", n)
	rm, ok := ap[n]
	if !ok {
		log.Warn("Field not found: ", n)
		return
	}
	log.Info("Incoming: ", string(rm))

	//value := reflect.New(fv.Type())
	//p := unsafe.Pointer(fv.Addr().Pointer())
	//p := unsafe.Pointer(value.Addr().Pointer())
	value, err := valueInstance(fv)
	if err != nil {
		return
	}
	_ = json.Unmarshal(rm, &value)
	log.Info("Unmarshaled raw message: ", value)

	log.Info("Mutable: ", fv.CanSet())
	fv.Set(reflect.ValueOf(value))

	log.Info("Deleting additional property: ", n)
	delete(ap, n)
}

func valueInstance(value reflect.Value) (interface{}, error) {
	switch value.Type().Kind() {
	case reflect.Bool:
		var b bool
		return b, nil
	case reflect.Float32, reflect.Float64:
		var f float64
		return f, nil
	case reflect.String:
		var s string
		return s, nil
	}

	return nil, errors.New("Couldn't do it")
}

// func (r CommonAttributes) MarshalJSON() ([]byte, error) {
// 	u, err := json.Marshal(JSONResource(r))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var um map[string]json.RawMessage
// 	err = json.Unmarshal(u, &um)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, v := range r.additionalProperties {
// 		um[k] = v
// 	}

// 	return json.Marshal(um)
// }

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

func removeKnownProperties(additionalProperties map[string]json.RawMessage, s reflect.Value) {
	for i := 0; i < s.NumField(); i++ {
		v := s.Field(i)
		f := v.Type().Field(i)
		n := jsonName(f)
		log.Infof("Type -- Name: %s, Type: %s, Kind: %s, Field: %v", v.Type().Name(), v.Type(), v.Type().Kind(), f)
		log.Infof("Value -- Name: %s, Type: %s, Kind: %s, Value: %v", n, v.Type(), v.Kind(), v)
		rm, ok := additionalProperties[n]
		if ok && v.Kind() == reflect.Struct {
			log.Infof("Recursing into: %s", n)
			removeKnownProperties(additionalProperties, v)
			delete(additionalProperties, n)
			continue
		}
		if ok {
			value := reflect.New(f.Type)
			_ = json.Unmarshal(rm, &value)
			v.Set(value)
		}

		delete(additionalProperties, n)
	}
	// for i := 0; i < t.NumField(); i++ {
	// 	f := t.Field(i)
	// 	n := jsonName(f)
	// 	if strings.HasSuffix(f.Type.Name(), n) && f.Type.Kind() == reflect.Struct {
	// 		log.Debugf("Recursing into: %s", n)
	// 		removeKnownProperties(additionalProperties, f.Type)
	// 	} else {
	// 		log.Debugf("Name: %s, Type: %s, Kind: %s", n, f.Type, f.Type.Kind())
	// 		delete(additionalProperties, n)
	// 	}
	// }
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
