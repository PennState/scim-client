package cprclient

//Entitlement may be an additional right to a thing, object, or service
type Entitlement struct {
	//Type is a label indicating the attribute's function.
	Type string `json:"type"`

	//Value is the value of an entitlement.
	Value string `json:"value"`

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.")
	Display string `json:"display"`

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.")
	Primary bool `json:"primary"`
}
