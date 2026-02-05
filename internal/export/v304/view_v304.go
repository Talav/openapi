package v304

import (
	"github.com/talav/openapi/internal/export/util"
)

// ViewV304 represents an OpenAPI 3.0.4 specification
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#openapi-object
type ViewV304 struct {
	// This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses.
	OpenAPI string `json:"openapi"`

	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info *InfoV30 `json:"info"`

	// An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of "/".
	Servers []*ServerV30 `json:"servers,omitempty"`

	// The available paths and operations for the API.
	Paths PathsV30 `json:"paths"`

	// An element to hold various schemas for the specification.
	Components *ComponentsV30 `json:"components,omitempty"`

	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition.
	Security []SecurityRequirementV30 `json:"security,omitempty"`

	// A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	Tags []*TagV30 `json:"tags,omitempty"`

	// Additional external documentation.
	ExternalDocs *ExternalDocsV30 `json:"externalDocs,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ViewV304 to inline extensions.
func (s *ViewV304) MarshalJSON() ([]byte, error) {
	type viewV304 ViewV304

	return util.MarshalWithExtensions(viewV304(*s), s.Extensions)
}

// InfoV30 provides metadata about the API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#info-object
type InfoV30 struct {
	// The title of the API.
	Title string `json:"title"`

	// A short description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	TermsOfService string `json:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *ContactV30 `json:"contact,omitempty"`

	// The license information for the exposed API.
	License *LicenseV30 `json:"license,omitempty"`

	// The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for InfoV30 to inline extensions.
func (i *InfoV30) MarshalJSON() ([]byte, error) {
	type infoV30 InfoV30

	return util.MarshalWithExtensions(infoV30(*i), i.Extensions)
}

// ContactV30 information for the exposed API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#contact-object
type ContactV30 struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty"`

	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`

	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `json:"email,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ContactV30 to inline extensions.
func (c *ContactV30) MarshalJSON() ([]byte, error) {
	type contactV30 ContactV30

	return util.MarshalWithExtensions(contactV30(*c), c.Extensions)
}

// LicenseV30 information for the exposed API
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#license-object
type LicenseV30 struct {
	// The license name used for the API.
	Name string `json:"name"`

	// A URL to the license used for the API. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for LicenseV30 to inline extensions.
func (l *LicenseV30) MarshalJSON() ([]byte, error) {
	type licenseV30 LicenseV30

	return util.MarshalWithExtensions(licenseV30(*l), l.Extensions)
}

// ServerV30 represents a server
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#server-object
type ServerV30 struct {
	// A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
	URL string `json:"url"`

	// An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A map between a variable name and its value. The value is used for substitution in the server's URL template.
	Variables map[string]*ServerVariableV30 `json:"variables,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ServerV30 to inline extensions.
func (s *ServerV30) MarshalJSON() ([]byte, error) {
	type serverV30 ServerV30

	return util.MarshalWithExtensions(serverV30(*s), s.Extensions)
}

// ServerVariableV30 represents a server variable for server URL template substitution
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#server-variable-object
type ServerVariableV30 struct {
	// An enumeration of string values to be used if the substitution options are from a limited set.
	Enum []string `json:"enum,omitempty"`

	// The default value to use for substitution, which SHALL be sent if an alternate value is not supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional.
	Default string `json:"default"`

	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ServerVariableV30 to inline extensions.
func (s *ServerVariableV30) MarshalJSON() ([]byte, error) {
	type serverVariableV30 ServerVariableV30

	return util.MarshalWithExtensions(serverVariableV30(*s), s.Extensions)
}

// PathsV30 is a map of paths to PathItem objects
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#paths-object
type PathsV30 map[string]*PathItemV30

// PathItemV30 describes the operations available on a single path
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#path-item-object
type PathItemV30 struct {
	// Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. If there are conflicts between the referenced definition and this Path Item's definition, the behavior is undefined.
	Ref string `json:"$ref,omitempty"`

	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`

	// An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A definition of a GET operation on this path.
	Get *OperationV30 `json:"get,omitempty"`

	// A definition of a PUT operation on this path.
	Put *OperationV30 `json:"put,omitempty"`

	// A definition of a POST operation on this path.
	Post *OperationV30 `json:"post,omitempty"`

	// A definition of a DELETE operation on this path.
	Delete *OperationV30 `json:"delete,omitempty"`

	// A definition of a OPTIONS operation on this path.
	Options *OperationV30 `json:"options,omitempty"`

	// A definition of a HEAD operation on this path.
	Head *OperationV30 `json:"head,omitempty"`

	// A definition of a PATCH operation on this path.
	Patch *OperationV30 `json:"patch,omitempty"`

	// A definition of a TRACE operation on this path.
	Trace *OperationV30 `json:"trace,omitempty"`

	// An alternative server array to service all operations in this path.
	Servers []*ServerV30 `json:"servers,omitempty"`

	// A list of parameters that are applicable to all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []*ParameterV30 `json:"parameters,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for PathItemV30 to inline extensions.
func (p *PathItemV30) MarshalJSON() ([]byte, error) {
	type pathItemV30 PathItemV30

	return util.MarshalWithExtensions(pathItemV30(*p), p.Extensions)
}

// OperationV30 describes a single API operation on a path
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#operation-object
type OperationV30 struct {
	// A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Tags []string `json:"tags,omitempty"`

	// A short summary of what the operation does.
	Summary string `json:"summary,omitempty"`

	// A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocsV30 `json:"externalDocs,omitempty"`

	// Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
	OperationID string `json:"operationId,omitempty"`

	// A list of parameters that are applicable to this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []*ParameterV30 `json:"parameters,omitempty"`

	// The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
	RequestBody *RequestBodyV30 `json:"requestBody,omitempty"`

	// The list of possible responses as they are returned from executing this operation.
	Responses ResponsesV30 `json:"responses,omitempty"`

	// A map of possible out-of band callbacks related to the parent operation. The key value used to identify the callback object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation.
	Callbacks map[string]*CallbackV30 `json:"callbacks,omitempty"`

	// Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
	Security []SecurityRequirementV30 `json:"security,omitempty"`

	// An alternative server array to service this operation. If an alternative server object is specified at the Path Item Object or Root level, it will be overridden by this value.
	Servers []*ServerV30 `json:"servers,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for OperationV30 to inline extensions.
func (o *OperationV30) MarshalJSON() ([]byte, error) {
	type operationV30 OperationV30

	return util.MarshalWithExtensions(operationV30(*o), o.Extensions)
}

// ParameterV30 describes a single operation parameter
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#parameter-object
type ParameterV30 struct {
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
	Schema *SchemaV30 `json:"schema,omitempty"`

	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example, the example value SHALL override the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV30 `json:"examples,omitempty"`

	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry. This field is mutually exclusive with the schema field.
	Content map[string]*MediaTypeV30 `json:"content,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ParameterV30 to inline extensions.
func (p *ParameterV30) MarshalJSON() ([]byte, error) {
	type parameterV30 ParameterV30

	return util.MarshalWithExtensions(parameterV30(*p), p.Extensions)
}

// RequestBodyV30 describes a single request body
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#request-body-object
type RequestBodyV30 struct {
	// A reference to a request body defined in components/requestBodies
	Ref string `json:"$ref,omitempty"`

	// A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*MediaTypeV30 `json:"content"`

	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for RequestBodyV30 to inline extensions.
func (r *RequestBodyV30) MarshalJSON() ([]byte, error) {
	type requestBodyV30 RequestBodyV30

	return util.MarshalWithExtensions(requestBodyV30(*r), r.Extensions)
}

// MediaTypeV30 provides schema and examples for the media type identified by its key
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#media-type-object
type MediaTypeV30 struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *SchemaV30 `json:"schema,omitempty"`

	// Example of the media type. The example object SHOULD be in the correct format as specified by the media type. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema which contains an example, the example value SHALL override the example provided by the schema.
	Example any `json:"example,omitempty"`

	// Examples of the media type. Each example object SHOULD match the media type and specified schema if present. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV30 `json:"examples,omitempty"`

	// A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
	Encoding map[string]*EncodingV30 `json:"encoding,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for MediaTypeV30 to inline extensions.
func (m *MediaTypeV30) MarshalJSON() ([]byte, error) {
	type mediaTypeV30 MediaTypeV30

	return util.MarshalWithExtensions(mediaTypeV30(*m), m.Extensions)
}

// EncodingV30 describes a single encoding definition applied to a single schema property
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#encoding-object
type EncodingV30 struct {
	// The Content-Type for encoding a specific property. Default value depends on the property type: for string with format being binary – application/octet-stream; for other primitive types – text/plain; for object - application/json; for array – the default is defined based on the inner type. The value can be a specific media type (e.g. application/json), a wildcard media type (e.g. image/*), or a comma-separated list of the two types.
	ContentType string `json:"contentType,omitempty"`

	// A map allowing additional information to be provided as headers, for example Content-Disposition. Content-Type is described separately and SHALL be ignored in this section. This property SHALL be ignored if the request body media type is not a multipart.
	Headers map[string]*HeaderV30 `json:"headers,omitempty"`

	// Describes how a specific property value will be serialized depending on its type. See Parameter Object for details on the style property. The behavior follows the same values as query parameters, including default values. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Style string `json:"style,omitempty"`

	// When this is true, property values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. The default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for EncodingV30 to inline extensions.
func (e *EncodingV30) MarshalJSON() ([]byte, error) {
	type encodingV30 EncodingV30

	return util.MarshalWithExtensions(encodingV30(*e), e.Extensions)
}

// ResponsesV30 is a container for the expected responses of an operation
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#responses-object
// ResponsesV30 represents the responses for an operation.
type ResponsesV30 map[string]*ResponseV30

// ResponseV30 describes a single response from an API Operation
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#response-object
type ResponseV30 struct {
	// A reference to a response defined in components/responses
	Ref string `json:"$ref,omitempty"`

	// A short description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`

	// Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Headers map[string]*HeaderV30 `json:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*MediaTypeV30 `json:"content,omitempty"`

	// Links to operations based on the response.
	Links map[string]*LinkV30 `json:"links,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ResponseV30 to inline extensions.
func (r *ResponseV30) MarshalJSON() ([]byte, error) {
	type responseV30 ResponseV30

	return util.MarshalWithExtensions(responseV30(*r), r.Extensions)
}

// SchemaV30 represents a JSON Schema
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#schema-object
type SchemaV30 struct {
	// A reference to a schema defined in components/schemas
	Ref string `json:"$ref,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but SHOULD NOT be sent as part of the request. If the property is marked as readOnly being true and is in the required list, the required will take effect on the response only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	ReadOnly bool `json:"readOnly,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "write only". Therefore, it MAY be sent as part of a request but SHOULD NOT be sent as part of the response. If the property is marked as writeOnly being true and is in the required list, the required will take effect on the request only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	WriteOnly bool `json:"writeOnly,omitempty"`

	// This MAY be used only on properties schemas. It has no effect on root schemas. Adds additional metadata to describe the XML representation of this property.
	XML *XMLV30 `json:"xml,omitempty"`

	// Additional external documentation for this schema.
	ExternalDocs *ExternalDocsV30 `json:"externalDocs,omitempty"`

	// A free-form property to include an example of an instance for this schema. To represent examples that cannot be naturally represented in JSON or YAML, a string value can be used to contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Allows sending a null value for the defined schema. Default value is false.
	Nullable bool `json:"nullable,omitempty"`

	// Adds support for polymorphism. The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description. See Composition and Inheritance for more details.
	Discriminator *DiscriminatorV30 `json:"discriminator,omitempty"`

	// Specifies that a schema is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// A true value adds "null" to the allowed type specified by the type keyword, only if type is explicitly defined within the same Schema Object. Other Schema Object constraints retain their defined behavior, and therefore may disallow the use of null as a value. A false value leaves the specified or default type unmodified. The default value is false.
	NullableInType bool `json:"x-nullable,omitempty"`

	// Any: any type
	AnyOf []*SchemaV30 `json:"anyOf,omitempty"`
	// All: all of the subschemas must be valid
	AllOf []*SchemaV30 `json:"allOf,omitempty"`
	// One: one of the subschemas must be valid
	OneOf []*SchemaV30 `json:"oneOf,omitempty"`
	// Not: the given schema must not be valid
	Not *SchemaV30 `json:"not,omitempty"`

	// Items: for array types, the schema of array items
	Items *SchemaV30 `json:"items,omitempty"`
	// Properties: for object types, the properties of the object
	Properties map[string]*SchemaV30 `json:"properties,omitempty"`
	// AdditionalProperties: for object types, allows additional properties beyond those specified
	AdditionalProperties any `json:"additionalProperties,omitempty"`

	// Description: a description of the schema
	Description string `json:"description,omitempty"`
	// Format: the format of the value (e.g., date, email, etc.)
	Format string `json:"format,omitempty"`
	// Default: the default value for the schema
	Default any `json:"default,omitempty"`
	// Title: a title for the schema
	Title string `json:"title,omitempty"`

	// MultipleOf: for numeric types, the number must be a multiple of this value
	MultipleOf *float64 `json:"multipleOf,omitempty"`
	// Maximum: for numeric types, the maximum allowed value
	Maximum *float64 `json:"maximum,omitempty"`
	// ExclusiveMaximum: whether the maximum is exclusive
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty"`
	// Minimum: for numeric types, the minimum allowed value
	Minimum *float64 `json:"minimum,omitempty"`
	// ExclusiveMinimum: whether the minimum is exclusive
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty"`

	// MaxLength: for string types, the maximum length
	MaxLength *int `json:"maxLength,omitempty"`
	// MinLength: for string types, the minimum length
	MinLength *int `json:"minLength,omitempty"`
	// Pattern: for string types, a regex pattern
	Pattern string `json:"pattern,omitempty"`

	// MaxItems: for array types, the maximum number of items
	MaxItems *int `json:"maxItems,omitempty"`
	// MinItems: for array types, the minimum number of items
	MinItems *int `json:"minItems,omitempty"`
	// UniqueItems: for array types, whether items must be unique
	UniqueItems bool `json:"uniqueItems,omitempty"`

	// MaxProperties: for object types, the maximum number of properties
	MaxProperties *int `json:"maxProperties,omitempty"`
	// MinProperties: for object types, the minimum number of properties
	MinProperties *int `json:"minProperties,omitempty"`
	// Required: for object types, the required properties
	Required []string `json:"required,omitempty"`

	// Enum: the allowed values for the schema
	Enum []any `json:"enum,omitempty"`
	// Type: the type of the schema
	Type string `json:"type,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for SchemaV30 to inline extensions.
func (s *SchemaV30) MarshalJSON() ([]byte, error) {
	type schemaV30 SchemaV30

	return util.MarshalWithExtensions(schemaV30(*s), s.Extensions)
}

// DiscriminatorV30 discriminates types for OneOf, AnyOf, AllOf
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#discriminator-object
type DiscriminatorV30 struct {
	// The name of the property in the payload that will hold the discriminator value.
	PropertyName string `json:"propertyName"`

	// An object to hold mappings between payload values and schema names or references.
	Mapping map[string]string `json:"mapping,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for DiscriminatorV30 to inline extensions.
func (d *DiscriminatorV30) MarshalJSON() ([]byte, error) {
	type discriminatorV30 DiscriminatorV30

	return util.MarshalWithExtensions(discriminatorV30(*d), d.Extensions)
}

// XMLV30 information for XML serialization
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#xml-object
type XMLV30 struct {
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

// MarshalJSON implements json.Marshaler for XMLV30 to inline extensions.
func (x *XMLV30) MarshalJSON() ([]byte, error) {
	type xMLV30 XMLV30

	return util.MarshalWithExtensions(xMLV30(*x), x.Extensions)
}

// ComponentsV30 holds a set of reusable objects for different aspects of the OAS
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#components-object
type ComponentsV30 struct {
	// An object to hold reusable Schema Objects.
	Schemas map[string]*SchemaV30 `json:"schemas,omitempty"`

	// An object to hold reusable Response Objects.
	Responses map[string]*ResponseV30 `json:"responses,omitempty"`

	// An object to hold reusable Parameter Objects.
	Parameters map[string]*ParameterV30 `json:"parameters,omitempty"`

	// An object to hold reusable Example Objects.
	Examples map[string]*ExampleV30 `json:"examples,omitempty"`

	// An object to hold reusable Request Body Objects.
	RequestBodies map[string]*RequestBodyV30 `json:"requestBodies,omitempty"`

	// An object to hold reusable Header Objects.
	Headers map[string]*HeaderV30 `json:"headers,omitempty"`

	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes map[string]*SecuritySchemeV30 `json:"securitySchemes,omitempty"`

	// An object to hold reusable Link Objects.
	Links map[string]*LinkV30 `json:"links,omitempty"`

	// An object to hold reusable Callback Objects.
	Callbacks map[string]*CallbackV30 `json:"callbacks,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ComponentsV30 to inline extensions.
func (c *ComponentsV30) MarshalJSON() ([]byte, error) {
	type componentsV30 ComponentsV30

	return util.MarshalWithExtensions(componentsV30(*c), c.Extensions)
}

// SecurityRequirementV30 lists the required security schemes
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#security-requirement-object
type SecurityRequirementV30 map[string][]string

// SecuritySchemeV30 defines a security scheme that can be used by the operations
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#security-scheme-object
type SecuritySchemeV30 struct {
	// A reference to a security scheme defined in components/securitySchemes
	Ref string `json:"$ref,omitempty"`

	// The type of the security scheme. Valid values are "apiKey", "http", "oauth2", "openIdConnect".
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
	Flows *OAuthFlowsV30 `json:"flows,omitempty"`

	// OpenId Connect URL to discover OAuth2 configuration values. This MUST be in the form of a URL.
	OpenIDConnectURL string `json:"openIdConnectUrl,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for SecuritySchemeV30 to inline extensions.
func (s *SecuritySchemeV30) MarshalJSON() ([]byte, error) {
	type securitySchemeV30 SecuritySchemeV30

	return util.MarshalWithExtensions(securitySchemeV30(*s), s.Extensions)
}

// OAuthFlowsV30 allows configuration of the supported OAuth Flows
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#oauth-flows-object
type OAuthFlowsV30 struct {
	// Configuration for the OAuth Implicit flow
	Implicit *OAuthFlowV30 `json:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow
	Password *OAuthFlowV30 `json:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow (previously called application in OAuth 2.0)
	ClientCredentials *OAuthFlowV30 `json:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow (previously called accessCode in OAuth 2.0)
	AuthorizationCode *OAuthFlowV30 `json:"authorizationCode,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for OAuthFlowsV30 to inline extensions.
func (o *OAuthFlowsV30) MarshalJSON() ([]byte, error) {
	type oAuthFlowsV30 OAuthFlowsV30

	return util.MarshalWithExtensions(oAuthFlowsV30(*o), o.Extensions)
}

// OAuthFlowV30 configuration details for a supported OAuth Flow
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#oauth-flow-object
type OAuthFlowV30 struct {
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

// MarshalJSON implements json.Marshaler for OAuthFlowV30 to inline extensions.
func (o *OAuthFlowV30) MarshalJSON() ([]byte, error) {
	type oAuthFlowV30 OAuthFlowV30

	return util.MarshalWithExtensions(oAuthFlowV30(*o), o.Extensions)
}

// TagV30 adds metadata to a single tag that is used by the Operation Object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#tag-object
type TagV30 struct {
	// The name of the tag.
	Name string `json:"name"`

	// A short description for the tag. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocsV30 `json:"externalDocs,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for TagV30 to inline extensions.
func (t *TagV30) MarshalJSON() ([]byte, error) {
	type tagV30 TagV30

	return util.MarshalWithExtensions(tagV30(*t), t.Extensions)
}

// ExternalDocsV30 allows referencing an external resource for extended documentation
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#external-documentation-object
type ExternalDocsV30 struct {
	// A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// The URL for the target documentation. Value MUST be in the format of a URL.
	URL string `json:"url"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for ExternalDocsV30 to inline extensions.
func (e *ExternalDocsV30) MarshalJSON() ([]byte, error) {
	type externalDocsV30 ExternalDocsV30

	return util.MarshalWithExtensions(externalDocsV30(*e), e.Extensions)
}

// ExampleV30 object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#example-object
type ExampleV30 struct {
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

// MarshalJSON implements json.Marshaler for ExampleV30 to inline extensions.
func (e *ExampleV30) MarshalJSON() ([]byte, error) {
	type exampleV30 ExampleV30

	return util.MarshalWithExtensions(exampleV30(*e), e.Extensions)
}

// HeaderV30 follows the structure of the Parameter Object
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#header-object
type HeaderV30 struct {
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
	Schema *SchemaV30 `json:"schema,omitempty"`

	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example, the example value SHALL override the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*ExampleV30 `json:"examples,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for HeaderV30 to inline extensions.
func (h *HeaderV30) MarshalJSON() ([]byte, error) {
	type headerV30 HeaderV30

	return util.MarshalWithExtensions(headerV30(*h), h.Extensions)
}

// LinkV30 represents a possible design-time link for a response
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#link-object
type LinkV30 struct {
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
	Server *ServerV30 `json:"server,omitempty"`

	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for LinkV30 to inline extensions.
func (l *LinkV30) MarshalJSON() ([]byte, error) {
	type linkV30 LinkV30

	return util.MarshalWithExtensions(linkV30(*l), l.Extensions)
}

// CallbackV30 represents a callback object that can be referenced or defined inline
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.4.md#callback-object
type CallbackV30 struct {
	// A reference to a callback defined in components/callbacks
	Ref string `json:"$ref,omitempty"`
	// A map of possible out-of-band callbacks related to the parent operation.
	// The key value used to identify the callback object is an expression,
	// evaluated at runtime, that identifies a URL to use for the callback operation.
	PathItems map[string]*PathItemV30 `json:"-"`
	// Extensions contains specification extensions (fields prefixed with x-).
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements json.Marshaler for CallbackV30 to inline extensions.
func (c *CallbackV30) MarshalJSON() ([]byte, error) {
	type callbackV30 CallbackV30

	return util.MarshalWithExtensions(callbackV30(*c), c.Extensions)
}
