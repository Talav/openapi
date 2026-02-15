package openapi

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/example"
	"github.com/talav/openapi/internal/build"
	"github.com/talav/openapi/internal/export"
	v304 "github.com/talav/openapi/internal/export/v304"
	v312 "github.com/talav/openapi/internal/export/v312"
	"github.com/talav/openapi/internal/model"
)

// API holds OpenAPI configuration and defines an API specification.
// All fields are public for functional options, but direct modification after creation
// is not recommended. Use functional options to configure.
//
// Create instances using [New] or [MustNew].
type API struct {
	// Info contains API metadata (title, version, description, contact, license).
	Info model.Info

	// Servers lists available server URLs for the API.
	Servers []model.Server

	// Tags provides additional metadata for operations.
	Tags []model.Tag

	// SecuritySchemes defines available authentication/authorization schemes.
	SecuritySchemes map[string]*model.SecurityScheme

	// DefaultSecurity applies security requirements to all operations by default.
	DefaultSecurity []model.SecurityRequirement

	// ExternalDocs provides external documentation links.
	ExternalDocs *model.ExternalDocs

	// Extensions contains specification extensions (fields prefixed with x-).
	// Extensions are added to the root of the OpenAPI specification.
	//
	// Direct mutation of this map after New/MustNew bypasses API.Validate().
	// However, projection-time filtering via copyExtensions still applies:
	// - Keys must start with "x-"
	// - In OpenAPI 3.1.x, keys starting with "x-oai-" or "x-oas-" are reserved and will be filtered
	//
	// Prefer using WithExtension() option instead of direct mutation.
	Extensions map[string]any

	// Version is the target OpenAPI version.
	Version string

	// StrictDownlevel causes projection to error (instead of warn) when
	// 3.1-only features are used with a 3.0 target.
	// Default: false
	StrictDownlevel bool

	// ValidateSpec enables JSON Schema validation of generated specs.
	// When enabled, Generate validates the output against the official
	// OpenAPI meta-schema (3.0.x or 3.1.x based on target version).
	// This catches specification errors early but adds ~1-5ms overhead.
	// Default: false
	ValidateSpec bool

	// SchemaPrefix is the prefix for the OpenAPI schema.
	SchemaPrefix string

	// TagConfig configures struct tag names used for OpenAPI schema generation.
	// If not set, uses default tag names (schema, body, openapi, validate, default, requires).
	TagConfig config.TagConfig

	generator       *build.SchemaGenerator
	requestBuilder  build.RequestBuilder
	responseBuilder build.ResponseBuilder
	exporter        export.Exporter
}

// Option configures OpenAPI behavior using the functional options pattern.
// Options are applied in order, with later options potentially overriding earlier ones.
type Option func(*API)

// NewAPI creates a new OpenAPI [API].
//
// Example:
//
//	api := openapi.NewAPI(
//	    openapi.WithInfoTitle("My API"),
//	    openapi.WithInfoVersion("1.0.0"),
//	    openapi.WithInfoDescription("API description"),
//	)
func NewAPI(opts ...Option) *API {
	api := &API{
		Info: model.Info{
			Title:   "API",
			Version: "1.0.0",
		},
	}
	api.TagConfig = config.DefaultTagConfig()
	api.SchemaPrefix = "#/components/schemas/"
	for _, opt := range opts {
		opt(api)
	}

	// Create metadata with tag configuration
	metadata := build.NewMetadata(api.TagConfig)

	// Create schema generator
	api.generator = build.NewSchemaGenerator(api.SchemaPrefix, metadata, api.TagConfig)

	// Create request and response builders
	api.requestBuilder = build.NewRequestBuilder(api.generator, metadata, api.TagConfig)
	api.responseBuilder = build.NewResponseBuilder(api.generator, metadata, api.TagConfig)
	api.exporter = export.NewExporter([]export.ViewAdapter{
		&v304.AdapterV304{},
		&v312.AdapterV312{},
	})

	return api
}

// WithInfoTitle sets the API title.
//
// Example:
//
//	openapi.WithInfoTitle("User Management API")
func WithInfoTitle(title string) Option {
	return func(a *API) {
		a.Info.Title = title
	}
}

// WithInfoVersion sets the API version.
//
// Example:
//
//	openapi.WithInfoVersion("2.1.0")
func WithInfoVersion(version string) Option {
	return func(a *API) {
		a.Info.Version = version
	}
}

