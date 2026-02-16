package build

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/internal/metadata"
	"github.com/talav/openapi/internal/model"
	"github.com/talav/schema"
)

// BaseRoute represents route information needed for response extraction.
type BaseRoute struct {
	Operation     *model.Operation
	DefaultStatus int
	Errors        []int
}

type ResponseBuilder interface {
	BuildOperationResponses(op *model.Operation, responses map[int]reflect.Type) error
}

// ContentTypeProvider allows you to override the content type for responses,
// allowing you to return a different content type like
// `application/problem+json` after using the `application/json` marshaller.
// This should be implemented by the response body struct.
type ContentTypeProvider interface {
	ContentType(string) string
}

// ResponseSchemaExtractor extracts OpenAPI response schemas from output struct types.
type responseBuilder struct {
	generator *SchemaGenerator
	metadata  *schema.Metadata
	tagCfg    config.TagConfig
}

// NewResponseBuilder creates a new response builder.
func NewResponseBuilder(generator *SchemaGenerator, metadata *schema.Metadata, tagCfg config.TagConfig) ResponseBuilder {
	return &responseBuilder{
		generator: generator,
		metadata:  metadata,
		tagCfg:    tagCfg,
	}
}

func (rb *responseBuilder) BuildOperationResponses(op *model.Operation, responses map[int]reflect.Type) error {
	// Initialize response
	if op.Responses == nil {
		op.Responses = make(map[string]*model.Response)
	}
	for status, response := range responses {
		if err := rb.buildOperationResponse(op, status, response); err != nil {
			return err
		}
	}

	return nil
}

func (rb *responseBuilder) buildOperationResponse(op *model.Operation, status int, response reflect.Type) error {
	structMeta, err := rb.metadata.GetStructMetadata(response)
	if err != nil {
		return fmt.Errorf("failed to get struct metadata for type %s: %w", response, err)
	}

	resp := getResponse(op, status)

	// Extract body schema - handles both tagged fields and plain structs
	if err := rb.extractBodySchema(structMeta, resp, op.OperationID); err != nil {
		return err
	}

	// Extract headers only when using wrapper pattern
	rb.buildResponseHeaders(structMeta, resp)

	return nil
}

// extractBodySchema extracts the body schema and adds it to the response.
// Supports both wrapper pattern (bodyField != nil) and plain struct pattern (bodyField == nil).
func (rb *responseBuilder) extractBodySchema(structMeta *schema.StructMetadata, resp *model.Response, operationID string) error {
	var bodyType reflect.Type
	var bodyMeta *schema.BodyMetadata
	var hint string
	var schemaBodyType schema.BodyType

	// Find body field with body tag
	bodyField := findBodyField(structMeta, rb.tagCfg)
	if bodyField != nil {
		// Wrapper pattern: extract from tagged field
		var ok bool
		bodyMeta, ok = schema.GetTagMetadata[*schema.BodyMetadata](bodyField, rb.tagCfg.Body)
		if !ok {
			return fmt.Errorf("body field missing body metadata")
		}
		bodyType = bodyField.Type
		schemaBodyType = bodyMeta.BodyType
		hint = getSchemaHint(structMeta.Type, bodyField.StructFieldName, operationID)
	} else {
		// Plain struct pattern: use entire struct as body
		bodyType = structMeta.Type
		schemaBodyType = schema.BodyTypeStructured
		hint = getSchemaHint(structMeta.Type, "Response", operationID)
	}

	// Determine content type
	ct := rb.determineContentType(bodyType, schemaBodyType)

	// Generate schema
	bodySchema := rb.generator.schema(bodyType, true, hint)
	if bodyMeta != nil && bodyMeta.BodyType == schema.BodyTypeFile {
		bodySchema = transformSchemaForFileResponse(bodySchema)
	}

	// Set response content
	resp.Content[ct] = &model.MediaType{
		Schema: bodySchema,
	}

	return nil
}

