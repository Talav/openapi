package model

// Spec represents a version-agnostic OpenAPI specification.
//
// This model supports all features from both OpenAPI 3.0.x and 3.1.x. Version-specific
// differences are handled by views.
type Spec struct {
	// Info contains API metadata (title, version, description, contact, license).
	Info Info

	// Servers lists available server URLs for the API.
	Servers []Server

	// Paths maps path patterns to PathItem objects containing operations.
	Paths map[string]*PathItem

	// Components holds reusable schemas, security schemes, etc.
	Components *Components

	// Webhooks defines webhook endpoints (3.1 feature).
	// In 3.0, this will be dropped with a warning.
	Webhooks map[string]*PathItem

	// Tags provides additional metadata for operations.
	Tags []Tag

	// Security defines global security requirements applied to all operations.
	Security []SecurityRequirement

	// ExternalDocs provides external documentation links.
	ExternalDocs *ExternalDocs

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any
}

// Info object that provides metadata about the API. The metadata MAY be used by
// the clients if needed, and MAY be presented in editing or documentation
// generation tools for convenience.
//
//	title: Sample Pet Store App
//	summary: A pet store manager.
//	description: This is a sample server for a pet store.
//	termsOfService: https://example.com/terms/
//	contact:
//	  name: API Support
//	  url: https://www.example.com/support
//	  email: support@example.com
//	license:
//	  name: Apache 2.0
//	  url: https://www.apache.org/licenses/LICENSE-2.0.html
//	version: 1.0.1
type Info struct {
	// Title of the API.
	Title string

	// Summary of the API. 3.1+ only: short summary of the API
	Summary string

	// Description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string

	// TermsOfService URL for the API.
	TermsOfService string

	// Contact information to get support for the API.
	Contact *Contact

	// License name & link for using the API.
	License *License

	// Version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string

	// Extensions (user-defined properties), if any. Values in this map will
	// be marshalled as siblings of the other properties above.
	Extensions map[string]any
}

// Contact information to get support for the API.
//
//	name: API Support
//	url: https://www.example.com/support
//	email: support@example.com
type Contact struct {
	// Name of the contact person/organization.
	Name string

	// URL pointing to the contact information.
	URL string

	// Email address of the contact person/organization.
	Email string

	// Extensions (user-defined properties), if any. Values in this map will
	// be marshalled as siblings of the other properties above.
	Extensions map[string]any
}

// License name & link for using the API.
//
//	name: Apache 2.0
//	identifier: Apache-2.0
type License struct {
	// Name of the license.
	Name string

	// Identifier SPDX license expression for the API. This field is mutually
	// exclusive with the URL field.
	Identifier string

	// URL pointing to the license. This field is mutually exclusive with the
	// Identifier field.
	URL string

	// Extensions (user-defined properties), if any. Values in this map will
	// be marshalled as siblings of the other properties above.
	Extensions map[string]any
}

