package scim

const SearchRequestURN = "urn:ietf:params:scim:api:messages:2.0:SearchRequest"

type searchRequest struct {
	schemas            []string  `json:"schemas" validate:"required"`
	Attributes         []string  `json:"attributes"`
	ExcludedAttributes []string  `json:"ExcludedAttributes"`
	Filter             string    `json:"filter"`
	SortBy             string    `json:"sortBy"`
	SortOrder          sortOrder `json:"sortOrder"`
	StartIndex         int       `json:"startIndex"`
	Count              int       `json:"count"`
}

type sortOrder string

const (
	Ascending  sortOrder = "ascending"
	Descending sortOrder = "descending"
)

func NewSearchRequest() searchRequest {
	var schemas []string
	schemas = append(schemas, SearchRequestURN)
	return searchRequest{
		schemas: schemas,
	}
}

func (sr searchRequest) Schemas() []string {
	return sr.schemas
}
