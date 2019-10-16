package scim

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

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

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("HTTP status: %v, Type: %v, Detail: %v", er.Status, er.ScimType, er.Detail)
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

type listResponse struct {
	Schemas      []string          `json:"schemas"`      //Schemas identifies the response as a ListResponse.
	ItemsPerPage int               `json:"itemsPerPage"` //ItemsPerPage is the number of resources returned in a list response page.
	Resources    []json.RawMessage `json:"Resources"`    //Resources is a multi-valued list of complex objects containing the requested resources.  This MAY be a subset of the full set of resources if pagination (Section 3.4.2.4) is requested.
	StartIndex   int               `json:"startIndex"`   //StartIndex is the 1-based index of the first result in the current set of list results.  REQUIRED when partial results are returned due to pagination.
	TotalResults int               `json:"totalResults"` //TotalResults is the total number of results returned by the list or query operation.  The value may be larger than the number of resources returned, such as when returning a single page (see Section 3.4.2.4) of results where multiple pages are available.
}

const SearchRequestURN = "urn:ietf:params:scim:api:messages:2.0:SearchRequest"

type SearchRequest struct {
	Schemas            []string  `json:"schemas" validate:"required"`
	Attributes         []string  `json:"attributes,omitempty"`
	ExcludedAttributes []string  `json:"excludedAttributes,omitempty"`
	Filter             string    `json:"filter" validate:"required"`
	SortBy             string    `json:"sortBy,omitempty"`
	SortOrder          sortOrder `json:"sortOrder,omitempty"`
	StartIndex         int       `json:"startIndex,omitempty"`
	Count              int       `json:"count,omitempty"`
}

func URN() string {
	return SearchRequestURN
}

type sortOrder string

const (
	Ascending    sortOrder = "ascending"
	Descending   sortOrder = "descending"
	NotSpecified sortOrder = ""
)

func (lro *ListResponse) UnmarshalJSON(data []byte) error {
	log.Trace("-> (ListResponse) UnmarshalJSON([]byte) error")
	log.Debug("ListResponse raw data: ", string(data))
	var lri listResponse
	err := json.Unmarshal(data, &lri)
	if err != nil {
		return err
	}
	log.Debug("lri: ", lri)

	lro.Schemas = lri.Schemas
	lro.ItemsPerPage = lri.ItemsPerPage
	lro.StartIndex = lri.StartIndex
	lro.TotalResults = lri.TotalResults

	for _, rm := range lri.Resources {
		var ca CommonAttributes
		err := json.Unmarshal(rm, &ca)
		if err != nil {
			return err
		}
		log.Debug("CA schemas: ", ca.Schemas)
		log.Debug("CA meta: ", ca.Meta.ResourceType)

		// TODO: Replace this block with code to unmarshall by resource type.
		var user User
		err = Unmarshal(rm, &user)
		if err != nil {
			return err
		}
		log.Debug("CA as user: ", user)
		lro.Resources = append(lro.Resources, &user)
	}

	log.Trace("(ListResponse) UnmarshalJSON([]byte) error ->")
	return nil
}
