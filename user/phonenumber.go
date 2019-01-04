package user

//PhoneNumber is a number for the user.
//The value SHOULD be specified according to the format defined in [RFC3966], e.g., 'tel:+1-201-555-0123'.  Service providers SHOULD canonicalize the value according to [RFC3966] format, when appropriate.
type PhoneNumber struct {

	//Value is phone number of the User
	Value string `json:"value"`

	//Type A label indicating the attribute's function; e.g., 'work' or 'home' or 'mobile' etc.
	Type string `json:"type"`

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.")
	Display string

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred phone number or primary phone number. The primary attribute value 'true' MUST appear no more than once.")
	Primary bool
}
