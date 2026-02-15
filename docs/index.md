# OpenAPI - Automatic Specification Generation

Generate OpenAPI 3.0.4 and 3.1.2 specifications from Go structs.

## Features

- **Type-Driven** - Define request/response structures, get specs automatically
- **Multiple Versions** - Support for OpenAPI 3.0.4 and 3.1.2
- **Rich Metadata** - Six tag systems give you complete control
- **Validation Integration** - Validation rules become schema constraints
- **Security Schemes** - Built-in support for all OpenAPI auth types
- **Examples** - Auto-generated or custom examples
- **Extensible** - Hooks for custom schema transformations

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/talav/openapi"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

type CreateUserRequest struct {
    Body User `body:"structured"`
}

func main() {
    api := openapi.NewAPI(
        openapi.WithInfoTitle("User API"),
        openapi.WithInfoVersion("1.0.0"),
    )
    
    result, _ := api.Generate(context.Background(),
        openapi.POST("/users",
            openapi.WithRequest(CreateUserRequest{}),
            openapi.WithResponse(201, User{}),
        ),
    )
    
    fmt.Println(string(result.JSON))
}
```

## Documentation

- **[Getting Started](getting-started/installation.md)** - Installation and first steps
- **[Guides](guides/validation.md)** - Feature guides and examples
- **[Advanced](advanced/hooks.md)** - Hooks, custom tags, debugging

## Why OpenAPI Generation?

Hand-writing OpenAPI specs is tedious and error-prone. This library:

1. **Eliminates duplication** - Your Go structs are the source of truth
2. **Catches errors early** - Type safety at compile time
3. **Keeps code and docs in sync** - They can't drift apart
4. **Reduces maintenance** - Update structs, regenerate specs

## Tag System

The library extracts metadata from struct tags:

- `json` - Property names
- `schema` - Parameter names, locations, and styles
- `body` - Request/response body marker
- `validate` - Validation rules become constraints
- `openapi` - Descriptions, examples, and metadata
- `default` - Default values
- `requires` - Conditional required fields

Learn more in [Core Concepts](getting-started/concepts.md).

## Installation

```bash
go get github.com/talav/openapi
```

Continue to [Quick Start](getting-started/quick-start.md) →
