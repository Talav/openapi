Automatic OpenAPI 3.0.4 and 3.1.2 specification generation for Go applications.

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

    "rivaas.dev/openapi"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

func main() {
    api := openapi.MustNew(
        openapi.WithTitle("My API", "1.0.0"),
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