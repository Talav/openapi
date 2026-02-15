# OpenAPI Versions

This library generates specifications for:

- `3.0.4`
- `3.1.2`

## Choosing a Version

Set the target version when creating the API:

```go
// OpenAPI 3.0.4 (default)
api := openapi.NewAPI(
    openapi.WithVersion("3.0.4"),
)

// OpenAPI 3.1.2
api := openapi.NewAPI(
    openapi.WithVersion("3.1.2"),
)
```

If `WithVersion` is omitted, the default target is `3.0.4`.

## Library Behavior

- The same Go structs and operation definitions are used for both versions.
- The library projects output to the requested target version.
- If a feature cannot be represented in the target version, behavior depends on configuration (degrade with warnings vs strict errors).

## External Specification References

For authoritative version semantics and compatibility details, use the official specs:

- [OpenAPI 3.0.4 Specification](https://spec.openapis.org/oas/v3.0.4)
- [OpenAPI 3.1.0/3.1.x Specification](https://spec.openapis.org/oas/v3.1.0)
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12)

## Next Steps

- [Security](security.md)
- [Advanced Hooks](../advanced/hooks.md)
