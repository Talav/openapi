package openapi

import (
	"net/http"
	"reflect"

	"github.com/talav/openapi/example"
)

// Operation represents an OpenAPI operation (HTTP method + path + metadata).
// Create operations using the HTTP method constructors: GET, POST, PUT, PATCH, DELETE, etc.
type Operation struct {
	Method string       // HTTP method (GET, POST, etc.)
	Path   string       // URL path with parameters (e.g. "/users/:id")
	doc    operationDoc // Operation documentation (private)
}

// OperationDocOption configures an OpenAPI operation.
// Use with HTTP method constructors like GET, POST, PUT, etc.
type OperationDocOption func(*operationDoc)

// operationDoc holds OpenAPI documentation for a single operation.
// This is private - users interact through Operation and OperationDocOption.
//
// This struct maps to the OpenAPI 3.1.0 Operation Object:
// https://spec.openapis.org/oas/v3.1.0#operation-object
type operationDoc struct {
	// Summary is a short summary of what the operation does.
	// Maps to the "summary" field in the Operation Object.
	Summary string

	// Description provides a verbose explanation of the operation behavior.
	// CommonMark syntax MAY be used for rich text representation.
	// Maps to the "description" field in the Operation Object.
	Description string

	// OperationID is a unique string used to identify the operation.
	// The id MUST be unique among all operations described in the API.
	// The operationId value is case-sensitive. Tools and libraries MAY use
	// the operationId to uniquely identify an operation, therefore, it is
	// RECOMMENDED to follow common programming naming conventions.
	// Maps to the "operationId" field in the Operation Object.
	OperationID string

	// Tags is a list of tags for API documentation control.
	// Tags can be used for logical grouping of operations by resources
	// or any other qualifier. Each tag name in the list MUST be unique.
	// Maps to the "tags" field in the Operation Object.
	Tags []string

	// Deprecated declares this operation to be deprecated.
	// Consumers SHOULD refrain from usage of the declared operation.
	// Default value is false.
	// Maps to the "deprecated" field in the Operation Object.
	Deprecated bool

	// Consumes specifies the MIME types that the operation can consume.
	// This is used to generate the requestBody content map.
	// Defaults to ["application/json"].
	// Implementation detail: not directly in spec, but used to construct
	// the requestBody.content map in the Operation Object.
	Consumes []string

	// Produces specifies the MIME types that the operation can produce.
	// This is used to generate the responses content map.
	// Defaults to ["application/json"].
	// Implementation detail: not directly in spec, but used to construct
	// the responses[statusCode].content map in the Operation Object.
	Produces []string

	// RequestType is the Go type for the request body.
	// Used to generate the requestBody schema in the Operation Object.
	// Implementation detail: not directly in spec, but used to construct
	// the requestBody field in the Operation Object.
	RequestType reflect.Type

	// RequestNamedExamples contains named examples for the request body.
	// These examples are placed in the Media Type Object's "examples" field
	// within requestBody.content[mediaType].examples.
	// Maps to requestBody.content[mediaType].examples in the Operation Object.
	// https://spec.openapis.org/oas/v3.1.0#media-type-object
	RequestNamedExamples []example.Example

	// ResponseTypes maps HTTP status codes to their response Go types.
	// Used to generate the responses map in the Operation Object.
	// Implementation detail: not directly in spec, but used to construct
	// the responses field in the Operation Object.
	ResponseTypes map[int]reflect.Type

	// ResponseNamedExamples maps HTTP status codes to named examples.
	// These examples are placed in the Media Type Object's "examples" field
	// within responses[statusCode].content[mediaType].examples.
	// Maps to responses[statusCode].content[mediaType].examples in the Operation Object.
	// https://spec.openapis.org/oas/v3.1.0#media-type-object
	ResponseNamedExamples map[int][]example.Example

	// Security is a declaration of which security mechanisms can be used
	// for this operation. The list of values includes alternative security
	// requirement objects that can be used. Only one of the security
	// requirement objects need to be satisfied to authorize a request.
	// Maps to the "security" field in the Operation Object.
	Security []SecurityReq

	// Extensions contains specification extensions (x-* fields).
	// Extension keys MUST start with "x-". In OpenAPI 3.1.x, keys starting
	// with "x-oai-" or "x-oas-" are reserved for the OpenAPI Initiative.
	// Maps to extension fields in the Operation Object.
	// https://spec.openapis.org/oas/v3.1.0#specification-extensions
	Extensions map[string]any
}

