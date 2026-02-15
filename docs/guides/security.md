# Security

Configure OpenAPI security schemes for authentication and authorization.

## Overview

The library supports all OpenAPI 3.0/3.1 security scheme types:

- HTTP Authentication (Basic, Bearer, etc.)
- API Keys (Header, Query, Cookie)
- OAuth 2.0
- OpenID Connect

## HTTP Authentication

### Bearer Authentication (JWT)

Most common for modern APIs:

```go
api := openapi.NewAPI(
    openapi.WithInfoTitle("My API"),
    openapi.WithInfoVersion("1.0.0"),
    openapi.WithBearerAuth("bearerAuth", "JWT authentication"),
)

result, _ := api.Generate(context.Background(),
    openapi.POST("/users",
        openapi.WithRequest(CreateUserRequest{}),
        openapi.WithResponse(201, User{}),
        openapi.WithSecurity("bearerAuth"), // Requires auth
    ),
    
    // Public endpoint (no auth)
    openapi.GET("/public/info",
        openapi.WithResponse(200, Info{}),
    ),
)
```

Generates:

```json
{
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "JWT authentication"
      }
    }
  }
}
```

### Basic Authentication

```go
api.AddSecurityScheme("basicAuth", openapi.SecurityScheme{
    Type:        "http",
    Scheme:      "basic",
    Description: "Basic HTTP authentication",
})
```

### Custom HTTP Schemes

```go
api.AddSecurityScheme("digestAuth", openapi.SecurityScheme{
    Type:        "http",
    Scheme:      "digest",
    Description: "Digest authentication",
})
```

## API Key Authentication

### Header-Based API Keys

```go
api.AddSecurityScheme("apiKey", openapi.SecurityScheme{
    Type:        "apiKey",
    In:          "header",
    Name:        "X-API-Key",
    Description: "API key for external integrations",
})
```

Usage:

```
GET /users
X-API-Key: your-api-key-here
```

### Query Parameter API Keys

```go
api.AddSecurityScheme("apiKeyQuery", openapi.SecurityScheme{
    Type:        "apiKey",
    In:          "query",
    Name:        "api_key",
    Description: "API key in query string",
})
```

Usage:

```
GET /users?api_key=your-api-key-here
```

### Cookie-Based API Keys

```go
api.AddSecurityScheme("cookieAuth", openapi.SecurityScheme{
    Type:        "apiKey",
    In:          "cookie",
    Name:        "session_id",
    Description: "Session cookie",
})
```

## OAuth 2.0

### Authorization Code Flow

For web applications with backend:

```go
api.WithOAuth2Auth("oauth2", "OAuth 2.0 authentication",
    openapi.WithAuthorizationCodeFlow(
        "https://auth.example.com/oauth/authorize",
        "https://auth.example.com/oauth/token",
        map[string]string{
            "read:users":  "Read user data",
            "write:users": "Modify user data",
            "admin":       "Administrative access",
        },
    ),
)
```

### Implicit Flow

For single-page applications (deprecated in OAuth 2.1):

```go
api.WithOAuth2Auth("oauth2Implicit", "OAuth 2.0 Implicit",
    openapi.WithImplicitFlow(
        "https://auth.example.com/oauth/authorize",
        map[string]string{
            "read": "Read access",
        },
    ),
)
```

### Client Credentials Flow

For server-to-server authentication:

```go
api.WithOAuth2Auth("oauth2Client", "OAuth 2.0 Client Credentials",
    openapi.WithClientCredentialsFlow(
        "https://auth.example.com/oauth/token",
        map[string]string{
            "api:access": "API access",
        },
    ),
)
```

### Password Flow

For trusted applications (deprecated in OAuth 2.1):

```go
api.WithOAuth2Auth("oauth2Password", "OAuth 2.0 Password",
    openapi.WithPasswordFlow(
        "https://auth.example.com/oauth/token",
        map[string]string{
            "read": "Read access",
            "write": "Write access",
        },
    ),
)
```

### Multiple Flows

Combine flows for different client types:

```go
api.AddSecurityScheme("oauth2Multi", openapi.SecurityScheme{
    Type:        "oauth2",
    Description: "OAuth 2.0",
    Flows: &openapi.OAuthFlows{
        AuthorizationCode: &openapi.OAuthFlow{
            AuthorizationURL: "https://auth.example.com/authorize",
            TokenURL:         "https://auth.example.com/token",
            Scopes: map[string]string{
                "read":  "Read access",
                "write": "Write access",
            },
        },
        ClientCredentials: &openapi.OAuthFlow{
            TokenURL: "https://auth.example.com/token",
            Scopes: map[string]string{
                "api": "API access",
            },
        },
    },
})
```

## OpenID Connect

```go
api.WithOpenIDConnectAuth("oidc", "OpenID Connect",
    "https://auth.example.com/.well-known/openid-configuration",
)
```

