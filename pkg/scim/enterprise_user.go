package scim

//EnterpriseUser defines attributes commonly used in representing users
//that belong to, or act on behalf of, a business or enterprise.
//https://tools.ietf.org/html/rfc7643#section-4.3
type EnterpriseUser struct {
	EmployeeNumber string  `json:"employeeNumber"` //EmployeeNumber is astring identifier, typically numeric or alphanumeric, assigned to a person, typically based on order of hire or association with an organization.
	CostCenter     string  `json:"costCenter"`     //CostCenter identifies the name of a cost center.
	Organization   string  `json:"organization"`   //Organization identifies the name of an organization.
	Division       string  `json:"division"`       //Division identifies the name of a division
	Department     string  `json:"department"`     //Department identifies the name of a department
	Manager        Manager `json:"manager"`        //The user's manager.  A complex type that optionally allows service providers to represent organizational hierarchy by referencing the "id" attribute of another User.
}

//Manager is a reference to the user's manager with a small amount of
//cargo data for convenience.
type Manager struct {
	Value       string `json:"value"`       //Value is the "id" of the SCIM resource representing the user's manager.
	Reference   string `json:"$ref"`        //Reference is the URI of the SCIM resource representing the User's manager.
	DisplayName string `json:"displayName"` //DisplayName is the displayName of the user's manager.
}

//EnterpriseUserURN is the IANA registered name that identifies SCIM
//enterprise user extension.
const EnterpriseUserURN = "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"

//URN returns the EnterpriseUserURN but more importantly implements
//the Extension interface that identifies this struct as a SCIM
//extension.
func (eu EnterpriseUser) URN() string {
	return EnterpriseUserURN
}
