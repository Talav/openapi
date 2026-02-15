# Examples

Generate realistic examples for your API documentation.

## How Examples Are Defined

This library does not auto-generate example payloads from Go types or validation rules.
Examples are defined explicitly:

- Field-level examples via the `openapi:"examples=..."` tag
- Operation-level examples via `openapi/example` helpers (`WithRequestExample`, `WithResponseExample`)

## Field-Level Examples

Define examples using the `openapi` tag:

### Single Example

```go
type User struct {
    Name  string `json:"name" openapi:"examples=John Doe"`
    Email string `json:"email" openapi:"examples=john@example.com"`
    Age   int    `json:"age" openapi:"examples=30"`
}
```

### Multiple Examples

Separate examples with `|`:

```go
type Product struct {
    Name  string  `json:"name" openapi:"examples=Widget|Gadget|Tool"`
    Price float64 `json:"price" openapi:"examples=9.99|19.99|29.99"`
}
```

The first example is used as the default.

## Operation-Level Examples

Add examples to entire requests/responses:

```go
import "github.com/talav/openapi/example"

// Create example instance
userExample := example.New("john_doe",
    example.WithSummary("Example user"),
    example.WithDescription("A typical user account"),
    example.WithValue(map[string]any{
        "id":    42,
        "name":  "John Doe",
        "email": "john@example.com",
        "role":  "admin",
    }),
)

// Apply to operation
openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithRequestExample(userExample),
    openapi.WithResponse(201, User{}),
)
```

### Multiple Request Examples

```go
adminExample := example.New("admin_user",
    example.WithSummary("Admin user"),
    example.WithValue(map[string]any{
        "name":  "Admin User",
        "email": "admin@example.com",
        "role":  "admin",
    }),
)

regularExample := example.New("regular_user",
    example.WithSummary("Regular user"),
    example.WithValue(map[string]any{
        "name":  "Regular User",
        "email": "user@example.com",
        "role":  "user",
    }),
)

openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithRequestExample(adminExample),
    openapi.WithRequestExample(regularExample),
    openapi.WithResponse(201, User{}),
)
```

## Response Examples

### Success Example

```go
successExample := example.New("success",
    example.WithSummary("Successful creation"),
    example.WithValue(map[string]any{
        "id":         123,
        "name":       "New User",
        "email":      "user@example.com",
        "created_at": "2024-03-15T10:30:00Z",
    }),
)

openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{},
        openapi.WithResponseExample(successExample),
    ),
)
```

### Error Examples

```go
notFoundExample := example.New("not_found",
    example.WithSummary("User not found"),
    example.WithValue(map[string]any{
        "error":   "not_found",
        "message": "User with ID 999 does not exist",
    }),
)

validationExample := example.New("validation_error",
    example.WithSummary("Validation failed"),
    example.WithValue(map[string]any{
        "error": "validation_failed",
        "details": []map[string]string{
            {"field": "email", "message": "invalid email format"},
            {"field": "age", "message": "must be at least 18"},
        },
    }),
)

openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{}),
    openapi.WithResponse(400, Error{},
        openapi.WithResponseExample(validationExample),
    ),
    openapi.WithResponse(404, Error{},
        openapi.WithResponseExample(notFoundExample),
    ),
)
```

## Complex Examples

### Nested Structures

```go
type Order struct {
    ID    int     `json:"id"`
    Items []Item  `json:"items"`
    User  User    `json:"user"`
    Total float64 `json:"total"`
}

orderExample := example.New("sample_order",
    example.WithValue(map[string]any{
        "id": 789,
        "items": []map[string]any{
            {"id": 1, "name": "Widget", "price": 9.99, "quantity": 2},
            {"id": 2, "name": "Gadget", "price": 19.99, "quantity": 1},
        },
        "user": map[string]any{
            "id":    42,
            "name":  "John Doe",
            "email": "john@example.com",
        },
        "total": 39.97,
    }),
)
```

### Arrays

