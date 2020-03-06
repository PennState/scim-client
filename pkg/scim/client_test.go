package scim

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientOptsParsing(t *testing.T) {
	type booleanOpt func(bool) ClientOpt
	opts := []booleanOpt{DisableDiscovery, DisableEtag, IgnoreRedirects}
	count := int(math.Pow(2, float64(len(opts))))
	for i := 0; i < count; i++ {
		dd := i&1 != 0
		de := i&2 != 0
		ir := i&4 != 0
		name := fmt.Sprintf("DisableDiscovery: %t, DisableEtag: %t, IgnoreRedirects: %t", dd, de, ir)
		t.Run(name, func(t *testing.T) {
			c, err := NewClient(
				nil,
				"https://example.com/scim",
				DisableDiscovery(dd),
				DisableEtag(de),
				IgnoreRedirects(ir),
			)
			assert.NoError(t, err)
			assert.Equal(t, dd, c.cfg.DisableDiscovery)
			assert.Equal(t, de, c.cfg.DisableEtag)
			assert.Equal(t, ir, c.cfg.IgnoreRedirects)
		})
	}
}

func TestNewClientFromEnv(t *testing.T) {
	url := "https://example.com/scim"
	os.Setenv("SCIM_SERVICE_URL", url)
	os.Setenv("SCIM_IGNORE_REDIRECTS", "true")
	os.Setenv("SCIM_DISABLE_DISCOVERY", "true")
	os.Setenv("SCIM_DISABLE_ETAG", "true")
	c, err := NewClientFromEnv(nil)
	assert.NoError(t, err)
	assert.Equal(t, c.cfg.ServiceURL, url)
	assert.True(t, c.cfg.IgnoreRedirects)
	assert.True(t, c.cfg.DisableDiscovery)
	assert.True(t, c.cfg.DisableEtag)
}

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
