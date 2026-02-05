package debug

import "fmt"

// Warning represents an informational, non-fatal issue during spec generation.
//
// Warnings are ADVISORY ONLY and never break execution.
// Use errors for issues that must stop the process.
//
// Common scenarios that produce warnings:
//   - Targeting OpenAPI 3.0 when using 3.1-only features (downlevel)
//   - Using deprecated API features
type Warning interface {
	// Code returns the warning identifier.
	// Compare with Warn* constants for type-safe checks.
	Code() WarningCode

	// Path returns the JSON pointer to the affected spec element.
	// Example: "#/webhooks", "#/info/summary"
	Path() string

	// Message returns a human-readable description.
	Message() string

	// String returns a formatted representation.
	String() string
}

// WarningCode identifies a specific warning type.
// Use the Warn* constants for type-safe comparisons.
type WarningCode string

// String returns the code as a string.
func (c WarningCode) String() string {
	return string(c)
}

// Schema degradation Warnings (3.1 → 3.0 schema feature losses).
const (
	// WarnDegradationWebhooks indicates webhooks were dropped (3.0 doesn't support them).
	WarnDegradationWebhooks WarningCode = "DEGRADATION_WEBHOOKS"

	// WarnDegradationInfoSummary indicates info.summary was dropped (3.0 doesn't support it).
	WarnDegradationInfoSummary WarningCode = "DEGRADATION_INFO_SUMMARY"

	// WarnDegradationLicenseIdentifier indicates license.identifier was dropped.
	WarnDegradationLicenseIdentifier WarningCode = "DEGRADATION_LICENSE_IDENTIFIER"

	// WarnDegradationMutualTLS indicates mutualTLS security scheme was dropped.
	WarnDegradationMutualTLS WarningCode = "DEGRADATION_MUTUAL_TLS"

	// WarnDegradationConstToEnum indicates JSON Schema const was converted to enum.
	WarnDegradationConstToEnum WarningCode = "DEGRADATION_CONST_TO_ENUM"

	// WarnDegradationConstToEnumConflict indicates const conflicted with existing enum.
	WarnDegradationConstToEnumConflict WarningCode = "DEGRADATION_CONST_TO_ENUM_CONFLICT"

	// WarnDegradationPathItems indicates $ref in pathItems was expanded.
	WarnDegradationPathItems WarningCode = "DEGRADATION_PATH_ITEMS"

	// WarnDegradationPatternProperties indicates patternProperties was dropped.
	WarnDegradationPatternProperties WarningCode = "DEGRADATION_PATTERN_PROPERTIES"

	// WarnDegradationUnevaluatedProperties indicates unevaluatedProperties was dropped.
	WarnDegradationUnevaluatedProperties WarningCode = "DEGRADATION_UNEVALUATED_PROPERTIES"

	// WarnDegradationContentEncoding indicates contentEncoding was dropped.
	WarnDegradationContentEncoding WarningCode = "DEGRADATION_CONTENT_ENCODING"

	// WarnDegradationContentMediaType indicates contentMediaType was dropped.
	WarnDegradationContentMediaType WarningCode = "DEGRADATION_CONTENT_MEDIA_TYPE"

	// WarnDegradationMultipleExamples indicates multiple examples were collapsed to one.
	WarnDegradationMultipleExamples WarningCode = "DEGRADATION_MULTIPLE_EXAMPLES"
)

// Spec violation warnings (invalid OpenAPI constructs).
const (
	// WarnInvalidExampleMutualExclusivity indicates both value and externalValue were set.
	WarnInvalidExampleMutualExclusivity WarningCode = "INVALID_EXAMPLE_MUTUAL_EXCLUSIVITY"
)

// Warnings is a collection of Warning with helper methods.
// Warnings are informational and never break execution.
type Warnings []Warning

// Has returns true if any warning matches the given code.
func (ws Warnings) Has(code WarningCode) bool {
	for _, w := range ws {
		if w.Code() == code {
			return true
		}
	}

	return false
}

// Append adds a warning to the collection.
func (ws *Warnings) Append(w Warning) {
	*ws = append(*ws, w)
}

// warning is the concrete implementation of Warning interface.
type warning struct {
	code    WarningCode
	path    string
	message string
}

func (w *warning) Code() WarningCode {
	return w.code
}

func (w *warning) Path() string {
	return w.path
}

func (w *warning) Message() string {
	return w.message
}

func (w *warning) String() string {
	return fmt.Sprintf("[%s] %s", w.code, w.message)
}

// NewWarning creates a new Warning instance.
// This is the primary way to create warnings from internal packages.
func NewWarning(code WarningCode, path, message string) Warning {
	return &warning{
		code:    code,
		path:    path,
		message: message,
	}
}