// WithInfoDescription sets the API description in the Info object.
//
// The description supports Markdown formatting and appears in the OpenAPI spec
// and Swagger UI.
//
// Example:
//
//	openapi.WithInfoDescription("A RESTful API for managing users and their profiles.")
func WithInfoDescription(desc string) Option {
	return func(a *API) {
		a.Info.Description = desc
	}
}

// WithInfoSummary sets the API summary in the Info object (OpenAPI 3.1+ only).
// In 3.0 targets, this will be dropped with a warning.
//
// Example:
//
//	openapi.WithInfoSummary("User Management API")
func WithInfoSummary(summary string) Option {
	return func(a *API) {
		a.Info.Summary = summary
	}
}

// WithTermsOfService sets the Terms of Service URL/URI.
func WithTermsOfService(url string) Option {
	return func(a *API) {
		a.Info.TermsOfService = url
	}
}

// WithInfoExtension adds a specification extension to the Info object.
//
// Extension keys must start with "x-". In OpenAPI 3.1.x, keys starting with
// "x-oai-" or "x-oas-" are reserved and cannot be used.
//
// Example:
//
//	openapi.WithInfoExtension("x-api-category", "public")
func WithInfoExtension(key string, value any) Option {
	return func(a *API) {
		if a.Info.Extensions == nil {
			a.Info.Extensions = make(map[string]any)
		}
		a.Info.Extensions[key] = value
	}
}

// WithExternalDocs sets external documentation URL and optional description.
func WithExternalDocs(url, description string) Option {
	return func(a *API) {
		a.ExternalDocs = &model.ExternalDocs{
			URL:         url,
			Description: description,
		}
	}
}

// WithContact sets contact information for the API.
//
// All parameters are optional. Empty strings are omitted from the specification.
//
// Example:
//
//	openapi.WithContact("API Support", "https://example.com/support", "support@example.com")
func WithContact(name, url, email string) Option {
	return func(a *API) {
		a.Info.Contact = &model.Contact{
			Name:  name,
			URL:   url,
			Email: email,
		}
	}
}

// WithLicense sets license information for the API using a URL (OpenAPI 3.0 style).
//
// The name is required. URL is optional.
// This is mutually exclusive with identifier - use WithLicenseIdentifier for SPDX identifiers.
// Validation occurs when New() is called.
//
// Example:
//
//	openapi.WithLicense("MIT", "https://opensource.org/licenses/MIT")
func WithLicense(name, url string) Option {
	return func(a *API) {
		a.Info.License = &model.License{
			Name: name,
			URL:  url,
		}
	}
}

// WithLicenseIdentifier sets license information for the API using an SPDX identifier (OpenAPI 3.1+).
//
// The name is required. Identifier is an SPDX license expression (e.g., "Apache-2.0").
// This is mutually exclusive with URL - use WithLicense for URL-based licenses.
// Validation occurs when New() is called.
//
// Example:
//
//	openapi.WithLicenseIdentifier("Apache 2.0", "Apache-2.0")
func WithLicenseIdentifier(name, identifier string) Option {
	return func(a *API) {
		a.Info.License = &model.License{
			Name:       name,
			Identifier: identifier,
		}
	}
}

// ServerOption configures a Server using the functional options pattern.
type ServerOption func(*model.Server)

// WithServer adds a server URL to the specification.
//
// Multiple servers can be added by calling this option multiple times.
// Use server options to configure description, variables, and extensions.
//
// Example:
//
//	openapi.WithServer("https://api.example.com",
//	    openapi.WithServerDescription("Production"),
//	),
//	openapi.WithServer("https://{env}.example.com",
//	    openapi.WithServerDescription("Multi-tenant API"),
//	    openapi.WithServerVariable("env", "prod", []string{"prod", "staging"}, "Environment"),
//	    openapi.WithServerExtension("x-region", "us-east-1"),
//	),
func WithServer(url string, opts ...ServerOption) Option {
	return func(a *API) {
		server := &model.Server{URL: url}
		for _, opt := range opts {
			opt(server)
		}
		a.Servers = append(a.Servers, *server)
	}
}