// SecurityReq represents a security requirement for an operation.
type SecurityReq struct {
	Scheme string
	Scopes []string
}

// newOperation creates an Operation from method, path, and options.
func newOperation(method, path string, opts ...OperationDocOption) Operation {
	doc := operationDoc{
		Consumes:              []string{"application/json"},
		Produces:              []string{"application/json"},
		ResponseTypes:         make(map[int]reflect.Type),
		ResponseNamedExamples: make(map[int][]example.Example),
	}
	for _, opt := range opts {
		opt(&doc)
	}

	return Operation{
		Method: method,
		Path:   path,
		doc:    doc,
	}
}

// GET creates an Operation for a GET request.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithSummary("Get user"),
//	    openapi.WithResponse(200, User{}),
//	)
func GET(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodGet, path, opts...)
}

// POST creates an Operation for a POST request.
//
// Example:
//
//	openapi.POST("/users",
//	    openapi.WithSummary("Create user"),
//	    openapi.WithRequest(CreateUserRequest{}),
//	    openapi.WithResponse(201, User{}),
//	)
func POST(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodPost, path, opts...)
}

// PUT creates an Operation for a PUT request.
//
// Example:
//
//	openapi.PUT("/users/:id",
//	    openapi.WithSummary("Update user"),
//	    openapi.WithRequest(UpdateUserRequest{}),
//	    openapi.WithResponse(200, User{}),
//	)
func PUT(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodPut, path, opts...)
}

// PATCH creates an Operation for a PATCH request.
//
// Example:
//
//	openapi.PATCH("/users/:id",
//	    openapi.WithSummary("Partially update user"),
//	    openapi.WithRequest(PatchUserRequest{}),
//	    openapi.WithResponse(200, User{}),
//	)
func PATCH(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodPatch, path, opts...)
}

// DELETE creates an Operation for a DELETE request.
//
// Example:
//
//	openapi.DELETE("/users/:id",
//	    openapi.WithSummary("Delete user"),
//	    openapi.WithResponse(204, nil),
//	)
func DELETE(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodDelete, path, opts...)
}

// HEAD creates an Operation for a HEAD request.
//
// Example:
//
//	openapi.HEAD("/users/:id",
//	    openapi.WithSummary("Check user exists"),
//	)
func HEAD(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodHead, path, opts...)
}

// OPTIONS creates an Operation for an OPTIONS request.
//
// Example:
//
//	openapi.OPTIONS("/users",
//	    openapi.WithSummary("Get supported methods"),
//	)
func OPTIONS(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodOptions, path, opts...)
}

// TRACE creates an Operation for a TRACE request.
//
// Example:
//
//	openapi.TRACE("/users/:id",
//	    openapi.WithSummary("Trace request"),
//	)
func TRACE(path string, opts ...OperationDocOption) Operation {
	return newOperation(http.MethodTrace, path, opts...)
}

// WithSummary sets the operation summary.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithSummary("Get user by ID"),
//	)
func WithSummary(s string) OperationDocOption {
	return func(d *operationDoc) { d.Summary = s }
}

// WithDescription sets the operation description.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithDescription("Retrieves a user by their unique identifier"),
//	)
func WithDescription(s string) OperationDocOption {
	return func(d *operationDoc) { d.Description = s }
}

// WithOperationID sets a custom operation ID.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithOperationID("getUserById"),
//	)
func WithOperationID(id string) OperationDocOption {
	return func(d *operationDoc) { d.OperationID = id }
}

// WithRequest sets the request type and optionally provides examples.
//
// Example:
//
//	openapi.POST("/users",
//	    openapi.WithRequest(CreateUserRequest{}),
//	)
func WithRequest(req any, examples ...example.Example) OperationDocOption {
	return func(d *operationDoc) {
		d.RequestType = reflect.TypeOf(req)
		if len(examples) > 0 {
			d.RequestNamedExamples = examples
		}
	}
}

