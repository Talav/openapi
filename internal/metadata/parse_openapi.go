package metadata

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/talav/tagparser"
)

// OpenAPIMetadata represents OpenAPI-specific schema metadata extracted from the openapi tag.
// Types match OpenAPI v3.0 specification for schema metadata.
// This metadata is used to generate OpenAPI schema properties that are not validation constraints
// but API contract metadata (e.g., readOnly, writeOnly, deprecated, title, description, examples).
//
// When used on a field (not the _ blank identifier), it represents field-level metadata.
// When used on the _ blank identifier field, it represents struct-level metadata
// (additionalProperties, nullable).
type OpenAPIMetadata struct {
	// Field-level API contract metadata (not validation constraints)
	// OpenAPI v3.0: readOnly, writeOnly, deprecated are booleans
	ReadOnly    *bool  // field is read-only
	WriteOnly   *bool  // field is write-only
	Deprecated  *bool  // field is deprecated
	Hidden      *bool  // field is hidden from schema (not included in properties)
	Required    *bool  // field is required (override for validate:"required")
	Title       string // title for the schema
	Description string // description for the schema
	Format      string // format for the schema (e.g., "date", "date-time", "time", "email", "uri")
	Examples    []any  // parsed example values

	// Struct-level metadata (only valid when used on _ blank identifier field)
	AdditionalProperties *bool // allow additional properties (struct-level)
	Nullable             *bool // struct is nullable (struct-level)

	// Extensions are OpenAPI specification extensions (x-* fields).
	// Keys must start with "x-" per OpenAPI spec requirement.
	Extensions map[string]any
}

// ParseOpenAPITag parses an openapi tag and returns OpenAPIMetadata.
// Tag format: openapi:"readOnly,writeOnly,deprecated,hidden,required,title=My Title,description=My description,examples=val1|val2|val3,x-custom=value"
//
// This parser:
// 1. Parses tag format (comma-separated, key=value pairs or flags)
// 2. Converts string values to proper OpenAPI types (bool for readOnly/writeOnly/deprecated/hidden/required)
// 3. Converts empty string to true for boolean flags (e.g., "readOnly" -> ReadOnly=true)
// 4. Routes x-* prefixed keys to Extensions map (OpenAPI spec requirement)
// 5. Detects struct-level vs field-level based on field name (blank identifier _ = struct-level)
// 6. Supports pipe-separated examples values: examples=val1|val2|val3
//
// Field-level options (for named fields):
//   - readOnly -> ReadOnly=true
//   - writeOnly -> WriteOnly=true
//   - deprecated -> Deprecated=true
//   - hidden -> Hidden=true (field excluded from schema properties)
//   - required -> Required=true (overrides validate:"required" for docs only)
//   - title=... -> Title="..."
//   - description=... -> Description="..."
//   - format=... -> Format="..." (e.g., "date", "date-time", "time", "email", "uri")
//   - examples=val1|val2|val3 -> Examples=[val1, val2, val3] (pipe-separated values)
//
// Struct-level options (for _ blank identifier field):
//   - additionalProperties=true/false -> AdditionalProperties=bool
//   - nullable=true/false -> Nullable=bool
//
// OpenAPI extensions (valid at both field and struct level):
//   - x-* -> Extensions["x-*"]="..." (MUST start with x-, minimum length 4)
func ParseOpenAPITag(field reflect.StructField, index int, tagValue string) (any, error) {
	om := &OpenAPIMetadata{}

	// Parse tag using tagparser (options mode - all items are options)
	tag, err := tagparser.Parse(tagValue)
	if err != nil {
		return nil, fmt.Errorf("field %s: failed to parse openapi tag: %w", field.Name, err)
	}

	// Detect if this is struct-level metadata (blank identifier field)
	isStructLevel := field.Name == "_"

	// Process all options
	for key, value := range tag.Options {
		if err := applyOpenAPIMapping(om, key, value, isStructLevel); err != nil {
			return nil, fmt.Errorf("field %s: failed to apply openapi mapping: %w", field.Name, err)
		}
	}

	return om, nil
}

// applyOpenAPIMapping maps a single openapi tag option to OpenAPIMetadata field.
// Extensions (x- prefix, length > 3) are processed first for both struct and field levels.
// isStructLevel indicates if this is struct-level metadata (on _ blank identifier field).
// Non-extension keys are routed to struct-level or field-level handlers based on isStructLevel.
// Supports pipe-separated examples values (e.g., examples=val1|val2|val3).
func applyOpenAPIMapping(om *OpenAPIMetadata, key, value string, isStructLevel bool) error {
	if isExtension(key) {
		applyExtension(om, key, value)

		return nil
	}

	if isStructLevel {
		return applyStructLevelOption(om, key, value)
	}

	return applyFieldLevelOption(om, key, value)
}

// isExtension checks if a key is a valid OpenAPI extension (x- prefix with length > 3).
func isExtension(key string) bool {
	return strings.HasPrefix(key, "x-") && len(key) > 3
}

// applyExtension adds an extension to the metadata.
func applyExtension(om *OpenAPIMetadata, key, value string) {
	if om.Extensions == nil {
		om.Extensions = make(map[string]any)
	}
	om.Extensions[key] = value
}

// applyStructLevelOption handles struct-level OpenAPI options.
func applyStructLevelOption(om *OpenAPIMetadata, key, value string) error {
	boolSetters := map[string]**bool{
		"additionalProperties": &om.AdditionalProperties,
		"nullable":             &om.Nullable,
	}

	if ptr, ok := boolSetters[key]; ok {
		b, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("invalid %s value: %w", key, err)
		}
		*ptr = b

		return nil
	}

	return fmt.Errorf("unknown struct-level option %q (valid: additionalProperties, nullable)", key)
}

// applyFieldLevelOption handles field-level OpenAPI options.
func applyFieldLevelOption(om *OpenAPIMetadata, key, value string) error {
	boolSetters := map[string]**bool{
		"readOnly":   &om.ReadOnly,
		"writeOnly":  &om.WriteOnly,
		"deprecated": &om.Deprecated,
		"hidden":     &om.Hidden,
		"required":   &om.Required,
	}

	if ptr, ok := boolSetters[key]; ok {
		b, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("invalid %s value: %w", key, err)
		}
		*ptr = b

		return nil
	}

	stringSetters := map[string]*string{
		"title":       &om.Title,
		"description": &om.Description,
		"format":      &om.Format,
	}

	if ptr, ok := stringSetters[key]; ok {
		*ptr = value

		return nil
	}

	if key == "examples" {
		om.Examples = append(om.Examples, parseExampleValues(value)...)

		return nil
	}

	return fmt.Errorf("unknown field-level option %q (valid: readOnly, writeOnly, deprecated, hidden, required, title, description, format, examples)", key)
}

// parseExampleValues parses pipe-separated example values.
// Numeric values are stored as float64, others as strings.
func parseExampleValues(value string) []any {
	var examples []any
	for part := range strings.SplitSeq(value, "|") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Try to detect if value is numeric
		if num, err := strconv.ParseFloat(part, 64); err == nil {
			examples = append(examples, num)
		} else {
			examples = append(examples, part)
		}
	}

	return examples
}
