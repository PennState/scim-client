package scim

//User describes a SCIM user based on the RFC7643 specification
type User struct {

	ScimResource

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
	Groups []Group `json:"groups"`

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

//Address is a street and country based addess for the identity
type Address struct {
	Multivalued

	//Country is the contry location of the address
	Country string `json:"country"`

	// Formatted is the full mailing address, formatted for display or use with a mailing label. This attribute MAY contain newlines.")
	Formatted string `json:"formatted"`

	//Locality is the city or locality component.
	Locality string `json:"locality"`

	//PostalCode is the zipcode or postal code component.
	PostalCode string

	//Region is the state or region component.
	Region string

	//StreetAddress is the full street address component, which may include house number, street name, PO BOX, and multi-line extended street address information. This attribute MAY contain newlines.")
	StreetAddress string
}

type Email = StringMultivalued
type Entitlement = StringMultivalued
type Group = StringMultivalued
type Im = StringMultivalued

//Name is the name of the user
type Name struct {
	//Formatted is the full name, including all middle names, titles, and suffixes as appropriate, formatted for display (e.g. Ms. Barbara J Jensen, III.)
	Formatted string `json:"formatted"`

	//FamilyName is the family name of the User, or Last Name in most Western languages (e.g. Jensen given the full name Ms. Barbara J Jensen, III.)
	FamilyName string `json:"familyName"`

	//GiveName is the given name of the User, or First Name in most Western languages (e.g. Barbara given the full name Ms. Barbara J Jensen, III.)
	GivenName string `json:"givenName"`

	//MiddleName is the middle name(s) of the User (e.g. Robert given the full name Ms. Barbara J Jensen, III.).
	MiddleName string `json:"middleName"`

	//HonorificPrevix is the honorific prefix(es) of the User, or Title in most Western languages (e.g. Ms. given the full name Ms. Barbara J Jensen, III.)
	HonorificPrefix string `json:"honorificPrefix"`

	//HonorificSuffix is the honorific suffix(es) of the User, or Suffix in most Western languages (e.g. III. given the full name Ms. Barbara J Jensen, III.)
	HonorificSuffix string `json:"honorificSuffix"`
}

type PhoneNumber = StringMultivalued
type Photo = StringMultivalued
type Role = StringMultivalued
type X509Certificate = StringMultivalued
