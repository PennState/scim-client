package examples

//Pantry is a SCIM Extension that allows employees and employers to track
//the amount owed to the pantry (or to the employee).
type Pantry struct {
	Building string  `json:"building"` //Building is the location housing the employee's office
	Office   string  `json:"office"`   //Office is the room number of the employee's office
	Balance  float64 `json:"balance"`  //Balance is the amount the employee owes to the pantry (if negative).  Credits can be represented by positive Balance values
}

//URN returns the SCIM Extension's URN (identifier) and, more importantly
//identifies the Pantry struct as a SCIM extension.
func (p Pantry) URN() string {
	return "urn:com:example:2.0:Pantry"
}
