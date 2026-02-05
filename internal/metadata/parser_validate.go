package metadata

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/talav/tagparser"
)

// ValidateMetadata represents validation constraints extracted from the validate tag.
// Types match OpenAPI v3.0 specification for schema validation constraints.
// This metadata is used to generate OpenAPI schema constraints by mapping
// go-playground/validator tags to OpenAPI/JSON Schema keywords.
type ValidateMetadata struct {
	// Numeric validation constraints (for number/integer types)
	// OpenAPI v3.0: minimum, maximum, exclusiveMinimum, exclusiveMaximum, multipleOf are numbers
	Minimum          *float64 // inclusive minimum value
	ExclusiveMinimum *float64 // exclusive minimum value (value must be > exclusiveMinimum)
	Maximum          *float64 // inclusive maximum value
	ExclusiveMaximum *float64 // exclusive maximum value (value must be < exclusiveMaximum)
	MultipleOf       *float64 // value must be a multiple of this number

	// String validation constraints (for string types)
	Pattern string // regular expression pattern that string must match
	Format  string // predefined format for string validation (e.g., "email", "date-time", "uri")

	// General validation constraints
	Enum     []any // parsed enum values
	Required *bool // field must be present
}

// ParseValidateTag parses a validate tag in go-playground/validator format and returns ValidateMetadata.
// Tag format: validate:"required,email,min=5,max=100,pattern=^[a-z]+$"
//
// This parser:
// 1. Parses go-playground/validator tag format (comma-separated, key=value pairs)
// 2. Maps validator tags to OpenAPI/JSON Schema constraints
// 3. Converts string values to proper OpenAPI types (int, float64, bool)
// 4. Returns error if value cannot be parsed to expected type
//
// Validator tag -> OpenAPI mapping:
//   - required -> Required=true
//   - min=N -> Minimum=N (as float64)
//   - max=N -> Maximum=N (as float64)
//   - len=N -> Minimum=N, Maximum=N (as float64, sets both to same value)
//   - email -> Format="email"
//   - url -> Format="uri"
//   - pattern=... -> Pattern="..."
//   - oneof=... -> Enum="[...]"
//   - etc.
func ParseValidateTag(field reflect.StructField, index int, tagValue string) (any, error) {
	vm := &ValidateMetadata{}

	// Parse go-playground/validator format using tagparser
	// Format: "required,email,min=5,max=100"
	// Use ParseFunc to handle all items, including flags without values
	allValidators := make(map[string]string)

	tag, err := tagparser.Parse(tagValue)
	if err != nil {
		return nil, fmt.Errorf("field %s: failed to parse validate tag: %w", field.Name, err)
	}

	for key, value := range tag.Options {
		if key == "" {
			// First item without equals sign (flag without value)
			allValidators[value] = ""
		} else {
			// Key=value pair
			allValidators[key] = value
		}
	}

	// Map validator tags to OpenAPI constraints
	for validator, value := range allValidators {
		if err := applyValidatorMapping(vm, validator, value); err != nil {
			return nil, fmt.Errorf("field %s: failed to apply validator %q: %w", field.Name, validator, err)
		}
	}

	return vm, nil
}

// applyValidatorMapping maps a single validator tag to OpenAPI constraint.
// Only includes validators actually supported by go-playground/validator v10.
// Reference: https://pkg.go.dev/github.com/go-playground/validator/v10
//
//nolint:cyclop // Map-based dispatch - acceptable complexity
func applyValidatorMapping(vm *ValidateMetadata, validator, value string) error {
	// Boolean flags
	boolSetters := map[string]**bool{
		"required": &vm.Required,
	}
	if ptr, ok := boolSetters[validator]; ok {
		b, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("invalid required value: %w", err)
		}
		*ptr = b

		return nil
	}

	// Numeric constraints (parse as float64 for OpenAPI)
	floatSetters := map[string]**float64{
		"min":         &vm.Minimum,
		"gte":         &vm.Minimum,
		"max":         &vm.Maximum,
		"lte":         &vm.Maximum,
		"gt":          &vm.ExclusiveMinimum,
		"lt":          &vm.ExclusiveMaximum,
		"multiple_of": &vm.MultipleOf,
	}
	if ptr, ok := floatSetters[validator]; ok {
		f, err := parseFloat64(value)
		if err != nil {
			return fmt.Errorf("invalid %s value %q: %w", validator, value, err)
		}
		*ptr = &f

		return nil
	}

	if validator == "len" {
		f, err := parseFloat64(value)
		if err != nil {
			return fmt.Errorf("invalid len value %q: %w", value, err)
		}
		vm.Minimum = &f
		vm.Maximum = &f

		return nil
	}

	// String format constraints (validator name -> OpenAPI format string)
	formatSetters := map[string]string{
		"email": "email",
		"url":   "uri",
	}
	if format, ok := formatSetters[validator]; ok {
		vm.Format = format

		return nil
	}

	// Fixed pattern constraints (validator name -> regex pattern)
	patternSetters := map[string]string{
		"alpha":           "^[a-zA-Z]+$",
		"alphanum":        "^[a-zA-Z0-9]+$",
		"alphaunicode":    "^[\\p{L}]+$",
		"alphanumunicode": "^[\\p{L}\\p{N}]+$",
	}
	if pattern, ok := patternSetters[validator]; ok {
		vm.Pattern = pattern

		return nil
	}

	if validator == "pattern" {
		vm.Pattern = value

		return nil
	}

	if validator == "oneof" {
		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("oneof requires at least one value")
		}
		var enumValues []any
		for _, part := range strings.Fields(value) {
			part = strings.TrimSpace(part)
			if part != "" {
				enumValues = append(enumValues, part)
			}
		}
		vm.Enum = enumValues

		return nil
	}

	return fmt.Errorf("unsupported validator %q (see go-playground/validator v10 docs)", validator)
}