// determineContentType determines the content type for a response body.
// Uses bodyMeta if available (wrapper pattern), otherwise defaults to JSON.
func (rb *responseBuilder) determineContentType(bodyType reflect.Type, bodySchemaType schema.BodyType) string {
	ct := contentTypeJSON

	switch bodySchemaType {
	case schema.BodyTypeStructured:
		ct = contentTypeJSON
	case schema.BodyTypeFile:
		ct = contentTypeOctetStream
	case schema.BodyTypeMultipart:
		// Multipart is not valid for responses, but we'll default to JSON
		// The validation will be caught elsewhere if needed
		ct = contentTypeJSON
	}

	// Check if type implements ContentTypeProvider interface
	contentTypeProviderType := reflect.TypeOf((*ContentTypeProvider)(nil)).Elem()
	if reflect.PointerTo(bodyType).Implements(contentTypeProviderType) {
		instance, ok := reflect.New(bodyType).Interface().(ContentTypeProvider)
		if ok {
			ct = instance.ContentType(ct)
		}
	}

	return ct
}

// buildResponseHeaders extracts header schemas from fields with "schema" tag and location=header
// and adds them to the success response.
func (rb *responseBuilder) buildResponseHeaders(structMeta *schema.StructMetadata, response *model.Response) {
	if response.Headers == nil {
		response.Headers = make(map[string]*model.Header)
	}

	// Iterate through metadata fields
	for _, fieldMeta := range structMeta.Fields {
		// Only process fields with schema tag and location=header
		schemaMeta, ok := schema.GetTagMetadata[*schema.SchemaMetadata](&fieldMeta, rb.tagCfg.Schema)
		if !ok {
			continue
		}

		if schemaMeta.Location != schema.LocationHeader {
			continue
		}

		headerName := schemaMeta.ParamName

		// Get field type for schema generation
		fieldType := fieldMeta.Type
		if fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
		}

		// Check if field implements fmt.Stringer (will be serialized as string)
		if reflect.PointerTo(fieldType).Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
			fieldType = reflect.TypeOf("")
		}

		// Generate schema for header
		hint := getSchemaHint(structMeta.Type, fieldMeta.StructFieldName, headerName)
		headerSchema := rb.generator.schema(fieldType, true, hint)

		// Get description from openapi metadata if available
		description := ""
		if openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](&fieldMeta, rb.tagCfg.OpenAPI); ok {
			description = openAPIMeta.Description
		}

		// Create header parameter
		response.Headers[headerName] = &model.Header{
			Schema:      headerSchema,
			Description: description,
		}
	}
}

// getResponse ensures a response exists for the given status code.
// If the response doesn't exist, it creates one with the provided description.
// If description is empty, it uses the HTTP status text.
// Returns the response (existing or newly created).
func getResponse(op *model.Operation, statusCode int) *model.Response {
	statusStr := strconv.Itoa(statusCode)
	if op.Responses[statusStr] == nil {
		op.Responses[statusStr] = &model.Response{
			Description: http.StatusText(statusCode),
		}
	}

	if op.Responses[statusStr].Content == nil {
		op.Responses[statusStr].Content = make(map[string]*model.MediaType)
	}

	return op.Responses[statusStr]
}

// transformSchemaForFileResponse transforms a schema for file/binary responses.
// For file responses, []byte should use format: binary (not contentEncoding: base64).
func transformSchemaForFileResponse(s *model.Schema) *model.Schema {
	if s == nil {
		return s
	}

	// For []byte fields, change from JSON Schema to OpenAPI binary format
	// In JSON: []byte -> {type: string, contentEncoding: base64, contentMediaType: application/octet-stream}
	// In OpenAPI file response: []byte -> {type: string, format: binary}
	if s.Type == TypeString && s.ContentEncoding == contentEncodingBase64 {
		sCopy := *s
		sCopy.ContentEncoding = ""
		sCopy.ContentMediaType = ""
		sCopy.Format = formatBinary

		return &sCopy
	}

	return s
}