// Server represents a server URL and optional description.
type Server struct {
	// URL to the target host. Supports Server Variables and MAY be relative.
	URL string

	// Description of the host designated by the URL.
	Description string

	// Variables for server URL template substitution.
	Variables map[string]*ServerVariable

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// ServerVariable represents a variable for server URL template substitution.
type ServerVariable struct {
	// Enumeration of allowed values (optional).
	Enum []string

	// REQUIRED. Default value for substitution.
	Default string

	// Description for the server variable.
	Description string

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// PathItem represents the operations available on a single path.
type PathItem struct {
	// If set, this PathItem is a reference ($ref).
	Ref string

	// Summary intended to apply to all operations in this path.
	Summary string

	// Description intended to apply to all operations in this path.
	Description string

	// HTTP method operations.
	Get     *Operation
	Put     *Operation
	Post    *Operation
	Delete  *Operation
	Options *Operation
	Head    *Operation
	Patch   *Operation
	Trace   *Operation

	// Alternative server array to service all operations in this path.
	Servers []Server

	// Parameters applicable to all operations described under this path.
	Parameters []Parameter

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Operation describes a single API operation on a path.
type Operation struct {
	// Tags for API documentation control.
	Tags []string

	// Short summary of what the operation does.
	Summary string

	// Verbose explanation of the operation behavior.
	Description string

	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocs

	// Unique string used to identify the operation.
	OperationID string

	// Parameters applicable to this operation.
	Parameters []Parameter

	// Request body applicable for this operation.
	RequestBody *RequestBody

	// REQUIRED. List of possible responses as they are returned from executing this operation.
	Responses map[string]*Response

	// Map of possible out-of-band callbacks related to the parent operation.
	Callbacks map[string]*Callback

	// Declares this operation to be deprecated.
	Deprecated bool

	// Security mechanisms that can be used for this operation.
	Security []SecurityRequirement

	// Alternative server array to service this operation.
	Servers []Server

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Parameter describes a single operation parameter.
type Parameter struct {
	// If set, this is a $ref.
	Ref string

	// REQUIRED. Name of the parameter.
	Name string

	// REQUIRED. Location of the parameter: query, path, header, cookie.
	In string

	// Description of the parameter.
	Description string

	// Determines whether this parameter is mandatory.
	Required bool

	// Specifies that a parameter is deprecated.
	Deprecated bool

	// Sets the ability to pass empty-valued parameters.
	AllowEmptyValue bool

	// Describes how the parameter value will be serialized.
	Style string

	// When true, parameter values of type array or object generate separate parameters for each value.
	Explode bool

	// Determines whether the parameter value SHOULD allow reserved characters.
	AllowReserved bool

	// Schema defining the type used for the parameter.
	Schema *Schema

	// Example of the parameter's potential value.
	Example any

	// Examples of the parameter's potential value (3.1 style).
	Examples map[string]*Example

	// Map containing the representations for the parameter (mutually exclusive with schema).
	Content map[string]*MediaType

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// RequestBody describes a single request body.
type RequestBody struct {
	// If set, this is a $ref.
	Ref string

	// Description of the request body.
	Description string

	// Determines if the request body is required in the request.
	Required bool

	// REQUIRED. Content of the request body.
	Content map[string]*MediaType

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Response describes a single response from an API operation.
type Response struct {
	// If set, this is a $ref.
	Ref string

	// REQUIRED. Description of the response.
	Description string

	// Map containing descriptions of potential response payloads.
	Content map[string]*MediaType

	// Map of headers that are supported with this response.
	Headers map[string]*Header

	// Map of operations links that can be followed from the response.
	Links map[string]*Link

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Header represents a response header.
type Header struct {
	// If set, this is a $ref.
	Ref string

	// Description of the header.
	Description string

	// Determines whether this header is mandatory.
	Required bool

	// Specifies that a header is deprecated.
	Deprecated bool

	// Sets the ability to pass empty-valued headers.
	AllowEmptyValue bool

	// Describes how the header value will be serialized.
	Style string

	// When true, header values of type array or object generate separate headers for each value.
	Explode bool

	// Schema defining the type used for the header.
	Schema *Schema

	// Example of the header's potential value.
	Example any

	// Examples of the header's potential value.
	Examples map[string]*Example

	// Map containing the representations for the header (mutually exclusive with schema).
	Content map[string]*MediaType

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// MediaType provides schema and examples for a specific content type.
type MediaType struct {
	// Schema defining the type used for the media type.
	Schema *Schema

	// Example of the media type.
	Example any

	// Examples of the media type.
	Examples map[string]*Example

	// Map between a property name and its encoding information.
	Encoding map[string]*Encoding

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Encoding describes encoding for a single schema property.
type Encoding struct {
	// Content type for encoding a specific property.
	ContentType string

	// Map of additional headers to be sent with the request.
	Headers map[string]*Header

	// Describes how a specific property value will be serialized.
	Style string

	// When true, property values of type array or object generate separate parameters for each value.
	Explode bool

	// Determines whether the parameter value SHOULD allow reserved characters.
	AllowReserved bool

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Example represents an example value with optional description.
type Example struct {
	// If set, this is a $ref.
	Ref string

	// Short description for the example.
	Summary string

	// Long description for the example.
	Description string

	// Embedded literal example.
	Value any

	// URI that points to a literal example.
	ExternalValue string

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Link represents a possible design-time link for a response.
type Link struct {
	// If set, this is a $ref.
	Ref string

	// Relative or absolute URI reference to an OAS operation.
	OperationRef string

	// Name of an existing, resolvable OAS operation.
	OperationID string

	// Map representing parameters to pass to an operation.
	Parameters map[string]any

	// Request body to pass to an operation.
	RequestBody any

	// Description of the link.
	Description string

	// Server object to be used by the target operation.
	Server *Server

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Callback represents a callback definition.
type Callback struct {
	// If set, this is a $ref.
	Ref string

	// Map of possible out-of-band callbacks related to the parent operation.
	PathItems map[string]*PathItem

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Components holds reusable components.
type Components struct {
	// Reusable Schema Objects.
	Schemas map[string]*Schema

	// Reusable Response Objects.
	Responses map[string]*Response

	// Reusable Parameter Objects.
	Parameters map[string]*Parameter

	// Reusable Example Objects.
	Examples map[string]*Example

	// Reusable Request Body Objects.
	RequestBodies map[string]*RequestBody

	// Reusable Header Objects.
	Headers map[string]*Header

	// Reusable Security Scheme Objects.
	SecuritySchemes map[string]*SecurityScheme

	// Reusable Link Objects.
	Links map[string]*Link

	// Reusable Callback Objects.
	Callbacks map[string]*Callback

	// Reusable Path Item Objects (3.1 feature).
	// In 3.0, this will be dropped with a warning.
	PathItems map[string]*PathItem

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// SecurityScheme defines a security scheme.
type SecurityScheme struct {
	// If set, this is a $ref.
	Ref string

	// REQUIRED. Type of the security scheme.
	Type string

	// Description for the security scheme.
	Description string

	// Name of the header, query or cookie parameter to be used (for apiKey).
	Name string

	// Location of the API key (for apiKey): query, header, cookie.
	In string

	// Name of the HTTP Authorization scheme (for http).
	Scheme string

	// Hint to the client to identify how the bearer token is formatted (for http bearer).
	BearerFormat string

	// Configuration for OAuth2 flows (for oauth2).
	Flows *OAuthFlows

	// OpenId Connect URL to discover OAuth2 configuration values (for openIdConnect).
	OpenIDConnectURL string

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow.
	Implicit *OAuthFlow

	// Configuration for the OAuth Resource Owner Password flow.
	Password *OAuthFlow

	// Configuration for the OAuth Client Credentials flow.
	ClientCredentials *OAuthFlow

	// Configuration for the OAuth Authorization Code flow.
	AuthorizationCode *OAuthFlow

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// OAuthFlow contains configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	// REQUIRED. Authorization URL for implicit and authorizationCode flows.
	AuthorizationURL string

	// REQUIRED. Token URL for password, clientCredentials, and authorizationCode flows.
	TokenURL string

	// URL for obtaining refresh tokens.
	RefreshURL string

	// REQUIRED. Map of scope name to description.
	Scopes map[string]string

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// SecurityRequirement lists required security schemes for an operation.
type SecurityRequirement map[string][]string

// Tag adds metadata to a tag.
type Tag struct {
	// REQUIRED. Name of the tag.
	Name string

	// Description of the tag.
	Description string

	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocs

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// ExternalDocs provides external documentation links.
type ExternalDocs struct {
	// Description of the target documentation.
	Description string

	// REQUIRED. URL for the target documentation.
	URL string

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Schema represents a version-agnostic JSON Schema.
// This IR supports all features from both OpenAPI 3.0.x and 3.1.x.
// Version-specific differences are handled by projectors in the export package.
type Schema struct {
	// Ref is a logical reference to a component schema.
	Ref string

	// Type of the schema (string representation).
	Type string

	// Nullable indicates if the value can be null.
	// In 3.0: represented as nullable: true
	// In 3.1: represented as type: ["T", "null"]
	Nullable bool

	// Title provides a title for the schema.
	Title string

	// Description provides documentation for the schema.
	Description string

	// Format provides additional type information.
	Format string

	// ContentEncoding specifies the encoding for binary data (3.1 feature).
	// In 3.0: converted to format: "byte"
	ContentEncoding string

	// ContentMediaType specifies the media type of binary content (3.1 feature).
	// In 3.0: omitted (not supported)
	ContentMediaType string

	// Deprecated marks the schema as deprecated.
	Deprecated bool

	// ReadOnly indicates the value is read-only.
	ReadOnly bool

	// WriteOnly indicates the value is write-only.
	WriteOnly bool

	// Example provides a single example value (3.0 style).
	Example any

	// Examples provides multiple example values (3.1 style).
	Examples []any

	// Pattern is a regex pattern for string validation.
	Pattern string

	// MinLength is the minimum string length.
	MinLength *int

	// MaxLength is the maximum string length.
	MaxLength *int

	// Minimum is the minimum numeric value (with exclusive flag).
	// Projectors will convert this to version-specific format.
	Minimum *Bound

	// Maximum is the maximum numeric value (with exclusive flag).
	// Projectors will convert this to version-specific format.
	Maximum *Bound

	// MultipleOf constrains numbers to be multiples of this value.
	MultipleOf *float64

	// Items defines the item schema for arrays.
	Items *Schema

	// MinItems is the minimum number of items in an array.
	MinItems *int

	// MaxItems is the maximum number of items in an array.
	MaxItems *int

	// UniqueItems indicates array items must be unique.
	UniqueItems bool

	// Properties defines object properties.
	Properties map[string]*Schema

	// Required lists required property names (for type "object").
	Required []string

	// DependentRequired specifies fields that become required when a given field is present
	// (JSON Schema 2019-09 / OpenAPI 3.1 feature).
	// Key: property name that, when present, makes other properties required.
	// Value: array of property names that become required when the key property is present.
	// In 3.0, this will be dropped with a warning.
	DependentRequired map[string][]string

	// Additional controls additionalProperties behavior.
	Additional *Additional

	// PatternProps defines pattern-based properties (3.1 feature).
	// In 3.0, this will be dropped with a warning.
	PatternProps map[string]*Schema

	// Unevaluated defines unevaluatedProperties schema (3.1 feature).
	// In 3.0, this will be dropped with a warning.
	Unevaluated *Schema

	// MinProperties is the minimum number of properties in an object.
	MinProperties *int

	// MaxProperties is the maximum number of properties in an object.
	MaxProperties *int

	// AllOf represents an allOf composition.
	AllOf []*Schema

	// AnyOf represents an anyOf composition.
	AnyOf []*Schema

	// OneOf represents a oneOf composition.
	OneOf []*Schema

	// Not represents a not composition.
	Not *Schema

	// Enum lists allowed values for the schema.
	Enum []any

	// Const is a constant value (3.1 feature).
	// In 3.0, this will be converted to enum: [const] with a warning.
	Const any

	// Default is the default value for the schema.
	Default any

	// Discriminator is used for polymorphism (optional).
	Discriminator *Discriminator

	// XML provides XML serialization hints (optional).
	XML *XML

	// ExternalDocs provides external documentation links (optional).
	ExternalDocs *ExternalDocs

	// Extensions (user-defined properties), if any.
	Extensions map[string]any
}

// Bound represents a numeric bound (minimum or maximum) with exclusive flag.
//
// In OpenAPI 3.0, exclusive bounds are represented as boolean flags:
//   - minimum: 10, exclusiveMinimum: true
//
// In OpenAPI 3.1, exclusive bounds are represented as numeric values:
//   - exclusiveMinimum: 10 (instead of minimum: 10)
type Bound struct {
	// Value is the numeric bound value.
	Value float64

	// Exclusive indicates if the bound is exclusive.
	Exclusive bool
}

// Additional represents additionalProperties configuration for objects.
//
// The semantics are:
//   - nil => not specified (JSON Schema default: true)
//   - Allow != nil && *Allow == false && Schema == nil => additionalProperties: false (strict)
//   - Allow != nil && *Allow == true && Schema == nil => additionalProperties: true (explicit allow-all)
//   - Schema != nil => additionalProperties: <schema> (takes precedence over Allow)
type Additional struct {
	// Allow controls whether additional properties are allowed.
	// If nil, the behavior is unspecified (defaults to true in JSON Schema).
	// If Schema is non-nil, it takes precedence over Allow.
	Allow *bool

	// Schema defines the schema for additional properties.
	// If set, this takes precedence over Allow.
	Schema *Schema
}

// Discriminator is used for polymorphism in oneOf/allOf compositions.
type Discriminator struct {
	// REQUIRED. Name of the property in the payload that will hold the discriminator value.
	PropertyName string

	// Object to hold mappings between payload values and schema names or references.
	Mapping map[string]string
}

// XML provides XML serialization hints.
type XML struct {
	// Replaces the name of the element/attribute used for the described schema property.
	Name string

	// URI of the namespace definition.
	Namespace string

	// Prefix to be used for the name.
	Prefix string

	// Declares whether the property definition translates to an attribute instead of an element.
	Attribute bool

	// MAY be used only for an array definition. Signifies whether the array is wrapped.
	Wrapped bool
}
