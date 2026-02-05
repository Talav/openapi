// Package example provides OpenAPI Example Object support.
//
// Examples demonstrate expected request/response formats for API documentation.
// They can contain inline values or reference external URLs.
//
// Basic usage:
//
//	import "github.com/talav/openapi/example"
//
//	// Create an inline example
//	example.New("user-found", map[string]any{"id": 42, "status": "active"})
//
//	// Add descriptive metadata
//	example.New("error-case", map[string]any{"error": "not found"},
//		example.WithSummary("Resource not found"),
//		example.WithDescription("Returned when the requested resource does not exist"),
//	)
//
//	// Reference an external example file
//	example.NewExternal("full-dataset", "https://example.com/data/full.json")
package example

// Example represents an OpenAPI Example Object.
// https://spec.openapis.org/oas/v3.1.0#example-object
//
// Examples contain either an inline value or an external URL reference.
// Per the spec, these two options are mutually exclusive.
type Example struct {
	// Unique key in the examples map (required). The name serves as the key in the examples map and must be unique within the parent request body or response object.
	name string

	// Short description for the example.
	summary string

	// Long description for the example. [CommonMark](https://commonmark.org/) syntax MAY be used for rich text representation.
	description string

	// Embedded literal example. The value field and externalValue field are mutually exclusive. To represent examples of media types that cannot naturally be represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
	value any

	// A URI that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents. The value field and externalValue field are mutually exclusive.
	externalValue string
}

// Option configures an Example using the functional options pattern.
type Option func(*Example)

// New creates an example with an inline value.
//
// The name serves as the key in the examples map and must be unique within
// the parent request body or response object.
//
// The value will be serialized to JSON in the generated OpenAPI document.
// Use functional options to add summary and description metadata.
//
// Examples:
//
//	example.New("success", map[string]any{"status": "ok"})
//	example.New("with-details", data, example.WithSummary("Detailed response"))
func New(name string, value any, opts ...Option) Example {
	example := Example{
		name:  name,
		value: value,
	}
	for _, opt := range opts {
		opt(&example)
	}

	return example
}

// NewExternal creates an example that references an external URL.
//
// This is useful for large examples that would bloat the specification,
// binary or XML content that cannot be embedded, or examples shared across
// multiple API specifications.
//
// Examples:
//
//	example.NewExternal("large-payload", "https://example.com/samples/large.json")
//	example.NewExternal("xml-sample", "https://example.com/samples/data.xml",
//		example.WithSummary("XML response format"))
func NewExternal(name, url string, opts ...Option) Example {
	example := Example{
		name:          name,
		externalValue: url,
	}
	for _, opt := range opts {
		opt(&example)
	}

	return example
}

// WithSummary adds a short description to the example.
// This typically appears as a title in documentation tools like Swagger UI.
func WithSummary(summary string) Option {
	return func(example *Example) {
		example.summary = summary
	}
}

// WithDescription adds a detailed explanation to the example.
// Supports CommonMark formatting for rich text documentation.
func WithDescription(description string) Option {
	return func(example *Example) {
		example.description = description
	}
}

// Name returns the example's unique identifier.
func (example Example) Name() string { return example.name }

// Value returns the inline example value, or nil for external examples.
func (example Example) Value() any { return example.value }

// ExternalValue returns the external URL, or empty string for inline examples.
func (example Example) ExternalValue() string { return example.externalValue }

// Summary returns the short description.
func (example Example) Summary() string { return example.summary }

// Description returns the detailed description.
func (example Example) Description() string { return example.description }

// IsExternal reports whether this example references an external URL.
func (example Example) IsExternal() bool { return example.externalValue != "" }
