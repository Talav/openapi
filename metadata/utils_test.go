package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        float64
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid integer string",
			input:   "42",
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "valid float string",
			input:   "3.14",
			want:    3.14,
			wantErr: false,
		},
		{
			name:    "valid negative float",
			input:   "-10.5",
			want:    -10.5,
			wantErr: false,
		},
		{
			name:    "valid zero",
			input:   "0",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "valid scientific notation",
			input:   "1e10",
			want:    1e10,
			wantErr: false,
		},
		{
			name:    "valid decimal with exponent",
			input:   "1.5e-2",
			want:    0.015,
			wantErr: false,
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			errContains: "empty value",
		},
		{
			name:        "invalid string",
			input:       "not a number",
			wantErr:     true,
			errContains: "invalid syntax",
		},
		{
			name:        "partial number",
			input:       "123abc",
			wantErr:     true,
			errContains: "invalid syntax",
		},
		{
			name:        "only decimal point",
			input:       ".",
			wantErr:     true,
			errContains: "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat64(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        int
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid positive integer",
			input:   "42",
			want:    42,
			wantErr: false,
		},
		{
			name:    "valid zero",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid large integer",
			input:   "2147483647",
			want:    2147483647,
			wantErr: false,
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			errContains: "empty value",
		},
		{
			name:        "negative number",
			input:       "-5",
			wantErr:     true,
			errContains: "value must be non-negative",
		},
		{
			name:        "invalid string",
			input:       "not a number",
			wantErr:     true,
			errContains: "invalid syntax",
		},
		{
			name:        "partial number",
			input:       "123abc",
			wantErr:     true,
			errContains: "invalid syntax",
		},
		{
			name:        "float string",
			input:       "3.14",
			wantErr:     true,
			errContains: "invalid syntax",
		},
		{
			name:        "only minus sign",
			input:       "-",
			wantErr:     true,
			errContains: "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInt(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *bool
		wantErr     bool
		errContains string
	}{
		{
			name:    "empty string returns true",
			input:   "",
			want:    boolPtr(true),
			wantErr: false,
		},
		{
			name:    "explicit true",
			input:   "true",
			want:    boolPtr(true),
			wantErr: false,
		},
		{
			name:    "explicit false",
			input:   "false",
			want:    boolPtr(false),
			wantErr: false,
		},
		{
			name:        "invalid value returns error",
			input:       "invalid",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "True (capitalized) returns error",
			input:       "True",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "FALSE (uppercase) returns error",
			input:       "FALSE",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "1 returns error",
			input:       "1",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "0 returns error",
			input:       "0",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "yes returns error",
			input:       "yes",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "no returns error",
			input:       "no",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
		{
			name:        "whitespace returns error",
			input:       " ",
			wantErr:     true,
			errContains: "invalid boolean value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBool(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}
