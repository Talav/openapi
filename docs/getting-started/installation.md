# Installation

## Install the Library

```bash
go get github.com/talav/openapi
```

## Verify Installation

Create a test file to confirm everything works:

```go
// main.go
package main

import (
    "context"
    "fmt"
    
    "github.com/talav/openapi"
)

func main() {
    api := openapi.NewAPI(
        openapi.WithInfoTitle("Test API"),
        openapi.WithInfoVersion("1.0.0"),
    )
    
    result, err := api.Generate(context.Background())
    if err != nil {
        panic(err)
    }
    
    fmt.Println("✓ OpenAPI library installed successfully")
    fmt.Printf("Generated %d bytes of OpenAPI spec\n", len(result.JSON))
}
```

Run it:

```bash
go run main.go
```

You should see:

```
✓ OpenAPI library installed successfully
Generated 106 bytes of OpenAPI spec
```

## Next Steps

- [Quick Start](quick-start.md) - Build your first API spec
- [Core Concepts](concepts.md) - Understand how it works
