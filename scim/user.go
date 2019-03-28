package scim

//UserURN is the IANA registered SCIM name for the standardized SCIM
//User.
const UserURN = "urn:ietf:params:scim:schemas:core:2.0:User"

//User describes a SCIM user based on the RFC7643 specification
//https://tools.ietf.org/html/rfc7643#section-4.1
type User struct {
	CommonAttributes
	Active            bool              `json:"active"`             //Active informs as to whether this User record is currently live in the system
	Addresses         []Address         `json:"addresses"`          //Addresses is the  physical mailing address for this User, as described in (address Element). Canonical Type Values of work, home, and other. The value attribute is a complex type with the following sub-attributes.
	DisplayName       string            `json:"displayName"`        //DisplayName is the name of the User, suitable for display to end-users. The name SHOULD be the full name of the User being described if known
	Emails            []Email           `json:"emails"`             //Emails are E-mail addresses for the user. The value SHOULD be canonicalized by the Service Provider, e.g. bjensen@example.com instead of bjensen@EXAMPLE.COM. Canonical Type values of work, home, and other.
	Entitlements      []Entitlement     `json:"entitlements"`       //Entitlements is a collection of entitlements
	Groups            []GroupRef        `json:"groups"`             //Groups is a list of groups that the user belongs to, either thorough direct membership, nested groups, or dynamically calculated")
	IMs               []IM              `json:"ims"`                //Ims are instant messaging addresses for the User.
	Locale            string            `json:"locale"`             //Locale is used to indicate the User's default location for purposes of localizing items such as currency, date time format, numerical representations, etc.
	Name              Name              `json:"name"`               //Name is the components of the user's real name. Providers MAY return just the full name as a single string in the formatted sub-attribute, or they MAY return just the individual component attributes using the other sub-attributes, or they MAY return both. If both variants are returned, they SHOULD be describing the same name, with the formatted name indicating how the component attributes should be combined.
	NickName          string            `json:"nickName"`           //NickName is the casual way to address the user in real life, e.g.'Bob' or 'Bobby' instead of 'Robert'. This attribute SHOULD NOT be used to represent a User's username (e.g. bjensen or mpepperidge)
	Password          string            `json:"password,omitempty"` //Password is the User's clear text password.  This attribute is intended to be used as a means to specify an initial password when creating a new User or to reset an existing User's password. -TODO - (This is a placeholder in case I implement the server in golang)
	PhoneNumbers      []PhoneNumber     `json:"phoneNumbers"`       //PhoneNumberrs are the phone numbers for the User.  The value SHOULD be canonicalized by the Service Provider according to format in RFC3966 e.g. 'tel:+1-201-555-0123'.  Canonical Type values of work, home, mobile, fax, pager and other.")
	Photos            []Photo           `json:"photos"`             //Photos are URLs of photos of the User.
	ProfileURL        string            `json:"profileUrl"`         //Profile URL is a fully qualified URL to a page representing the User's online profile
	PreferredLanguage string            `json:"preferredLanguage"`  //PreferredLanguage indicates the User's preferred written or spoken language.  Generally used for selecting a localized User interface. e.g., 'en_US' specifies the language English and country US.
	Roles             []Role            `json:"roles"`              //Roles are a list of roles for the User that collectively represent who the User is; e.g., 'Student', 'Faculty'.
	Timezone          string            `json:"timezone"`           //Timezone is the User's time zone in the 'Olson' timezone database format; e.g.,'America/Los_Angeles'
	Title             string            `json:"title"`              //Title is the user's title, such as "Vice President.
	UserName          string            `json:"userName"`           //UserName is a unique identifier for the User typically used by the user to directly authenticate to the service provider. Each User MUST include a non-empty userName value.  This identifier MUST be unique across the Service Consumer's entire set of Users.  REQUIRED
	UserType          string            `json:"userType"`           //UserType is used to identify the organization to user relationship. Typical values used might be 'Contractor', 'Employee', 'Intern', 'Temp', 'External', and 'Unknown' but any value may be used.
	X509Certificates  []X509Certificate `json:"x509Certificates"`   //X509Certificates is list of certificates issued to the User.
}

