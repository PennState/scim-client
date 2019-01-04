package user

//Address is a street and country based addess for the identity
type Address struct {
	//Type represents the type of address (i.e. home, work, other, etc..)
	Type string `json:"type"`

	//Display is a human readable representation of the address
	Display string `json:"display"`

	//Primary is a flag indicating whether this is the primary address for a user record
	Primary bool `json:"primary"`

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
