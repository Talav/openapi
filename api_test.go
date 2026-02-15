package openapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// normalizeJSON normalizes JSON by unmarshaling and remarshaling to ensure consistent formatting.
func normalizeJSON(jsonBytes []byte) (string, error) {
	var v any
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return "", err
	}

	normalized, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}

func TestGenerate_SimpleGET(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users", WithResponse(200, GetUsersResponse{})),
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.JSON)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_POST_WithRequestBody(t *testing.T) {
	type CreateUserRequest struct {
		Body struct {
			Name  string `json:"name" validate:"required,min=3"`
			Email string `json:"email" validate:"required,email"`
		} `body:"structured"`
	}

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type CreateUserResponse struct {
		Body User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		POST("/users",
			WithRequest(CreateUserRequest{}),
			WithResponse(201, CreateUserResponse{}),
		),
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "CreateUserRequestBody": {
        "properties": {
          "email": {
            "format": "email",
            "type": "string"
          },
          "name": {
            "minLength": 3,
            "type": "string"
          }
        },
        "required": [
          "name",
          "email"
        ],
        "type": "object"
      },
      "User": {
        "properties": {
          "email": {
            "type": "string"
          },
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateUserRequestBody"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_CRUD_Operations(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type CreateUserRequest struct {
		Body struct {
			Name string `json:"name" validate:"required"`
		} `body:"structured"`
	}

	type UpdateUserRequest struct {
		ID   int `schema:"id,location=path"`
		Body struct {
			Name string `json:"name" validate:"required"`
		} `body:"structured"`
	}

	type GetUserRequest struct {
		ID int `schema:"id,location=path"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	type CreateUserResponse struct {
		Body User `body:"structured"`
	}

	type UpdateUserResponse struct {
		Body User `body:"structured"`
	}

	type DeleteUserRequest struct {
		ID int `schema:"id,location=path"`
	}

	api := NewAPI(
		WithInfoTitle("CRUD API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		POST("/users", WithRequest(CreateUserRequest{}), WithResponse(201, CreateUserResponse{})),
		GET("/users", WithResponse(200, GetUsersResponse{})),
		GET("/users/:id", WithRequest(GetUserRequest{}), WithResponse(200, GetUserResponse{})),
		PUT("/users/:id", WithRequest(UpdateUserRequest{}), WithResponse(200, UpdateUserResponse{})),
		DELETE("/users/:id", WithRequest(DeleteUserRequest{})),
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "CreateUserRequestBody": {
        "properties": {
          "name": {
            "type": "string"
          }
        },
        "required": [
          "name"
        ],
        "type": "object"
      },
      "UpdateUserRequestBody": {
        "properties": {
          "name": {
            "type": "string"
          }
        },
        "required": [
          "name"
        ],
        "type": "object"
      },
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "CRUD API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      },
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateUserRequestBody"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    },
    "/users/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "simple"
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      },
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "simple"
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "simple"
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UpdateUserRequestBody"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithPathParameters(t *testing.T) {
	type GetUserRequest struct {
		ID int `schema:"id,location=path"`
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users/:id",
			WithRequest(GetUserRequest{}),
			WithResponse(200, GetUserResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users/{id}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "simple"
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithQueryParameters(t *testing.T) {
	type ListUsersRequest struct {
		Limit  int    `schema:"limit,location=query"`
		Offset int    `schema:"offset,location=query"`
		Search string `schema:"search,location=query"`
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users",
			WithRequest(ListUsersRequest{}),
			WithResponse(200, GetUsersResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "parameters": [
          {
            "explode": true,
            "in": "query",
            "name": "limit",
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "form"
          },
          {
            "explode": true,
            "in": "query",
            "name": "offset",
            "schema": {
              "format": "int64",
              "type": "integer"
            },
            "style": "form"
          },
          {
            "explode": true,
            "in": "query",
            "name": "search",
            "schema": {
              "type": "string"
            },
            "style": "form"
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_NestedStructs(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type User struct {
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users/:id", WithResponse(200, GetUserResponse{})),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "Address": {
        "properties": {
          "city": {
            "type": "string"
          },
          "street": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "User": {
        "properties": {
          "address": {
            "$ref": "#/components/schemas/Address"
          },
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithValidation(t *testing.T) {
	type CreateUserRequest struct {
		Body struct {
			Name  string `json:"name" validate:"required,min=3,max=50"`
			Email string `json:"email" validate:"required,email"`
			Age   int    `json:"age" validate:"min=0,max=150"`
		} `body:"structured"`
	}

	type CreateUserResponse struct {
		Body struct {
			ID int `json:"id"`
		} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
		WithValidation(true),
	)

	result, err := api.Generate(context.Background(),
		POST("/users", WithRequest(CreateUserRequest{}), WithResponse(201, CreateUserResponse{})),
	)

	// With validation enabled, the spec should be validated against OpenAPI schema
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.JSON)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "CreateUserRequestBody": {
        "properties": {
          "age": {
            "format": "int64",
            "maximum": 150,
            "minimum": 0,
            "type": "integer"
          },
          "email": {
            "format": "email",
            "type": "string"
          },
          "name": {
            "maxLength": 50,
            "minLength": 3,
            "type": "string"
          }
        },
        "required": [
          "name",
          "email"
        ],
        "type": "object"
      },
      "CreateUserResponseBody": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateUserRequestBody"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/CreateUserResponseBody"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_Version304(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.0.4"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users", WithResponse(200, GetUsersResponse{})),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.0.4",
  "paths": {
    "/users": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithServers(t *testing.T) {
	type HealthResponse struct {
		Body struct {
			Status string `json:"status"`
		} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
		WithServer("https://api.example.com", WithServerDescription("Production server")),
		WithServer("https://staging.example.com", WithServerDescription("Staging server")),
	)

	result, err := api.Generate(context.Background(),
		GET("/health", WithResponse(200, HealthResponse{})),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "HealthResponseBody": {
        "properties": {
          "status": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/health": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponseBody"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  },
  "servers": [
    {
      "description": "Production server",
      "url": "https://api.example.com"
    },
    {
      "description": "Staging server",
      "url": "https://staging.example.com"
    }
  ]
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithTags(t *testing.T) {
	type User struct {
		ID int `json:"id"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
		WithTag("users", "User operations"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users", WithTags("users"), WithResponse(200, GetUsersResponse{})),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        },
        "tags": [
          "users"
        ]
      }
    }
  },
  "tags": [
    {
      "description": "User operations",
      "name": "users"
    }
  ]
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_MultipleResponseCodes(t *testing.T) {
	type User struct {
		ID int `json:"id"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	type ErrorResponse struct {
		Body struct {
			Message string `json:"message"`
		} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users/:id",
			WithResponse(200, GetUserResponse{}),
			WithResponse(404, ErrorResponse{}),
			WithResponse(500, ErrorResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "ErrorResponseBody": {
        "properties": {
          "message": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          },
          "404": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponseBody"
                }
              }
            },
            "description": "Not Found"
          },
          "500": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponseBody"
                }
              }
            },
            "description": "Internal Server Error"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithHeaderParameters(t *testing.T) {
	type GetUsersRequest struct {
		APIKey      string `schema:"X-API-Key,location=header"`
		ContentLang string `schema:"Accept-Language,location=header"`
	}

	type User struct {
		ID int `json:"id"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users",
			WithRequest(GetUsersRequest{}),
			WithResponse(200, GetUsersResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "parameters": [
          {
            "in": "header",
            "name": "X-API-Key",
            "schema": {
              "type": "string"
            },
            "style": "simple"
          },
          {
            "in": "header",
            "name": "Accept-Language",
            "schema": {
              "type": "string"
            },
            "style": "simple"
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_WithCookieParameters(t *testing.T) {
	type GetUsersRequest struct {
		SessionID string `schema:"session_id,location=cookie"`
	}

	type User struct {
		ID int `json:"id"`
	}

	type GetUsersResponse struct {
		Body []User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/users",
			WithRequest(GetUsersRequest{}),
			WithResponse(200, GetUsersResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "get": {
        "parameters": [
          {
            "explode": true,
            "in": "cookie",
            "name": "session_id",
            "schema": {
              "type": "string"
            },
            "style": "form"
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/User"
                  },
                  "type": "array"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_MultipartFormData(t *testing.T) {
	type UploadRequest struct {
		Body struct {
			Name string `json:"name"`
			File []byte `json:"file"`
		} `body:"multipart"`
	}

	type UploadResponse struct {
		Body struct {
			ID string `json:"id"`
		} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		POST("/upload",
			WithRequest(UploadRequest{}),
			WithResponse(201, UploadResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "UploadResponseBody": {
        "properties": {
          "id": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/upload": {
      "post": {
        "requestBody": {
          "content": {
            "multipart/form-data": {
              "encoding": {
                "file": {
                  "contentType": "application/octet-stream"
                }
              },
              "schema": {
                "properties": {
                  "file": {
                    "format": "binary",
                    "type": "string"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UploadResponseBody"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_FileUpload(t *testing.T) {
	type FileUploadRequest struct {
		Body []byte `body:"file"`
	}

	type FileUploadResponse struct {
		Body struct {
			FileID string `json:"file_id"`
		} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		POST("/files",
			WithRequest(FileUploadRequest{}),
			WithResponse(201, FileUploadResponse{}),
		),
	)

	require.NoError(t, err)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "FileUploadResponseBody": {
        "properties": {
          "file_id": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/files": {
      "post": {
        "requestBody": {
          "content": {
            "application/octet-stream": {
              "schema": {
                "format": "binary",
                "type": "string"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/FileUploadResponseBody"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_Version312_FullFeatures(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithInfoDescription("Test API description"),
		WithVersion("3.1.2"),
		WithValidation(true),
	)

	result, err := api.Generate(context.Background(),
		GET("/users/:id", WithResponse(200, GetUserResponse{})),
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "description": "Test API description",
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_Version304_WithValidation(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.0.4"),
		WithValidation(true),
	)

	result, err := api.Generate(context.Background(),
		GET("/users/:id", WithResponse(200, GetUserResponse{})),
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.0.4",
  "paths": {
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected, normalized)
}

func TestGenerate_CompareVersions_SameAPI(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type GetUserResponse struct {
		Body User `body:"structured"`
	}

	// Generate with 3.1.2
	api312 := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result312, err := api312.Generate(context.Background(),
		GET("/users/:id", WithResponse(200, GetUserResponse{})),
		POST("/users", WithResponse(201, GetUserResponse{})),
	)

	require.NoError(t, err)

	normalized312, err := normalizeJSON(result312.JSON)
	require.NoError(t, err)

	// Generate with 3.0.4
	api304 := NewAPI(
		WithInfoTitle("Test API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.0.4"),
	)

	result304, err := api304.Generate(context.Background(),
		GET("/users/:id", WithResponse(200, GetUserResponse{})),
		POST("/users", WithResponse(201, GetUserResponse{})),
	)

	require.NoError(t, err)

	normalized304, err := normalizeJSON(result304.JSON)
	require.NoError(t, err)

	expected312 := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {
    "/users": {
      "post": {
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	expected304 := `{
  "components": {
    "schemas": {
      "User": {
        "properties": {
          "id": {
            "format": "int64",
            "type": "integer"
          },
          "name": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "openapi": "3.0.4",
  "paths": {
    "/users": {
      "post": {
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "Created"
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "OK"
          }
        }
      }
    }
  }
}`

	assert.Equal(t, expected312, normalized312)
	assert.Equal(t, expected304, normalized304)
}

func TestGenerate_EmptyAPI(t *testing.T) {
	api := NewAPI(
		WithInfoTitle("Empty API"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)

	normalized, err := normalizeJSON(result.JSON)
	require.NoError(t, err)

	expected := `{
  "components": {},
  "info": {
    "title": "Empty API",
    "version": "1.0.0"
  },
  "openapi": "3.1.2",
  "paths": {}
}`

	assert.Equal(t, expected, normalized)
}
