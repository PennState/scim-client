package schema

//ResourceReference is the url reference for the ScimValue
type ResourceReference struct {

	//Value is the reference Element Identifier
	Value string `json:"value"`

	//Ref is the URI of the corresponding resource ", referenceTypes={"User", "Group"})
	Ref string `json:"ref"`

	//Display is a human readable name, primarily used for display purposes. READ-ONLY.
	Display string `json:"display"`

	//ReferenceType is a label indicating the attribute's function; e.g., 'direct' or 'indirect'.", canonicalValueList={"direct", "indirect"}
	Type string `json:"type"`
}
