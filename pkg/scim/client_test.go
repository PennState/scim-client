package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceURLParsing(t *testing.T) {
	tests := []struct {
		Name        string
		InputURL    string
		ExpectedURL string
		ErrMessage  string
	}{
		{"Valid URL", "http://example.com/scim", "http://example.com/scim", ""},
		{"Valid URL with trailing slash", "http://example.com/scim/", "http://example.com/scim", ""},
		{"Empty URL", "", "", noServiceURLMessage},
		{"Invalid URL", ":", "", invalidServiceURLMessage},
	}
	for idx := range tests {
		test := tests[idx]
		t.Run(test.Name, func(t *testing.T) {
			c, err := NewClient(nil, test.InputURL)
			if err != nil {
				assert.EqualError(t, err, test.ErrMessage)
				return
			}
			assert.Equal(t, test.ExpectedURL, c.cfg.ServiceURL)
		})
	}
}
