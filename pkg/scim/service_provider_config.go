package scim

const ServiceProviderConfigURN = "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"

//https://tools.ietf.org/html/rfc7643#section-7
type ServiceProviderConfig struct {
	CommonAttributes
	DocumentationURI      string                       `json:"documentationUri,omitempty"`
	PatchConfig           PatchConfig                  `json:"patch" validation:"required"`
	BulkConfig            BulkConfig                   `json:"bulk" validation:"required"`
	FilterConfig          FilterConfig                 `json:"filter" validation:"required"`
	ChangePasswordConfig  ChangePasswordConfig         `json:"changePassword" validation:"required"`
	SortConfig            SortConfig                   `json:"sort" validation:"required"`
	ETagConfig            ETagConfig                   `json:"etag" validation:"required"`
	AuthenticationSchemes []AuthenticationSchemeConfig `json:"authenticationSchemes" validation:"required"`
}

type supportedConfig struct {
	Supported bool `json:"supported" validation:"required"`
}

type PatchConfig supportedConfig

type BulkConfig struct {
	supportedConfig
	MaxOperations  int `json:"maxOperations" validation:"required"`
	MaxPayloadSize int `json:"maxPayloadSize" validation:"required"`
}

type FilterConfig struct {
	supportedConfig
	MaxResults int `json:"maxResults" validation:"required"`
}

type ChangePasswordConfig supportedConfig

type SortConfig supportedConfig

type ETagConfig supportedConfig

type AuthenticationSchemeConfig struct {
	Type             AuthenticationSchemeType `json:"type" validation:"required"`
	Name             string                   `json:"name" validation:"required"`
	Description      string                   `json:"description" validation:"required"`
	SpecURI          string                   `json:"specUri,omitempty"`
	DocumentationURI string                   `json:"documentationUri,omitempty"`
}

type AuthenticationSchemeType string

const (
	OAuth            AuthenticationSchemeType = "oauth"
	OAuth2           AuthenticationSchemeType = "oauth2"
	OAuthBearerToken AuthenticationSchemeType = "oauthbearertoken"
	HTTPBasic        AuthenticationSchemeType = "httpbasic"
	HTTPDigest       AuthenticationSchemeType = "httpdigest"
)

var ServiceProviderConfigResourceType = ResourceType{
	CommonAttributes: CommonAttributes{
		Schemas: []string{
			ResourceTypeURN,
		},
		ID: "ResourceType",
	},
	Name:        "ServiceProviderConfig",
	Endpoint:    "/ServiceProviderConfig",
	Description: "SCIM Service Provider Config - See https://tools.ietf.org/html/rfc7643#section-5",
	Schema:      ServiceProviderConfigURN,
}

func (spc ServiceProviderConfig) URN() string {
	return ServiceProviderConfigURN
}

func (spc ServiceProviderConfig) ResourceType() ResourceType {
	return ServiceProviderConfigResourceType
}
