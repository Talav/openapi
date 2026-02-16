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

Response structs can use either pattern:

**Simple pattern** - the struct is the response body:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

openapi.WithResponse(200, User{}) // User becomes the response schema
```

**Advanced pattern** - use body tag for headers:

```go
type UserResponse struct {
    ETag string `schema:"ETag,location=header"` // Response header
    Body User   `body:"structured"`              // Response body
}

openapi.WithResponse(200, UserResponse{}) // Includes headers
```

## Response Patterns

The library supports two patterns for responses:

### Simple Pattern (Most Common)

Pass your response struct directly - the entire struct becomes the response body:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type ErrorResponse struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

openapi.GET("/users/:id",
    openapi.WithResponse(200, User{}),           // ✅ Simple
    openapi.WithResponse(404, ErrorResponse{}),  // ✅ Simple
)
```

### Advanced Pattern (With Response Headers)

Use the `body` tag wrapper when you need to include response headers:

```go
type UserWithHeaders struct {
    ETag         string `schema:"ETag,location=header"`
    CacheControl string `schema:"Cache-Control,location=header"`
    Body         User   `body:"structured"`
}

openapi.GET("/users/:id",
    openapi.WithResponse(200, UserWithHeaders{}), // Includes headers
)
```

Use the wrapper pattern only when you need response headers. For standard responses, use the simple pattern.

### Custom Content Types

By default, responses use `application/json`. To return a different content type, implement the `ContentTypeProvider` interface on your response struct:

```go
type ProblemDetail struct {
    Type   string `json:"type"`
    Title  string `json:"title"`
    Status int    `json:"status"`
    Detail string `json:"detail"`
}

// ContentTypeProvider interface
func (ProblemDetail) ContentType(defaultType string) string {
    return "application/problem+json"
}

openapi.GET("/users/:id",
    openapi.WithResponse(404, ProblemDetail{}),
)
```

The generated OpenAPI spec will use `application/problem+json` instead of `application/json` for this response.

**Common use cases:**

- RFC 7807 Problem Details: `application/problem+json`
- HAL responses: `application/hal+json`
- JSON:API: `application/vnd.api+json`
- Custom vendor types: `application/vnd.myapp+json`

The method receives the default content type as a parameter, allowing conditional logic:

```go
func (r APIResponse) ContentType(defaultType string) string {
    if r.IsError {
        return "application/problem+json"
    }
    return defaultType // Use default application/json
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

// User appears in both operations
openapi.POST("/users",
    openapi.WithResponse(201, User{}),
),
openapi.GET("/users/:id",
    openapi.WithResponse(200, User{}),
)
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
