package v312

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talav/openapi/debug"
	"github.com/talav/openapi/internal/model"
)

func normalizeJSON(jsonStr string) string {
	var m any
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal JSON in normalizeJSON: %v", err))
	}
	normalized, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON in normalizeJSON: %v", err))
	}

	return string(normalized)
}

func TestView_ComprehensiveSpec(t *testing.T) {
	spec := createComprehensiveSpec()

	adapter := &AdapterV312{}
	result, warnings, err := adapter.View(spec)

	// Check errors
	require.NoError(t, err)

	// Check warnings - 3.1.2 should NOT warn about 3.1-only features
	assert.Empty(t, warnings, "3.1.2 adapter should not generate warnings for 3.1 features")

	// Create JSON from result and compare with expected
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)

	// Expected JSON for comprehensive spec (OpenAPI 3.1.2 format)
	expectedJSON := `{
  "openapi": "3.1.2",
  "info": {
    "title": "Test API",
    "summary": "This is a summary (3.1-only feature)",
    "description": "A test API",
    "license": {
      "name": "MIT",
      "identifier": "MIT"
    },
    "version": "1.0.0",
    "x-custom-info": "custom info extension",
    "x-api-version": "1.0.0-beta"
  },
  "tags": [
    {
      "name": "Users",
      "description": "User management operations"
    }
  ],
  "paths": {
    "/users": {
      "get": {
        "summary": "Get users",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "Maximum number of users to return",
            "schema": {
              "type": "integer",
              "default": 10
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/User"
                  }
                },
                "example": [
                  {
                    "id": "1",
                    "name": "John"
                  }
                ],
                "examples": {
                  "single": {
                    "summary": "Single user example",
                    "description": "An example of a single user",
                    "value": {
                      "id": "1",
                      "name": "John"
                    }
                  },
                  "external": {
                    "summary": "External example",
                    "description": "Example loaded from external source",
                    "externalValue": "https://api.example.com/examples/user.json"
                  }
                }
              }
            }
          }
        },
        "x-operation-type": "list",
        "x-cache-ttl": 300
      },
      "post": {
        "summary": "Create user",
        "requestBody": {
          "description": "User data",
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/User"
              },
              "example": {
                "id": "123",
                "name": "New User"
              },
              "examples": {
                "newUser": {
                  "summary": "New user creation",
                  "description": "Example of creating a new user",
                  "value": {
                    "id": "123",
                    "name": "New User"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                },
                "example": {
                  "id": "123",
                  "name": "Created User"
                }
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "servers": [
          {
            "url": "https://api.example.com"
          }
        ]
      }
    },
    "/users/{userId}": {
      "get": {
        "summary": "Get user by ID",
        "parameters": [
          {
            "name": "userId",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "examples": {
              "example1": {
                "value": "123"
              },
              "example2": {
                "value": "456"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "User found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                },
                "example": {
                  "id": "123",
                  "name": "Retrieved User"
                }
              }
            },
            "headers": {
              "X-Rate-Limit": {
                "$ref": "#/components/headers/X-Rate-Limit"
              }
            }
          },
          "404": {
            "description": "User not found"
          }
        }
      }
    }
  },
  "webhooks": {
    "userCreated": {
      "post": {
        "summary": "User created webhook"
      }
    }
  },
  "components": {
    "schemas": {
      "StatusConst": {
        "title": "Status Constant",
        "const": "active"
      },
      "User": {
        "type": "object",
        "title": "User Schema",
        "properties": {
          "id": {
            "type": "string",
            "description": "Unique user identifier"
          },
          "name": {
            "type": "string",
            "description": "User name"
          }
        },
        "required": [
          "id",
          "name"
        ],
        "examples": [
          {
            "id": "1",
            "name": "Example 1"
          },
          {
            "id": "2",
            "name": "Example 2"
          }
        ],
        "contentEncoding": "gzip",
        "contentMediaType": "application/json",
        "unevaluatedProperties": {
          "type": "string"
        }
      }
    },
    "parameters": {
      "limit": {
        "name": "limit",
        "in": "query",
        "description": "Maximum number of users to return",
        "schema": {
          "type": "integer",
          "default": 10
        }
      }
    },
    "responses": {
      "UserResponse": {
        "description": "User response",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/User"
            }
          }
        }
      },
      "NotFound": {
        "description": "Resource not found"
      }
    },
    "examples": {
      "userExample": {
        "summary": "User example",
        "description": "An example user object",
        "value": {
          "id": "123",
          "name": "John Doe"
        }
      }
    },
    "requestBodies": {
      "UserRequest": {
        "description": "User request body",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/User"
            }
          }
        }
      }
    },
    "headers": {
      "X-Custom-Header": {
        "description": "Custom header with all properties",
        "required": true,
        "deprecated": true,
        "style": "simple",
        "schema": {
          "type": "string"
        },
        "example": "custom-value"
      },
      "X-Rate-Limit": {
        "description": "Rate limit information",
        "schema": {
          "type": "string"
        },
        "example": "100/hour",
        "examples": {
          "normal": {
            "summary": "Normal rate limit",
            "value": "100/hour"
          },
          "throttled": {
            "summary": "Throttled rate limit",
            "value": "10/hour"
          }
        },
        "style": "simple"
      }
    },
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "description": "JWT Bearer token",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "links": {
      "userOrders": {
        "operationId": "getUserOrders",
        "description": "Link to user orders",
        "parameters": {
          "userId": "$response.body#/id"
        },
        "server": {
          "url": "https://api.example.com"
        }
      }
    },
    "callbacks": {
      "userCreated": {
        "{$request.body#/callbackUrl}": {
          "post": {
            "summary": "User created callback",
            "description": "Called when a user is successfully created",
            "requestBody": {
              "description": "Callback payload",
              "content": {
                "application/json": {
                  "schema": {
                    "$ref": "#/components/schemas/User"
                  }
                }
              }
            },
            "responses": {
              "200": {
                "description": "Callback processed successfully"
              }
            }
          }
        }
      }
    },
    "pathItems": {
      "common": {
        "parameters": [
          {
            "name": "apiKey",
            "in": "header",
            "schema": {
              "type": "string"
            }
          }
        ]
      }
    }
  }
}`

	// Compare actual JSON with expected (ignoring whitespace differences)
	actualNormalized := normalizeJSON(string(jsonBytes))
	expectedNormalized := normalizeJSON(expectedJSON)
	assert.Equal(t, expectedNormalized, actualNormalized, "Generated JSON does not match expected")
}

