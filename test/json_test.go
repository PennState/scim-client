package blackbox

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/PennState/golang_scimclient/pkg/scim"
	"github.com/PennState/proctor/pkg/goldenfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

var enterpriseUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		ID:         "2819c223-7f76-453a-919d-413861904646",
		ExternalID: "701984",
		AdditionalProperties: map[string]json.RawMessage{
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": json.RawMessage{0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x3a, 0x20, 0x22, 0x37, 0x30, 0x31, 0x39, 0x38, 0x34, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x63, 0x6f, 0x73, 0x74, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x3a, 0x20, 0x22, 0x34, 0x31, 0x33, 0x30, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x20, 0x22, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x61, 0x6c, 0x20, 0x53, 0x74, 0x75, 0x64, 0x69, 0x6f, 0x73, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x20, 0x22, 0x54, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x50, 0x61, 0x72, 0x6b, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x54, 0x6f, 0x75, 0x72, 0x20, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0x3a, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x32, 0x36, 0x31, 0x31, 0x38, 0x39, 0x31, 0x35, 0x2d, 0x36, 0x30, 0x39, 0x30, 0x2d, 0x34, 0x36, 0x31, 0x30, 0x2d, 0x38, 0x37, 0x65, 0x34, 0x2d, 0x34, 0x39, 0x64, 0x38, 0x63, 0x61, 0x39, 0x66, 0x38, 0x30, 0x38, 0x64, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x24, 0x72, 0x65, 0x66, 0x22, 0x3a, 0x20, 0x22, 0x2e, 0x2e, 0x2f, 0x55, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x32, 0x36, 0x31, 0x31, 0x38, 0x39, 0x31, 0x35, 0x2d, 0x36, 0x30, 0x39, 0x30, 0x2d, 0x34, 0x36, 0x31, 0x30, 0x2d, 0x38, 0x37, 0x65, 0x34, 0x2d, 0x34, 0x39, 0x64, 0x38, 0x63, 0x61, 0x39, 0x66, 0x38, 0x30, 0x38, 0x64, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x4a, 0x6f, 0x68, 0x6e, 0x20, 0x53, 0x6d, 0x69, 0x74, 0x68, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7d},
		},
		Meta: scim.Meta{
			ResourceType: "User",
			Created:      parseTime("2010-01-23T04:56:22Z"),
			LastModified: parseTime("2011-05-13T04:42:34Z"),
			Version:      "W/\"a330bc54f0671c9\"",
			Location:     "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
		},
	},
	UserName: "bjensen@example.com",
}

var fullUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
		},
		ID:         "2819c223-7f76-453a-919d-413861904646",
		ExternalID: "701984",
		Meta: scim.Meta{
			ResourceType: "User",
			Created:      parseTime("2010-01-23T04:56:22Z"),
			LastModified: parseTime("2011-05-13T04:42:34Z"),
			Version:      "W/\"a330bc54f0671c9\"",
			Location:     "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
		},
	},
	Active: true,
	Addresses: []scim.Address{
		scim.Address{
			Multivalued: scim.Multivalued{
				Type:    "work",
				Primary: true,
			},
			StreetAddress: "100 Universal City Plaza",
			Locality:      "Hollywood",
			Region:        "CA",
			PostalCode:    "91608",
			Country:       "USA",
			Formatted:     "100 Universal City Plaza\nHollywood, CA 91608 USA",
		},
		scim.Address{
			Multivalued: scim.Multivalued{
				Type: "home",
			},
			StreetAddress: "456 Hollywood Blvd",
			Locality:      "Hollywood",
			Region:        "CA",
			PostalCode:    "91608",
			Country:       "USA",
			Formatted:     "456 Hollywood Blvd\nHollywood, CA 91608 USA",
		},
	},
	DisplayName: "Babs Jensen",
	Emails: []scim.Email{
		scim.Email{
			Multivalued: scim.Multivalued{
				Type:    "work",
				Primary: true,
			},
			Value: "bjensen@example.com",
		},
		scim.Email{
			Multivalued: scim.Multivalued{
				Type: "home",
			},
			Value: "babs@jensen.org",
		},
	},
	Groups: []scim.GroupRef{
		scim.GroupRef{
			Multivalued: scim.Multivalued{
				Reference: "https://example.com/v2/Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
				Display:   "Tour Guides",
			},
			Value: "e9e30dba-f08f-4109-8486-d5c6a331660a",
		},
		scim.GroupRef{
			Multivalued: scim.Multivalued{
				Reference: "https://example.com/v2/Groups/fc348aa8-3835-40eb-a20b-c726e15c55b5",
				Display:   "Employees",
			},
			Value: "fc348aa8-3835-40eb-a20b-c726e15c55b5",
		},
		scim.GroupRef{
			Multivalued: scim.Multivalued{
				Reference: "https://example.com/v2/Groups/71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
				Display:   "US Employees",
			},
			Value: "71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
		},
	},
	IMs: []scim.IM{
		scim.IM{
			Multivalued: scim.Multivalued{
				Type: "aim",
			},
			Value: "someaimhandle",
		},
	},
	Locale: "en-US",
	Name: scim.Name{
		Formatted:       "Ms. Barbara J Jensen, III",
		FamilyName:      "Jensen",
		GivenName:       "Barbara",
		MiddleName:      "Jane",
		HonorificPrefix: "Ms.",
		HonorificSuffix: "III",
	},
	NickName: "Babs",
	Password: "t1meMa$heen",
	PhoneNumbers: []scim.PhoneNumber{
		scim.PhoneNumber{
			Multivalued: scim.Multivalued{
				Type: "work",
			},
			Value: "555-555-5555",
		},
		scim.PhoneNumber{
			Multivalued: scim.Multivalued{
				Type: "mobile",
			},
			Value: "555-555-4444",
		},
	},
	Photos: []scim.Photo{
		scim.Photo{
			Multivalued: scim.Multivalued{
				Type: "photo",
			},
			Value: "https://photos.example.com/profilephoto/72930000000Ccne/F",
		},
		scim.Photo{
			Multivalued: scim.Multivalued{
				Type: "thumbnail",
			},
			Value: "https://photos.example.com/profilephoto/72930000000Ccne/T",
		},
	},
	PreferredLanguage: "en-US",
	ProfileURL:        "https://login.example.com/bjensen",
	Timezone:          "America/Los_Angeles",
	Title:             "Tour Guide",
	UserName:          "bjensen@example.com",
	UserType:          "Employee",
	X509Certificates: []scim.X509Certificate{
		scim.X509Certificate{
			Value: "MIIDQzCCAqygAwIBAgICEAAwDQYJKoZIhvcNAQEFBQAwTjELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFDASBgNVBAoMC2V4YW1wbGUuY29tMRQwEgYDVQQDDAtleGFtcGxlLmNvbTAeFw0xMTEwMjIwNjI0MzFaFw0xMjEwMDQwNjI0MzFaMH8xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRQwEgYDVQQKDAtleGFtcGxlLmNvbTEhMB8GA1UEAwwYTXMuIEJhcmJhcmEgSiBKZW5zZW4gSUlJMSIwIAYJKoZIhvcNAQkBFhNiamVuc2VuQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7Kr+Dcds/JQ5GwejJFcBIP682X3xpjis56AK02bc1FLgzdLI8auoR+cC9/Vrh5t66HkQIOdA4unHh0AaZ4xL5PhVbXIPMB5vAPKpzz5iPSi8xO8SL7I7SDhcBVJhqVqr3HgllEG6UClDdHO7nkLuwXq8HcISKkbT5WFTVfFZzidPl8HZ7DhXkZIRtJwBweq4bvm3hM1Os7UQH05ZS6cVDgweKNwdLLrT51ikSQG3DYrl+ft781UQRIqxgwqCfXEuDiinPh0kkvIi5jivVu1Z9QiwlYEdRbLJ4zJQBmDrSGTMYn4lRc2HgHO4DqB/bnMVorHB0CC6AV1QoFK4GPe1LwIDAQABo3sweTAJBgNVHRMEAjAAMCwGCWCGSAGG+EIBDQQfFh1PcGVuU1NMIEdlbmVyYXRlZCBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQU8pD0U0vsZIsaA16lL8En8bx0F/gwHwYDVR0jBBgwFoAUdGeKitcaF7gnzsNwDx708kqaVt0wDQYJKoZIhvcNAQEFBQADgYEAA81SsFnOdYJtNg5Tcq+/ByEDrBgnusx0jloUhByPMEVkoMZ3J7j1ZgI8rAbOkNngX8+pKfTiDz1RC4+dx8oU6Za+4NJXUjlL5CvV6BEYb1+QAEJwitTVvxB/A67g42/vzgAtoRUeDov1+GFiBZ+GNF/cAYKcMtGcrs2i97ZkJMo=",
		},
	},
}

