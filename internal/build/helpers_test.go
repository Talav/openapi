package build

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talav/openapi/config"
)

func TestGetSchemaHint(t *testing.T) {
	tests := []struct {
		name        string
		typ         reflect.Type
		fieldName   string
		operationID string
		want        string
	}{
		{
			name:        "type with name",
			typ:         reflect.TypeOf(struct{ Name string }{}),
			fieldName:   "Email",
			operationID: "",
			want:        "Email",
		},
		{
			name:        "anonymous type with operationID",
			typ:         reflect.TypeOf(struct{ Name string }{}),
			fieldName:   "Email",
			operationID: "createUser",
			want:        "createUserEmail", // Anonymous type has empty Name(), so uses operationID
		},
		{
			name:        "named type",
			typ:         reflect.TypeOf((*User)(nil)).Elem(),
			fieldName:   "Email",
			operationID: "",
			want:        "UserEmail",
		},
		{
			name:        "named type with operationID",
			typ:         reflect.TypeOf((*User)(nil)).Elem(),
			fieldName:   "Email",
			operationID: "createUser",
			want:        "UserEmail", // type.Name() takes priority
		},
		{
			name:        "no type name, no operationID",
			typ:         reflect.TypeOf(struct{}{}),
			fieldName:   "Email",
			operationID: "",
			want:        "Email",
		},
		{
			name:        "no type name, with operationID",
			typ:         reflect.TypeOf(struct{}{}),
			fieldName:   "Email",
			operationID: "createUser",
			want:        "createUserEmail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSchemaHint(tt.typ, tt.fieldName, tt.operationID)
			assert.Equal(t, tt.want, got)
		})
	}
}

// User is a test type for named type tests.
type User struct {
	Name string
}

func TestDeref(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want reflect.Type
	}{
		{
			name: "non-pointer type",
			typ:  reflect.TypeOf(""),
			want: reflect.TypeOf(""),
		},
		{
			name: "single pointer",
			typ:  reflect.TypeOf((*string)(nil)),
			want: reflect.TypeOf(""),
		},
		{
			name: "double pointer",
			typ:  reflect.TypeOf((**string)(nil)),
			want: reflect.TypeOf(""),
		},
		{
			name: "triple pointer",
			typ:  reflect.TypeOf((***string)(nil)),
			want: reflect.TypeOf(""),
		},
		{
			name: "pointer to struct",
			typ:  reflect.TypeOf((*User)(nil)),
			want: reflect.TypeOf(User{}),
		},
		{
			name: "pointer to int",
			typ:  reflect.TypeOf((*int)(nil)),
			want: reflect.TypeOf(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deref(tt.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{
			name:  "true bool",
			input: true,
			want:  true,
		},
		{
			name:  "false bool",
			input: false,
			want:  false,
		},
		{
			name:  "pointer to true",
			input: boolPtr(true),
			want:  true,
		},
		{
			name:  "pointer to false",
			input: boolPtr(false),
			want:  false,
		},
		{
			name:  "nil pointer",
			input: (*bool)(nil),
			want:  false,
		},
		{
			name:  "other type",
			input: "not a bool",
			want:  false,
		},
		{
			name:  "int",
			input: 42,
			want:  false,
		},
		{
			name:  "nil interface",
			input: nil,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toBool(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindBodyField(t *testing.T) {
	cfg := config.DefaultTagConfig()

	// Create test structs using reflection
	type WithBody struct {
		Name  string `json:"name"`
		Body  string `body:""`
		Email string `json:"email"`
	}

	type WithoutBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type CustomBodyTag struct {
		Name    string `json:"name"`
		Payload string `request:"" json:"payload"`
	}

	type EmptyStruct struct{}

	tests := []struct {
		name       string
		structType reflect.Type
		cfg        config.TagConfig
		wantIndex  int
		wantNil    bool
	}{
		{
			name:       "field with body tag",
			structType: reflect.TypeOf(WithBody{}),
			cfg:        cfg,
			wantIndex:  1,
			wantNil:    false,
		},
		{
			name:       "no body field",
			structType: reflect.TypeOf(WithoutBody{}),
			cfg:        cfg,
			wantIndex:  -1,
			wantNil:    true,
		},
		{
			name:       "empty struct",
			structType: reflect.TypeOf(EmptyStruct{}),
			cfg:        cfg,
			wantIndex:  -1,
			wantNil:    true,
		},
		{
			name:       "custom body tag name",
			structType: reflect.TypeOf(CustomBodyTag{}),
			cfg:        config.NewTagConfig("schema", "request", "openapi", "validate", "default", "requires"),
			wantIndex:  1,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create metadata with the test config
			testMetadata := NewMetadata(tt.cfg)
			structMeta, err := testMetadata.GetStructMetadata(tt.structType)
			require.NoError(t, err)

			got := findBodyField(structMeta, tt.cfg)

			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.wantIndex, got.Index)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
