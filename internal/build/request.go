package build

import (
	"fmt"
	"reflect"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/internal/metadata"
	"github.com/talav/openapi/internal/model"
	"github.com/talav/schema"
)

type RequestBuilder interface {
	BuildRequest(op *model.Operation, inputType reflect.Type) error
}

// requestBuilder extracts OpenAPI request schemas from input struct types.
// It populates operation Parameters and RequestBody based on struct field tags.
type requestBuilder struct {
	generator *SchemaGenerator
	metadata  *schema.Metadata
	tagCfg    config.TagConfig
}

// NewRequestBuilder creates a new request builder.
func NewRequestBuilder(generator *SchemaGenerator, metadata *schema.Metadata, tagCfg config.TagConfig) RequestBuilder {
	return &requestBuilder{
		generator: generator,
		metadata:  metadata,
		tagCfg:    tagCfg,
	}
}

// BuildRequest extracts OpenAPI request schemas from an input struct type
// and populates the operation's Parameters and RequestBody.
//
// This handles:
// - Parameters: Generated from fields with "schema" tag (path, query, header, cookie)
// - Request body: Generated from field with "body" tag
// - Content type: Determined from body field type (defaults to application/json)
// - Required fields: Set based on field type (non-pointer = required) and metadata.
func (rb *requestBuilder) BuildRequest(op *model.Operation, inputType reflect.Type) error {
	// Get struct metadata (parsed and cached)
	structMeta, err := rb.metadata.GetStructMetadata(inputType)
	if err != nil {
		return fmt.Errorf("failed to get struct metadata for type %s: %w", inputType, err)
	}

	// Process parameters (fields with "schema" tag, excluding body)
	// Parameters can be in path, query, header, or cookie locations
	rb.buildParameters(op, structMeta, inputType)

	// Process request body (field with "body" tag)
	// Body is handled separately as it's not a parameter
	if err := rb.buildRequestBody(op, structMeta, inputType); err != nil {
		return fmt.Errorf("failed to build request body schema: %w", err)
	}

	return nil
}

// buildParameters extracts OpenAPI parameters from struct fields with "schema" tag.
// Skips fields with "body" tag (handled separately).
// Only processes valid parameter locations: path, query, header, cookie.
func (rb *requestBuilder) buildParameters(op *model.Operation, structMeta *schema.StructMetadata, inputType reflect.Type) {
	if op.Parameters == nil {
		op.Parameters = make([]model.Parameter, 0, len(structMeta.Fields))
	}

	for i := range structMeta.Fields {
		field := &structMeta.Fields[i]

		// Get schema metadata (must have schema tag)
		schemaMeta, ok := schema.GetTagMetadata[*schema.SchemaMetadata](field, rb.tagCfg.Schema)
		if !ok {
			continue
		}

		// Generate schema for parameter type
		hint := getSchemaHint(inputType, field.StructFieldName, op.OperationID+"Request")
		paramSchema := rb.generator.schema(field.Type, true, hint)
		if paramSchema == nil {
			continue
		}

		// Create and add parameter using values from schema parser
		op.Parameters = append(op.Parameters, model.Parameter{
			Name:        schemaMeta.ParamName,
			Description: rb.getDescription(field),
			In:          string(schemaMeta.Location),
			Required:    schemaMeta.Required,
			Schema:      paramSchema,
			Style:       string(schemaMeta.Style),
			Explode:     schemaMeta.Explode,
		})
	}
}

// getDescription returns the description from openapi metadata for the field, or "" if unset.
func (rb *requestBuilder) getDescription(field *schema.FieldMetadata) string {
	if openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](field, rb.tagCfg.OpenAPI); ok {
		return openAPIMeta.Description
	}

	return ""
}

// buildRequestBody extracts OpenAPI request body from struct field with body tag.
// Initializes RequestBody if needed and sets content type and schema.
func (rb *requestBuilder) buildRequestBody(op *model.Operation, structMeta *schema.StructMetadata, inputType reflect.Type) error {
	// Find body field by checking for body tag
	bodyField := findBodyField(structMeta, rb.tagCfg)
	// No body field - nothing to do
	if bodyField == nil {
		return nil
	}

	// Get body metadata
	bodyMeta, ok := schema.GetTagMetadata[*schema.BodyMetadata](bodyField, rb.tagCfg.Body)
	if !ok {
		return fmt.Errorf("body field missing body metadata")
	}

	// Initialize RequestBody if needed
	if op.RequestBody == nil {
		op.RequestBody = &model.RequestBody{
			Content: make(map[string]*model.MediaType),
		}
	}

	op.RequestBody.Required = isRequestBodyRequired(bodyMeta, bodyField)

	// Determine content type based on BodyType
	contentType := getRequestContentType(bodyMeta.BodyType)

	// Initialize content map entry if needed
	if op.RequestBody.Content[contentType] == nil {
		op.RequestBody.Content[contentType] = &model.MediaType{}
	}

	// Generate and transform body schema based on body type
	hint := getSchemaHint(inputType, bodyField.StructFieldName, op.OperationID+"Request")
	bodySchema, encoding := rb.generateBodySchema(bodyField, bodyMeta, hint)
	op.RequestBody.Content[contentType].Schema = bodySchema
	op.RequestBody.Content[contentType].Encoding = encoding

	return nil
}

