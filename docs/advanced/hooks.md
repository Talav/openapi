# Schema Hooks

Transform schemas programmatically using hooks.

## Overview

Schema hooks let you customize how types are converted to OpenAPI schemas. Use them for:

- Custom type mappings
- Adding computed properties
- Modifying existing schemas
- Supporting domain-specific types

## The SchemaTransformer Interface

Implement this interface to modify schemas after generation:

```go
type SchemaTransformer interface {
    TransformSchema(schema *Schema) (*Schema, error)
}
```

### Basic Example

```go
import "github.com/talav/openapi/hook"

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func (User) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add a custom property
    schema.Extensions = map[string]any{
        "x-resource-type": "user",
        "x-searchable":    true,
    }
    
    return schema, nil
}
```

Now when `User` is used, the generated schema includes your extensions.

### Modifying Properties

```go
type Product struct {
    ID    int     `json:"id"`
    Price float64 `json:"price"`
}

func (Product) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add currency information to price field
    if priceSchema, ok := schema.Properties["price"]; ok {
        priceSchema.Extensions = map[string]any{
            "x-currency": "USD",
        }
    }
    
    // Add computed field
    schema.Properties["price_formatted"] = &hook.Schema{
        Type:        "string",
        Description: "Formatted price with currency symbol",
        Examples:    []any{"$9.99"},
        Extensions: map[string]any{
            "x-computed": true,
        },
    }
    
    return schema, nil
}
```

## The SchemaProvider Interface

Implement this to completely replace schema generation:

```go
type SchemaProvider interface{
    ProvideSchema() (*Schema, error)
}
```

### Custom Schema

```go
type CustomID string

func (CustomID) ProvideSchema() (*hook.Schema, error) {
    return &hook.Schema{
        Type:        "string",
        Format:      "uuid",
        Pattern:     "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
        Description: "UUID v4 identifier",
        Examples:    []any{"550e8400-e29b-41d4-a716-446655440000"},
    }, nil
}

type User struct {
    ID   CustomID `json:"id"`
    Name string   `json:"name"`
}
```

The `ID` field will use your custom schema instead of the default string schema.

## Real-World Examples

### Enum Types

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusClosed
)

func (Status) ProvideSchema() (*hook.Schema, error) {
    return &hook.Schema{
        Type:        "string",
        Enum:        []any{"pending", "active", "closed"},
        Description: "Status of the resource",
    }, nil
}
```

### Time Types

```go
type Timestamp time.Time

func (Timestamp) ProvideSchema() (*hook.Schema, error) {
    return &hook.Schema{
        Type:        "string",
        Format:      "date-time",
        Description: "ISO 8601 timestamp",
        Examples:    []any{"2024-03-15T10:30:00Z"},
    }, nil
}
```

### Money Types

```go
type Money struct {
    Amount   int    `json:"amount"`   // cents
    Currency string `json:"currency"`
}

func (Money) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add description
    schema.Description = "Money amount in smallest currency unit (e.g., cents for USD)"
    
    // Add example
    schema.Example = map[string]any{
        "amount":   999,
        "currency": "USD",
    }
    
    // Mark currency as required
    schema.Required = []string{"amount", "currency"}
    
    // Add constraints
    if amountSchema, ok := schema.Properties["amount"]; ok {
        amountSchema.Minimum = &hook.Bound{Value: 0}
        amountSchema.Description = "Amount in smallest currency unit (cents)"
    }
    
    if currencySchema, ok := schema.Properties["currency"]; ok {
        currencySchema.Pattern = "^[A-Z]{3}$"
        currencySchema.Description = "ISO 4217 currency code"
        currencySchema.Examples = []any{"USD", "EUR", "GBP"}
    }
    
    return schema, nil
}
```

### Polymorphic Types

```go
type Event struct {
    Type string `json:"type"`
    Data any    `json:"data"`
}

func (Event) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add discriminator for polymorphism
    schema.Discriminator = &hook.Discriminator{
        PropertyName: "type",
        Mapping: map[string]string{
            "user_created":  "#/components/schemas/UserCreatedEvent",
            "user_updated":  "#/components/schemas/UserUpdatedEvent",
            "user_deleted":  "#/components/schemas/UserDeletedEvent",
        },
    }
    
    return schema, nil
}
```

### Database IDs

```go
type DatabaseID int64

func (DatabaseID) ProvideSchema() (*hook.Schema, error) {
    return &hook.Schema{
        Type:        "integer",
        Format:      "int64",
        Minimum:     &hook.Bound{Value: 1},
        Description: "Database primary key",
        Examples:    []any{1, 42, 1000},
        Extensions: map[string]any{
            "x-database-type": "BIGINT",
        },
    }, nil
}
```

## Error Handling

Return errors from hooks when transformation fails:

```go
func (User) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    if schema == nil {
        return nil, fmt.Errorf("schema cannot be nil")
    }
    
    // Validate schema before transforming
    if schema.Type != "object" {
        return nil, fmt.Errorf("User schema must be an object, got %s", schema.Type)
    }
    
    // Transform...
    
    return schema, nil
}
```

## Combining Hooks with Tags

Hooks work alongside struct tags:

```go
type Product struct {
    ID    int     `json:"id" openapi:"readOnly"`
    Price float64 `json:"price" validate:"min=0"`
}

func (Product) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Tags are already applied, now add more
    schema.Extensions = map[string]any{
        "x-resource": "product",
    }
    
    return schema, nil
}
```

Execution order:
1. Generate base schema from Go type
2. Apply struct tag metadata (json, validate, openapi, etc.)
3. Call `TransformSchema` if implemented

## Advanced Patterns

### Conditional Properties

```go
func (Order) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add conditional validation in OpenAPI 3.1
    schema.If = &hook.Schema{
        Properties: map[string]*hook.Schema{
            "payment_method": {
                Const: "credit_card",
            },
        },
    }
    
    schema.Then = &hook.Schema{
        Required: []string{"card_number", "cvv"},
    }
    
    return schema, nil
}
```

### Computed Fields

```go
func (User) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add virtual/computed fields
    schema.Properties["full_name"] = &hook.Schema{
        Type:        "string",
        Description: "Computed from first_name and last_name",
        ReadOnly:    true,
        Extensions: map[string]any{
            "x-computed": true,
        },
    }
    
    return schema, nil
}
```

### Schema Versioning

```go
func (User) TransformSchema(schema *hook.Schema) (*hook.Schema, error) {
    // Add version information
    schema.Extensions = map[string]any{
        "x-schema-version": "2.0",
        "x-deprecated-fields": []string{"username"},
    }
    
    // Mark deprecated fields
    if usernameSchema, ok := schema.Properties["username"]; ok {
        usernameSchema.Deprecated = true
        usernameSchema.Description += " (deprecated: use email instead)"
    }
    
    return schema, nil
}
```

## Testing Hooks

```go
func TestUserSchemaTransform(t *testing.T) {
    user := User{}
    
    // Generate base schema (you'll need access to internal build package)
    schema := &hook.Schema{
        Type: "object",
        Properties: map[string]*hook.Schema{
            "id":   {Type: "integer"},
            "name": {Type: "string"},
        },
    }
    
    // Apply transformation
    result, err := user.TransformSchema(schema)
    require.NoError(t, err)
    
    // Verify extensions
    assert.Equal(t, "user", result.Extensions["x-resource-type"])
    assert.True(t, result.Extensions["x-searchable"].(bool))
}
```

## Next Steps

- [Custom Tags](custom-tags.md) - Extend the tag system
