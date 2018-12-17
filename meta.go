package cprclient

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//Meta holds meta data about a Scim Resource
type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created,string"`
	LastModified time.Time `json:"lastModified,string"`

	Location string `json:"location"`
	Version  string `json:"version"`
}

//UnmarshalJSON turns the meta string values on the wire into the correct types
func (m *Meta) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		key := strings.ToLower(k)
		if key == "resourceType" {
			m.ResourceType = v
		} else if key == "created" {
			t, err := time.Parse(time.RFC3339, v)

			if err != nil {
				log.Warnf("Unable to parse \"created\" %s", v)
				return err
			}
			m.Created = t
		} else if key == "lastModified" {
			t, err := time.Parse(time.RFC3339, v)

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

//MarshalJSON converts a Meta into a json representation
func (m Meta) MarshalJSON() ([]byte, error) {
	meta := struct {
		ResourceType string `json:"resourceType"`
		Created      string `json:"created,string"`
		LastModified string `json:"lastModified,string"`

		Location string `json:"location"`
		Version  string `json:"version"`
	}{
		ResourceType: m.ResourceType,
		Created:      m.Created.Format(time.RFC3339),
		LastModified: m.LastModified.Format(time.RFC3339),
		Location:     m.Location,
		Version:      m.Version,
	}

	return json.Marshal(meta)
}
