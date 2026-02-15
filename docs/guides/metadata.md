# Metadata

Add rich OpenAPI metadata using the `openapi`, `default`, and `requires` tags.

## The `openapi` Tag

The `openapi` tag adds OpenAPI-specific properties that don't affect validation.

### Basic Properties

```go
type Product struct {
    ID    string `json:"id" openapi:"readOnly"`
    Name  string `json:"name" openapi:"title=Product Name"`
    Price float64 `json:"price" openapi:"description=Price in USD"`
}
```

### Field-Level Options

| Option | Description | Example |
|--------|-------------|---------|
| `readOnly` | Field appears in responses only | `openapi:"readOnly"` |
| `writeOnly` | Field appears in requests only | `openapi:"writeOnly"` |
| `deprecated` | Field is deprecated | `openapi:"deprecated"` |
| `hidden` | Exclude from schema | `openapi:"hidden"` |
| `required` | Override required status | `openapi:"required"` |
| `title` | Schema title | `openapi:"title=User ID"` |
| `description` | Field description | `openapi:"description=Unique identifier"` |
| `format` | Data format hint | `openapi:"format=date-time"` |
| `examples` | Example values | `openapi:"examples=val1|val2"` |

### ReadOnly and WriteOnly

Useful for fields that appear in only one direction:

```go
type User struct {
    // Server-generated, appears only in responses
    ID        int       `json:"id" openapi:"readOnly"`
    CreatedAt time.Time `json:"created_at" openapi:"readOnly"`
    
    // Sent in requests, never returned
    Password string `json:"password" openapi:"writeOnly"`
    
    // Appears in both
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### Deprecated Fields

Mark fields as deprecated to signal future removal:

```go
type User struct {
    Name     string `json:"name"`
    FullName string `json:"full_name"` // New field
    Username string `json:"username" openapi:"deprecated"` // Old field
}
```

### Hidden Fields

Exclude fields from the generated schema:

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Internal string `json:"internal" openapi:"hidden"` // Not in spec
}
```

### Examples

Provide example values for documentation:

```go
type Product struct {
    Name  string  `json:"name" openapi:"examples=Widget|Gadget|Tool"`
    Price float64 `json:"price" openapi:"examples=9.99|19.99|29.99"`
    SKU   string  `json:"sku" openapi:"examples=WDG-001"`
}
```

Multiple examples separated by `|`.

### Format Hints

Specify data formats beyond standard JSON Schema formats:

```go
type Event struct {
    Date     string `json:"date" openapi:"format=date"`          // yyyy-mm-dd
    Time     string `json:"time" openapi:"format=time"`          // hh:mm:ss
    DateTime string `json:"datetime" openapi:"format=date-time"` // RFC3339
    Duration string `json:"duration" openapi:"format=duration"`   // ISO 8601
}
```

### Custom Extensions

Add vendor-specific extensions (must start with `x-`):

```go
type User struct {
    ID int `json:"id" openapi:"x-internal-id=users,x-searchable=true"`
}
```

Generates:

```json
{
  "id": {
    "type": "integer",
    "x-internal-id": "users",
    "x-searchable": true
  }
}
```

## The `default` Tag

Specify default values for optional fields:

```go
type Config struct {
    Host    string `json:"host" default:"localhost"`
    Port    int    `json:"port" default:"8080"`
    Debug   bool   `json:"debug" default:"false"`
    Timeout int    `json:"timeout" default:"30"`
}
```

Generates:

```json
{
  "host": {
    "type": "string",
    "default": "localhost"
  },
  "port": {
    "type": "integer",
    "default": 8080
  }
}
```

### Type Conversion

Default values are automatically converted to the correct type:

```go
type Settings struct {
    Count   int     `json:"count" default:"10"`      // String → int
    Enabled bool    `json:"enabled" default:"true"`  // String → bool
    Rate    float64 `json:"rate" default:"0.5"`      // String → float64
}
```

## The `requires` Tag

Specify conditional required fields (OpenAPI 3.1 `dependentRequired`):

```go
type Payment struct {
    Method     string `json:"method"`
    CardNumber string `json:"card_number" requires:"cvv,expiry"`
    CVV        string `json:"cvv"`
    Expiry     string `json:"expiry"`
}
```

When `card_number` is present, `cvv` and `expiry` become required.

Generates:

```json
{
  "type": "object",
  "properties": {
    "method": { "type": "string" },
    "card_number": { "type": "string" },
    "cvv": { "type": "string" },
    "expiry": { "type": "string" }
  },
  "dependentRequired": {
    "card_number": ["cvv", "expiry"]
  }
}
```

### Multiple Dependencies

```go
type Form struct {
    Type         string `json:"type"`
    PersonName   string `json:"person_name" requires:"person_age"`
    PersonAge    int    `json:"person_age"`
    CompanyName  string `json:"company_name" requires:"company_tax_id"`
    CompanyTaxID string `json:"company_tax_id"`
}
```

## Struct-Level Metadata

Use the blank identifier `_` for struct-level properties:

```go
type User struct {
    _ struct{} `openapi:"additionalProperties=false,nullable=true"`
    
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

- `additionalProperties=false`: Disallow extra properties
- `nullable=true`: Allow null values for the entire object

## Combining Tags

Use multiple tags together for complete schemas:

```go
type User struct {
    ID    int    `json:"id" openapi:"readOnly,title=User ID,description=Unique identifier,examples=1|42|100"`
    Name  string `json:"name" validate:"required,min=3,max=50" openapi:"title=Full Name,description=User's display name,examples=John Doe|Jane Smith"`
    Email string `json:"email" validate:"required,email" openapi:"description=Contact email,examples=user@example.com" default:""`
    Age   int    `json:"age" validate:"min=18" openapi:"description=Age in years" default:"18"`
}
```

## Next Steps

- [Examples](examples.md) - Generate examples automatically
- [Security](security.md) - Add authentication schemes
- [Advanced Hooks](../advanced/hooks.md) - Transform schemas programmatically
