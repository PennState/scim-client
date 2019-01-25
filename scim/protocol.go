package scim

const ErrorResponseURN = "urn:ietf:params:scim:api:messages:2.0:Error"

//ErrorResponse is a SCIM standard JSON response body used by a SCIM
//server when an error must be returned to the client.
//https://tools.ietf.org/html/rfc7644#section-3.12
type ErrorResponse struct {
	Schemas  []string `json:"schemas"`  //Schemas identifies the response as an ErrorResponse.
	ScimType string   `json:"scimType"` //ScimType is a detail error keyword.  See Table 9.
	Detail   string   `json:"detail"`   //Detail is a human-readable message.
	Status   string   `json:"status"`   //Status is the HTTP status code expressed as a JSON string.
}

const ListResponseURN = "urn:ietf:params:scim:api:messages:2.0:ListResponse"

//ListResponse defines the SCIM standard response to a valid search query
//(which might return zero or more results).
//https://tools.ietf.org/html/rfc7644#section-3.4.2
type ListResponse struct {
	Schemas      []string   `json:"schemas"`      //Schemas identifies the response as a ListResponse.
	ItemsPerPage int        `json:"itemsPerPage"` //ItemsPerPage is the number of resources returned in a list response page.
	Resources    []Resource `json:"Resources"`    //Resources is a multi-valued list of complex objects containing the requested resources.  This MAY be a subset of the full set of resources if pagination (Section 3.4.2.4) is requested.
	StartIndex   int        `json:"startIndex"`   //StartIndex is the 1-based index of the first result in the current set of list results.  REQUIRED when partial results are returned due to pagination.
	TotalResults int        `json:"totalResults"` //TotalResults is the total number of results returned by the list or query operation.  The value may be larger than the number of resources returned, such as when returning a single page (see Section 3.4.2.4) of results where multiple pages are available.
}
