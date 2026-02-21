package metadata

import (
	"fmt"
	"reflect"

	"github.com/talav/tagparser"
)

// RequiresMetadata represents fields that become required when this field is present.
// Extracted from the requires tag for OpenAPI schema generation (JSON Schema dependentRequired keyword).
type RequiresMetadata struct {
	Fields []string // List of field names that become required when this field is present
}

// ParseRequiresTag parses a requires tag and returns RequiresMetadata.
// Tag format: requires:"field1,field2,field3"
//
// Parses comma-separated list of field names that become required when this field is present.
// Empty strings and whitespace are filtered out.
//
// Example:
//   - requires:"billing_address,cvv" -> Fields=["billing_address", "cvv"]
//   - requires:"field1" -> Fields=["field1"]
//   - requires:"" -> Fields=[] (empty, will be ignored)
func ParseRequiresTag(field reflect.StructField, index int, tagValue string) (any, error) {
	tag, err := tagparser.Parse(tagValue)
	if err != nil {
		return nil, fmt.Errorf("field %s: failed to parse requires tag: %w", field.Name, err)
	}

	fields := make([]string, 0, len(tag.Options))
	for key := range tag.Options {
		fields = append(fields, key)
	}

	return &RequiresMetadata{
		Fields: fields,
	}, nil
}
