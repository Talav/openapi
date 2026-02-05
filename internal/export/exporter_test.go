package export

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talav/openapi/debug"
	v304 "github.com/talav/openapi/internal/export/v304"
	v312 "github.com/talav/openapi/internal/export/v312"
	"github.com/talav/openapi/internal/model"
)

// mockAdapter is a mock ViewAdapter for testing error cases.
type mockAdapter struct {
	version    string
	schemaJSON []byte
	viewFunc   func(*model.Spec) (any, debug.Warnings, error)
}

func (m *mockAdapter) Version() string {
	return m.version
}

func (m *mockAdapter) SchemaJSON() []byte {
	return m.schemaJSON
}

func (m *mockAdapter) View(spec *model.Spec) (any, debug.Warnings, error) {
	if m.viewFunc != nil {
		return m.viewFunc(spec)
	}

	return nil, nil, nil
}

func TestExport_NilSpec(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	ctx := context.Background()
	result, err := exporter.Export(ctx, nil, ExporterConfig{Version: "3.0.4"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "nil spec")
}

func TestExport_UnknownVersion(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	spec := createMinimalSpec()
	ctx := context.Background()
	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "2.0.0"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unknown version: 2.0.0")
}

func TestExport_AdapterViewError(t *testing.T) {
	expectedError := errors.New("view error")
	mock := &mockAdapter{
		version:    "3.0.4",
		schemaJSON: []byte(`{"$schema":"http://json-schema.org/draft-04/schema#"}`),
		viewFunc: func(*model.Spec) (any, debug.Warnings, error) {
			return nil, nil, expectedError
		},
	}

	exporter := NewExporter([]ViewAdapter{mock})
	spec := createMinimalSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create a view of the spec")
	assert.ErrorIs(t, err, expectedError)
}

func TestExport_JSONMarshalError(t *testing.T) {
	// Create a view with un-marshalable data (channel)
	unmarshalableView := struct {
		Channel chan int `json:"channel"`
	}{
		Channel: make(chan int),
	}

	mock := &mockAdapter{
		version:    "3.0.4",
		schemaJSON: []byte(`{"$schema":"http://json-schema.org/draft-04/schema#"}`),
		viewFunc: func(*model.Spec) (any, debug.Warnings, error) {
			return unmarshalableView, nil, nil
		},
	}

	exporter := NewExporter([]ViewAdapter{mock})
	spec := createMinimalSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to marshal spec to JSON")
}

func TestExport_ValidatorCreationError(t *testing.T) {
	invalidSchemaJSON := []byte(`{invalid json}`)

	mock := &mockAdapter{
		version:    "3.0.4",
		schemaJSON: invalidSchemaJSON,
		viewFunc: func(*model.Spec) (any, debug.Warnings, error) {
			return map[string]string{"openapi": "3.0.4"}, nil, nil
		},
	}

	exporter := NewExporter([]ViewAdapter{mock})
	spec := createMinimalSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4", ShouldValidate: true})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create validator")
}

func TestExport_ValidationFailure(t *testing.T) {
	// Use a valid schema but return a view that doesn't match it
	// We'll use the 3.0.4 schema but return invalid OpenAPI JSON
	adapter := &v304.AdapterV304{}
	mock := &mockAdapter{
		version:    "3.0.4",
		schemaJSON: adapter.SchemaJSON(),
		viewFunc: func(*model.Spec) (any, debug.Warnings, error) {
			// Return invalid OpenAPI structure (missing required fields)
			return map[string]any{
				"openapi": "3.0.4",
				// Missing required "info" field
			}, nil, nil
		},
	}

	exporter := NewExporter([]ViewAdapter{mock})
	spec := createMinimalSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4", ShouldValidate: true})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestExport_Success_V304(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	spec := createComprehensiveSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Result)
	// Warnings can be nil or empty - both are valid

	// Verify JSON contains correct version
	var jsonData map[string]any
	err = json.Unmarshal(result.Result, &jsonData)
	require.NoError(t, err)
	assert.Equal(t, "3.0.4", jsonData["openapi"])
}

