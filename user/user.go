package user

import (
	"github.com/PennState/golang_scimclient/resource"
	"github.com/PennState/golang_scimclient/schema"
)

//User describes a SCIM user based on the RFC7643 specification
type User struct {

	resource.ScimResource

	//Active informs as to whether this User record is currently live in the system
	Active bool `json:"active"`

	//Addresses is the  physical mailing address for this User, as described in (address Element). Canonical Type Values of work, home, and other. The value attribute is a complex type with the following sub-attributes.
	Addresses []Address

	//DisplayName is the name of the User, suitable for display to end-users. The name SHOULD be the full name of the User being described if known
	DisplayName string

	//Emails are E-mail addresses for the user. The value SHOULD be canonicalized by the Service Provider, e.g. bjensen@example.com instead of bjensen@EXAMPLE.COM. Canonical Type values of work, home, and other.
	Emails []Email

	//Entitlements is a collection of entitlements
	Entitlements []Entitlement

	//Groups is a list of groups that the user belongs to, either thorough direct membership, nested groups, or dynamically calculated")
	Groups []schema.ResourceReference `json:"groups"`

	//Ims are instant messaging addresses for the User.
	Ims []Im

	//Locale is used to indicate the User's default location for purposes of localizing items such as currency, date time format, numerical representations, etc.
	Locale string `json:"locale"`

	//Name is the components of the user's real name. Providers MAY return just the full name as a single string in the formatted sub-attribute, or they MAY return just the individual component attributes using the other sub-attributes, or they MAY return both. If both variants are returned, they SHOULD be describing the same name, with the formatted name indicating how the component attributes should be combined.
	Name Name

	//NickName is the casual way to address the user in real life, e.g.'Bob' or 'Bobby' instead of 'Robert'. This attribute SHOULD NOT be used to represent a User's username (e.g. bjensen or mpepperidge)
	NickName string `json:"nickName"`

	//Password ishe User's clear text password.  This attribute is intended to be used as a means to specify an initial password when creating a new User or to reset an existing User's password.
	//(This is a placeholder in case I implement the server in golang)
	Password string `json:"_,omitempty"`

	//PhoneNumberrs are the phone numbers for the User.  The value SHOULD be canonicalized by the Service Provider according to format in RFC3966 e.g. 'tel:+1-201-555-0123'.  Canonical Type values of work, home, mobile, fax, pager and other.")
	PhoneNumbers []PhoneNumber

	//Photos are URLs of photos of the User.
	Photos []Photo

	//Profile URL is a fully qualified URL to a page representing the User's online profile
	ProfileURL string `json:"profileUrl"`

	//PreferredLanguage indicates the User's preferred written or spoken language.  Generally used for selecting a localized User interface. e.g., 'en_US' specifies the language English and country US.
	PreferredLanguage string `json:"preferredLanguage"`

	//Rols are a list of roles for the User that collectively represent who the User is; e.g., 'Student', 'Faculty'.
	Roles []Role

	//Timezone is the User's time zone in the 'Olson' timezone database format; e.g.,'America/Los_Angeles'
	Timezone string `json:"timezone"`

	//Title is the user's title, such as "Vice President.
	Title string `json:"title"`

	//UserName is a unique identifier for the User typically used by the user to directly authenticate to the service provider. Each User MUST include a non-empty userName value.  This identifier MUST be unique across the Service Consumer's entire set of Users.  REQUIRED
	UserName string `json:"userName"`

	//UserType is used to identify the organization to user relationship. Typical values used might be 'Contractor', 'Employee', 'Intern', 'Temp', 'External', and 'Unknown' but any value may be used.
	UserType string `json:"userType"`

	//X509Certificates is list of certificates issued to the User.
	X509Certificates []X509Certificate
}