## Applying Security to Operations

### Single Security Scheme

```go
openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{}),
    openapi.WithSecurity("bearerAuth"),
)
```

### Multiple Security Schemes (AND)

Requires ALL schemes:

```go
openapi.POST("/admin/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{}),
    openapi.WithSecurity("bearerAuth", "apiKey"), // Both required
)
```

### Alternative Security Schemes (OR)

Accepts ANY scheme:

```go
// Method 1: Multiple WithSecurity calls
openapi.GET("/users",
    openapi.WithResponse(200, []User{}),
    openapi.WithSecurity("bearerAuth"),
    openapi.WithSecurity("apiKey"), // Alternative
)

// Method 2: Pass empty strings to separate alternatives
openapi.GET("/users",
    openapi.WithResponse(200, []User{}),
    openapi.WithSecurity("bearerAuth"),
    openapi.WithSecurity(""), // OR separator
    openapi.WithSecurity("apiKey"),
)
```

## OAuth Scopes

Specify required OAuth scopes per operation:

```go
api.WithOAuth2Auth("oauth2", "OAuth 2.0",
    openapi.WithAuthorizationCodeFlow(
        "https://auth.example.com/authorize",
        "https://auth.example.com/token",
        map[string]string{
            "users:read":  "Read users",
            "users:write": "Write users",
            "admin":       "Admin access",
        },
    ),
)

// Read operation - requires read scope
openapi.GET("/users",
    openapi.WithResponse(200, []User{}),
    openapi.WithOAuthScopes("oauth2", "users:read"),
)

// Write operation - requires write scope
openapi.POST("/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{}),
    openapi.WithOAuthScopes("oauth2", "users:write"),
)

// Admin operation - requires both scopes
openapi.DELETE("/users/:id",
    openapi.WithRequest(DeleteUserRequest{}),
    openapi.WithResponse(204, nil),
    openapi.WithOAuthScopes("oauth2", "users:write", "admin"),
)
```

## Global Security

Apply security to all operations by default:

```go
api := openapi.NewAPI(
    openapi.WithInfoTitle("My API"),
    openapi.WithInfoVersion("1.0.0"),
    openapi.WithBearerAuth("bearerAuth", "JWT auth"),
    openapi.WithGlobalSecurity("bearerAuth"), // All operations require auth
)

result, _ := api.Generate(context.Background(),
    // Inherits global security
    openapi.GET("/users",
        openapi.WithResponse(200, []User{}),
    ),
    
    // Override to make public
    openapi.GET("/public/status",
        openapi.WithResponse(200, Status{}),
        openapi.WithoutSecurity(), // No auth required
    ),
)
```

## Common Patterns

### Multi-Tenant API Keys

```go
api.AddSecurityScheme("tenantAuth", openapi.SecurityScheme{
    Type:        "apiKey",
    In:          "header",
    Name:        "X-Tenant-ID",
    Description: "Tenant identifier",
})

api.AddSecurityScheme("apiKey", openapi.SecurityScheme{
    Type:        "apiKey",
    In:          "header",
    Name:        "X-API-Key",
    Description: "API key",
})

// Requires both tenant ID and API key
openapi.GET("/data",
    openapi.WithResponse(200, Data{}),
    openapi.WithSecurity("tenantAuth", "apiKey"),
)
```

### Mixed Public and Private Endpoints

```go
api := openapi.NewAPI(
    openapi.WithBearerAuth("bearerAuth", "JWT auth"),
)

result, _ := api.Generate(context.Background(),
    // Public endpoints
    openapi.GET("/public/products",
        openapi.WithResponse(200, []Product{}),
    ),
    openapi.GET("/public/pricing",
        openapi.WithResponse(200, Pricing{}),
    ),
    
    // Protected endpoints
    openapi.POST("/users",
        openapi.WithRequest(CreateUserRequest{}),
        openapi.WithResponse(201, User{}),
        openapi.WithSecurity("bearerAuth"),
    ),
    openapi.GET("/users/me",
        openapi.WithResponse(200, User{}),
        openapi.WithSecurity("bearerAuth"),
    ),
)
```

### Admin Endpoints with Extra Security

```go
api.WithBearerAuth("userAuth", "User JWT")
api.WithAPIKeyAuth("adminKey", "X-Admin-Key", "header", "Admin API key")

// Regular user endpoints
openapi.GET("/profile",
    openapi.WithResponse(200, Profile{}),
    openapi.WithSecurity("userAuth"),
)

// Admin endpoints require BOTH auth methods
openapi.POST("/admin/users",
    openapi.WithRequest(CreateUserRequest{}),
    openapi.WithResponse(201, User{}),
    openapi.WithSecurity("userAuth", "adminKey"),
)
```

## Next Steps

- [Examples](examples.md) - Add example requests with auth headers
- [Versions](versions.md) - Security differences between 3.0 and 3.1