func TestExport_Success_V312(t *testing.T) {
	adapter := &v312.AdapterV312{}
	exporter := NewExporter([]ViewAdapter{adapter})

	spec := createComprehensiveSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.1.2"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Result)
	// Warnings can be nil or empty - both are valid

	// Verify JSON contains correct version
	var jsonData map[string]any
	err = json.Unmarshal(result.Result, &jsonData)
	require.NoError(t, err)
	assert.Equal(t, "3.1.2", jsonData["openapi"])
}

func TestExport_Success_MinimalSpec(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	spec := createMinimalSpec()
	ctx := context.Background()

	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Result)

	// Verify JSON is valid and minimal
	var jsonData map[string]any
	err = json.Unmarshal(result.Result, &jsonData)
	require.NoError(t, err)
	assert.Equal(t, "3.0.4", jsonData["openapi"])

	info, ok := jsonData["info"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Test API", info["title"])
	assert.Equal(t, "1.0.0", info["version"])
}

func TestExport_Success_WithWarnings(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	// Create spec with 3.1-only features (webhooks) to trigger warnings
	spec := &model.Spec{
		Info: model.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Webhooks: map[string]*model.PathItem{
			"userCreated": {
				Post: &model.Operation{
					Summary: "User created webhook",
				},
			},
		},
	}

	ctx := context.Background()
	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Result)

	// Verify warnings are present
	assert.NotEmpty(t, result.Warnings)
	assert.True(t, result.Warnings.Has(debug.WarnDegradationWebhooks))
}

func TestExport_Success_WithExtensions(t *testing.T) {
	adapter := &v304.AdapterV304{}
	exporter := NewExporter([]ViewAdapter{adapter})

	spec := &model.Spec{
		Info: model.Info{
			Title:   "Test API",
			Version: "1.0.0",
			Extensions: map[string]any{
				"x-custom-info": "custom value",
				"x-api-version": "1.0.0-beta",
			},
		},
		Extensions: map[string]any{
			"x-top-level": "top level extension",
		},
	}

	ctx := context.Background()
	result, err := exporter.Export(ctx, spec, ExporterConfig{Version: "3.0.4"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Result)

	// Verify extensions appear in generated JSON
	var jsonData map[string]any
	err = json.Unmarshal(result.Result, &jsonData)
	require.NoError(t, err)

	info, ok := jsonData["info"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "custom value", info["x-custom-info"])
	assert.Equal(t, "1.0.0-beta", info["x-api-version"])

	assert.Equal(t, "top level extension", jsonData["x-top-level"])
}

// Helper function to create a minimal spec.
func createMinimalSpec() *model.Spec {
	return &model.Spec{
		Info: model.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}
}

// Helper function to create a comprehensive spec (reused from adapter tests).
func createComprehensiveSpec() *model.Spec {
	return &model.Spec{
		Info: model.Info{
			Title:       "Test API",
			Description: "A test API",
			Version:     "1.0.0",
			License: &model.License{
				Name: "MIT",
			},
			Extensions: map[string]any{
				"x-custom-info": "custom info extension",
			},
		},
		Tags: []model.Tag{
			{
				Name:        "Users",
				Description: "User management operations",
			},
		},
		Paths: map[string]*model.PathItem{
			"/users": {
				Get: &model.Operation{
					Summary: "Get users",
					Parameters: []model.Parameter{
						{
							Name:        "limit",
							In:          "query",
							Schema:      &model.Schema{Type: "integer", Default: 10},
							Description: "Maximum number of users to return",
						},
					},
					Responses: map[string]*model.Response{
						"200": {
							Description: "Success",
							Content: map[string]*model.MediaType{
								"application/json": {
									Schema: &model.Schema{
										Type:  "array",
										Items: &model.Schema{Ref: "#/components/schemas/User"},
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &model.Components{
			Schemas: map[string]*model.Schema{
				"User": {
					Type:  "object",
					Title: "User Schema",
					Properties: map[string]*model.Schema{
						"id": {
							Type:        "string",
							Description: "Unique user identifier",
						},
						"name": {
							Type:        "string",
							Description: "User name",
						},
					},
					Required: []string{"id", "name"},
				},
			},
		},
	}
}
