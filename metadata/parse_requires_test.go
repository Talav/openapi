package metadata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequiresTag(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		tagValue  string
		want      *RequiresMetadata
		wantErr   bool
	}{
		{
			name:      "single required field",
			fieldName: "CreditCard",
			tagValue:  "billing_address",
			want: &RequiresMetadata{
				Fields: []string{"billing_address"},
			},
		},
		{
			name:      "multiple required fields",
			fieldName: "CreditCard",
			tagValue:  "billing_address,cardholder_name",
			want: &RequiresMetadata{
				Fields: []string{"billing_address", "cardholder_name"},
			},
		},
		{
			name:      "multiple required fields with spaces",
			fieldName: "CreditCard",
			tagValue:  "billing_address, cardholder_name, phone_number",
			want: &RequiresMetadata{
				Fields: []string{"billing_address", "cardholder_name", "phone_number"},
			},
		},
		{
			name:      "empty tag",
			fieldName: "Field",
			tagValue:  "",
			want: &RequiresMetadata{
				Fields: []string{},
			},
		},
		{
			name:      "whitespace only",
			fieldName: "Field",
			tagValue:  "   ",
			wantErr:   true, // tagparser doesn't allow whitespace-only keys
		},
		{
			name:      "quoted field names",
			fieldName: "Field",
			tagValue:  "'field1','field2 with spaces','field3'",
			want: &RequiresMetadata{
				Fields: []string{"field1", "field2 with spaces", "field3"},
			},
		},
		{
			name:      "quoted field with comma",
			fieldName: "Field",
			tagValue:  "field1,'field,with,comma',field2",
			want: &RequiresMetadata{
				Fields: []string{"field1", "field,with,comma", "field2"},
			},
		},
		{
			name:      "empty values filtered",
			fieldName: "Field",
			tagValue:  "field1,,field2,  ,field3",
			wantErr:   true, // tagparser doesn't allow empty keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: tt.fieldName,
			}

			result, err := ParseRequiresTag(field, 0, tt.tagValue)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			rm, ok := result.(*RequiresMetadata)
			require.True(t, ok, "result should be *RequiresMetadata")

			// Compare as sets (order may vary due to map iteration)
			assert.ElementsMatch(t, tt.want.Fields, rm.Fields, "Fields mismatch")
		})
	}
}