```go
usersExample := example.New("user_list",
    example.WithSummary("List of users"),
    example.WithValue([]map[string]any{
        {"id": 1, "name": "Alice", "email": "alice@example.com"},
        {"id": 2, "name": "Bob", "email": "bob@example.com"},
        {"id": 3, "name": "Charlie", "email": "charlie@example.com"},
    }),
)

openapi.GET("/users",
    openapi.WithResponse(200, []User{},
        openapi.WithResponseExample(usersExample),
    ),
)
```

## Parameter Examples

Examples for query/path/header parameters:

```go
type SearchRequest struct {
    Query    string   `schema:"q" openapi:"examples=golang|kubernetes|docker"`
    Category string   `schema:"category" openapi:"examples=tutorial|reference"`
    Tags     []string `schema:"tags" openapi:"examples=beginner|advanced"`
    Page     int      `schema:"page" default:"1" openapi:"examples=1|2|5"`
}
```

## External Examples

Reference examples from external URLs:

```go
externalExample := example.New("external_user",
    example.WithSummary("User from external source"),
    example.WithExternalValue("https://example.com/examples/user.json"),
)
```

## Best Practices

### Realistic Data

Use realistic data that helps developers understand your API:

```go
// ✗ Bad: generic placeholders
example.WithValue(map[string]any{
    "name":  "string",
    "email": "user@example.com",
})

// ✓ Good: realistic data
example.WithValue(map[string]any{
    "name":  "Sarah Johnson",
    "email": "sarah.johnson@company.com",
})
```

### Cover Edge Cases

Provide examples for different scenarios:

```go
// Normal case
normalExample := example.New("normal",
    example.WithValue(map[string]any{"age": 30}),
)

// Edge cases
minimumExample := example.New("minimum_age",
    example.WithValue(map[string]any{"age": 18}),
)

seniorExample := example.New("senior",
    example.WithValue(map[string]any{"age": 75}),
)
```

### Consistent Naming

Use clear, descriptive example names:

```go
// ✗ Bad
example.New("ex1")
example.New("example2")

// ✓ Good
example.New("successful_registration")
example.New("duplicate_email_error")
example.New("invalid_format_error")
```

### Document with Summaries

Add summaries to explain examples:

```go
example.New("admin_user",
    example.WithSummary("User with admin privileges"),
    example.WithDescription("This user has full access to all admin features"),
    example.WithValue(...),
)
```

## Complete Example

```go
package main

import (
    "context"
    "github.com/talav/openapi"
    "github.com/talav/openapi/example"
)

func main() {
    api := openapi.NewAPI(
        openapi.WithInfoTitle("User API"),
        openapi.WithInfoVersion("1.0.0"),
    )
    
    // Success examples
    userCreated := example.New("user_created",
        example.WithSummary("Successfully created user"),
        example.WithValue(map[string]any{
            "id":         123,
            "name":       "Alice Johnson",
            "email":      "alice@company.com",
            "role":       "user",
            "created_at": "2024-03-15T10:30:00Z",
        }),
    )
    
    // Error examples
    invalidEmail := example.New("invalid_email",
        example.WithSummary("Invalid email format"),
        example.WithValue(map[string]any{
            "error":   "validation_error",
            "message": "email must be a valid email address",
            "field":   "email",
        }),
    )
    
    duplicateEmail := example.New("duplicate_email",
        example.WithSummary("Email already registered"),
        example.WithValue(map[string]any{
            "error":   "conflict",
            "message": "A user with this email already exists",
            "field":   "email",
        }),
    )
    
    result, _ := api.Generate(context.Background(),
        openapi.POST("/users",
            openapi.WithRequest(CreateUserRequest{}),
            openapi.WithResponse(201, User{},
                openapi.WithResponseExample(userCreated),
            ),
            openapi.WithResponse(400, Error{},
                openapi.WithResponseExample(invalidEmail),
            ),
            openapi.WithResponse(409, Error{},
                openapi.WithResponseExample(duplicateEmail),
            ),
        ),
    )
}
```

## Next Steps

- [OpenAPI Versions](versions.md) - Example format differences
- [Metadata](metadata.md) - Field-level examples with `openapi` tag
