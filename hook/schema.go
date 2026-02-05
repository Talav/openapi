package hook

import (
	"reflect"

	"github.com/talav/openapi/internal/model"
)

// SchemaProvider is an interface that can be implemented by types to provide
// a custom schema for themselves, overriding the built-in schema generation.
// This can be used by custom types with their own special serialization rules.
type SchemaProvider interface {
	Schema(r SchemaRegistry) *model.Schema
}

// SchemaTransformer is an interface that can be implemented by types
// to transform the generated schema as needed.
// This can be used to leverage the default schema generation for a type,
// and arbitrarily modify parts of it.
type SchemaTransformer interface {
	TransformSchema(r SchemaRegistry, s *model.Schema) *model.Schema
}

// SchemaRegistry is a minimal interface for schema generation.
// It's used by SchemaProvider and SchemaTransformer implementations.
type SchemaRegistry interface {
	Schema(t reflect.Type) *model.Schema
}
