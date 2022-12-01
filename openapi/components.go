package openapi

// Components is a programmatic representation of the Components object defined here: https://swagger.io/specification/#components-object
type Components struct {
	Schemas         map[string]Schema         // can also be a $ref
	Responses       map[string]Response       // can also be a $ref
	Parameters      map[string]Parameter      // can also be a $ref
	Examples        map[string]Example        // can also be a $ref
	RequestBodies   map[string]RequestBody    // can also be a $ref
	Headers         map[string]Header         // can also be a $ref
	SecuritySchemes map[string]SecurityScheme // can also be a $ref
	Links           map[string]Link           // can also be a $ref
	Callbacks       map[string]Callback       // can also be a $ref
}

// Example is a programmatic representation of the Example object defined here: https://swagger.io/specification/#components-object
type Example struct {
	Summary       string      `yaml:"summary"`
	Description   string      `yaml:"description"`
	Value         interface{} `yaml:"value"` // TODO: does this imply an explicit need for reflective interpretation of data while marshalling?
	ExternalValue string      `yaml:"externalValue"`
}

// TODO: work out ABNF expressions for the resolution of $ref strings for below structs

// Header is a programmatic representation of the Header object defined here:https://swagger.io/specification/#header-object
type Header struct{}

type OAuthFlows struct {
	node
	Implicit          OAuthFlow `yaml:"implicit"`          // Configuration for the OAuth Implicit flow
	Password          OAuthFlow `yaml:"password"`          // Configuration for the OAuth Resource Owner Password flow
	ClientCredentials OAuthFlow `yaml:"clientCredentials"` // Configuration for the OAuth Client Credentials flow. Previously called application in OpenAPI 2.0.
	AuthorizationCode OAuthFlow `yaml:"authorizationCode"` // Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.
}
type OAuthFlow struct {
	node
	AuthorizationUrl map[string]string `yaml:"authorizationUrl"`
	TokenUrl         string            `yaml:"tokenUrl"`
	RefreshUrl       string            `yaml:"refreshUrl"`
	Scopes           map[string]string `yaml:"scopes"`
}

// SecurityScheme is a programmatic representation of the SecurityScheme object defined here: https://swagger.io/specification/#security-scheme-object
type SecurityScheme struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`

	// only applies to API key
	Name string `yaml:"name"`
	In   string `yaml:"in"`

	// only applies to http
	Scheme string `yaml:"scheme"`

	// only applies to http bearer tokens
	BearerFormat string `yaml:"bearerFormat"`

	// only applies to oauth2
	Flows OAuthFlows `yaml:"flows"`

	// only applies to openIdConnect
	OpenIdConnectUrl string `yaml:"openIdConnectUrl"`
}

type SecurityRequirement map[string][]string // Contains a named set of string arrays eg. { "petstore_auth": ["write:pets", "read:pets"] }

// Link is a programmatic representation of the Link object defined here: https://swagger.io/specification/#link-object
type Link struct{}

// Callback is a programmatic representation of the Callback object defined here: https://swagger.io/specification/#callback-object
type Callback struct{}
