# Validation

Transform Go validation tags into OpenAPI schema constraints.

## Overview

The `validate` tag (from go-playground/validator) is automatically converted to OpenAPI/JSON Schema validation rules. This keeps your validation logic and API documentation in sync.

## String Validation

### Required Fields

```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}
```

Generates:

```json
{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "email": { "type": "string", "format": "email" }
  },
  "required": ["name", "email"]
}
```

### Length Constraints

```go
type User struct {
    Username string `json:"username" validate:"min=3,max=20"`
    Bio      string `json:"bio" validate:"max=500"`
    Code     string `json:"code" validate:"len=6"` // Exact length
}
```

Becomes:

```json
{
  "username": {
    "type": "string",
    "minLength": 3,
    "maxLength": 20
  },
  "bio": {
    "type": "string",
    "maxLength": 500
  },
  "code": {
    "type": "string",
    "minLength": 6,
    "maxLength": 6
  }
}
```

### Format Validation

```go
type Contact struct {
    Email string `json:"email" validate:"email"`
    URL   string `json:"url" validate:"url"`
    Phone string `json:"phone" validate:"e164"` // E.164 phone format
}
```

Maps to OpenAPI formats:

```json
{
  "email": {
    "type": "string",
    "format": "email"
  },
  "url": {
    "type": "string",
    "format": "uri"
  }
}
```

### Pattern Matching

```go
type Product struct {
    SKU  string `json:"sku" validate:"pattern=^[A-Z]{3}-[0-9]{4}$"`
    Slug string `json:"slug" validate:"pattern=^[a-z0-9-]+$"`
}
```

## Numeric Validation

### Range Constraints

```go
type Product struct {
    Price    float64 `json:"price" validate:"min=0,max=10000"`
    Quantity int     `json:"quantity" validate:"min=1"`
    Rating   float64 `json:"rating" validate:"min=0,max=5"`
    Age      int     `json:"age" validate:"min=18,max=120"`
}
```

Generates:

```json
{
  "price": {
    "type": "number",
    "minimum": 0,
    "maximum": 10000
  },
  "quantity": {
    "type": "integer",
    "minimum": 1
  }
}
```

### Exclusive Bounds

```go
type Range struct {
    Value int `json:"value" validate:"gt=0,lt=100"` // 0 < value < 100
}
```

Becomes:

```json
{
  "value": {
    "type": "integer",
    "exclusiveMinimum": 0,
    "exclusiveMaximum": 100
  }
}
```

### Multiple Of

```go
type Measurement struct {
    Increment float64 `json:"increment" validate:"multipleOf=0.5"`
    Count     int     `json:"count" validate:"multipleOf=10"`
}
```

## Array Validation

### Size Constraints

```go
type Form struct {
    Tags       []string `json:"tags" validate:"min=1,max=5"`
    Signatures []string `json:"signatures" validate:"min=2"`
    Options    []string `json:"options" validate:"len=3"` // Exactly 3
}
```

Generates:

```json
{
  "tags": {
    "type": "array",
    "items": { "type": "string" },
    "minItems": 1,
    "maxItems": 5
  }
}
```

### Unique Items

```go
type List struct {
    UniqueIDs []int `json:"unique_ids" validate:"unique"`
}
```

## Enum Validation

### String Enums

```go
type Order struct {
    Status string `json:"status" validate:"oneof=pending processing shipped delivered"`
}
```

Becomes:

```json
{
  "status": {
    "type": "string",
    "enum": ["pending", "processing", "shipped", "delivered"]
  }
}
```

### Numeric Enums

```go
type Priority struct {
    Level int `json:"level" validate:"oneof=1 2 3 4 5"`
}
```

## Combining Validators

Multiple validation rules work together:

```go
type User struct {
    Email    string `json:"email" validate:"required,email,max=100"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
    Username string `json:"username" validate:"required,min=3,max=20,pattern=^[a-z0-9_]+$"`
}
```

## Validation Mapping Reference

| Validator | OpenAPI Property | Example |
|-----------|-----------------|---------|
| `required` | `required` array | `validate:"required"` |
| `min=N` (string) | `minLength` | `validate:"min=3"` |
| `max=N` (string) | `maxLength` | `validate:"max=100"` |
| `len=N` (string) | `minLength` + `maxLength` | `validate:"len=6"` |
| `min=N` (number) | `minimum` | `validate:"min=0"` |
| `max=N` (number) | `maximum` | `validate:"max=100"` |
| `gt=N` | `exclusiveMinimum` | `validate:"gt=0"` |
| `lt=N` | `exclusiveMaximum` | `validate:"lt=100"` |
| `multipleOf=N` | `multipleOf` | `validate:"multipleOf=0.5"` |
| `email` | `format: email` | `validate:"email"` |
| `url` | `format: uri` | `validate:"url"` |
| `pattern=...` | `pattern` | `validate:"pattern=^[a-z]+$"` |
| `oneof=...` | `enum` | `validate:"oneof=a b c"` |
| `min=N` (array) | `minItems` | `validate:"min=1"` |
| `max=N` (array) | `maxItems` | `validate:"max=10"` |
| `unique` | `uniqueItems` | `validate:"unique"` |

## Overriding with `openapi` Tag

You can override validation-derived required status:

```go
type User struct {
    // Required by validate, but optional in OpenAPI docs
    Email string `json:"email" validate:"required" openapi:"required=false"`
    
    // Optional in code, but required in OpenAPI docs
    Name string `json:"name" openapi:"required"`
}
```

## Nested Validation

Validation works on nested structs:

```go
type User struct {
    Name    string  `json:"name" validate:"required"`
    Address Address `json:"address" validate:"required"`
}

type Address struct {
    Street string `json:"street" validate:"required"`
    City   string `json:"city" validate:"required"`
    Zip    string `json:"zip" validate:"required,len=5"`
}
```

All validation rules propagate to the generated schema.

## Complex Validation

### Conditional Required

Use the `requires` tag for fields that become required when another field is present:

```go
type Payment struct {
    Method     string `json:"method"`
    CardNumber string `json:"card_number" requires:"cvv,expiry"`
    CVV        string `json:"cvv"`
    Expiry     string `json:"expiry"`
}
```

In OpenAPI 3.1, this uses `dependentRequired`:

```json
{
  "dependentRequired": {
    "card_number": ["cvv", "expiry"]
  }
}
```

## Next Steps

- [Metadata](metadata.md) - Add descriptions, examples, and more
- [Tag Reference (talav/schema)](https://talav.github.io/schema/) - Parameter and body tag semantics
