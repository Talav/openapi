package openapi

import "github.com/talav/openapi/debug"

type Result struct {
	JSON []byte

	// Warnings contains informational, non-fatal issues.
	// These are advisory only and do not indicate failure.
	Warnings debug.Warnings
}