// WithServerDescription sets the server description.
//
// Example:
//
//	openapi.WithServer("https://api.example.com",
//	    openapi.WithServerDescription("Production environment"),
//	),
func WithServerDescription(desc string) ServerOption {
	return func(s *model.Server) {
		s.Description = desc
	}
}

// WithServerVariable adds a variable to the server URL template.
//
// The variable name should match a placeholder in the server URL (e.g., {username}).
// Default is required. Enum and description are optional.
//
// Example:
//
//	openapi.WithServer("https://{username}.example.com:{port}/v1",
//	    openapi.WithServerVariable("username", "demo", []string{"demo", "prod"}, "User subdomain"),
//	    openapi.WithServerVariable("port", "8443", []string{"8443", "443"}, "Server port"),
//	),
func WithServerVariable(name, defaultValue string, enum []string, description string) ServerOption {
	return func(s *model.Server) {
		if s.Variables == nil {
			s.Variables = make(map[string]*model.ServerVariable)
		}
		s.Variables[name] = &model.ServerVariable{
			Enum:        enum,
			Default:     defaultValue,
			Description: description,
		}
	}
}

// WithServerExtension adds a specification extension to the server.
//
// Extension keys MUST start with "x-". In OpenAPI 3.1.x, keys starting with
// "x-oai-" or "x-oas-" are reserved for the OpenAPI Initiative.
//
// Example:
//
//	openapi.WithServer("https://api.example.com",
//	    openapi.WithServerExtension("x-region", "us-east-1"),
//	    openapi.WithServerExtension("x-deployment", "production"),
//	),
func WithServerExtension(key string, value any) ServerOption {
	return func(s *model.Server) {
		if s.Extensions == nil {
			s.Extensions = make(map[string]any)
		}
		s.Extensions[key] = value
	}
}

// WithTag adds a tag to the specification.
//
// Tags are used to group operations in Swagger UI. Operations can be assigned
// tags using RouteWrapper.Tags(). Multiple tags can be added by calling this
// option multiple times.
//
// Example:
//
//	openapi.WithTag("users", "User management operations"),
//	openapi.WithTag("orders", "Order processing operations"),
func WithTag(name, desc string) Option {
	return func(a *API) {
		a.Tags = append(a.Tags, model.Tag{
			Name:        name,
			Description: desc,
		})
	}
}

// WithBearerAuth adds Bearer (JWT) authentication scheme.
//
// The name is used to reference this scheme in security requirements.
// The description appears in Swagger UI to help users understand the authentication.
//
// Example:
//
//	openapi.WithBearerAuth("bearerAuth", "JWT token authentication. Format: Bearer <token>")
//
// Then use in routes:
//
//	app.GET("/protected", handler).Bearer()
func WithBearerAuth(name, desc string) Option {
	return func(a *API) {
		if a.SecuritySchemes == nil {
			a.SecuritySchemes = make(map[string]*model.SecurityScheme)
		}
		a.SecuritySchemes[name] = &model.SecurityScheme{
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  desc,
		}
	}
}

// ParameterLocation represents where an API parameter can be located.
type ParameterLocation string

const (
	// InHeader indicates the parameter is passed in the HTTP header.
	InHeader ParameterLocation = "header"

	// InQuery indicates the parameter is passed as a query string parameter.
	InQuery ParameterLocation = "query"

	// InCookie indicates the parameter is passed as a cookie.
	InCookie ParameterLocation = "cookie"
)

// WithAPIKey adds API key authentication scheme.
//
// Parameters:
//   - name: Scheme name used in security requirements
//   - paramName: Name of the header/query parameter (e.g., "X-API-Key")
//   - in: Location of the API key - use InHeader, InQuery, or InCookie
//   - desc: Description shown in Swagger UI
//
// Example:
//
//	openapi.WithAPIKey("apiKey", "X-API-Key", openapi.InHeader, "API key in X-API-Key header")
func WithAPIKey(name, paramName string, in ParameterLocation, desc string) Option {
	return func(a *API) {
		if a.SecuritySchemes == nil {
			a.SecuritySchemes = make(map[string]*model.SecurityScheme)
		}
		a.SecuritySchemes[name] = &model.SecurityScheme{
			Type:        "apiKey",
			Name:        paramName,
			In:          string(in),
			Description: desc,
		}
	}
}

// OAuthFlowType represents the type of OAuth2 flow.
type OAuthFlowType string

