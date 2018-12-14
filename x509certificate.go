package cprclient

//X509Certificate DER-encoded X.509 certificate (see Section 4 of [RFC5280]), which MUST be base64 encoded per Section 4 of [RFC4648].
type X509Certificate struct {
	//Type is a label indicating the attribute's function.
	Type string `json:"type,omitempty"`

	//Value is the value of a X509 certificate.
	Value string `json:"value"`

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.
	Display string `json:"display"`

	//Primary is a boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g. the preferred mailing address or primary e-mail address. The primary attribute value 'true' MUST appear no more than once.
	Primary bool `json:"primary"`
}
