# Core Concepts

Understanding how the library transforms Go structs into OpenAPI specifications.

## The Tag System

The library uses multiple struct tags, each with a specific purpose:

### 1. `json` - Property Names

Standard Go tag for JSON field names, used to name properties in schemas:

```go
type User struct {
    ID   int    `json:"id"`           // Property: "id"
    Name string `json:"name"`         // Property: "name"
    Age  int    `json:"age,omitempty"` // Property: "age" (omitempty ignored)
}
```

### 2. `schema` - Parameter Metadata

Defines where parameters come from and how they're serialized:

```go
type Request struct {
    ID      string `schema:"id,location=path"`        // Path parameter
    Search  string `schema:"q,location=query"`        // Query parameter
    APIKey  string `schema:"X-API-Key,location=header"` // Header
    Session string `schema:"session,location=cookie"`   // Cookie
}
```

Learn more: [Tag Reference (talav/schema)](https://talav.github.io/schema/)

### 3. `body` - Request/Response Bodies

Marks which field contains the body:

```go
type CreateUserRequest struct {
    Body User `body:"structured"` // JSON/XML/form body
}

type FileUploadRequest struct {
    Body []byte `body:"file"` // Binary file upload
}
```

### 4. `validate` - Validation Rules

Validation tags transform into OpenAPI schema constraints:

```go
type User struct {
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}
```

Becomes:

```json
{
  "name": {
    "type": "string",
    "minLength": 3,
    "maxLength": 50
  },
  "email": {
    "type": "string",
    "format": "email"
  },
  "age": {
    "type": "integer",
    "minimum": 18,
    "maximum": 120
  },
  "required": ["name", "email"]
}
```

Learn more: [Validation](../guides/validation.md)

### 5. `openapi` - OpenAPI Metadata

Additional OpenAPI-specific properties:

```go
type Product struct {
    ID    string  `json:"id" openapi:"readOnly"`
    Name  string  `json:"name" openapi:"title=Product Name,description=Display name"`
    Price float64 `json:"price" openapi:"examples=9.99|19.99|29.99"`
}
```

Learn more: [Metadata](../guides/metadata.md)

### 6. `default` - Default Values

Specify defaults for optional fields:

```go
type Config struct {
    Host string `json:"host" default:"localhost"`
    Port int    `json:"port" default:"8080"`
}
```

### 7. `requires` - Conditional Requirements

Fields that become required when another field is present:

```go
type PaymentRequest struct {
    Method      string `json:"method"`
    CardNumber  string `json:"card_number" requires:"cvv,expiry"`
    CVV         string `json:"cvv"`
    Expiry      string `json:"expiry"`
}
```

When `card_number` is provided, both `cvv` and `expiry` become required.

## How It Works

The generation process:

1. **Parse struct metadata** - Extract tags and reflect on types
2. **Build operations** - Create paths, parameters, request bodies
3. **Generate schemas** - Transform validation rules, apply metadata
4. **Validate spec** - Check against OpenAPI 3.0/3.1 schema
5. **Return JSON/YAML** - Serialized specification

## Request vs Response Mapping

### Request Structs

Request structs can have parameters and a body:

```go
type CreatePostRequest struct {
    OrgID  string `schema:"org_id,location=path"` // Parameter
    APIKey string `schema:"X-API-Key,location=header"` // Parameter
    Body struct {
        Title   string `json:"title" validate:"required"`
        Content string `json:"content" validate:"required"`
    } `body:"structured"` // Request body
}
```

### Response Structs

Response structs typically have a body and optional headers:

```go
type CreatePostResponse struct {
    Location string `schema:"Location,location=header"` // Response header
    Body struct {
        ID        int    `json:"id"`
        CreatedAt string `json:"created_at"`
    } `body:"structured"` // Response body
}
```

## Schema Generation

### Type Mapping

Go types map to OpenAPI/JSON Schema types:

| Go Type | OpenAPI Type | Format |
|---------|-------------|---------|
| `string` | string | - |
| `int`, `int32` | integer | int32 |
| `int64` | integer | int64 |
| `float32` | number | float |
| `float64` | number | double |
| `bool` | boolean | - |
| `[]T` | array | items: T |
| `map[string]T` | object | additionalProperties: T |
| `struct` | object | properties: fields |
| `[]byte` | string | binary (files) or base64 (JSON) |

### Nested Structures

Nested structs become nested schemas:

```go
type User struct {
    Name    string  `json:"name"`
    Address Address `json:"address"`
}

type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
}
```

The `Address` type becomes a reusable component referenced via `$ref`.

## Component Reuse

Types used in multiple places are generated once:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type CreateUserRequest struct {
    Body User `body:"structured"`
}

type GetUserResponse struct {
    Body User `body:"structured"`
}
```

`User` appears once in `components/schemas` and is referenced from both operations.

## OpenAPI Versions

Choose your target version:

```go
// OpenAPI 3.0.4 (maximum compatibility)
api := openapi.NewAPI(
    openapi.WithVersion("3.0.4"),
)

// OpenAPI 3.1.2 (modern features)
api := openapi.NewAPI(
    openapi.WithVersion("3.1.2"),
)
```

Learn more: [OpenAPI Versions](../guides/versions.md)

## Next Steps

- [Tag Reference (talav/schema)](https://talav.github.io/schema/) - Master `schema`/`body` tag semantics
- [Validation](../guides/validation.md) - Transform validation rules
- [Metadata](../guides/metadata.md) - Add rich documentation
