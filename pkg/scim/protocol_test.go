package scim

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	errorResponse = `{
		"detail": "Error converting Scim user (id: 9991533) to Person: Full or partial date of birth must be provided",
		"status": "400",
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:Error"]
	}`
)

func TestErrorResponseUnmarshaling(t *testing.T) {
	var er ErrorResponse
	err := json.Unmarshal([]byte(errorResponse), &er)
	assert.NoError(t, err)
}
