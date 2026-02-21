package metadata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx // Comprehensive table-driven test with many cases
func TestParseOpenAPI(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		tagValue    string
		want        *OpenAPIMetadata
		wantErr     bool
		errContains string
	}{
		{
			name:      "empty tag",
			fieldName: "Field",
			tagValue:  "",
			want:      &OpenAPIMetadata{},
		},
		{
			name:      "readOnly flag",
			fieldName: "ID",
			tagValue:  "readOnly",
			want: &OpenAPIMetadata{
				ReadOnly: boolPtr(true),
			},
		},
		{
			name:      "writeOnly flag",
			fieldName: "Password",
			tagValue:  "writeOnly",
			want: &OpenAPIMetadata{
				WriteOnly: boolPtr(true),
			},
		},
		{
			name:      "deprecated flag",
			fieldName: "OldField",
			tagValue:  "deprecated",
			want: &OpenAPIMetadata{
				Deprecated: boolPtr(true),
			},
		},
		{
			name:      "required flag",
			fieldName: "Email",
			tagValue:  "required",
			want: &OpenAPIMetadata{
				Required: boolPtr(true),
			},
		},
		{
			name:      "required with explicit true",
			fieldName: "Name",
			tagValue:  "required=true",
			want: &OpenAPIMetadata{
				Required: boolPtr(true),
			},
		},
		{
			name:      "required with explicit false",
			fieldName: "OptionalField",
			tagValue:  "required=false",
			want: &OpenAPIMetadata{
				Required: boolPtr(false),
			},
		},
		{
			name:      "readOnly with explicit true",
			fieldName: "ID",
			tagValue:  "readOnly=true",
			want: &OpenAPIMetadata{
				ReadOnly: boolPtr(true),
			},
		},
		{
			name:      "readOnly with explicit false",
			fieldName: "ID",
			tagValue:  "readOnly=false",
			want: &OpenAPIMetadata{
				ReadOnly: boolPtr(false),
			},
		},
		{
			name:      "hidden with explicit true",
			fieldName: "InternalField",
			tagValue:  "hidden=true",
			want: &OpenAPIMetadata{
				Hidden: boolPtr(true),
			},
		},
		{
			name:      "hidden with explicit false",
			fieldName: "PublicField",
			tagValue:  "hidden=false",
			want: &OpenAPIMetadata{
				Hidden: boolPtr(false),
			},
		},
		{
			name:      "title",
			fieldName: "Name",
			tagValue:  "title=User Name",
			want: &OpenAPIMetadata{
				Title: "User Name",
			},
		},
		{
			name:      "description",
			fieldName: "Email",
			tagValue:  "description=User email address",
			want: &OpenAPIMetadata{
				Description: "User email address",
			},
		},
		{
			name:      "single example value",
			fieldName: "Age",
			tagValue:  "examples=25",
			want: &OpenAPIMetadata{
				Examples: []any{25.0},
			},
		},
		{
			name:      "multiple examples values (pipe-separated)",
			fieldName: "Status",
			tagValue:  "examples=active|inactive|pending",
			want: &OpenAPIMetadata{
				Examples: []any{"active", "inactive", "pending"},
			},
		},
		{
			name:      "examples with numeric values (pipe-separated)",
			fieldName: "Score",
			tagValue:  "examples=0|50|100",
			want: &OpenAPIMetadata{
				Examples: []any{0.0, 50.0, 100.0},
			},
		},
		{
			name:      "single extension",
			fieldName: "Field",
			tagValue:  "x-custom=value",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom": "value",
				},
			},
		},
		{
			name:      "extension with hyphen",
			fieldName: "Field",
			tagValue:  "x-custom-feature=enabled",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom-feature": "enabled",
				},
			},
		},
		{
			name:      "multiple extensions",
			fieldName: "Field",
			tagValue:  "x-custom=value,x-vendor-tool=test,x-rate-limit=100",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom":      "value",
					"x-vendor-tool": "test",
					"x-rate-limit":  "100",
				},
			},
		},
		{
			name:      "extension with quoted value",
			fieldName: "Field",
			tagValue:  "x-custom='value with spaces'",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom": "value with spaces",
				},
			},
		},
		{
			name:      "mixed standard and extension fields",
			fieldName: "ID",
			tagValue:  "readOnly,deprecated,title=User ID,x-custom=value,x-vendor=test",
			want: &OpenAPIMetadata{
				ReadOnly:   boolPtr(true),
				Deprecated: boolPtr(true),
				Title:      "User ID",
				Extensions: map[string]any{
					"x-custom": "value",
					"x-vendor": "test",
				},
			},
		},
		{
			name:      "all standard fields",
			fieldName: "Field",
			tagValue:  "readOnly,writeOnly,deprecated,hidden,required,title=Title,description=Description,examples=val1|val2",
			want: &OpenAPIMetadata{
				ReadOnly:    boolPtr(true),
				WriteOnly:   boolPtr(true),
				Deprecated:  boolPtr(true),
				Hidden:      boolPtr(true),
				Required:    boolPtr(true),
				Title:       "Title",
				Description: "Description",
				Examples:    []any{"val1", "val2"},
			},
		},
		{
			name:      "extension with empty value",
			fieldName: "Field",
			tagValue:  "x-custom=",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom": "",
				},
			},
		},
		{
			name:        "invalid extension key - too short (error)",
			fieldName:   "Field",
			tagValue:    "x-=",
			wantErr:     true,
			errContains: "unknown field-level option",
		},
		{
			name:        "non-extension key returns error",
			fieldName:   "Field",
			tagValue:    "x=value",
			wantErr:     true,
			errContains: "unknown field-level option",
		},
		{
			name:        "invalid tag parsing",
			fieldName:   "Field",
			tagValue:    "readOnly,'unclosed quote",
			wantErr:     true,
			errContains: "failed to parse openapi tag",
		},
		{
			name:      "extension with special characters",
			fieldName: "Field",
			tagValue:  "x-custom-value=test-value-123",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom-value": "test-value-123",
				},
			},
		},
		{
			name:      "description with comma",
			fieldName: "Field",
			tagValue:  "description='Description, with comma'",
			want: &OpenAPIMetadata{
				Description: "Description, with comma",
			},
		},
		{
			name:      "multiple boolean flags",
			fieldName: "Field",
			tagValue:  "readOnly,writeOnly,deprecated,hidden,required",
			want: &OpenAPIMetadata{
				ReadOnly:   boolPtr(true),
				WriteOnly:  boolPtr(true),
				Deprecated: boolPtr(true),
				Hidden:     boolPtr(true),
				Required:   boolPtr(true),
			},
		},
		{
			name:      "extension overwrites previous value",
			fieldName: "Field",
			tagValue:  "x-custom=first,x-custom=second",
			want: &OpenAPIMetadata{
				Extensions: map[string]any{
					"x-custom": "second",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: tt.fieldName,
			}

			result, err := ParseOpenAPITag(field, 0, tt.tagValue)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			om, ok := result.(*OpenAPIMetadata)
			require.True(t, ok, "result should be *OpenAPIMetadata")

			assert.Equal(t, tt.want.ReadOnly, om.ReadOnly, "ReadOnly mismatch")
			assert.Equal(t, tt.want.WriteOnly, om.WriteOnly, "WriteOnly mismatch")
			assert.Equal(t, tt.want.Deprecated, om.Deprecated, "Deprecated mismatch")
			assert.Equal(t, tt.want.Hidden, om.Hidden, "Hidden mismatch")
			assert.Equal(t, tt.want.Required, om.Required, "Required mismatch")
			assert.Equal(t, tt.want.Title, om.Title, "Title mismatch")
			assert.Equal(t, tt.want.Description, om.Description, "Description mismatch")
			assert.Equal(t, tt.want.Examples, om.Examples, "Examples mismatch")

			if tt.want.Extensions != nil {
				require.NotNil(t, om.Extensions, "Extensions should not be nil")
				assert.Equal(t, tt.want.Extensions, om.Extensions, "Extensions mismatch")
			} else {
				assert.Nil(t, om.Extensions, "Extensions should be nil")
			}
		})
	}
}

