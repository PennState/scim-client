package scim

import (
	"encoding/json"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Marshal(r resource) ([]byte, error) {
	//Unmarshal everything but struct's using the default (json.Marshal)
	t := dereference(reflect.TypeOf(r))
	if t.Kind() != reflect.Struct {
		return json.Marshal(r)
	}

	ap := r.getAdditionalProperties()

	//Iterate over the individual fields
	st := reflect.TypeOf(r).Elem()
	sv := reflect.ValueOf(r).Elem()
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		log.Debug("Field type: ", ft)
		n := name(ft)

		//Skip fields that are tagged with "-"
		if n == "-" {
			continue
		}

		fv := sv.Field(i)
		log.Debug("Value: ", fv)

		//Skip fields that are not both addressable and interfaceable
		log.Debug("Addressable: ", fv.CanAddr())
		log.Debug("Interfacable: ", fv.CanInterface())
		if !fv.CanAddr() || !fv.CanInterface() {
			continue
		}

		//Marshal all the other fields
		m, err := json.Marshal(fv.Interface())
		if err != nil {
			log.Error(err)
			continue
		}

		//Add them to the additional properties map as json.RawMessages
		log.Debug("Marshalled: ", string(m))
		ap[n] = json.RawMessage(m)
	}

	//Marshal the map that now contains all the struct's fields plus the
	//original additional properties
	return json.Marshal(r.getAdditionalProperties())
}

func Unmarshal(data []byte, v resource) error {
	t := dereference(reflect.TypeOf(v))
	if t.Kind() != reflect.Struct {
		return json.Unmarshal(data, v)
	}
	return unmarshalResource(data, v)
}

func dereference(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return dereference(t.Elem())
	}
	return t
}

func name(sf reflect.StructField) string {
	t := sf.Tag.Get("json")
	log.Debug("Tag: ", t)

	if t != "" {
		if idx := strings.Index(t, ","); idx != -1 {
			return t[:idx]
		}
		return t
	}

	return sf.Name
}

func unmarshalResource(data []byte, resource resource) error {
	var ap map[string]json.RawMessage
	err := json.Unmarshal(data, &ap)
	if err != nil {
		return err
	}

	err = unmarshalStruct(data, resource, ap)
	if err != nil {
		return err
	}

	resource.addAdditionalProperties(ap)
	return nil
}

func unmarshalStruct(data []byte, v interface{}, ap map[string]json.RawMessage) error {
	st := reflect.TypeOf(v).Elem()
	sv := reflect.ValueOf(v).Elem()

	//Iterate through the struct's fields
	for i := 0; i < st.NumField(); i++ {
		log.Debug("RawMessage count: ", len(ap))

		//Get the field's JSON name
		ft := st.Field(i)
		log.Debug("Field type: ", ft)
		n := name(ft)

		//Fields tagged with "-" should not be marshaled/unmarshaled so go
		//to the next field
		if n == "-" {
			continue
		}

		//Fields that can't be addressed or interfaced can't be set through
		//reflection so go to the next field
		fv := sv.Field(i)
		log.Debug("CanAddr: ", fv.CanAddr())
		log.Debug("CanInterface: ", fv.CanInterface())
		if !fv.CanAddr() || !fv.CanInterface() {
			continue
		}

		//Get a pointer to the value to pass so that we can either recurse
		//into anonymous structures or unmarshal it outright
		log.Debug("Field value before: ", fv)
		pv := fv.Addr().Interface()
		log.Debug("Field pointer: ", reflect.TypeOf(pv).Elem())

		//Anonymous struct's fields are part of the current JSON object
		//but we have to recurse into them to set their fields
		if ft.Anonymous {
			err := unmarshalStruct(data, pv, ap)
			if err != nil {
				return err
			}
			log.Debug("Anonymous field value after: ", fv)
			continue
		}

		//If a RawMessage doesn't exist for a given field name there's no
		//point wasting resources trying to unmarshal it
		rm, ok := ap[n]
		if !ok {
			log.Debug("No raw message found: ", n)
			continue
		}

		//Use the encoding/json version of Unmarshal to turn each
		//RawMessage into the individual fields
		err := json.Unmarshal(rm, pv)
		log.Debug("Field value after: ", fv)
		if err != nil {
			return err
		}

		//As fields are unmarshaled the struct's values will be filled
		//in and the RawMessageCount will decrease
		log.Debug("Struct value now: ", sv)
		delete(ap, n)
	}

	return nil
}
