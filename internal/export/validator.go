package export

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validator validates OpenAPI specifications against a specific meta-schema.
// Each validator instance is tied to a specific OpenAPI version.
type Validator struct {
	schema *jsonschema.Schema
}

// New creates a new version-specific Validator with the provided meta-schema JSON.
//
// The validator uses santhosh-tekuri/jsonschema which supports both
// JSON Schema draft-04 (for OpenAPI 3.0) and draft-2020-12 (for OpenAPI 3.1).
//
// Example:
//
//	validator, err := validator.New(schemaJSON)
//	if err != nil {
//	    log.Fatalf("Failed to create validator: %v", err)
//	}
func NewValidator(schemaJSON []byte) (*Validator, error) {
	// Unmarshal the schema JSON into a document
	var schemaDoc any
	if err := json.Unmarshal(schemaJSON, &schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema JSON: %w", err)
	}

	compiler := jsonschema.NewCompiler()

	// Use a simple resource name
	resourceName := "openapi-schema.json"
	if err := compiler.AddResource(resourceName, schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}

	schema, err := compiler.Compile(resourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &Validator{
		schema: schema,
	}, nil
}

// Validate validates an OpenAPI specification JSON against the meta-schema.
func (v *Validator) Validate(ctx context.Context, specJSON []byte) error {
	// Unmarshal JSON first, then validate the unmarshaled data
	var data any
	if err := json.Unmarshal(specJSON, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return v.schema.Validate(data)
}
