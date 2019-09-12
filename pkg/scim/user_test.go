package scim

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// https://tools.ietf.org/html/rfc7643#section-8.1
func TestMinimalUserUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	const minimalUser = `
		{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
		"id": "2819c223-7f76-453a-919d-413861904646",
		"userName": "bjensen@example.com",
		"meta": {
			"resourceType": "User",
			"created": "2010-01-23T04:56:22Z",
			"lastModified": "2011-05-13T04:42:34Z",
			"version": "W/3694e05e9dff590",
			"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
		}
	}`

	var user User
	error := Unmarshal([]byte(minimalUser), &user)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(user.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
	assert.Equal(user.UserName, "bjensen@example.com")

	assert.Equal(user.Meta.ResourceType, "User")
	assert.Equal(user.Meta.Created, time.Date(2010, time.January, 23, 4, 56, 22, 0, time.UTC))
	assert.Equal(user.Meta.LastModified, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC))
	assert.Equal(user.Meta.Version, "W/3694e05e9dff590")
	assert.Equal(user.Meta.Location, "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646")
}

// https://tools.ietf.org/html/rfc7643#section-8.2
func TestFullUserUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	const fullUser = `
	{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
		"id": "2819c223-7f76-453a-919d-413861904646",
		"externalId": "701984",
		"userName": "bjensen@example.com",
		"name": {
			"formatted": "Ms. Barbara J Jensen, III",
			"familyName": "Jensen",
			"givenName": "Barbara",
			"middleName": "Jane",
			"honorificPrefix": "Ms.",
			"honorificSuffix": "III"
		},
		"displayName": "Babs Jensen",
		"nickName": "Babs",
		"profileUrl": "https://login.example.com/bjensen",
		"emails": [
			{
				"value": "bjensen@example.com",
				"type": "work",
				"primary": true
			},
			{
				"value": "babs@jensen.org",
				"type": "home"
			}
		],
		"addresses": [
			{
				"type": "work",
				"streetAddress": "100 Universal City Plaza",
				"locality": "Hollywood",
				"region": "CA",
				"postalCode": "91608",
				"country": "USA",
				"formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
				"primary": true
			},
			{
				"type": "home",
				"streetAddress": "456 Hollywood Blvd",
				"locality": "Hollywood",
				"region": "CA",
				"postalCode": "91608",
				"country": "USA",
				"formatted": "456 Hollywood Blvd\nHollywood, CA 91608 USA"
			}
		],
		"phoneNumbers": [
			{
				"value": "555-555-5555",
				"type": "work"
			},
			{
				"value": "555-555-4444",
				"type": "mobile"
			}
		],
		"ims": [
			{
				"value": "someaimhandle",
				"type": "aim"
			}
		],
		"photos": [
			{
				"value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
				"type": "photo"
			},
			{
				"value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
				"type": "thumbnail"
			}
		],
		"userType": "Employee",
		"title": "Tour Guide",
		"preferredLanguage": "en-US",
		"locale": "en-US",
		"timezone": "America/Los_Angeles",
		"active": true,
		"password": "t1meMa$heen",
		"groups": [
			{
				"value": "e9e30dba-f08f-4109-8486-d5c6a331660a",
				"$ref": "https://example.com/v2/Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
				"display": "Tour Guides"
			},
			{
				"value": "fc348aa8-3835-40eb-a20b-c726e15c55b5",
				"$ref": "https://example.com/v2/Groups/fc348aa8-3835-40eb-a20b-c726e15c55b5",
				"display": "Employees"
			},
			{
				"value": "71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
				"$ref": "https://example.com/v2/Groups/71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
				"display": "US Employees"
			}
		],
		"x509Certificates": [
			{
				"value": "MIIDQzCCAqygAwIBAgICEAAwDQYJKoZIhvcNAQEFBQAwTjELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFDASBgNVBAoMC2V4YW1wbGUuY29tMRQwEgYDVQQDDAtleGFtcGxlLmNvbTAeFw0xMTEwMjIwNjI0MzFaFw0xMjEwMDQwNjI0MzFaMH8xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRQwEgYDVQQKDAtleGFtcGxlLmNvbTEhMB8GA1UEAwwYTXMuIEJhcmJhcmEgSiBKZW5zZW4gSUlJMSIwIAYJKoZIhvcNAQkBFhNiamVuc2VuQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7Kr+Dcds/JQ5GwejJFcBIP682X3xpjis56AK02bc1FLgzdLI8auoR+cC9/Vrh5t66HkQIOdA4unHh0AaZ4xL5PhVbXIPMB5vAPKpzz5iPSi8xO8SL7I7SDhcBVJhqVqr3HgllEG6UClDdHO7nkLuwXq8HcISKkbT5WFTVfFZzidPl8HZ7DhXkZIRtJwBweq4bvm3hM1Os7UQH05ZS6cVDgweKNwdLLrT51ikSQG3DYrl+ft781UQRIqxgwqCfXEuDiinPh0kkvIi5jivVu1Z9QiwlYEdRbLJ4zJQBmDrSGTMYn4lRc2HgHO4DqB/bnMVorHB0CC6AV1QoFK4GPe1LwIDAQABo3sweTAJBgNVHRMEAjAAMCwGCWCGSAGG+EIBDQQfFh1PcGVuU1NMIEdlbmVyYXRlZCBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQU8pD0U0vsZIsaA16lL8En8bx0F/gwHwYDVR0jBBgwFoAUdGeKitcaF7gnzsNwDx708kqaVt0wDQYJKoZIhvcNAQEFBQADgYEAA81SsFnOdYJtNg5Tcq+/ByEDrBgnusx0jloUhByPMEVkoMZ3J7j1ZgI8rAbOkNngX8+pKfTiDz1RC4+dx8oU6Za+4NJXUjlL5CvV6BEYb1+QAEJwitTVvxB/A67g42/vzgAtoRUeDov1+GFiBZ+GNF/cAYKcMtGcrs2i97ZkJMo="
			}
		],
		"meta": {
			"resourceType": "User",
			"created": "2010-01-23T04:56:22Z",
			"lastModified": "2011-05-13T04:42:34Z",
			"version": "W/a330bc54f0671c9",
			"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
		}
	}`

	var user User
	error := Unmarshal([]byte(fullUser), &user)

	assert.Nil(error, "Error unmarshaling the User object - %q", error)
	assert.Equal(user.ID, "2819c223-7f76-453a-919d-413861904646", "Missing or incorrect id attribute")
}