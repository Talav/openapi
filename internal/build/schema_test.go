package build

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talav/openapi/config"
	"github.com/talav/openapi/internal/model"
)

func TestSchemaGenerator_PrimitiveTypes(t *testing.T) {
	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("", metadata, config.DefaultTagConfig())

	tests := []struct {
		name       string
		typ        any
		wantType   string
		wantFormat string
	}{
		{"string", "", "string", ""},
		{"int", 0, "integer", "int64"},
		{"int8", int8(0), "integer", "int32"},
		{"int16", int16(0), "integer", "int32"},
		{"int32", int32(0), "integer", "int32"},
		{"int64", int64(0), "integer", "int64"},
		{"uint", uint(0), "integer", "int64"},
		{"uint8", uint8(0), "integer", "int32"},
		{"uint16", uint16(0), "integer", "int32"},
		{"uint32", uint32(0), "integer", "int32"},
		{"uint64", uint64(0), "integer", "int64"},
		{"float32", float32(0), "number", "float"},
		{"float64", float64(0), "number", "double"},
		{"bool", false, "boolean", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := gen.Schema(reflect.TypeOf(tt.typ))
			require.NotNil(t, schema)
			assert.Equal(t, tt.wantType, schema.Type)
			if tt.wantFormat != "" {
				assert.Equal(t, tt.wantFormat, schema.Format)
			}
		})
	}
}

func TestSchemaGenerator_NestedStruct(t *testing.T) {
	type Address struct {
		City string `json:"city"`
	}
	type User struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf(User{}))

	require.NotNil(t, schema)
	assert.Equal(t, "#/components/schemas/User", schema.Ref)

	// Check both schemas are generated
	schemas := gen.Schemas()
	assert.Contains(t, schemas, "User")
	assert.Contains(t, schemas, "Address")
}

func TestSchemaGenerator_Array(t *testing.T) {
	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf([]string{}))

	require.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "string", schema.Items.Type)
}

func TestSchemaGenerator_ArrayOfStructs(t *testing.T) {
	type Item struct {
		ID int `json:"id"`
	}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf([]Item{}))

	require.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "#/components/schemas/Item", schema.Items.Ref)

	schemas := gen.Schemas()
	assert.Contains(t, schemas, "Item")
}

func TestSchemaGenerator_Map(t *testing.T) {
	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf(map[string]string{}))

	require.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.NotNil(t, schema.Additional)
	assert.NotNil(t, schema.Additional.Schema)
	assert.Equal(t, "string", schema.Additional.Schema.Type)
}

func TestSchemaGenerator_Pointer(t *testing.T) {
	type User struct {
		ID int `json:"id"`
	}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf(&User{}))

	require.NotNil(t, schema)
	// Pointer should resolve to the underlying type
	assert.Equal(t, "#/components/schemas/User", schema.Ref)
}

func TestSchemaGenerator_StructFeatures(t *testing.T) {
	type WithValidation struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"min=0,max=150"`
	}

	type WithEnum struct {
		Value string `json:"value" validate:"oneof=active inactive pending"`
	}

	type WithOpenAPI struct {
		ID   int    `json:"id" openapi:"readOnly,title=User ID"`
		Name string `json:"name" openapi:"title=Name"`
	}

	type WithTime struct {
		CreatedAt time.Time `json:"created_at"`
	}

	tests := []struct {
		name        string
		structType  any
		schemaName  string
		checkSchema func(t *testing.T, s *model.Schema)
	}{
		{
			name:       "validation tags",
			structType: WithValidation{},
			schemaName: "WithValidation",
			checkSchema: func(t *testing.T, s *model.Schema) {
				t.Helper()
				assert.Equal(t, "object", s.Type)
				assert.NotEmpty(t, s.Properties)
			},
		},
		{
			name:       "enum from validation",
			structType: WithEnum{},
			schemaName: "WithEnum",
			checkSchema: func(t *testing.T, s *model.Schema) {
				t.Helper()
				assert.NotEmpty(t, s.Properties)
			},
		},
		{
			name:       "openapi tags",
			structType: WithOpenAPI{},
			schemaName: "WithOpenAPI",
			checkSchema: func(t *testing.T, s *model.Schema) {
				t.Helper()
				assert.NotEmpty(t, s.Properties)
			},
		},
		{
			name:       "time type",
			structType: WithTime{},
			schemaName: "WithTime",
			checkSchema: func(t *testing.T, s *model.Schema) {
				t.Helper()
				assert.NotEmpty(t, s.Properties)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := NewMetadata(config.DefaultTagConfig())
			gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

			schema := gen.Schema(reflect.TypeOf(tt.structType))
			require.NotNil(t, schema)

			schemas := gen.Schemas()
			assert.Contains(t, schemas, tt.schemaName)

			if tt.checkSchema != nil {
				tt.checkSchema(t, schemas[tt.schemaName])
			}
		})
	}
}

func TestSchemaGenerator_ComplexStructJSON(t *testing.T) {
	type Address struct {
		Street  string `json:"street" validate:"required"`
		City    string `json:"city" validate:"required"`
		ZipCode string `json:"zip_code" validate:"required,len=5"`
	}

	type User struct {
		ID      int      `json:"id" openapi:"readOnly"`
		Name    string   `json:"name" validate:"required,min=3,max=100"`
		Email   string   `json:"email" validate:"required,email"`
		Age     int      `json:"age" validate:"min=0,max=150"`
		Tags    []string `json:"tags"`
		Address Address  `json:"address" validate:"required"`
	}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	gen.Schema(reflect.TypeOf(User{}))
	schemas := gen.Schemas()

	jsonBytes, err := json.MarshalIndent(schemas, "", "  ")
	require.NoError(t, err)

	// Verify JSON structure
	var parsed map[string]any
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	// Check both schemas exist
	assert.Contains(t, parsed, "User")
	assert.Contains(t, parsed, "Address")
}

func TestSchemaGenerator_EmptyStruct(t *testing.T) {
	type Empty struct{}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf(Empty{}))
	require.NotNil(t, schema)

	schemas := gen.Schemas()
	emptySchema := schemas["Empty"]
	assert.Equal(t, "object", emptySchema.Type)
}

func TestSchemaGenerator_InterfaceType(t *testing.T) {
	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("", metadata, config.DefaultTagConfig())

	schema := gen.Schema(reflect.TypeOf((*any)(nil)).Elem())
	require.NotNil(t, schema)

	// Interface should be treated as any type
	assert.Empty(t, schema.Type)
}

func TestSchemaGenerator_Caching(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	metadata := NewMetadata(config.DefaultTagConfig())
	gen := NewSchemaGenerator("#/components/schemas/", metadata, config.DefaultTagConfig())

	// Generate schema twice
	schema1 := gen.Schema(reflect.TypeOf(User{}))
	schema2 := gen.Schema(reflect.TypeOf(User{}))

	// Should return same reference
	assert.Equal(t, schema1.Ref, schema2.Ref)

	// Should only have one schema in cache
	schemas := gen.Schemas()
	assert.Len(t, schemas, 1)
}