const (
	// FlowImplicit represents the OAuth2 implicit flow.
	FlowImplicit OAuthFlowType = "implicit"
	// FlowPassword represents the OAuth2 resource owner password flow.
	FlowPassword OAuthFlowType = "password"
	// FlowClientCredentials represents the OAuth2 client credentials flow.
	FlowClientCredentials OAuthFlowType = "clientCredentials"
	// FlowAuthorizationCode represents the OAuth2 authorization code flow.
	FlowAuthorizationCode OAuthFlowType = "authorizationCode"
)

// OAuth2Flow configures a single OAuth2 flow with explicit type.
type OAuth2Flow struct {
	// Type specifies the OAuth2 flow type (implicit, password, clientCredentials, authorizationCode).
	Type OAuthFlowType

	// AuthorizationURL is required for implicit and authorizationCode flows.
	AuthorizationURL string

	// TokenURL is required for password, clientCredentials, and authorizationCode flows.
	TokenURL string

	// RefreshURL is optional for all flows.
	RefreshURL string

	// Scopes maps scope names to descriptions (required, can be empty).
	Scopes map[string]string
}

// WithOAuth2 adds OAuth2 authentication scheme.
//
// At least one flow must be configured. Use OAuth2Flow to configure each flow type.
// Multiple flows can be provided to support different OAuth2 flow types.
//
// Example:
//
//	openapi.WithOAuth2("oauth2", "OAuth2 authentication",
//		openapi.OAuth2Flow{
//			Type:             openapi.FlowAuthorizationCode,
//			AuthorizationURL: "https://example.com/oauth/authorize",
//			TokenURL:         "https://example.com/oauth/token",
//			Scopes: map[string]string{
//				"read":  "Read access",
//				"write": "Write access",
//			},
//		},
//		openapi.OAuth2Flow{
//			Type:     openapi.FlowClientCredentials,
//			TokenUrl: "https://example.com/oauth/token",
//			Scopes:   map[string]string{"read": "Read access"},
//		},
//	)
func WithOAuth2(name, desc string, flows ...OAuth2Flow) Option {
	return func(a *API) {
		if a.SecuritySchemes == nil {
			a.SecuritySchemes = make(map[string]*model.SecurityScheme)
		}
		oauthFlows := &model.OAuthFlows{}
		for _, flow := range flows {
			flowConfig := &model.OAuthFlow{
				AuthorizationURL: flow.AuthorizationURL,
				TokenURL:         flow.TokenURL,
				RefreshURL:       flow.RefreshURL,
				Scopes:           flow.Scopes,
			}
			switch flow.Type {
			case FlowImplicit:
				oauthFlows.Implicit = flowConfig
			case FlowPassword:
				oauthFlows.Password = flowConfig
			case FlowClientCredentials:
				oauthFlows.ClientCredentials = flowConfig
			case FlowAuthorizationCode:
				oauthFlows.AuthorizationCode = flowConfig
			}
		}
		a.SecuritySchemes[name] = &model.SecurityScheme{
			Type:        "oauth2",
			Description: desc,
			Flows:       oauthFlows,
		}
	}
}

// WithOpenIDConnect adds OpenID Connect authentication scheme.
//
// Parameters:
//   - name: Scheme name used in security requirements
//   - url: Well-known URL to discover OpenID Connect provider metadata
//   - desc: Description shown in Swagger UI
//
// Example:
//
//	openapi.WithOpenIDConnect("oidc", "https://example.com/.well-known/openid-configuration", "OpenID Connect authentication")
func WithOpenIDConnect(name, url, desc string) Option {
	return func(a *API) {
		if a.SecuritySchemes == nil {
			a.SecuritySchemes = make(map[string]*model.SecurityScheme)
		}
		a.SecuritySchemes[name] = &model.SecurityScheme{
			Type:             "openIdConnect",
			Description:      desc,
			OpenIDConnectURL: url,
		}
	}
}

// WithDefaultSecurity sets default security requirements applied to all operations.
//
// Operations can override this by specifying their own security requirements
// using RouteWrapper.Security() or RouteWrapper.Bearer().
//
// Example:
//
//	// Apply Bearer auth to all operations by default
//	openapi.WithDefaultSecurity("bearerAuth")
//
//	// Apply OAuth with specific scopes
//	openapi.WithDefaultSecurity("oauth2", "read", "write")
func WithDefaultSecurity(scheme string, scopes ...string) Option {
	return func(a *API) {
		// Ensure scopes is always an empty slice, never nil, per OpenAPI spec
		if scopes == nil {
			scopes = []string{}
		}
		a.DefaultSecurity = append(a.DefaultSecurity, model.SecurityRequirement{
			scheme: scopes,
		})
	}
}

