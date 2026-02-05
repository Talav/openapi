package config

// TagConfig configures struct tag names used for OpenAPI schema generation.
type TagConfig struct {
	// Schema is the tag name for schema metadata (e.g., "schema").
	Schema string

	// Body is the tag name for body field identification (e.g., "body").
	Body string

	// OpenAPI is the tag name for OpenAPI-specific metadata (e.g., "openapi").
	OpenAPI string

	// Validate is the tag name for validation constraints (e.g., "validate").
	Validate string

	// Default is the tag name for default values (e.g., "default").
	Default string

	// Requires is the tag name for dependent required fields (e.g., "requires").
	Requires string
}

// DefaultTagConfig returns the default tag configuration with standard tag names.
func DefaultTagConfig() TagConfig {
	return TagConfig{
		Schema:   "schema",
		Body:     "body",
		OpenAPI:  "openapi",
		Validate: "validate",
		Default:  "default",
		Requires: "requires",
	}
}

// MergeInto merges cfg into current, preserving current values when cfg fields are empty.
// Non-empty values in cfg override corresponding fields in current.
// This is useful for chaining multiple partial configurations.
func MergeTagConfig(current, cfg TagConfig) TagConfig {
	result := current

	if cfg.Schema != "" {
		result.Schema = cfg.Schema
	}
	if cfg.Body != "" {
		result.Body = cfg.Body
	}
	if cfg.OpenAPI != "" {
		result.OpenAPI = cfg.OpenAPI
	}
	if cfg.Validate != "" {
		result.Validate = cfg.Validate
	}
	if cfg.Default != "" {
		result.Default = cfg.Default
	}
	if cfg.Requires != "" {
		result.Requires = cfg.Requires
	}

	return result
}

// NewTagConfig creates a TagConfig with explicit values for all fields.
// Use this when you want to specify all tag names explicitly.
func NewTagConfig(schema, body, openapi, validate, default_, requires string) TagConfig {
	return TagConfig{
		Schema:   schema,
		Body:     body,
		OpenAPI:  openapi,
		Validate: validate,
		Default:  default_,
		Requires: requires,
	}
}
