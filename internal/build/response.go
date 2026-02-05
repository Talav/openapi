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

	if resp.Content == nil {
		resp.Content = make(map[string]*model.MediaType)
	}

	// Find and process body field
	bodyField := findBodyField(structMeta, rb.tagCfg)
	if bodyField == nil {
		return nil
	}

	// Extract body schema and add to response
	if err := rb.extractBodySchema(bodyField, resp, structMeta.Type, op); err != nil {
		return err
	}

	// Extract header schemas and add to success response
	rb.buildResponseHeaders(structMeta, resp)

	return nil
}

// extractBodySchema extracts the body schema from a body field and adds it to the response.
func (rb *responseBuilder) extractBodySchema(
	bodyField *schema.FieldMetadata,
	resp *model.Response,
	structType reflect.Type,
	op *model.Operation,
) error {
	// Get body metadata to determine content type based on body tag
	bodyMeta, ok := schema.GetTagMetadata[*schema.BodyMetadata](bodyField, rb.tagCfg.Body)
	if !ok {
		return fmt.Errorf("body field missing body metadata")
	}

	// Determine content type
	ct, err := rb.determineContentType(bodyField, bodyMeta)
	if err != nil {
		return err
	}

	// Initialize media type if needed (only if Content is empty)
	if len(resp.Content) == 0 {
		resp.Content[ct] = &model.MediaType{}
	}

	// Generate and transform body schema based on body type
	hint := getSchemaHint(structType, bodyField.StructFieldName, op.OperationID)
	bodySchema := rb.generateBodySchema(bodyField, bodyMeta, hint)
	if bodySchema != nil && resp.Content[ct] != nil && resp.Content[ct].Schema == nil {
		resp.Content[ct].Schema = bodySchema
	}

	return nil
}

// generateBodySchema generates and transforms the response body schema based on body type.
func (rb *responseBuilder) generateBodySchema(bodyField *schema.FieldMetadata, bodyMeta *schema.BodyMetadata, hint string) *model.Schema {
	bodySchema := rb.generator.schema(bodyField.Type, true, hint)

	if bodyMeta.BodyType == schema.BodyTypeFile {
		return transformSchemaForFileResponse(bodySchema)
	}

	return bodySchema
}

// determineContentType determines the content type for a body field.
// Returns an error if the body type is invalid for responses.
func (rb *responseBuilder) determineContentType(bodyField *schema.FieldMetadata, bodyMeta *schema.BodyMetadata) (string, error) {
	// Determine content type based on BodyType (validates that multipart is not used)
	ct, err := getResponseContentType(bodyMeta.BodyType)
	if err != nil {
		return "", fmt.Errorf("field %q: %w", bodyField.StructFieldName, err)
	}

	// Fallback to ContentTypeProvider interface if needed
	if ct == contentTypeJSON && reflect.PointerTo(bodyField.Type).Implements(reflect.TypeOf((*ContentTypeProvider)(nil)).Elem()) {
		instance, ok := reflect.New(bodyField.Type).Interface().(ContentTypeProvider)
		if ok {
			ct = instance.ContentType(ct)
		}
	}

	return ct, nil
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

// getResponseContentType maps BodyType to HTTP content-type for responses.
// Returns an error if the body type is invalid for responses.
// Valid types: BodyTypeStructured (JSON), BodyTypeFile (octet-stream).
func getResponseContentType(bodyType schema.BodyType) (string, error) {
	switch bodyType {
	case schema.BodyTypeStructured:
		return contentTypeJSON, nil
	case schema.BodyTypeFile:
		return contentTypeOctetStream, nil
	case schema.BodyTypeMultipart:
		return "", fmt.Errorf("invalid body type for response: multipart is not supported, use %q or %q", schema.BodyTypeStructured, schema.BodyTypeFile)
	default:
		return "", fmt.Errorf("invalid body type for response: empty body tag, must explicitly specify body type (e.g., body:\"structured\" or body:\"file\")")
	}
}
