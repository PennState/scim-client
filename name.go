package cprclient

//Name is the name of the user
type Name struct {
	//Formatted is the full name, including all middle names, titles, and suffixes as appropriate, formatted for display (e.g. Ms. Barbara J Jensen, III.)
	Formatted string `json:"formatted"`

	//FamilyName is the family name of the User, or Last Name in most Western languages (e.g. Jensen given the full name Ms. Barbara J Jensen, III.)
	FamilyName string `json:"familyName"`

	//GiveName is the given name of the User, or First Name in most Western languages (e.g. Barbara given the full name Ms. Barbara J Jensen, III.)
	GivenName string `json:"givenName"`

	//MiddleName is the middle name(s) of the User (e.g. Robert given the full name Ms. Barbara J Jensen, III.).
	MiddleName string `json:"middleName"`

	//HonorificPrevix is the honorific prefix(es) of the User, or Title in most Western languages (e.g. Ms. given the full name Ms. Barbara J Jensen, III.)
	HonorificPrefix string `json:"honorificPrefix"`

	//HonorificSuffix is the honorific suffix(es) of the User, or Suffix in most Western languages (e.g. III. given the full name Ms. Barbara J Jensen, III.)
	HonorificSuffix string `json:"honorificSuffix"`
}