var group scim.Group = scim.Group{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:Group",
		},
		ID: "e9e30dba-f08f-4109-8486-d5c6a331660a",
		Meta: scim.Meta{
			ResourceType: "Group",
			Created:      parseTime("2010-01-23T04:56:22Z"),
			LastModified: parseTime("2011-05-13T04:42:34Z"),
			Version:      "W/\"3694e05e9dff592\"",
			Location:     "https://example.com/v2/Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
		},
	},
	DisplayName: "Tour Guides",
	Members: []scim.MemberRef{
		scim.MemberRef{
			Multivalued: scim.Multivalued{
				Reference: "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
				Display:   "Babs Jensen",
			},
			Value: "2819c223-7f76-453a-919d-413861904646",
		},
		scim.MemberRef{
			Multivalued: scim.Multivalued{
				Reference: "https://example.com/v2/Users/902c246b-6245-4190-8e05-00816be7344a",
				Display:   "Mandy Pepperidge",
			},
			Value: "902c246b-6245-4190-8e05-00816be7344a",
		},
	},
}

var minimalUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		ID:      "2819c223-7f76-453a-919d-413861904646",
		Meta: scim.Meta{
			ResourceType: "User",
			Created:      parseTime("2010-01-23T04:56:22Z"),
			LastModified: parseTime("2011-05-13T04:42:34Z"),
			Version:      "W/\"3694e05e9dff590\"",
			Location:     "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
		},
	},
	UserName: "bjensen@example.com",
}

func TestResourceMarshalingAndUnmarshaling(t *testing.T) {
	tests := []struct {
		Name         string
		GoldenFile   string
		ZeroResource scim.Resource
		TestResource interface{}
	}{
		{"Enterprise user", "enterpriseuser.json", &scim.User{}, &enterpriseUser},
		{"Group", "group.json", &scim.Group{}, &group},
		{"Full user", "fulluser.json", &scim.User{}, &fullUser},
		{"Minimal user", "minimaluser.json", &scim.User{}, &minimalUser},
	}
	for _, test := range tests {
		fp := goldenfile.GetDefaultFilePath(test.GoldenFile)
		t.Run("Marshaling "+test.Name, func(t *testing.T) {
			data, err := json.Marshal(test.TestResource)
			require.NoError(t, err)
			goldenfile.AssertJSONEq(t, fp, string(data))
		})
		t.Run("Unmarshaling "+test.Name, func(t *testing.T) {
			data, err := ioutil.ReadFile(fp)
			require.NoError(t, err)
			res := test.ZeroResource
			err = json.Unmarshal(data, res)
			require.NoError(t, err)
			assert.Equal(t, test.TestResource, res)
		})
	}
}
