package v312

import (
	"encoding/json"
	"maps"

	"github.com/talav/openapi/internal/export/util"
)

// ViewV312 represents an OpenAPI 3.1.2 specification
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#openapi-object
type ViewV312 struct {
	// This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses.
	OpenAPI string `json:"openapi"`

	// The JSON Schema dialect used by the API. This field is OPTIONAL and defaults to the JSON Schema 2020-12 dialect.
	JSONSchemaDialect string `json:"jsonSchemaDialect,omitempty"`

	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info *InfoV31 `json:"info"`

	// An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of "/".
	Servers []*ServerV31 `json:"servers,omitempty"`

	// The available paths and operations for the API.
	Paths PathsV31 `json:"paths"`

	// An element to hold various schemas for the specification.
	Components *ComponentsV31 `json:"components,omitempty"`

	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition.
	Security []SecurityRequirementV31 `json:"security,omitempty"`

	// A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	Tags []*TagV31 `json:"tags,omitempty"`

	// Additional external documentation.
	ExternalDocs *ExternalDocsV31 `json:"externalDocs,omitempty"`

	// A map of named webhook definitions available in the API. Webhooks are event-driven interactions initiated by the API provider to registered webhook listeners.
	Webhooks map[string]*PathItemV31 `json:"webhooks,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ViewV312 to inline extensions.
func (s *ViewV312) MarshalJSON() ([]byte, error) {
	type viewV312 ViewV312

	return util.MarshalWithExtensions(viewV312(*s), s.Extensions)
}

// InfoV31 provides metadata about the API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#info-object
type InfoV31 struct {
	// The title of the API.
	Title string `json:"title"`

	// A short summary of the API.
	Summary string `json:"summary,omitempty"`

	// A short description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	TermsOfService string `json:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *ContactV31 `json:"contact,omitempty"`

	// The license information for the exposed API.
	License *LicenseV31 `json:"license,omitempty"`

	// The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for InfoV31 to inline extensions.
func (i *InfoV31) MarshalJSON() ([]byte, error) {
	type infoV31 InfoV31

	return util.MarshalWithExtensions(infoV31(*i), i.Extensions)
}

// ContactV31 information for the exposed API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#contact-object
type ContactV31 struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty"`

	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`

	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `json:"email,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ContactV31 to inline extensions.
func (c *ContactV31) MarshalJSON() ([]byte, error) {
	type contactV31 ContactV31

	return util.MarshalWithExtensions(contactV31(*c), c.Extensions)
}

// LicenseV31 information for the exposed API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#license-object
type LicenseV31 struct {
	// The license name used for the API.
	Name string `json:"name"`

	// An SPDX license expression for the API. The identifier field is mutually exclusive with the url field. The value is case sensitive and SHOULD be a valid SPDX license expression.
	Identifier string `json:"identifier,omitempty"`

	// A URL to the license used for the API. MUST be in the format of a URL. The url field is mutually exclusive with the identifier field.
	URL string `json:"url,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for LicenseV31 to inline extensions.
func (l *LicenseV31) MarshalJSON() ([]byte, error) {
	type licenseV31 LicenseV31

	return util.MarshalWithExtensions(licenseV31(*l), l.Extensions)
}

// ServerV31 represents a server
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#server-object
type ServerV31 struct {
	// A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
	URL string `json:"url"`

	// An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A map between a variable name and its value. The value is used for substitution in the server's URL template.
	Variables map[string]*ServerVariableV31 `json:"variables,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ServerV31 to inline extensions.
func (s *ServerV31) MarshalJSON() ([]byte, error) {
	type serverV31 ServerV31

	return util.MarshalWithExtensions(serverV31(*s), s.Extensions)
}

// ServerVariableV31 represents a server variable for server URL template substitution
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#server-variable-object
type ServerVariableV31 struct {
	// An enumeration of string values to be used if the substitution options are from a limited set.
	Enum []string `json:"enum,omitempty"`

	// The default value to use for substitution, which SHALL be sent if an alternate value is not supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional.
	Default string `json:"default"`

	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ServerVariableV31 to inline extensions.
func (s *ServerVariableV31) MarshalJSON() ([]byte, error) {
	type serverVariableV31 ServerVariableV31

	return util.MarshalWithExtensions(serverVariableV31(*s), s.Extensions)
}

// PathsV31 is a map of paths to PathItem objects
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#paths-object
type PathsV31 map[string]*PathItemV31

// PathItemV31 describes the operations available on a single path
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#path-item-object
type PathItemV31 struct {
	// Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. If there are conflicts between the referenced definition and this Path Item's definition, the behavior is undefined.
	Ref string `json:"$ref,omitempty"`

	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`

	// An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A definition of a GET operation on this path.
	Get *OperationV31 `json:"get,omitempty"`

	// A definition of a PUT operation on this path.
	Put *OperationV31 `json:"put,omitempty"`

	// A definition of a POST operation on this path.
	Post *OperationV31 `json:"post,omitempty"`

	// A definition of a DELETE operation on this path.
	Delete *OperationV31 `json:"delete,omitempty"`

	// A definition of a OPTIONS operation on this path.
	Options *OperationV31 `json:"options,omitempty"`

	// A definition of a HEAD operation on this path.
	Head *OperationV31 `json:"head,omitempty"`

	// A definition of a PATCH operation on this path.
	Patch *OperationV31 `json:"patch,omitempty"`

	// A definition of a TRACE operation on this path.
	Trace *OperationV31 `json:"trace,omitempty"`

	// An alternative server array to service all operations in this path.
	Servers []*ServerV31 `json:"servers,omitempty"`

	// A list of parameters that are applicable to all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []*ParameterV31 `json:"parameters,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for PathItemV31 to inline extensions.
func (p *PathItemV31) MarshalJSON() ([]byte, error) {
	type pathItemV31 PathItemV31

	return util.MarshalWithExtensions(pathItemV31(*p), p.Extensions)
}

// OperationV31 describes a single API operation on a path
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#operation-object
type OperationV31 struct {
	// A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Tags []string `json:"tags,omitempty"`

	// A short summary of what the operation does.
	Summary string `json:"summary,omitempty"`

	// A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocsV31 `json:"externalDocs,omitempty"`

	// Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
	OperationID string `json:"operationId,omitempty"`

	// A list of parameters that are applicable to this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []*ParameterV31 `json:"parameters,omitempty"`

	// The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
	RequestBody *RequestBodyV31 `json:"requestBody,omitempty"`

	// The list of possible responses as they are returned from executing this operation.
	Responses map[string]*ResponseV31 `json:"responses,omitempty"`

	// A map of possible out-of band callbacks related to the parent operation. The key value used to identify the callback object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation.
	Callbacks map[string]*CallbackV31 `json:"callbacks,omitempty"`

	// Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
	Security []SecurityRequirementV31 `json:"security,omitempty"`

	// An alternative server array to service this operation. If an alternative server object is specified at the Path Item Object or Root level, it will be overridden by this value.
	Servers []*ServerV31 `json:"servers,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for OperationV31 to inline extensions.
func (o *OperationV31) MarshalJSON() ([]byte, error) {
	type operationV31 OperationV31

	return util.MarshalWithExtensions(operationV31(*o), o.Extensions)
}

// ParameterV31 describes a single operation parameter
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#parameter-object
type ParameterV31 struct {
	// A reference to a parameter defined in components/parameters
	Ref string `json:"$ref,omitempty"`

	// The name of the parameter. Parameter names are case sensitive.
	Name string `json:"name"`

	// The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	In string `json:"in"`

	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Required bool `json:"required,omitempty"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage.
	Deprecated bool `json:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value. Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`

	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for query - form; for path - simple; for header - simple; for cookie - form.
	Style string `json:"style,omitempty"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. This property only applies to parameters with an in value of query. The default value is false.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// The schema defining the type used for the parameter.
	Schema *SchemaV31 `json:"schema,omitempty"`

	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example, the example value SHALL override the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV31 `json:"examples,omitempty"`

	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry. This field is mutually exclusive with the schema field.
	Content map[string]*MediaTypeV31 `json:"content,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ParameterV31 to inline extensions.
func (p *ParameterV31) MarshalJSON() ([]byte, error) {
	type parameterV31 ParameterV31

	return util.MarshalWithExtensions(parameterV31(*p), p.Extensions)
}

// RequestBodyV31 describes a single request body
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#request-body-object
type RequestBodyV31 struct {
	// A reference to a request body defined in components/requestBodies
	Ref string `json:"$ref,omitempty"`

	// A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*MediaTypeV31 `json:"content"`

	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for RequestBodyV31 to inline extensions.
func (r *RequestBodyV31) MarshalJSON() ([]byte, error) {
	type requestBodyV31 RequestBodyV31

	return util.MarshalWithExtensions(requestBodyV31(*r), r.Extensions)
}

// MediaTypeV31 provides schema and examples for the media type identified by its key
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#media-type-object
type MediaTypeV31 struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *SchemaV31 `json:"schema,omitempty"`

	// Example of the media type. The example object SHOULD be in the correct format as specified by the media type. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema which contains an example, the example value SHALL override the example provided by the schema.
	Example any `json:"example,omitempty"`

	// Examples of the media type. Each example object SHOULD match the media type and specified schema if present. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV31 `json:"examples,omitempty"`

	// A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
	Encoding map[string]*EncodingV31 `json:"encoding,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for MediaTypeV31 to inline extensions.
func (m *MediaTypeV31) MarshalJSON() ([]byte, error) {
	type mediaTypeV31 MediaTypeV31

	return util.MarshalWithExtensions(mediaTypeV31(*m), m.Extensions)
}

// EncodingV31 describes a single encoding definition applied to a single schema property
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#encoding-object
type EncodingV31 struct {
	// The Content-Type for encoding a specific property. Default value depends on the property type: for string with format being binary – application/octet-stream; for other primitive types – text/plain; for object - application/json; for array – the default is defined based on the inner type. The value can be a specific media type (e.g. application/json), a wildcard media type (e.g. image/*), or a comma-separated list of the two types.
	ContentType string `json:"contentType,omitempty"`

	// A map allowing additional information to be provided as headers, for example Content-Disposition. Content-Type is described separately and SHALL be ignored in this section. This property SHALL be ignored if the request body media type is not a multipart.
	Headers map[string]*HeaderV31 `json:"headers,omitempty"`

	// Describes how a specific property value will be serialized depending on its type. See Parameter Object for details on the style property. The behavior follows the same values as query parameters, including default values. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Style string `json:"style,omitempty"`

	// When this is true, property values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. The default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for EncodingV31 to inline extensions.
func (e *EncodingV31) MarshalJSON() ([]byte, error) {
	type encodingV31 EncodingV31

	return util.MarshalWithExtensions(encodingV31(*e), e.Extensions)
}

// ResponseV31 describes a single response from an API Operation
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#response-object
type ResponseV31 struct {
	// A reference to a response defined in components/responses
	Ref string `json:"$ref,omitempty"`

	// A short description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`

	// Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Headers map[string]*HeaderV31 `json:"headers,omitempty"`

	// A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*MediaTypeV31 `json:"content,omitempty"`

	// Links to operations based on the response.
	Links map[string]*LinkV31 `json:"links,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ResponseV31 to inline extensions.
func (r *ResponseV31) MarshalJSON() ([]byte, error) {
	type responseV31 ResponseV31

	return util.MarshalWithExtensions(responseV31(*r), r.Extensions)
}

// SchemaV31 represents a JSON Schema (Draft 2020-12)
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#schema-object
type SchemaV31 struct {
	// A reference to a schema defined in components/schemas
	Ref string `json:"$ref,omitempty"`

	// The type of the schema
	Type any `json:"type,omitempty"`

	// Title of the schema
	Title string `json:"title,omitempty"`

	// Format constraint
	Format string `json:"format,omitempty"`

	// Content encoding for binary data
	ContentEncoding string `json:"contentEncoding,omitempty"`

	// Content media type for encoded content
	ContentMediaType string `json:"contentMediaType,omitempty"`

	// Description of the schema
	Description string `json:"description,omitempty"`

	// Default value
	Default any `json:"default,omitempty"`

	// Example value
	Example any `json:"example,omitempty"`

	// Examples array
	Examples []any `json:"examples,omitempty"`

	// Read-only flag
	ReadOnly bool `json:"readOnly,omitempty"`

	// Write-only flag
	WriteOnly bool `json:"writeOnly,omitempty"`

	// Deprecated flag
	Deprecated bool `json:"deprecated,omitempty"`

	// Discriminator for polymorphism
	Discriminator *DiscriminatorV31 `json:"discriminator,omitempty"`

	// XML serialization hints
	XML *XMLV31 `json:"xml,omitempty"`

	// Enum values
	Enum []any `json:"enum,omitempty"`

	// Const value constraint
	Const any `json:"const,omitempty"`

	// All of composition
	AllOf []*SchemaV31 `json:"allOf,omitempty"`

	// Any of composition
	AnyOf []*SchemaV31 `json:"anyOf,omitempty"`

	// One of composition
	OneOf []*SchemaV31 `json:"oneOf,omitempty"`

	// Not composition
	Not *SchemaV31 `json:"not,omitempty"`

	// Items for arrays
	Items *SchemaV31 `json:"items,omitempty"`

	// Prefix items for tuple schemas
	PrefixItems []*SchemaV31 `json:"prefixItems,omitempty"`

	// Contains validation for arrays
	Contains *SchemaV31 `json:"contains,omitempty"`

	// Minimum contains count
	MinContains *int `json:"minContains,omitempty"`

	// Maximum contains count
	MaxContains *int `json:"maxContains,omitempty"`

	// Properties for objects
	Properties map[string]*SchemaV31 `json:"properties,omitempty"`

	// Pattern properties for objects
	PatternProperties map[string]*SchemaV31 `json:"patternProperties,omitempty"`

	// Additional properties for objects
	AdditionalProperties any `json:"additionalProperties,omitempty"`

	// Property names constraint
	PropertyNames *SchemaV31 `json:"propertyNames,omitempty"`

	// Unevaluated properties
	UnevaluatedProperties any `json:"unevaluatedProperties,omitempty"`

	// Required properties for objects
	Required []string `json:"required,omitempty"`

	// Maximum value for numbers
	Maximum *float64 `json:"maximum,omitempty"`

	// Exclusive maximum value for numbers
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`

	// Minimum value for numbers
	Minimum *float64 `json:"minimum,omitempty"`

	// Exclusive minimum value for numbers
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`

	// Multiple of constraint for numbers
	MultipleOf *float64 `json:"multipleOf,omitempty"`

	// Maximum length for strings
	MaxLength *int `json:"maxLength,omitempty"`

	// Minimum length for strings
	MinLength *int `json:"minLength,omitempty"`

	// Pattern for strings
	Pattern string `json:"pattern,omitempty"`

	// Maximum items for arrays
	MaxItems *int `json:"maxItems,omitempty"`

	// Minimum items for arrays
	MinItems *int `json:"minItems,omitempty"`

	// Unique items for arrays
	UniqueItems bool `json:"uniqueItems,omitempty"`

	// Maximum properties for objects
	MaxProperties *int `json:"maxProperties,omitempty"`

	// Minimum properties for objects
	MinProperties *int `json:"minProperties,omitempty"`

	// Additional external documentation for this schema.
	ExternalDocs *ExternalDocsV31 `json:"externalDocs,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for SchemaV31 to inline extensions.
func (s *SchemaV31) MarshalJSON() ([]byte, error) {
	type schemaV31 SchemaV31

	return util.MarshalWithExtensions(schemaV31(*s), s.Extensions)
}

// DiscriminatorV31 discriminates types for OneOf, AnyOf, AllOf
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#discriminator-object
type DiscriminatorV31 struct {
	// The name of the property in the payload that will hold the discriminator value.
	PropertyName string `json:"propertyName"`

	// An object to hold mappings between payload values and schema names or references.
	Mapping map[string]string `json:"mapping,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for DiscriminatorV31 to inline extensions.
func (d *DiscriminatorV31) MarshalJSON() ([]byte, error) {
	type discriminatorV31 DiscriminatorV31

	return util.MarshalWithExtensions(discriminatorV31(*d), d.Extensions)
}

// XMLV31 information for XML serialization
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#xml-object
type XMLV31 struct {
	// Replaces the name of the element/attribute used for the described schema property. When defined within items, it will affect the name of the individual XML elements within the list. When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true. If wrapped is false, it will be ignored.
	Name string `json:"name,omitempty"`

	// The URI of the namespace definition. Value MUST be in the form of an absolute URI.
	Namespace string `json:"namespace,omitempty"`

	// The prefix to be used for the name.
	Prefix string `json:"prefix,omitempty"`

	// Declares whether the property definition translates to an attribute instead of an element. Default value is false.
	Attribute bool `json:"attribute,omitempty"`

	// MAY be used only for an array definition. Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>). Default value is true. The definition takes effect only when defined alongside type being array (outside the items).
	Wrapped bool `json:"wrapped,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for XMLV31 to inline extensions.
func (x *XMLV31) MarshalJSON() ([]byte, error) {
	type xMLV31 XMLV31

	return util.MarshalWithExtensions(xMLV31(*x), x.Extensions)
}

// ComponentsV31 holds a set of reusable objects for different aspects of the OAS
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#components-object
type ComponentsV31 struct {
	// An object to hold reusable Schema Objects.
	Schemas map[string]*SchemaV31 `json:"schemas,omitempty"`

	// An object to hold reusable Response Objects.
	Responses map[string]*ResponseV31 `json:"responses,omitempty"`

	// An object to hold reusable Parameter Objects.
	Parameters map[string]*ParameterV31 `json:"parameters,omitempty"`

	// An object to hold reusable Example Objects.
	Examples map[string]*ExampleV31 `json:"examples,omitempty"`

	// An object to hold reusable Request Body Objects.
	RequestBodies map[string]*RequestBodyV31 `json:"requestBodies,omitempty"`

	// An object to hold reusable Header Objects.
	Headers map[string]*HeaderV31 `json:"headers,omitempty"`

	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes map[string]*SecuritySchemeV31 `json:"securitySchemes,omitempty"`

	// An object to hold reusable Link Objects.
	Links map[string]*LinkV31 `json:"links,omitempty"`

	// An object to hold reusable Callback Objects.
	Callbacks map[string]*CallbackV31 `json:"callbacks,omitempty"`

	// An object to hold reusable Path Item Objects.
	PathItems map[string]*PathItemV31 `json:"pathItems,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ComponentsV31 to inline extensions.
func (c *ComponentsV31) MarshalJSON() ([]byte, error) {
	type componentsV31 ComponentsV31

	return util.MarshalWithExtensions(componentsV31(*c), c.Extensions)
}

// SecurityRequirementV31 lists the required security schemes
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#security-requirement-object
type SecurityRequirementV31 map[string][]string

// SecuritySchemeV31 defines a security scheme that can be used by the operations
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#security-scheme-object
type SecuritySchemeV31 struct {
	// A reference to a security scheme defined in components/securitySchemes
	Ref string `json:"$ref,omitempty"`

	// The type of the security scheme. Valid values are "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect".
	Type string `json:"type"`

	// A short description for security scheme. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// The name of the header, query or cookie parameter to be used.
	Name string `json:"name,omitempty"`

	// The location of the API key. Valid values are "query", "header" or "cookie".
	In string `json:"in,omitempty"`

	// The name of the HTTP Authorization scheme to be used in the Authorization header as defined in RFC7235.
	Scheme string `json:"scheme,omitempty"`

	// A hint to the client to identify how the bearer token is formatted. Bearer tokens are usually generated by an authorization server, so this information is primarily for documentation purposes.
	BearerFormat string `json:"bearerFormat,omitempty"`

	// An object containing configuration information for the flow types supported.
	Flows *OAuthFlowsV31 `json:"flows,omitempty"`

	// OpenId Connect URL to discover OAuth2 configuration values. This MUST be in the form of a URL.
	OpenIDConnectURL string `json:"openIdConnectUrl,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for SecuritySchemeV31 to inline extensions.
func (s *SecuritySchemeV31) MarshalJSON() ([]byte, error) {
	type securitySchemeV31 SecuritySchemeV31

	return util.MarshalWithExtensions(securitySchemeV31(*s), s.Extensions)
}

// OAuthFlowsV31 allows configuration of the supported OAuth Flows
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#oauth-flows-object
type OAuthFlowsV31 struct {
	// Configuration for the OAuth Implicit flow
	Implicit *OAuthFlowV31 `json:"implicit,omitempty"`

	// Configuration for the OAuth Resource Owner Password flow
	Password *OAuthFlowV31 `json:"password,omitempty"`

	// Configuration for the OAuth Client Credentials flow (previously called application in OAuth 2.0)
	ClientCredentials *OAuthFlowV31 `json:"clientCredentials,omitempty"`

	// Configuration for the OAuth Authorization Code flow (previously called accessCode in OAuth 2.0)
	AuthorizationCode *OAuthFlowV31 `json:"authorizationCode,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for OAuthFlowsV31 to inline extensions.
func (o *OAuthFlowsV31) MarshalJSON() ([]byte, error) {
	type oAuthFlowsV31 OAuthFlowsV31

	return util.MarshalWithExtensions(oAuthFlowsV31(*o), o.Extensions)
}

// OAuthFlowV31 configuration details for a supported OAuth Flow
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#oauth-flow-object
type OAuthFlowV31 struct {
	// The authorization URL to be used for this flow. This MUST be in the form of a URL.
	AuthorizationURL string `json:"authorizationUrl,omitempty"`

	// The token URL to be used for this flow. This MUST be in the form of a URL.
	TokenURL string `json:"tokenUrl,omitempty"`

	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL.
	RefreshURL string `json:"refreshUrl,omitempty"`

	// The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it.
	Scopes map[string]string `json:"scopes"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for OAuthFlowV31 to inline extensions.
func (o *OAuthFlowV31) MarshalJSON() ([]byte, error) {
	type oAuthFlowV31 OAuthFlowV31

	return util.MarshalWithExtensions(oAuthFlowV31(*o), o.Extensions)
}

// TagV31 adds metadata to a single tag that is used by the Operation Object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#tag-object
type TagV31 struct {
	// The name of the tag.
	Name string `json:"name"`

	// A short description for the tag. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocsV31 `json:"externalDocs,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for TagV31 to inline extensions.
func (t *TagV31) MarshalJSON() ([]byte, error) {
	type tagV31 TagV31

	return util.MarshalWithExtensions(tagV31(*t), t.Extensions)
}

// ExternalDocsV31 allows referencing an external resource for extended documentation
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#external-documentation-object
type ExternalDocsV31 struct {
	// A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// The URL for the target documentation. Value MUST be in the format of a URL.
	URL string `json:"url"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ExternalDocsV31 to inline extensions.
func (e *ExternalDocsV31) MarshalJSON() ([]byte, error) {
	type externalDocsV31 ExternalDocsV31

	return util.MarshalWithExtensions(externalDocsV31(*e), e.Extensions)
}

// ExampleV31 object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#example-object
type ExampleV31 struct {
	// A reference to an example defined in components/examples
	Ref string `json:"$ref,omitempty"`

	// Short description for the example.
	Summary string `json:"summary,omitempty"`

	// Long description for the example. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Any example value - either a primitive, an object or an array.
	Value any `json:"value,omitempty"`

	// A URL that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents. The value field and externalValue field are mutually exclusive. To represent examples of media types that cannot naturally be represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
	ExternalValue string `json:"externalValue,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ExampleV31 to inline extensions.
func (e *ExampleV31) MarshalJSON() ([]byte, error) {
	type exampleV31 ExampleV31

	return util.MarshalWithExtensions(exampleV31(*e), e.Extensions)
}

// HeaderV31 follows the structure of the Parameter Object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#header-object
type HeaderV31 struct {
	// A reference to a header defined in components/headers
	Ref string `json:"$ref,omitempty"`

	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Required bool `json:"required,omitempty"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage.
	Deprecated bool `json:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value. Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`

	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for query - form; for path - simple; for header - simple; for cookie - form.
	Style string `json:"style,omitempty"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. This property only applies to parameters with an in value of query. The default value is false.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// The schema defining the type used for the parameter.
	Schema *SchemaV31 `json:"schema,omitempty"`

	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example, the example value SHALL override the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV31 `json:"examples,omitempty"`

	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry. This field is mutually exclusive with the schema field.
	Content map[string]*MediaTypeV31 `json:"content,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for HeaderV31 to inline extensions.
func (h *HeaderV31) MarshalJSON() ([]byte, error) {
	type headerV31 HeaderV31

	return util.MarshalWithExtensions(headerV31(*h), h.Extensions)
}

// LinkV31 represents a possible design-time link for a response
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#link-object
type LinkV31 struct {
	// A reference to a link defined in components/links
	Ref string `json:"$ref,omitempty"`

	// A relative or absolute URI reference to an OAS operation. This field is mutually exclusive of the operationId field, and MUST point to an Operation Object. Relative operationRef values MAY be used to locate an existing Operation Object in the OpenAPI definition.
	OperationRef string `json:"operationRef,omitempty"`

	// The name of an existing, resolvable OAS operation, as defined with a unique operationId. This field is mutually exclusive of the operationRef field.
	OperationID string `json:"operationId,omitempty"`

	// A map representing parameters to pass to an operation as specified with operationId or identified via operationRef. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation. The parameter name can be qualified using the parameter location [{in}.]{name} for operations that use the same parameter name in different locations (e.g. path.id).
	Parameters map[string]any `json:"parameters,omitempty"`

	// A literal value or {expression} to use as a request body when calling the target operation.
	RequestBody any `json:"requestBody,omitempty"`

	// A description of the link. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A server object to be used by the target operation.
	Server *ServerV31 `json:"server,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for LinkV31 to inline extensions.
func (l *LinkV31) MarshalJSON() ([]byte, error) {
	type linkV31 LinkV31

	return util.MarshalWithExtensions(linkV31(*l), l.Extensions)
}

// CallbackV31 represents a callback object that can be referenced or defined inline
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.2.md#callback-object
type CallbackV31 struct {
	// A reference to a callback defined in components/callbacks
	Ref string `json:"$ref,omitempty"`

	// A map of possible out-of-band callbacks related to the parent operation.
	// The key value used to identify the callback object is an expression,
	// evaluated at runtime, that identifies a URL to use for the callback operation.
	PathItems map[string]*PathItemV31 `json:"-"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for CallbackV31.
// Callbacks are maps of path expressions to PathItems, so PathItems become the top-level keys.
func (c *CallbackV31) MarshalJSON() ([]byte, error) {
	// Build the map - start with $ref if present, otherwise use PathItems
	m := make(map[string]any, len(c.PathItems)+len(c.Extensions)+1)

	if c.Ref != "" {
		m["$ref"] = c.Ref
	} else {
		for k, v := range c.PathItems {
			m[k] = v
		}
	}

	// Merge extensions
	if len(c.Extensions) > 0 {
		maps.Copy(m, c.Extensions)
	}

	return json.Marshal(m)
}