// WithVersion sets the target OpenAPI version.
//
// Example:
//
//	openapi.WithVersion("3.1.2")
func WithVersion(version string) Option {
	return func(a *API) {
		a.Version = version
	}
}

// WithStrictDownlevel causes projection to error (instead of warn) when
// 3.1-only features are used with a 3.0 target.
//
// Default: false (warnings only)
//
// Example:
//
//	openapi.WithStrictDownlevel(true)
func WithStrictDownlevel(strict bool) Option {
	return func(a *API) {
		a.StrictDownlevel = strict
	}
}

// WithValidation enables or disables JSON Schema validation of the generated OpenAPI spec.
//
// When enabled, Generate() validates the output against the official
// OpenAPI meta-schema and returns an error if the spec is invalid.
//
// This is useful for:
//   - Development: Catch spec generation bugs early
//   - CI/CD: Ensure generated specs are valid before deployment
//   - Testing: Verify spec correctness in tests
//
// Performance: Adds ~1-5ms overhead per generation. The default is false
// for backward compatibility. Enable for development and testing to catch
// errors early.
//
// Default: false
//
// Example:
//
//	openapi.WithValidation(false) // Disable for performance
func WithValidation(enabled bool) Option {
	return func(a *API) {
		a.ValidateSpec = enabled
	}
}

// WithExtension adds a specification extension to the root OpenAPI specification.
//
// Extension keys MUST start with "x-". In OpenAPI 3.1.x, keys starting with
// "x-oai-" or "x-oas-" are reserved for the OpenAPI Initiative.
//
// The value can be any valid JSON value (null, primitive, array, or object).
// Validation of extension keys happens during API.Validate().
//
// Example:
//
//	openapi.WithExtension("x-internal-id", "api-v2")
//	openapi.WithExtension("x-code-samples", []map[string]any{
//	    {"lang": "curl", "source": "curl https://api.example.com/users"},
//	})
func WithExtension(key string, value any) Option {
	return func(a *API) {
		if a.Extensions == nil {
			a.Extensions = make(map[string]any)
		}
		a.Extensions[key] = value
	}
}

// WithTagConfig configures struct tag names used for OpenAPI schema generation.
//
// By default, the following tag names are used:
//   - schema: for parameter location/style metadata
//   - body: for request/response body metadata
//   - openapi: for OpenAPI-specific metadata
//   - validate: for validation constraints
//   - default: for default values
//   - requires: for dependent required fields
//
// Use this option to customize tag names for compatibility with other libraries
// or to match your existing codebase conventions.
//
// Partial configurations are supported - only provide the tag names you want
// to customize. Unspecified tags will use their defaults.
//
// Example:
//
//	import "github.com/talav/openapi/config"
//
//	// Customize only specific tags - others use defaults
//	openapi.NewAPI(
//	    openapi.WithTagConfig(config.TagConfig{
//	        Schema:   "param",  // Custom
//	        OpenAPI:  "api",    // Custom
//	        Validate: "rules",  // Custom
//	        // Body, Default, Requires will use defaults ("body", "default", "requires")
//	    }),
//	)
func WithTagConfig(cfg config.TagConfig) Option {
	return func(a *API) {
		a.TagConfig = config.MergeTagConfig(a.TagConfig, cfg)
	}
}

// WithSchemaPrefix sets the prefix for OpenAPI schema references.
// The prefix is used when generating $ref references to schemas in components/schemas.
//
// Default: "#/components/schemas/"
//
// Example:
//
//	openapi.WithSchemaPrefix("#/definitions/")
func WithSchemaPrefix(prefix string) Option {
	return func(a *API) {
		a.SchemaPrefix = prefix
	}
}

