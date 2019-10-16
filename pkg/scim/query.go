package scim

import "fmt"

// NewSearchRequestFromFormat returns a search request with the Filter
// field constructed from the provided format string and arguments.
func NewSearchRequestFromFormat(format string, a ...interface{}) SearchRequest {
	return SearchRequest{
		Filter: fmt.Sprintf(format, a...),
	}
}
