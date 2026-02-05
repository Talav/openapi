package export

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/talav/openapi/debug"
	"github.com/talav/openapi/internal/model"
)

type Exporter interface {
	Export(ctx context.Context, spec *model.Spec, cfg ExporterConfig) (*ExporterResult, error)
	IsSupportedVersion(version string) bool
}

type ExporterConfig struct {
	Version        string
	ShouldValidate bool
}

// Result contains the output of spec projection.
type ExporterResult struct {
	// Result is the marshaled OpenAPI specification as bytes.
	Result []byte

	// Warnings contains any warnings generated during projection.
	// Warnings are generated when features are not supported by the target version.
	Warnings debug.Warnings
}

type ViewAdapter interface {
	View(spec *model.Spec) (any, debug.Warnings, error)
	Version() string
	SchemaJSON() []byte
}

type exporter struct {
	adapters map[string]ViewAdapter
}

func NewExporter(adapters []ViewAdapter) Exporter {
	adaptersMap := make(map[string]ViewAdapter)
	for _, adapter := range adapters {
		adaptersMap[adapter.Version()] = adapter
	}

	return &exporter{adapters: adaptersMap}
}

func (e *exporter) IsSupportedVersion(version string) bool {
	_, ok := e.adapters[version]

	return ok
}

func (e *exporter) Export(ctx context.Context, spec *model.Spec, cfg ExporterConfig) (*ExporterResult, error) {
	if spec == nil {
		return nil, errors.New("nil spec")
	}

	adapter, ok := e.adapters[cfg.Version]
	if !ok {
		return nil, fmt.Errorf("unknown version: %s", cfg.Version)
	}
	out, warns, err := adapter.View(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create a view of the spec: %w", err)
	}

	result, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec to JSON: %w", err)
	}

	if cfg.ShouldValidate {
		schemaJSON := adapter.SchemaJSON()

		validator, err := NewValidator(schemaJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to create validator: %w", err)
		}
		if err := validator.Validate(ctx, result); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	return &ExporterResult{
		Result:   result,
		Warnings: warns,
	}, nil
}
