package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_InlineExample(t *testing.T) {
	value := map[string]any{"id": 42, "status": "active"}
	ex := New("user-found", value)

	assert.Equal(t, "user-found", ex.Name())
	assert.Equal(t, value, ex.Value())
	assert.Empty(t, ex.ExternalValue())
	assert.Empty(t, ex.Summary())
	assert.Empty(t, ex.Description())
	assert.False(t, ex.IsExternal())
}

func TestNew_WithOptions(t *testing.T) {
	value := map[string]any{"error": "not found"}
	ex := New("error-case", value,
		WithSummary("Resource not found"),
		WithDescription("Returned when the requested resource does not exist"),
	)

	assert.Equal(t, "error-case", ex.Name())
	assert.Equal(t, value, ex.Value())
	assert.Empty(t, ex.ExternalValue())
	assert.Equal(t, "Resource not found", ex.Summary())
	assert.Equal(t, "Returned when the requested resource does not exist", ex.Description())
	assert.False(t, ex.IsExternal())
}

func TestNewExternal_BasicURL(t *testing.T) {
	ex := NewExternal("full-dataset", "https://example.com/data/full.json")

	assert.Equal(t, "full-dataset", ex.Name())
	assert.Nil(t, ex.Value())
	assert.Equal(t, "https://example.com/data/full.json", ex.ExternalValue())
	assert.Empty(t, ex.Summary())
	assert.Empty(t, ex.Description())
	assert.True(t, ex.IsExternal())
}

func TestNewExternal_WithOptions(t *testing.T) {
	ex := NewExternal("large-payload", "https://example.com/samples/large.json",
		WithSummary("Large payload example"),
		WithDescription("Example with large data that would bloat the spec"),
	)

	assert.Equal(t, "large-payload", ex.Name())
	assert.Nil(t, ex.Value())
	assert.Equal(t, "https://example.com/samples/large.json", ex.ExternalValue())
	assert.Equal(t, "Large payload example", ex.Summary())
	assert.Equal(t, "Example with large data that would bloat the spec", ex.Description())
	assert.True(t, ex.IsExternal())
}

func TestIsExternal_WithInlineValue(t *testing.T) {
	ex := New("inline", map[string]any{"test": "value"})
	assert.False(t, ex.IsExternal())
}

func TestIsExternal_WithExternalURL(t *testing.T) {
	ex := NewExternal("external", "https://example.com/example.json")
	assert.True(t, ex.IsExternal())
}

func TestIsExternal_WithNilValue(t *testing.T) {
	ex := New("nil-value", nil)
	assert.False(t, ex.IsExternal())
}

func TestAccessors(t *testing.T) {
	tests := []struct {
		name     string
		example  Example
		wantName string
		wantExt  bool
	}{
		{
			name:     "inline example",
			example:  New("test", "value"),
			wantName: "test",
			wantExt:  false,
		},
		{
			name:     "external example",
			example:  NewExternal("ext", "https://example.com"),
			wantName: "ext",
			wantExt:  true,
		},
		{
			name: "with all options",
			example: New("complete", map[string]any{"data": 123},
				WithSummary("Test summary"),
				WithDescription("Test description"),
			),
			wantName: "complete",
			wantExt:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantName, tt.example.Name())
			assert.Equal(t, tt.wantExt, tt.example.IsExternal())

			// Ensure Name() returns the same value on multiple calls
			assert.Equal(t, tt.example.Name(), tt.example.Name())
		})
	}
}

func TestWithSummary(t *testing.T) {
	ex := New("test", "value", WithSummary("Test Summary"))
	assert.Equal(t, "Test Summary", ex.Summary())
}

func TestWithDescription(t *testing.T) {
	ex := New("test", "value", WithDescription("Test Description"))
	assert.Equal(t, "Test Description", ex.Description())
}

func TestMultipleOptions(t *testing.T) {
	ex := New("test", "value",
		WithSummary("Summary"),
		WithDescription("Description"),
	)
	assert.Equal(t, "Summary", ex.Summary())
	assert.Equal(t, "Description", ex.Description())
}

func TestExampleValue_DifferentTypes(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"string", "test"},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"map", map[string]any{"key": "value"}},
		{"slice", []any{1, 2, 3}},
		{"struct", struct{ Name string }{Name: "test"}},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := New("test", tt.value)
			assert.Equal(t, tt.value, ex.Value())
			assert.False(t, ex.IsExternal())
		})
	}
}