// isRequestBodyRequired reports whether the request body must be present (OpenAPI required: true).
//
// Two inputs are considered, in order:
//
//  1. Explicit metadata: bodyMeta.Required from the body tag (e.g. body:"structured,required").
//     When set, the body is required regardless of type.
//
//  2. Type-based inference when metadata does not require it:
//     - Concrete types (struct, string, int, slice, etc.) are required:
//     the field cannot be nil in Go, so the body is always sent.
//     - Pointer types (*T) are optional: the field can be nil, so the body may be omitted.
//     - Interface types (interface{}) are optional: the field can be nil.
//
// Examples:
//
//	Body MyStruct  `body:"structured"`   -> required (non-pointer)
//	Body *MyStruct `body:"structured"`   -> optional (pointer)
//	Body *MyStruct `body:"structured,required"` -> required (explicit flag)
func isRequestBodyRequired(bodyMeta *schema.BodyMetadata, bodyField *schema.FieldMetadata) bool {
	if bodyMeta.Required {
		return true
	}
	// Non-pointer, non-interface types cannot be nil; treat as required.
	k := bodyField.Type.Kind()

	return k != reflect.Pointer && k != reflect.Interface
}

// generateBodySchema generates and transforms the request body schema based on body type.
// Returns the schema and optional encoding map (for multipart).
func (rb *requestBuilder) generateBodySchema(bodyField *schema.FieldMetadata, bodyMeta *schema.BodyMetadata, hint string) (*model.Schema, map[string]*model.Encoding) {
	// Multipart schemas must be inline and excluded from components
	allowRef := bodyMeta.BodyType != schema.BodyTypeMultipart
	if !allowRef {
		rb.generator.markInlineOnly(bodyField.Type, hint)
	}

	bodySchema := rb.generator.schema(bodyField.Type, allowRef, hint)

	// Apply content-type-specific transformations
	switch bodyMeta.BodyType {
	case schema.BodyTypeMultipart:
		bodySchema = transformSchemaForMultipart(bodySchema)

		return bodySchema, extractMultipartEncoding(bodySchema)
	case schema.BodyTypeFile:
		return transformSchemaForBinary(bodySchema), nil
	case schema.BodyTypeStructured:
		return bodySchema, nil
	default:
		return bodySchema, nil
	}
}

// transformSchemaForMultipart transforms a JSON schema for multipart/form-data.
// For multipart:
// - Binary fields ([]byte) should use format: binary (not contentEncoding: base64).
// - This represents raw octet-stream upload, not base64-encoded JSON.
func transformSchemaForMultipart(s *model.Schema) *model.Schema {
	// Create a copy to avoid modifying the original cached schema
	transformed := *s
	transformed.Properties = make(map[string]*model.Schema, len(s.Properties))

	for name, prop := range s.Properties {
		transformed.Properties[name] = transformSchemaForBinary(prop)
	}

	return &transformed
}

// transformSchemaForBinary transforms a schema for file/binary request bodies.
// For file requests, []byte should use format: binary (not contentEncoding: base64).
func transformSchemaForBinary(s *model.Schema) *model.Schema {
	// For []byte fields, change from JSON Schema to OpenAPI binary format
	// In JSON: []byte -> {type: string, contentEncoding: base64, contentMediaType: application/octet-stream}
	// In OpenAPI file request: []byte -> {type: string, format: binary}
	if s.Type == TypeString && s.ContentEncoding == contentEncodingBase64 {
		sCopy := *s
		sCopy.ContentEncoding = ""
		sCopy.ContentMediaType = ""
		sCopy.Format = formatBinary

		return &sCopy
	}

	return s
}

// extractMultipartEncoding creates an encoding object for multipart/form-data.
// Per OpenAPI spec, the encoding object specifies content-type for each part.
func extractMultipartEncoding(s *model.Schema) map[string]*model.Encoding {
	encoding := make(map[string]*model.Encoding)

	for name, prop := range s.Properties {
		// Only add encoding for binary fields (format: binary)
		if prop.Type == TypeString && prop.Format == formatBinary {
			encoding[name] = &model.Encoding{
				ContentType: contentTypeOctetStream,
			}
		}
	}

	if len(encoding) == 0 {
		return nil
	}

	return encoding
}

// getRequestContentType maps BodyType to HTTP content-type for requests.
func getRequestContentType(bodyType schema.BodyType) string {
	switch bodyType {
	case schema.BodyTypeMultipart:
		return contentTypeMultipart
	case schema.BodyTypeFile:
		return contentTypeOctetStream
	case schema.BodyTypeStructured:
		fallthrough
	default:
		return contentTypeJSON
	}
}
