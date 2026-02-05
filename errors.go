package openapi

import "errors"

// Configuration Errors (returned by [New]).
var (
	// ErrTitleRequired indicates the API title was not provided.
	ErrTitleRequired = errors.New("openapi: title is required")

	// ErrVersionRequired indicates the API version was not provided.
	ErrVersionRequired = errors.New("openapi: version is required")

	// ErrLicenseMutuallyExclusive indicates both license identifier and URL were set.
	ErrLicenseMutuallyExclusive = errors.New("openapi: license identifier and url are mutually exclusive")

	// ErrServerVariablesNeedURL indicates server variables were set without a server URL.
	ErrServerVariablesNeedURL = errors.New("openapi: server variables require a server URL")

	// ErrInvalidVersion indicates an unsupported OpenAPI version was specified.
	ErrInvalidVersion = errors.New("openapi: invalid OpenAPI version")
)
