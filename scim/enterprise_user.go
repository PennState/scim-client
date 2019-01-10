package scim

//https://tools.ietf.org/html/rfc7643#section-4.3
type EnterpriseUser struct {
	//EmployeeNumber is astring identifier, typically numeric or alphanumeric, assigned to a person, typically based on order of hire or association with an organization.
	EmployeeNumber string `json:"employeeNumber"`

	//CostCenter identifies the name of a cost center.
	CostCenter string `json:"costCenter"`

	//Organization identifies the name of an organization.
	Organization string `json:"organization"`

	//Division identifies the name of a division
	Division string `json:"division"`

	//Department identifies the name of a department
	Department string `json:"department"`

	//The user's manager.  A complex type that optionally allows service providers to represent organizational hierarchy by referencing the "id" attribute of another User.
	Manager Manager `json:"manager"`
}

type Manager struct {
	//Value is the "id" of the SCIM resource representing the user's manager.
	Value string `json:"value"`

	//Reference is the URI of the SCIM resource representing the User's manager.
	Reference string `json:"$ref"`

	//DisplayName is the displayName of the user's manager.
	DisplayName string `json:"displayName"`
}

func (eu EnterpriseUser) GetUrn() string {
	return "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
}