// Generate produces an OpenAPI specification from operations.
//
// This is a pure function with no side effects. It takes configuration and operations
// as input and produces JSON/YAML bytes as output. Caching and state management are
// the caller's responsibility.
//
// Example:
//
//	api := openapi.MustNew(
//	    openapi.WithTitle("My API", "1.0.0"),
//	    openapi.WithBearerAuth("bearerAuth", "JWT"),
//	)
//
//	result, err := api.Generate(ctx,
//	    openapi.GET("/users/:id",
//	        openapi.Summary("Get user"),
//	        openapi.Response(200, UserResponse{}),
//	    ),
//	    openapi.POST("/users",
//	        openapi.Summary("Create user"),
//	        openapi.Request(CreateUserRequest{}),
//	        openapi.Response(201, UserResponse{}),
//	    ),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(result.JSON))
func (a *API) Generate(ctx context.Context, ops ...Operation) (*Result, error) {
	spec := a.generateSpec()

	// Process operations and add them to the spec
	if err := a.processOperations(spec, ops); err != nil {
		return nil, fmt.Errorf("failed to process operations: %w", err)
	}

	// Update schemas after operations are processed (they're populated during operation building)
	spec.Components.Schemas = a.generator.Schemas()

	sortSpec(spec)

	if !a.exporter.IsSupportedVersion(a.Version) {
		return nil, fmt.Errorf("unsupported OpenAPI version: %s", a.Version)
	}

	// Export spec
	exportCfg := export.ExporterConfig{
		Version:        a.Version,
		ShouldValidate: a.ValidateSpec,
	}

	result, err := a.exporter.Export(ctx, spec, exportCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to export OpenAPI spec: %w", err)
	}

	return &Result{
		JSON:     result.Result,
		Warnings: result.Warnings,
	}, nil
}

// convertOperationToModel converts a public Operation to model.Operation.
// This uses RequestBuilder and ResponseBuilder to generate the structure,
// then adds examples and customizes content types.
func (a *API) convertOperationToModel(op Operation) (*model.Operation, error) {
	doc := op.doc

	// Convert security requirements
	security := make([]model.SecurityRequirement, 0, len(doc.Security))
	for _, s := range doc.Security {
		security = append(security, model.SecurityRequirement{
			s.Scheme: s.Scopes,
		})
	}

	modelOp := &model.Operation{
		Summary:     doc.Summary,
		Description: doc.Description,
		OperationID: doc.OperationID,
		Tags:        doc.Tags,
		Deprecated:  doc.Deprecated,
		Security:    security,
		Extensions:  copyExtensions(doc.Extensions),
		Responses:   map[string]*model.Response{},
		Parameters:  []model.Parameter{},
	}

	// Build request using RequestBuilder
	if doc.RequestType != nil {
		if err := a.requestBuilder.BuildRequest(modelOp, doc.RequestType); err != nil {
			return nil, fmt.Errorf("failed to build request: %w", err)
		}

		// Add examples to request body if present
		if modelOp.RequestBody != nil && len(doc.RequestNamedExamples) > 0 {
			a.addRequestExamples(modelOp.RequestBody, doc.RequestNamedExamples)
		}
	}

	// Build responses using ResponseBuilder
	if len(doc.ResponseTypes) > 0 {
		if err := a.responseBuilder.BuildOperationResponses(modelOp, doc.ResponseTypes); err != nil {
			return nil, fmt.Errorf("failed to build responses: %w", err)
		}

		// Add examples to responses if present
		if len(doc.ResponseNamedExamples) > 0 {
			a.addResponseExamples(modelOp.Responses, doc.ResponseNamedExamples)
		}
	}

	// Ensure at least one response exists
	if len(modelOp.Responses) == 0 {
		modelOp.Responses[strconv.Itoa(http.StatusOK)] = &model.Response{Description: "OK"}
	}

	return modelOp, nil
}

// addRequestExamples adds named examples to request body media types.
func (a *API) addRequestExamples(reqBody *model.RequestBody, examples []example.Example) {
	for _, content := range reqBody.Content {
		if content.Examples == nil {
			content.Examples = make(map[string]*model.Example)
		}
		for _, ex := range examples {
			m := &model.Example{Summary: ex.Summary(), Description: ex.Description()}
			if ex.IsExternal() {
				m.ExternalValue = ex.ExternalValue()
			} else {
				m.Value = ex.Value()
			}
			content.Examples[ex.Name()] = m
		}
	}
}

