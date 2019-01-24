package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSearchRequestHasSchema(t *testing.T) {
	assert = assert.New(t)
	searchRequest := NewSearchRequest()
	assert.Len(searchRequest.Schemas(), 1)
	assert.Equals(SearchRequestURN, searchRequest.Schemas()[0])
}
