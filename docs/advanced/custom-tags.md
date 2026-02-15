# Custom Tags

Customize built-in tag names with `WithTagConfig`.

This page documents the real public extensibility in this package: renaming the built-in tags.  
It does **not** expose a public API for registering arbitrary custom tag parsers.

## What Is Supported

`WithTagConfig` lets you rename these built-in tags:

- `schema`
- `body`
- `openapi`
- `validate`
- `default`
- `requires`

```go
import "github.com/talav/openapi/config"

api := openapi.NewAPI(
    openapi.WithTagConfig(config.TagConfig{
        Schema:   "param",
        Body:     "payload",
        OpenAPI:  "api",
        Validate: "rules",
        Default:  "def",
        Requires: "needs",
    }),
)
```

Then your structs use those names:

```go
type Request struct {
    ID   int  `param:"id,location=path"`
    Body User `payload:"structured"`

    Name string `json:"name" rules:"required" api:"description=User name"`
}
```

## Partial Configuration

You only need to set fields you want to override:

```go
api := openapi.NewAPI(
    openapi.WithTagConfig(config.TagConfig{
        Validate: "rules",
    }),
)
```

Unspecified tag names keep their defaults.

## Default Tag Names

By default:

```go
config.DefaultTagConfig()
// Schema: "schema"
// Body: "body"
// OpenAPI: "openapi"
// Validate: "validate"
// Default: "default"
// Requires: "requires"
```

## When Custom Tags Make Sense

- You already use a conflicting `validate` tag format.
- You need naming consistency with an existing codebase.

Example conflict resolution:

```go
type User struct {
    Email string `validate:"custom_email_rule"` // Existing library usage
}

api := openapi.NewAPI(
    openapi.WithTagConfig(config.TagConfig{
        Validate: "openapi_validate",
    }),
)

type User struct {
    Email string `validate:"custom_email_rule" openapi_validate:"required,email"`
}
```

## What Is Not Publicly Exposed

Registering arbitrary custom tag parsers is not part of the stable public API of this package.

If you need deeper schema customization, use documented hooks:

- [Hooks](hooks.md)
