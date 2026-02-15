# OpenAPI

[![tag](https://img.shields.io/github/tag/talav/openapi.svg)](https://github.com/talav/openapi/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/talav/openapi.svg)](https://pkg.go.dev/github.com/talav/openapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/talav/openapi)](https://goreportcard.com/report/github.com/talav/openapi)
[![CI](https://github.com/talav/openapi/actions/workflows/openapi-ci.yml/badge.svg)](https://github.com/talav/openapi/actions)
[![codecov](https://codecov.io/gh/talav/openapi/graph/badge.svg)](https://codecov.io/gh/talav/openapi)
[![License](https://img.shields.io/github/license/talav/openapi)](./LICENSE)

Automatic OpenAPI 3.0.4 and 3.1.2 specification generation for Go applications.

## Features

- Type-Driven - Define structs, get OpenAPI specs automatically
- OpenAPI 3.0.4 and 3.1.2 - Support for both major versions
- Rich Metadata - Six tag systems for complete control
- Validation Integration - Transform validation rules into schema constraints
- Security Schemes - Built-in support for all OpenAPI auth types
- Examples - Auto-generate or provide custom examples
- Extensible - Hooks for custom schema transformations

## Installation

```bash
go get github.com/talav/openapi
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/talav/openapi"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Body struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required,email"`
    } `body:"structured"`
}

func main() {
    api := openapi.NewAPI(
        openapi.WithInfoTitle("My API"),
        openapi.WithInfoVersion("1.0.0"),
        openapi.WithInfoDescription("API for managing users"),
        openapi.WithServer("http://localhost:8080", "Local development"),
        openapi.WithBearerAuth("bearerAuth", "JWT authentication"),
    )

    result, err := api.Generate(context.Background(),
        openapi.GET("/users/:id",
            openapi.WithSummary("Get user"),
            openapi.WithResponse(200, User{}),
            openapi.WithSecurity("bearerAuth"),
        ),
        openapi.POST("/users",
            openapi.WithSummary("Create user"),
            openapi.WithRequest(CreateUserRequest{}),
            openapi.WithResponse(201, User{}),
        ),
        openapi.DELETE("/users/:id",
            openapi.WithSummary("Delete user"),
            openapi.WithResponse(204, nil),
            openapi.WithSecurity("bearerAuth"),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Check for warnings (optional)
    if len(result.Warnings) > 0 {
        fmt.Printf("Generated with %d warnings\n", len(result.Warnings))
    }

    fmt.Println(string(result.JSON))
}
```

## Documentation

**[Full Documentation →](https://talav.github.io/openapi)**

### Getting Started

- [Installation](https://talav.github.io/openapi/getting-started/installation/) - Setup and verification
- [Quick Start](https://talav.github.io/openapi/getting-started/quick-start/) - Build your first spec
- [Core Concepts](https://talav.github.io/openapi/getting-started/concepts/) - Understand how it works

### Guides

- [Tag Reference (talav/schema)](https://talav.github.io/schema/) - Complete `schema`/`body` tag semantics
- [Validation](https://talav.github.io/openapi/guides/validation/) - Transform validation rules to constraints
- [Metadata](https://talav.github.io/openapi/guides/metadata/) - Add rich OpenAPI metadata
- [Security](https://talav.github.io/openapi/guides/security/) - Configure authentication schemes
- [Examples](https://talav.github.io/openapi/guides/examples/) - Generate realistic examples
- [OpenAPI Versions](https://talav.github.io/openapi/guides/versions/) - Choose between 3.0.4 and 3.1.2

### Advanced

- [Schema Hooks](https://talav.github.io/openapi/advanced/hooks/) - Transform schemas programmatically
- [Custom Tags](https://talav.github.io/openapi/advanced/custom-tags/) - Extend the tag system

### API Reference

- [pkg.go.dev](https://pkg.go.dev/github.com/talav/openapi) - Complete API documentation

## Key Concepts

### Struct Tags

The library uses six struct tags:

```go
type Request struct {
    ID      int    `schema:"id,location=path"`        // Parameter metadata
    APIKey  string `schema:"X-API-Key,location=header"` // Header parameter
    
    Body struct {
        Name  string `json:"name" validate:"required,min=3"`  // Validation
        Email string `json:"email" validate:"required,email" openapi:"description=Contact email,examples=user@example.com"` // Metadata
        Age   int    `json:"age" default:"18"` // Default value
    } `body:"structured"` // Request body
}
```

- **`json`** - Property names
- **`schema`** - Parameter location and style
- **`body`** - Request/response body
- **`validate`** - Validation rules → schema constraints
- **`openapi`** - OpenAPI metadata (descriptions, examples)
- **`default`** - Default values
- **`requires`** - Conditional required fields

### Tag Semantics

OpenAPI consumes the same `schema` and `body` tag semantics as [talav/schema](https://talav.github.io/schema/).  
Use the schema docs as the canonical reference for parameter locations, styles, body types, and serialization behavior.

## Testing

```bash
# Run tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test -race ./...`
5. Run linter: `golangci-lint run`
6. Commit with clear messages
7. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.
