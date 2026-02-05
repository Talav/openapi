package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTagConfig(t *testing.T) {
	cfg := DefaultTagConfig()

	assert.Equal(t, "schema", cfg.Schema)
	assert.Equal(t, "body", cfg.Body)
	assert.Equal(t, "openapi", cfg.OpenAPI)
	assert.Equal(t, "validate", cfg.Validate)
	assert.Equal(t, "default", cfg.Default)
	assert.Equal(t, "requires", cfg.Requires)
}

func TestNewTagConfig(t *testing.T) {
	cfg := NewTagConfig("s", "b", "o", "v", "d", "r")

	assert.Equal(t, "s", cfg.Schema)
	assert.Equal(t, "b", cfg.Body)
	assert.Equal(t, "o", cfg.OpenAPI)
	assert.Equal(t, "v", cfg.Validate)
	assert.Equal(t, "d", cfg.Default)
	assert.Equal(t, "r", cfg.Requires)
}

func TestMergeTagConfig(t *testing.T) {
	tests := []struct {
		name     string
		base     TagConfig
		override TagConfig
		want     TagConfig
	}{
		{
			name: "empty override does not change base",
			base: DefaultTagConfig(),
			override: TagConfig{
				Schema: "custom-schema",
			},
			want: TagConfig{
				Schema:   "custom-schema",
				Body:     "body",
				OpenAPI:  "openapi",
				Validate: "validate",
				Default:  "default",
				Requires: "requires",
			},
		},
		{
			name: "non-empty override replaces all fields",
			base: DefaultTagConfig(),
			override: TagConfig{
				Schema:   "s",
				Body:     "b",
				OpenAPI:  "o",
				Validate: "v",
				Default:  "d",
				Requires: "r",
			},
			want: TagConfig{
				Schema:   "s",
				Body:     "b",
				OpenAPI:  "o",
				Validate: "v",
				Default:  "d",
				Requires: "r",
			},
		},
		{
			name: "partial override only replaces specified fields",
			base: DefaultTagConfig(),
			override: TagConfig{
				Schema:  "custom-schema",
				OpenAPI: "custom-openapi",
			},
			want: TagConfig{
				Schema:   "custom-schema",
				Body:     "body",
				OpenAPI:  "custom-openapi",
				Validate: "validate",
				Default:  "default",
				Requires: "requires",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeTagConfig(tt.base, tt.override)
			assert.Equal(t, tt.want, result)
		})
	}
}
