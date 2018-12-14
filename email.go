package cprclient

//Email describes a user email address
type Email struct {
	//Type is a label indicating the attribute's function; e.g., 'work' or 'home'.
	Type string `json:"type"`

	//Value is an E-mail addresses for the user. The value SHOULD be canonicalized by the Service Provider, e.g. bjensen@example.com instead of bjensen@EXAMPLE.COM. Canonical Type values of work, home, and other.
	Value string `json:"value"`

	//Display is a  human readable name, primarily used for display purposes. READ-ONLY.
	Display string `json:"display"`

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.
	Primary bool `json:"primary"`
}