func TestParseOpenAPI_ExtensionValidation(t *testing.T) {
	tests := []struct {
		name        string
		tagValue    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid extension - x-custom",
			tagValue: "x-custom=value",
			wantErr:  false,
		},
		{
			name:     "invalid extension - x- (error)",
			tagValue: "x-=value",
			wantErr:  true, // "x-=" is too short, returns error
		},
		{
			name:     "valid extension - x-123",
			tagValue: "x-123=value",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: "TestField",
			}

			_, err := ParseOpenAPITag(field, 0, tt.tagValue)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseOpenAPI_ComplexScenarios(t *testing.T) {
	t.Run("real world example - user field", func(t *testing.T) {
		field := reflect.StructField{
			Name: "UserID",
		}

		tagValue := "readOnly,deprecated,required,title=User Identifier,description='Unique identifier for the user',examples=12345,x-custom-feature=enabled,x-vendor-tool=test"

		result, err := ParseOpenAPITag(field, 0, tagValue)
		require.NoError(t, err)

		om, ok := result.(*OpenAPIMetadata)
		require.True(t, ok)
		assert.Equal(t, boolPtr(true), om.ReadOnly)
		assert.Equal(t, boolPtr(true), om.Deprecated)
		assert.Equal(t, boolPtr(true), om.Required)
		assert.Equal(t, "User Identifier", om.Title)
		assert.Equal(t, "Unique identifier for the user", om.Description)
		assert.Equal(t, []any{12345.0}, om.Examples)
		assert.Equal(t, "enabled", om.Extensions["x-custom-feature"])
		assert.Equal(t, "test", om.Extensions["x-vendor-tool"])
	})

	t.Run("extension with JSON value", func(t *testing.T) {
		field := reflect.StructField{
			Name: "Config",
		}

		tagValue := "x-config='{\"key\":\"value\"}'"

		result, err := ParseOpenAPITag(field, 0, tagValue)
		require.NoError(t, err)

		om, ok := result.(*OpenAPIMetadata)
		require.True(t, ok)
		assert.Equal(t, "{\"key\":\"value\"}", om.Extensions["x-config"])
	})

	t.Run("all boolean flags with explicit values", func(t *testing.T) {
		field := reflect.StructField{
			Name: "Field",
		}

		tagValue := "readOnly=true,writeOnly=false,deprecated=true,hidden=false,required=true"

		result, err := ParseOpenAPITag(field, 0, tagValue)
		require.NoError(t, err)

		om, ok := result.(*OpenAPIMetadata)
		require.True(t, ok)
		assert.Equal(t, boolPtr(true), om.ReadOnly)
		assert.Equal(t, boolPtr(false), om.WriteOnly)
		assert.Equal(t, boolPtr(true), om.Deprecated)
		assert.Equal(t, boolPtr(false), om.Hidden)
		assert.Equal(t, boolPtr(true), om.Required)
	})
}