// addResponseExamples adds named examples to response media types.
func (a *API) addResponseExamples(responses map[string]*model.Response, examples map[int][]example.Example) {
	for status, exList := range examples {
		statusStr := strconv.Itoa(status)
		if resp, ok := responses[statusStr]; ok && resp.Content != nil {
			for _, content := range resp.Content {
				if content.Examples == nil {
					content.Examples = make(map[string]*model.Example)
				}
				for _, ex := range exList {
					m := &model.Example{Summary: ex.Summary(), Description: ex.Description()}
					if ex.IsExternal() {
						m.ExternalValue = ex.ExternalValue()
					} else {
						m.Value = ex.Value()
					}
					content.Examples[ex.Name()] = m
				}
			}
		}
	}
}

// processOperations processes operations and adds them to the spec.
func (a *API) processOperations(spec *model.Spec, ops []Operation) error {
	// Group operations by path
	byPath := make(map[string][]Operation)
	for _, op := range ops {
		path := convertPathToOpenAPI(op.Path)
		byPath[path] = append(byPath[path], op)
	}

	// Process each path
	for path, pathOps := range byPath {
		pathItem := &model.PathItem{}

		for _, op := range pathOps {
			modelOp, err := a.convertOperationToModel(op)
			if err != nil {
				return fmt.Errorf("failed to convert operation %s %s: %w", op.Method, op.Path, err)
			}

			// Add operation to path item based on HTTP method
			if err := assignOperationToPathItem(pathItem, op.Method, modelOp); err != nil {
				return err
			}
		}

		spec.Paths[path] = pathItem
	}

	return nil
}

// assignOperationToPathItem assigns an operation to the appropriate HTTP method field on a PathItem.
func assignOperationToPathItem(pathItem *model.PathItem, method string, op *model.Operation) error {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		pathItem.Get = op
	case http.MethodPost:
		pathItem.Post = op
	case http.MethodPut:
		pathItem.Put = op
	case http.MethodDelete:
		pathItem.Delete = op
	case http.MethodPatch:
		pathItem.Patch = op
	case http.MethodOptions:
		pathItem.Options = op
	case http.MethodHead:
		pathItem.Head = op
	case http.MethodTrace:
		pathItem.Trace = op
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	return nil
}

// convertPathToOpenAPI converts router path format (/users/:id) to OpenAPI format (/users/{id}).
func convertPathToOpenAPI(path string) string {
	// Convert :param to {param}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if param, ok := strings.CutPrefix(part, ":"); ok {
			parts[i] = "{" + param + "}"
		}
	}

	return strings.Join(parts, "/")
}

// copyExtensions creates a deep copy of extensions map.
func copyExtensions(ext map[string]any) map[string]any {
	if ext == nil {
		return nil
	}
	copied := make(map[string]any, len(ext))
	maps.Copy(copied, ext)

	return copied
}

func (a *API) generateSpec() *model.Spec {
	spec := &model.Spec{
		Info:         a.Info,
		Servers:      a.Servers,
		Tags:         a.Tags,
		Paths:        make(map[string]*model.PathItem),
		Security:     a.DefaultSecurity,
		ExternalDocs: a.ExternalDocs,
		Components: &model.Components{
			Schemas:         a.generator.Schemas(),
			SecuritySchemes: a.SecuritySchemes,
		},
	}

	return spec
}

// sortSpec sorts paths, tags, and components for deterministic output.
func sortSpec(s *model.Spec) {
	// Sort paths
	paths := make([]string, 0, len(s.Paths))
	for p := range s.Paths {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	// Create sorted paths map
	sortedPaths := make(map[string]*model.PathItem, len(paths))
	for _, p := range paths {
		sortedPaths[p] = s.Paths[p]
	}
	s.Paths = sortedPaths

	// Sort tags
	sort.Slice(s.Tags, func(i, j int) bool {
		return s.Tags[i].Name < s.Tags[j].Name
	})

	// Sort component schemas
	if s.Components != nil && s.Components.Schemas != nil {
		schemaNames := make([]string, 0, len(s.Components.Schemas))
		for n := range s.Components.Schemas {
			schemaNames = append(schemaNames, n)
		}
		sort.Strings(schemaNames)

		sortedSchemas := make(map[string]*model.Schema, len(schemaNames))
		for _, n := range schemaNames {
			sortedSchemas[n] = s.Components.Schemas[n]
		}
		s.Components.Schemas = sortedSchemas
	}
}
