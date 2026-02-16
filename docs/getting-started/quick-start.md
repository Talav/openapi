# Quick Start

Let's build a complete OpenAPI specification for a user management API.

## Define Your Types

Start by defining request and response structures:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/talav/openapi"
)

// Domain models
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ErrorResponse struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

// Request types
type GetUserRequest struct {
    ID int `schema:"id,location=path"`
}

type CreateUserRequest struct {
    Body struct {
        Name  string `json:"name" validate:"required,min=3"`
        Email string `json:"email" validate:"required,email"`
    } `body:"structured"`
}

type UpdateUserRequest struct {
    ID   int `schema:"id,location=path"`
    Body struct {
        Name  string `json:"name"`
        Email string `json:"email" validate:"email"`
    } `body:"structured"`
}
```

## Generate the Specification

Create the API instance and define operations:

```go
func main() {
    // Create API with metadata
    api := openapi.NewAPI(
        openapi.WithInfoTitle("User Management API"),
        openapi.WithInfoVersion("1.0.0"),
        openapi.WithInfoDescription("CRUD operations for user management"),
        openapi.WithServer("https://api.example.com", "Production"),
        openapi.WithServer("http://localhost:8080", "Development"),
    )
    
    // Generate spec from operations
    result, err := api.Generate(context.Background(),
        // List users
        openapi.GET("/users",
            openapi.WithSummary("List all users"),
            openapi.WithResponse(200, []User{}),
        ),
        
        // Get user by ID
        openapi.GET("/users/:id",
            openapi.WithSummary("Get user by ID"),
            openapi.WithRequest(GetUserRequest{}),
            openapi.WithResponse(200, User{}),
            openapi.WithResponse(404, ErrorResponse{}),
        ),
        
        // Create user
        openapi.POST("/users",
            openapi.WithSummary("Create new user"),
            openapi.WithRequest(CreateUserRequest{}),
            openapi.WithResponse(201, User{}),
            openapi.WithResponse(400, ErrorResponse{}),
        ),
        
        // Update user
        openapi.PUT("/users/:id",
            openapi.WithSummary("Update user"),
            openapi.WithRequest(UpdateUserRequest{}),
            openapi.WithResponse(200, User{}),
            openapi.WithResponse(404, ErrorResponse{}),
        ),
        
        // Delete user
        openapi.DELETE("/users/:id",
            openapi.WithSummary("Delete user"),
            openapi.WithRequest(GetUserRequest{}),
            openapi.WithResponse(204, nil),
            openapi.WithResponse(404, ErrorResponse{}),
        ),
    )
    
    if err != nil {
        panic(err)
    }
    
    // Save to file
    fmt.Println(string(result.JSON))
}
```

## Understanding the Output

The generated specification includes everything you'd expect:

- **Path parameters** from `:id` in URLs
- **Request bodies** from `body:"structured"` tags
- **Response schemas** from plain structs or wrapped types
- **Schema constraints** from `validate` tags
- **Component schemas** for reusable types
- **Multiple responses** with status codes

Here's the structure:

```json
{
  "openapi": "3.1.2",
  "info": {
    "title": "User Management API",
    "version": "1.0.0",
    "description": "CRUD operations for user management"
  },
  "servers": [
    {
      "url": "https://api.example.com",
      "description": "Production"
    }
  ],
  "paths": {
    "/users": {
      "get": { ... },
      "post": { ... }
    },
    "/users/{id}": {
      "get": { ... },
      "put": { ... },
      "delete": { ... }
    }
  },
  "components": {
    "schemas": {
      "User": { ... },
      "CreateUserRequestBody": { ... }
    }
  }
}
```

## Adding Security

Protect endpoints with authentication:

```go
api := openapi.NewAPI(
    openapi.WithInfoTitle("User API"),
    openapi.WithInfoVersion("1.0.0"),
    openapi.WithBearerAuth("bearerAuth", "JWT authentication"),
)

result, _ := api.Generate(context.Background(),
    openapi.POST("/users",
        openapi.WithRequest(CreateUserRequest{}),
        openapi.WithResponse(201, User{}),
        openapi.WithSecurity("bearerAuth"), // Requires auth
    ),
)
```

## Next Steps

Now that you've generated your first spec:

- [Core Concepts](concepts.md) - Learn how the library works
- [Tag Reference (talav/schema)](https://talav.github.io/schema/) - Complete `schema`/`body` tag semantics
- [Validation](../guides/validation.md) - Transform validation rules to schemas
