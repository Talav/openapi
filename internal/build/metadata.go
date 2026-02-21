package build

import (
	"reflect"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/metadata"
	"github.com/talav/schema"
)

// NewMetadata creates a new schema metadata instance with the given tag configuration.
// Partial configs are merged with defaults using config.MergeTagConfig().
func NewMetadata(cfg config.TagConfig) *schema.Metadata {
	// Merge with defaults to handle partial configs
	cfg = config.MergeTagConfig(config.DefaultTagConfig(), cfg)

	return schema.NewMetadata(schema.NewTagParserRegistry(
		schema.WithTagParser(cfg.Schema, schema.ParseSchemaTag, func(field reflect.StructField, index int) any {
			return conditionalSchemaDefault(field, index, cfg)
		}),
		schema.WithTagParser(cfg.Body, schema.ParseBodyTag),
		schema.WithTagParser(cfg.OpenAPI, metadata.ParseOpenAPITag),
		schema.WithTagParser(cfg.Validate, metadata.ParseValidateTag),
		schema.WithTagParser(cfg.Default, metadata.ParseDefaultTag),
		schema.WithTagParser(cfg.Requires, metadata.ParseRequiresTag),
	))
}

// conditionalSchemaDefault applies schema default metadata only if the field doesn't have a body tag.
// Business rule: fields with body tags should not receive default schema metadata.
func conditionalSchemaDefault(field reflect.StructField, index int, cfg config.TagConfig) any {
	// Don't apply schema default if field has body tag
	if _, ok := field.Tag.Lookup(cfg.Body); ok {
		return nil
	}

	return schema.DefaultSchemaMetadata(field, index)
}
