package cprclient

//Photo is the url location of the users photo
type Photo struct {
	//Value is the URL of a photo of the User.
	Value string `json:"value"`

	//Type is a label indicating the attribute's function; e.g., 'photo' or 'thumbnail'.
	Type string `json:"type,omitempty"`

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.
	Display string `json:"display"`

	//Primary is a Boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.
	Primary bool `json:"primary"`
}