func TestView_NilSpec(t *testing.T) {
	adapter := &AdapterV312{}
	result, warnings, err := adapter.View(nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, warnings)
	assert.Contains(t, err.Error(), "nil spec")
}

func TestView_EmptySpec(t *testing.T) {
	// Create minimal valid spec with only Info containing Title and Version
	spec := &model.Spec{
		Info: model.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	// Call View method
	adapter := &AdapterV312{}
	result, warnings, err := adapter.View(spec)

	// Verify no error
	require.NoError(t, err)

	// Verify no warnings
	assert.Empty(t, warnings)

	// Verify result is not nil
	require.NotNil(t, result)

	// Generate JSON for visual inspection
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)

	// Expected JSON for empty spec
	expectedJSON := `{
  "openapi": "3.1.2",
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "paths": {}
}`

	// Compare actual JSON with expected
	actualNormalized := normalizeJSON(string(jsonBytes))
	expectedNormalized := normalizeJSON(expectedJSON)
	assert.Equal(t, expectedNormalized, actualNormalized, "Generated JSON does not match expected")
}

func TestTransformSchema_RefCases(t *testing.T) {
	adapter := &AdapterV312{}

	// Test nil schema
	result := adapter.transformSchema(nil, nil)
	assert.Nil(t, result)

	// Test schema with only ref
	schema := &model.Schema{Ref: "#/components/schemas/User"}
	result = adapter.transformSchema(schema, nil)
	require.NotNil(t, result)
	assert.Equal(t, "#/components/schemas/User", result.Ref)
	// Other fields should not be set
	assert.Nil(t, result.Type)
	assert.Equal(t, "", result.Title)
}

