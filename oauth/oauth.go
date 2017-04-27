package oauth

import (
	"proxy"
)

// AccessRequestType is the type for OAuth param `grant_type`
type AccessRequestType string

// AuthorizeRequestType is the type for OAuth param `response_type`
type AuthorizeRequestType string

// Spec Holds an api definition and basic options
type Spec struct {
	*OAuth
	Manager Manager
}

// OAuth holds the configuration for oauth proxies
type OAuth struct {
	Name                   string                 `bson:"name" json:"name" valid:"required"`
	Endpoints              Endpoints              `bson:"oauth_endpoints" json:"oauth_endpoints" mapstructure:"oauth_endpoints"`
	ClientEndpoints        ClientEndpoints        `bson:"oauth_client_endpoints" json:"oauth_client_endpoints" mapstructure:"oauth_client_endpoints"`
	AllowedAccessTypes     []AccessRequestType    `bson:"allowed_access_types" json:"allowed_access_types" mapstructure:"allowed_access_types" `
	AllowedAuthorizeTypes  []AuthorizeRequestType `bson:"allowed_authorize_types" json:"allowed_authorize_types" mapstructure:"allowed_authorize_types"`
	AuthorizeLoginRedirect string                 `bson:"auth_login_redirect" json:"auth_login_redirect" mapstructure:"auth_login_redirect"`
	Secrets                map[string]string      `bson:"secrets" json:"secrets"`
	CorsMeta               meta                   `bson:"cors_meta" json:"cors_meta" valid:"cors_meta" mapstructure:"cors_meta"`
	TokenStrategy          TokenStrategy          `bson:"token_strategy" json:"token_strategy" mapstructure:"token_strategy"`
}

// Endpoints defines the oauth endpoints that wil be proxied
type Endpoints struct {
	Authorize *proxy.Definition `bson:"authorize" json:"authorize"`
	Token     *proxy.Definition `bson:"token" json:"token"`
	Info      *proxy.Definition `bson:"info" json:"info"`
	Revoke    *proxy.Definition `bson:"revoke" json:"revoke"`
}

// ClientEndpoints defines the oauth client endpoints that wil be proxied
type ClientEndpoints struct {
	Create *proxy.Definition `bson:"create" json:"create"`
	Remove *proxy.Definition `bson:"remove" json:"remove"`
}

// TokenStrategy defines the token strategy fields
type TokenStrategy struct {
	Name     string                `bson:"name" json:"name"`
	Settings TokenStrategySettings `bson:"settings" json:"settings"`
}

// TokenStrategySettings represents the settings for the token strategy
type TokenStrategySettings map[string]string

type meta struct {
	Domains        []string `mapstructure:"domains" bson:"domains" json:"domains"`
	Methods        []string `mapstructure:"methods" bson:"methods" json:"methods"`
	RequestHeaders []string `mapstructure:"request_headers" bson:"request_headers" json:"request_headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers" bson:"exposed_headers" json:"exposed_headers"`
	Enabled        bool     `bson:"enabled" json:"enabled"`
}
