package build

import (
	"reflect"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/metadata"
	"github.com/talav/schema"
)

const (
	contentTypeMultipart   = "multipart/form-data"
	contentTypeOctetStream = "application/octet-stream"
	contentTypeJSON        = "application/json"
	formatBinary           = "binary"
)

// getSchemaHint generates a hint for schema naming from type and field name.
// Used by the schema registry to name schemas for anonymous/unnamed types.
// Priority:
//  1. type.Name() + fieldName (e.g., "CreateUserRequestName")
//  2. operationID + fieldName (e.g., "createUserRequestName")
//  3. fieldName (fallback)
func getSchemaHint(typ reflect.Type, fieldName, operationID string) string {
	if typ.Name() != "" {
		return typ.Name() + fieldName
	}

	if operationID != "" {
		return operationID + fieldName
	}

	return fieldName
}

// findBodyField finds the field with body tag in the struct metadata.
// Returns nil if no body field is found.
func findBodyField(structMeta *schema.StructMetadata, cfg config.TagConfig) *schema.FieldMetadata {
	for i := range structMeta.Fields {
		if structMeta.Fields[i].HasTag(cfg.Body) {
			return &structMeta.Fields[i]
		}
	}

	return nil
}

// deref removes all pointer indirections from a type.
func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

// toBool converts a bool or *bool to bool.
// If the input is a pointer, returns false if nil, otherwise the dereferenced value.
// If the input is a bool, returns it directly.
func toBool(b any) bool {
	switch v := b.(type) {
	case *bool:
		if v == nil {
			return false
		}

		return *v
	case bool:
		return v
	default:
		return false
	}
}

// isRequiredFromMetadata returns true if the field is marked required via openapi or validate tags.
func isRequiredFromMetadata(field *schema.FieldMetadata, tagCfg config.TagConfig) bool {
	if openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](field, tagCfg.OpenAPI); ok {
		if toBool(openAPIMeta.Required) {
			return true
		}
	}
	if validateMeta, ok := schema.GetTagMetadata[*metadata.ValidateMetadata](field, tagCfg.Validate); ok {
		if toBool(validateMeta.Required) {
			return true
		}
	}

	return false
}
