package scim

import (
	"encoding/json"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Marshal(v resource) ([]byte, error) {
	return nil, nil
}

func Unmarshal(data []byte, v resource) error {
	t := dereference(reflect.TypeOf(v))
	if t.Kind() != reflect.Struct {
		return json.Unmarshal(data, v)
	}
	log.Info("Kind: ", t.Kind())
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
	log.Debugf("Tag: %s", t)

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
	for i := 0; i < st.NumField(); i++ {
		log.Info("RawMessage count: ", len(ap))

		ft := st.Field(i)
		log.Info("Field type: ", ft)
		n := name(ft)

		if n == "-" {
			continue
		}

		fv := sv.Field(i)
		log.Info("CanAddr: ", fv.CanAddr())
		log.Info("CanInterface: ", fv.CanInterface())
		if !fv.CanAddr() || !fv.CanInterface() {
			continue
		}

		log.Info("Field value before: ", fv)
		pv := fv.Addr().Interface()
		log.Info("Field pointer: ", reflect.TypeOf(pv).Elem())

		// if ft.Type.Kind() == reflect.Struct && strings.HasSuffix(ft.Type.String(), ft.Name) {
		if ft.Anonymous {
			err := unmarshalStruct(data, pv, ap)
			if err != nil {
				return err
			}
			log.Info("Anonymous field value after: ", fv)
			continue
		}

		rm, ok := ap[n]
		if !ok {
			log.Info("No raw message found: ", n)
			continue
		}

		log.Info("PkgPath: ", ft.PkgPath)
		log.Info("Field name: ", ft.Name)
		if ft.PkgPath != "" || ft.Name == "additionalProperties" {
			continue
		}

		err := json.Unmarshal(rm, pv)
		log.Info("Field value after: ", fv)
		if err != nil {
			return err
		}
		log.Info("Struct value now: ", sv)
		delete(ap, n)
	}

	return nil
}
