package blackbox

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/PennState/proctor/pkg/goldenfile"
	"github.com/PennState/scim-client/examples"
	"github.com/PennState/scim-client/pkg/scim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//nolint:gochecknoglobals
var enterpriseUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		ID:         "2819c223-7f76-453a-919d-413861904646",
		ExternalID: "701984",
		AdditionalProperties: map[string]json.RawMessage{
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": json.RawMessage([]byte("{\"employeeNumber\":\"701984\",\"costCenter\":\"4130\",\"organization\":\"Universal Studios\",\"division\":\"Theme Park\",\"department\":\"Tour Operations\",\"manager\":{\"value\":\"26118915-6090-4610-87e4-49d8ca9f808d\",\"$ref\":\"../Users/26118915-6090-4610-87e4-49d8ca9f808d\",\"displayName\":\"John Smith\"}}")),
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

//nolint:gochecknoglobals
var fullUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
		},
		ID:                   "2819c223-7f76-453a-919d-413861904646",
		ExternalID:           "701984",
		AdditionalProperties: map[string]json.RawMessage{},
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
	Name: &scim.Name{
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

//nolint:gochecknoglobals
var group scim.Group = scim.Group{
	CommonAttributes: scim.CommonAttributes{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:Group",
		},
		ID:                   "e9e30dba-f08f-4109-8486-d5c6a331660a",
		AdditionalProperties: map[string]json.RawMessage{},
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

//nolint:gochecknoglobals
var minimalUser scim.User = scim.User{
	CommonAttributes: scim.CommonAttributes{
		Schemas:              []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		ID:                   "2819c223-7f76-453a-919d-413861904646",
		AdditionalProperties: map[string]json.RawMessage{},
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

//nolint:gochecknoglobals
var organization examples.Organization = examples.Organization{
	CommonAttributes: scim.CommonAttributes{
		Schemas:              []string{"urn:com:example:2.0:Organization"},
		ID:                   "430beb5c-a361-4c04-b308-2845789a496e",
		AdditionalProperties: map[string]json.RawMessage{},
		Meta: scim.Meta{
			ResourceType: "Organization",
			Created:      parseTime("2010-01-23T04:56:22Z"),
			LastModified: parseTime("2011-05-13T04:42:34Z"),
			Version:      "W/\"3694e05e9dff590\"",
			Location:     "https://example.com/v2/Organizations/430beb5c-a361-4c04-b308-2845789a496e",
		},
	},
	Name:   "Tour Promotion",
	Type:   "Department",
	Parent: "../Organizations/4a7741a3-a436-4a52-a6d5-149e6c1b9578",
	Children: []examples.OrganizationReference{
		"../Organizations/7eb59c46-35a4-4443-b8c1-5de8be88f973",
		"../Organizations/66506f29-8c44-414e-b52d-a993b94f370c",
		"../Organizations/0a365d4f-10e5-45c5-ae05-ee5184b59627",
	},
}

//nolint:gochecknoglobals
var cases = []struct {
	Name         string
	GoldenFile   string
	ZeroResource scim.Resource
	TestResource interface{}
}{
	{"Enterprise user", "enterpriseuser.json", &scim.User{}, &enterpriseUser},
	// {"Group", "group.json", &scim.Group{}, &group},
	{"Full user", "fulluser.json", &scim.User{}, &fullUser},
	{"Minimal user", "minimaluser.json", &scim.User{}, &minimalUser},
	// {"Organization", "organization.json", &examples.Organization{}, &organization},
}

func TestResourceMarshaling(t *testing.T) {
	for idx := range cases {
		c := cases[idx]
		t.Run(c.Name, func(t *testing.T) {
			fp := goldenfile.GetDefaultFilePath(c.GoldenFile)
			data, err := json.Marshal(c.TestResource)
			require.NoError(t, err)
			goldenfile.AssertJSONEq(t, fp, string(data))
		})
	}
}

func TestResourceUnmarshaling(t *testing.T) {
	for idx := range cases {
		c := cases[idx]
		t.Run(c.Name, func(t *testing.T) {
			fp := goldenfile.GetDefaultFilePath(c.GoldenFile)
			data, err := ioutil.ReadFile(fp)
			require.NoError(t, err)
			res := c.ZeroResource
			err = json.Unmarshal(data, res)
			require.NoError(t, err)
			assert.Equal(t, c.TestResource, res)
		})
	}
}
