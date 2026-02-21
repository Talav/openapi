package metadata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOpenAPITag_StructLevel(t *testing.T) {
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
			fieldName: "_",
			tagValue:  "",
			want:      &OpenAPIMetadata{},
		},
		{
			name:      "additionalProperties true",
			fieldName: "_",
			tagValue:  "additionalProperties=true",
			want: &OpenAPIMetadata{
				AdditionalProperties: boolPtr(true),
			},
		},
		{
			name:      "additionalProperties false",
			fieldName: "_",
			tagValue:  "additionalProperties=false",
			want: &OpenAPIMetadata{
				AdditionalProperties: boolPtr(false),
			},
		},
		{
			name:      "additionalProperties flag (true)",
			fieldName: "_",
			tagValue:  "additionalProperties",
			want: &OpenAPIMetadata{
				AdditionalProperties: boolPtr(true),
			},
		},
		{
			name:      "nullable true",
			fieldName: "_",
			tagValue:  "nullable=true",
			want: &OpenAPIMetadata{
				Nullable: boolPtr(true),
			},
		},
		{
			name:      "nullable false",
			fieldName: "_",
			tagValue:  "nullable=false",
			want: &OpenAPIMetadata{
				Nullable: boolPtr(false),
			},
		},
		{
			name:      "nullable flag (true)",
			fieldName: "_",
			tagValue:  "nullable",
			want: &OpenAPIMetadata{
				Nullable: boolPtr(true),
			},
		},
		{
			name:      "both options",
			fieldName: "_",
			tagValue:  "additionalProperties=false,nullable=true",
			want: &OpenAPIMetadata{
				AdditionalProperties: boolPtr(false),
				Nullable:             boolPtr(true),
			},
		},
		{
			name:      "both options reversed",
			fieldName: "_",
			tagValue:  "nullable=true,additionalProperties=false",
			want: &OpenAPIMetadata{
				AdditionalProperties: boolPtr(false),
				Nullable:             boolPtr(true),
			},
		},
		{
			name:        "unknown option returns error",
			fieldName:   "_",
			tagValue:    "additionalProperties=true,unknown=value",
			wantErr:     true,
			errContains: "unknown struct-level option",
		},
		{
			name:        "invalid tag parsing",
			fieldName:   "_",
			tagValue:    "additionalProperties=true,'unclosed quote",
			wantErr:     true,
			errContains: "failed to parse openapi tag",
		},
		{
			name:        "field-level options return error on struct level",
			fieldName:   "_",
			tagValue:    "readOnly,additionalProperties=false",
			wantErr:     true,
			errContains: "unknown struct-level option",
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

			assert.Equal(t, tt.want.AdditionalProperties, om.AdditionalProperties, "AdditionalProperties mismatch")
			assert.Equal(t, tt.want.Nullable, om.Nullable, "Nullable mismatch")
		})
	}
}
