// Package metadata extracts OpenAPI schema metadata from Go struct field tags.
//
// This package provides parsers for various struct tags that control OpenAPI schema generation.
// It bridges the gap between Go struct tags and OpenAPI/JSON Schema properties by parsing
// tag values and converting them to appropriate metadata structures.
//
// # Supported Tags
//
// The package supports three categories of tags:
//
// 1. OpenAPI Metadata (openapi tag):
//   - Schema documentation: title, description, format, examples
//   - Field modifiers: readOnly, writeOnly, deprecated, hidden, required
//   - Extensions: x-* prefixed custom fields (field or struct level)
//   - Struct-level only: additionalProperties, nullable (on _ field)
//
// 2. Validation Constraints (validate tag):
//   - Numeric: min, max, gt, gte, lt, lte, multipleOf
//   - String: pattern, email, url, uuid, etc.
//   - General: required, enum (oneof)
//
// 3. Specialized Tags:
//   - default: Default values for fields
//   - requires: Fields that become required when this field is present
//
// # Tag Precedence Rules
//
// When tags overlap, the following precedence applies:
//
// 1. Format: validate:"email" sets format, openapi:"format=custom" overrides it
// 2. Required: validate:"required" is source of truth, openapi:"required" can override for docs
// 3. Hidden: json:"-" hides from JSON, openapi:"hidden" hides from OpenAPI schema only
//
// Best practices:
//   - Use validate:"required" for actual validation
//   - Use openapi:"required" only to override docs (e.g., mark optional field as required in docs)
//   - Use validate:"email" for format, not openapi:"format=email" (unless overriding)
//
// # Complete Example
//
//	type User struct {
//	    _            struct{} `openapi:"additionalProperties=false,nullable=false"`
//	    ID           int      `json:"id" openapi:"readOnly,description=Unique identifier"`
//	    Email        string   `json:"email" validate:"required,email" openapi:"description=User email address,examples=user@example.com"`
//	    Name         string   `json:"name" validate:"required,min=3" openapi:"description=Full name,examples=John Doe|Jane Smith"`
//	    Age          int      `json:"age" validate:"min=0,max=150" openapi:"description=User age,examples=25|30|35"`
//	    Status       string   `json:"status" validate:"oneof=active inactive" openapi:"description=Account status,examples=active|inactive"`
//	    CreditCard   string   `json:"credit_card,omitempty" requires:"billing_address,cvv"`
//	    BillingAddr  string   `json:"billing_address,omitempty"`
//	    CVV          string   `json:"cvv,omitempty" openapi:"writeOnly"`
//	    Internal     bool     `json:"-" openapi:"hidden"`
//	    LegacyField  string   `json:"legacy_field,omitempty" openapi:"deprecated"`
//	}
//
// # OpenAPI Tag
//
// The openapi tag controls OpenAPI-specific schema properties:
//
//	// Boolean flags (no value needed)
//	openapi:"readOnly"              // Field is read-only (e.g., ID, created_at)
//	openapi:"writeOnly"             // Field is write-only (e.g., password, secret)
//	openapi:"deprecated"            // Field is deprecated
//	openapi:"hidden"                // Field excluded from OpenAPI schema (but in JSON)
//	openapi:"required"              // Override required status for docs only
//
//	// Documentation
//	openapi:"title=Field Title"
//	openapi:"description=Detailed description"
//	openapi:"format=date-time"      // OpenAPI format (date, date-time, email, uri, uuid, etc.)
//
//	// Examples (pipe-separated for multiple values)
//	openapi:"examples=value"        // Single example
//	openapi:"examples=val1|val2"    // Multiple examples
//
//	// Extensions (must start with x-, valid at both field and struct level)
//	openapi:"x-internal=true,x-category=admin"
//
//	// Struct-level options (on _ blank identifier field)
//	openapi:"additionalProperties=false"           // Disallow additional properties
//	openapi:"nullable=true"                        // Struct can be null
//	openapi:"additionalProperties=false,x-strict=true"  // Can combine with extensions
//
// # Validate Tag
//
// The validate tag uses go-playground/validator format and maps to OpenAPI constraints:
//
//	validate:"required"             -> Required=true, used as source of truth
//	validate:"min=5,max=100"        -> Minimum=5, Maximum=100
//	validate:"gt=0,lt=100"          -> ExclusiveMinimum=0, ExclusiveMaximum=100
//	validate:"email"                -> Format="email" (can be overridden by openapi:"format=...")
//	validate:"url"                  -> Format="uri"
//	validate:"oneof=red green blue" -> Enum=["red","green","blue"]
//
// Note: validate tag is the source of truth for constraints. Use openapi tag only to
// add documentation metadata or override specific values for documentation purposes.
//
// # Default Tag
//
// The default tag sets default values for fields:
//
//	default:"value"              // String default (no quotes needed)
//	default:"42"                 // Number default (parsed as JSON)
//	default:"true"               // Boolean default (parsed as JSON)
//	default:"[1,2,3]"           // Array default (parsed as JSON)
//	default:"{\"key\":\"value\"}" // Object default (parsed as JSON)
//
// # Requires Tag
//
// The requires tag specifies fields that become required when this field is present
// (JSON Schema dependentRequired keyword):
//
//	type Payment struct {
//	    CreditCard     string `json:"credit_card" requires:"billing_address,cvv"`
//	    BillingAddress string `json:"billing_address"`
//	    CVV            string `json:"cvv"`
//	}
//
// When credit_card is present, billing_address and cvv become required.
//
// # Hidden vs json:"-"
//
// Two ways to hide fields, with different purposes:
//
//	json:"-"                 // Field not serialized to JSON AND hidden from OpenAPI
//	openapi:"hidden"         // Field serialized to JSON BUT hidden from OpenAPI schema
//
// Use cases:
//   - json:"-" → Completely internal fields (never in API responses or docs)
//   - openapi:"hidden" → Runtime fields that appear in responses but not in documentation
//
// # Error Handling
//
// All parsers return descriptive errors that include:
//   - Field name where the error occurred
//   - Tag name that failed to parse
//   - Specific reason for the failure
//
// Example error: "field Email: failed to parse validate tag: invalid email format"
package metadata