func TestTransformSchema_NoWarnings(t *testing.T) {
	adapter := &AdapterV312{}

	tests := []struct {
		name   string
		schema *model.Schema
	}{
		{
			name: "multiple examples",
			schema: &model.Schema{
				Type:     "string",
				Examples: []any{"example1", "example2"},
			},
		},
		{
			name: "const value",
			schema: &model.Schema{
				Type:  "string",
				Const: "constant-value",
			},
		},
		{
			name: "content encoding",
			schema: &model.Schema{
				Type:            "string",
				ContentEncoding: "base64",
			},
		},
		{
			name: "content media type",
			schema: &model.Schema{
				Type:             "string",
				ContentMediaType: "application/json",
			},
		},
		{
			name: "unevaluated properties",
			schema: &model.Schema{
				Type: "object",
				Unevaluated: &model.Schema{
					Type: "string",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var warnings debug.Warnings
			result := adapter.transformSchema(tt.schema, &warnings)

			require.NotNil(t, result)
			assert.Empty(t, warnings, "3.1.2 should not generate warnings for 3.1 features")
		})
	}
}

func TestTransformPathItem_RefCase(t *testing.T) {
	adapter := &AdapterV312{}

	// Test nil path item
	result := adapter.transformPathItem(nil, nil)
	assert.Nil(t, result)

	// Test path item with only ref
	pathItem := &model.PathItem{Ref: "#/paths/users"}
	result = adapter.transformPathItem(pathItem, nil)
	require.NotNil(t, result)
	assert.Equal(t, "#/paths/users", result.Ref)
	// Other fields should not be processed
	assert.Nil(t, result.Get)
	assert.Nil(t, result.Post)
}

func TestTransformComponents_NilAndEmpty(t *testing.T) {
	adapter := &AdapterV312{}

	// Test nil components
	result := adapter.transformComponents(nil, nil)
	assert.Nil(t, result)

	// Test empty components
	emptyComponents := &model.Components{}
	var warnings debug.Warnings
	result = adapter.transformComponents(emptyComponents, &warnings)

	require.NotNil(t, result)
	assert.Nil(t, result.Schemas)
	assert.Nil(t, result.Responses)
	assert.Nil(t, result.Parameters)
	assert.Nil(t, result.Examples)
	assert.Nil(t, result.RequestBodies)
	assert.Nil(t, result.Headers)
	assert.Nil(t, result.SecuritySchemes)
	assert.Nil(t, result.Links)
	assert.Nil(t, result.Callbacks)
	assert.Nil(t, result.PathItems)
	assert.Empty(t, warnings)
}

// Helper function to create a comprehensive test spec.
func createComprehensiveSpec() *model.Spec {
	return &model.Spec{
		Info: model.Info{
			Title:       "Test API",
			Description: "A test API",
			Version:     "1.0.0",
			Summary:     "This is a summary (3.1-only feature)",
			License: &model.License{
				Name:       "MIT",
				Identifier: "MIT",
			},
			Extensions: map[string]any{
				"x-custom-info": "custom info extension",
				"x-api-version": "1.0.0-beta",
			},
		},
		Tags: []model.Tag{
			{
				Name:        "Users",
				Description: "User management operations",
			},
		},
		Paths: map[string]*model.PathItem{
			"/users": {
				Get: &model.Operation{
					Summary: "Get users",
					Parameters: []model.Parameter{
						{
							Name:        "limit",
							In:          "query",
							Schema:      &model.Schema{Type: "integer", Default: 10},
							Description: "Maximum number of users to return",
						},
					},
					Responses: map[string]*model.Response{
						"200": {
							Description: "Success",
							Content: map[string]*model.MediaType{
								"application/json": {
									Schema:  &model.Schema{Type: "array", Items: &model.Schema{Ref: "#/components/schemas/User"}},
									Example: []map[string]any{{"id": "1", "name": "John"}},
									Examples: map[string]*model.Example{
										"single": {
											Summary:     "Single user example",
											Description: "An example of a single user",
											Value:       map[string]any{"id": "1", "name": "John"},
										},
										"external": {
											Summary:       "External example",
											Description:   "Example loaded from external source",
											ExternalValue: "https://api.example.com/examples/user.json",
										},
									},
								},
							},
						},
					},
					Extensions: map[string]any{
						"x-operation-type": "list",
						"x-cache-ttl":      300,
					},
				},
				Post: &model.Operation{
					Summary: "Create user",
					RequestBody: &model.RequestBody{
						Description: "User data",
						Content: map[string]*model.MediaType{
							"application/json": {
								Schema:  &model.Schema{Ref: "#/components/schemas/User"},
								Example: map[string]any{"id": "123", "name": "New User"},
								Examples: map[string]*model.Example{
									"newUser": {
										Summary:     "New user creation",
										Description: "Example of creating a new user",
										Value:       map[string]any{"id": "123", "name": "New User"},
									},
								},
							},
						},
					},
					Responses: map[string]*model.Response{
						"201": {
							Description: "User created successfully",
							Content: map[string]*model.MediaType{
								"application/json": {
									Schema:  &model.Schema{Ref: "#/components/schemas/User"},
									Example: map[string]any{"id": "123", "name": "Created User"},
								},
							},
						},
					},
					Security: []model.SecurityRequirement{
						{"bearerAuth": []string{}},
					},
					Servers: []model.Server{
						{URL: "https://api.example.com"},
					},
				},
			},
			"/users/{userId}": {
				Get: &model.Operation{
					Summary: "Get user by ID",
					Parameters: []model.Parameter{
						{
							Name:     "userId",
							In:       "path",
							Required: true,
							Schema:   &model.Schema{Type: "string"},
							Examples: map[string]*model.Example{
								"example1": {Value: "123"},
								"example2": {Value: "456"},
							},
						},
					},
					Responses: map[string]*model.Response{
						"200": {
							Description: "User found",
							Content: map[string]*model.MediaType{
								"application/json": {
									Schema:  &model.Schema{Ref: "#/components/schemas/User"},
									Example: map[string]any{"id": "123", "name": "Retrieved User"},
								},
							},
							Headers: map[string]*model.Header{
								"X-Rate-Limit": {
									Ref: "#/components/headers/X-Rate-Limit",
								},
							},
						},
						"404": {
							Description: "User not found",
						},
					},
				},
			},
		},
		Components: &model.Components{
			Schemas: map[string]*model.Schema{
				"User":        createComplexSchema(),
				"StatusConst": createConstSchema(),
			},
			Parameters: map[string]*model.Parameter{
				"limit": {
					Name:        "limit",
					In:          "query",
					Description: "Maximum number of users to return",
					Schema:      &model.Schema{Type: "integer", Default: 10},
				},
			},
			Responses: map[string]*model.Response{
				"UserResponse": {
					Description: "User response",
					Content: map[string]*model.MediaType{
						"application/json": {
							Schema: &model.Schema{Ref: "#/components/schemas/User"},
						},
					},
				},
				"NotFound": {
					Description: "Resource not found",
				},
			},
			Examples: map[string]*model.Example{
				"userExample": {
					Summary:     "User example",
					Description: "An example user object",
					Value: map[string]any{
						"id":   "123",
						"name": "John Doe",
					},
				},
			},
			RequestBodies: map[string]*model.RequestBody{
				"UserRequest": {
					Description: "User request body",
					Content: map[string]*model.MediaType{
						"application/json": {
							Schema: &model.Schema{Ref: "#/components/schemas/User"},
						},
					},
				},
			},
			Headers: map[string]*model.Header{
				"X-Rate-Limit": {
					Description:     "Rate limit information",
					Schema:          &model.Schema{Type: "string"},
					Deprecated:      false,
					AllowEmptyValue: false,
					Style:           "simple",
					Explode:         false,
					Example:         "100/hour",
					Examples: map[string]*model.Example{
						"normal": {
							Summary: "Normal rate limit",
							Value:   "100/hour",
						},
						"throttled": {
							Summary: "Throttled rate limit",
							Value:   "10/hour",
						},
					},
				},
				"X-Custom-Header": {
					Description: "Custom header with all properties",
					Schema:      &model.Schema{Type: "string"},
					Required:    true,
					Deprecated:  true,
					Style:       "simple",
					Explode:     false,
					Example:     "custom-value",
				},
			},
			SecuritySchemes: map[string]*model.SecurityScheme{
				"bearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
					Description:  "JWT Bearer token",
				},
			},
			Links: map[string]*model.Link{
				"userOrders": {
					OperationID: "getUserOrders",
					Parameters:  map[string]any{"userId": "$response.body#/id"},
					Description: "Link to user orders",
					Server: &model.Server{
						URL: "https://api.example.com",
					},
				},
			},
			Callbacks: map[string]*model.Callback{
				"userCreated": {
					PathItems: map[string]*model.PathItem{
						"{$request.body#/callbackUrl}": {
							Post: &model.Operation{
								Summary:     "User created callback",
								Description: "Called when a user is successfully created",
								RequestBody: &model.RequestBody{
									Description: "Callback payload",
									Content: map[string]*model.MediaType{
										"application/json": {
											Schema: &model.Schema{Ref: "#/components/schemas/User"},
										},
									},
								},
								Responses: map[string]*model.Response{
									"200": {
										Description: "Callback processed successfully",
									},
								},
							},
						},
					},
				},
			},
			PathItems: map[string]*model.PathItem{
				"common": {
					Parameters: []model.Parameter{
						{Name: "apiKey", In: "header", Schema: &model.Schema{Type: "string"}},
					},
				},
			},
		},
		Webhooks: map[string]*model.PathItem{
			"userCreated": {
				Post: &model.Operation{
					Summary: "User created webhook",
				},
			},
		},
	}
}

func createComplexSchema() *model.Schema {
	return &model.Schema{
		Type:  "object",
		Title: "User Schema",
		Properties: map[string]*model.Schema{
			"id": {
				Type:        "string",
				Description: "Unique user identifier",
			},
			"name": {
				Type:        "string",
				Description: "User name",
			},
		},
		Required: []string{"id", "name"},
		// 3.1 features - should be supported without warnings
		Examples: []any{
			map[string]any{"id": "1", "name": "Example 1"},
			map[string]any{"id": "2", "name": "Example 2"},
		},
		ContentEncoding:  "gzip",
		ContentMediaType: "application/json",
		Unevaluated:      &model.Schema{Type: "string"},
	}
}

func createConstSchema() *model.Schema {
	return &model.Schema{
		Title: "Status Constant",
		Const: "active",
	}
}