//Address is a street and country based addess for the identity
type Address struct {
	Multivalued
	Country       string `json:"country"`       //Country is the contry location of the address
	Formatted     string `json:"formatted"`     //Formatted is the full mailing address, formatted for display or use with a mailing label. This attribute MAY contain newlines.")
	Locality      string `json:"locality"`      //Locality is the city or locality component.
	PostalCode    string `json:"postalCode"`    //PostalCode is the zipcode or postal code component.
	Region        string `json:"region"`        //Region is the state or region component.
	StreetAddress string `json:"streetAddress"` //StreetAddress is the full street address component, which may include house number, street name, PO BOX, and multi-line extended street address information. This attribute MAY contain newlines.")
}

//Email provides an email address in the StringMultivalued.Value field.
type Email StringMultivalued

//Entitlement provides an entitlement name in the StringMultivalued.Value field.
type Entitlement StringMultivalued

//GroupRef indicates membership in a scim.Group by providing a reference as well as a small amount of cargo data to the group.
type GroupRef StringMultivalued

//IM provides an instant message address in the StringMultivalued.Value field.
type IM StringMultivalued

//Name is the name of the user
type Name struct {
	Formatted       string `json:"formatted"`       //Formatted is the full name, including all middle names, titles, and suffixes as appropriate, formatted for display (e.g. Ms. Barbara J Jensen, III.)
	FamilyName      string `json:"familyName"`      //FamilyName is the family name of the User, or Last Name in most Western languages (e.g. Jensen given the full name Ms. Barbara J Jensen, III.)
	GivenName       string `json:"givenName"`       //GiveName is the given name of the User, or First Name in most Western languages (e.g. Barbara given the full name Ms. Barbara J Jensen, III.)
	MiddleName      string `json:"middleName"`      //MiddleName is the middle name(s) of the User (e.g. Robert given the full name Ms. Barbara J Jensen, III.).
	HonorificPrefix string `json:"honorificPrefix"` //HonorificPrevix is the honorific prefix(es) of the User, or Title in most Western languages (e.g. Ms. given the full name Ms. Barbara J Jensen, III.)
	HonorificSuffix string `json:"honorificSuffix"` //HonorificSuffix is the honorific suffix(es) of the User, or Suffix in most Western languages (e.g. III. given the full name Ms. Barbara J Jensen, III.)
}

//PhoneNumber provides an RFC3966 compliant phone number in the StringMultivalued.Value field.
type PhoneNumber StringMultivalued

//Photo provides a link to a user's photograph in the StringMultivalued.Value field.
type Photo StringMultivalued

//Role provides an identifier for a role in the StringMultivalued.Value field.
type Role StringMultivalued

//X509Certificate provides a DER-encoded X.509 certificate in the StringMultivalued.Value field.
type X509Certificate StringMultivalued

//UserResourceType provides the default structure which connects the User
//struct to its associated ResourceType.
var UserResourceType = ResourceType{
	CommonAttributes: CommonAttributes{
		Schemas: []string{
			ResourceTypeURN,
		},
		ID: "User",
	},
	Name:        "User",
	Endpoint:    "/Users",
	Description: "SCIM ResourceType - See https://tools.ietf.org/html/rfc7643#section-6",
	Schema:      ResourceTypeURN,
}

//URN returns the IANA registered SCIM name for the User data structure
//and, together with ResourceType() identifies this code as implementing
//the Resource interface.
func (u User) URN() string {
	return UserURN
}

//ResourceType returns the default structure describing the availability
//of the User resource and, together with URN() identifies this code as
//implementing the Resource interface.
func (u User) ResourceType() ResourceType {
	return UserResourceType
}