// WithResponse sets the response schema and examples for a status code.
//
// Supports two patterns:
//
// 1. Simple pattern - pass the response struct directly:
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	openapi.GET("/users/:id",
//	    openapi.WithResponse(200, User{}),
//	    openapi.WithResponse(404, ErrorModel{}),
//	)
//
// 2. Advanced pattern - wrap with body tag for headers:
//
//	type UserResponse struct {
//	    Body User   `body:"structured"`
//	    ETag string `schema:"ETag,location=header"`
//	}
//
//	openapi.GET("/users/:id",
//	    openapi.WithResponse(200, UserResponse{}),
//	)
//
// With named examples:
//
//	openapi.GET("/users/:id",
//	    openapi.WithResponse(200, User{},
//	        example.New("success", User{ID: 1, Name: "John"}),
//	        example.New("admin", User{ID: 2, Name: "Admin"}),
//	    ),
//	)
func WithResponse(status int, resp any, examples ...example.Example) OperationDocOption {
	return func(d *operationDoc) {
		if resp == nil {
			d.ResponseTypes[status] = nil

			return
		}

		d.ResponseTypes[status] = reflect.TypeOf(resp)
		if len(examples) > 0 {
			d.ResponseNamedExamples[status] = examples
		}
	}
}

// WithTags adds tags to the operation.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithTags("users", "authentication"),
//	)
func WithTags(tags ...string) OperationDocOption {
	return func(d *operationDoc) { d.Tags = append(d.Tags, tags...) }
}

// WithSecurity adds a security requirement.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithSecurity("bearerAuth"),
//	)
//
//	openapi.POST("/users",
//	    openapi.WithSecurity("oauth2", "read:users", "write:users"),
//	)
func WithSecurity(scheme string, scopes ...string) OperationDocOption {
	return func(d *operationDoc) {
		// Ensure scopes is always an empty slice, never nil, per OpenAPI spec
		if scopes == nil {
			scopes = []string{}
		}
		d.Security = append(d.Security, SecurityReq{
			Scheme: scheme,
			Scopes: scopes,
		})
	}
}

// WithDeprecated marks the operation as deprecated.
//
// Example:
//
//	openapi.GET("/old-endpoint",
//	    openapi.WithDeprecated(),
//	)
func WithDeprecated() OperationDocOption {
	return func(d *operationDoc) { d.Deprecated = true }
}

// WithConsumes sets the content types that this operation accepts.
//
// Example:
//
//	openapi.POST("/users",
//	    openapi.WithConsumes("application/xml", "application/json"),
//	)
func WithConsumes(contentTypes ...string) OperationDocOption {
	return func(d *operationDoc) { d.Consumes = contentTypes }
}

// WithProduces sets the content types that this operation returns.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithProduces("application/xml", "application/json"),
//	)
func WithProduces(contentTypes ...string) OperationDocOption {
	return func(d *operationDoc) { d.Produces = contentTypes }
}

// WithOperationExtension adds a specification extension to the operation.
//
// Extension keys MUST start with "x-". In OpenAPI 3.1.x, keys starting with
// "x-oai-" or "x-oas-" are reserved for the OpenAPI Initiative.
//
// Example:
//
//	openapi.GET("/users/:id",
//	    openapi.WithOperationExtension("x-rate-limit", 100),
//	    openapi.WithOperationExtension("x-internal", true),
//	)
func WithOperationExtension(key string, value any) OperationDocOption {
	return func(d *operationDoc) {
		if d.Extensions == nil {
			d.Extensions = make(map[string]any)
		}
		d.Extensions[key] = value
	}
}

// WithOptions composes multiple OperationDocOptions into a single option.
//
// This enables creating reusable option sets for common patterns across operations.
// Options are applied in the order they are provided, with later options potentially
// overriding values set by earlier options.
//
// Example:
//
//	// Define reusable option sets
//	var (
//	    CommonErrors = openapi.WithOptions(
//	        openapi.WithResponse(400, Error{}),
//	        openapi.WithResponse(401, Error{}),
//	        openapi.WithResponse(500, Error{}),
//	    )
//
//	    AuthRequired = openapi.WithOptions(
//	        openapi.WithSecurity("jwt"),
//	    )
//
//	    UserEndpoint = openapi.WithOptions(
//	        openapi.WithTags("users"),
//	        AuthRequired,
//	        CommonErrors,
//	    )
//	)
//
//	// Apply composed options to operations
//	openapi.GET("/users/:id",
//	    UserEndpoint,
//	    openapi.WithSummary("Get user"),
//	    openapi.WithResponse(200, User{}),
//	)
//
//	openapi.POST("/users",
//	    UserEndpoint,
//	    openapi.WithSummary("Create user"),
//	    openapi.WithRequest(CreateUser{}),
//	    openapi.WithResponse(201, User{}),
//	)
func WithOptions(opts ...OperationDocOption) OperationDocOption {
	return func(d *operationDoc) {
		for _, opt := range opts {
			opt(d)
		}
	}
}
