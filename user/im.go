package user

//Im is an instant messaging address for the user.
type Im struct {

	//Type is a label indicating the attribute's function; e.g., 'aim', 'gtalk', 'mobile' etc.
	Type string `json:"type"`

	//Value is an instant messaging address for the User.
	Value string

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.")
	Display string `json:"display"`

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.")
	Primary bool `json:"primary"`
}
