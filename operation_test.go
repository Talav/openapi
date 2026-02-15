package openapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talav/openapi/example"
)

// getOperation is a test helper that safely navigates spec JSON to retrieve an operation.
func getOperation(t *testing.T, spec map[string]any, method string) map[string]any {
	t.Helper()
	paths, ok := spec["paths"].(map[string]any)
	require.True(t, ok, "paths must exist in spec")
	pathItem, ok := paths["/test"].(map[string]any)
	require.True(t, ok, "path /test must exist")
	op, ok := pathItem[method].(map[string]any)
	require.True(t, ok, "%s operation must exist at /test", method)

	return op
}

func TestOperation_HTTPMethods(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		op     Operation
		method string
		path   string
	}{
		{"GET", GET("/x"), "GET", "/x"},
		{"POST", POST("/x"), "POST", "/x"},
		{"PUT", PUT("/x"), "PUT", "/x"},
		{"PATCH", PATCH("/x"), "PATCH", "/x"},
		{"DELETE", DELETE("/x"), "DELETE", "/x"},
		{"HEAD", HEAD("/x"), "HEAD", "/x"},
		{"OPTIONS", OPTIONS("/x"), "OPTIONS", "/x"},
		{"TRACE", TRACE("/x"), "TRACE", "/x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.method, tt.op.Method)
			assert.Equal(t, tt.path, tt.op.Path)
		})
	}
}

func TestGenerate_HTTPMethods(t *testing.T) {
	type emptyResp struct {
		Body struct{} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/get", WithResponse(200, emptyResp{})),
		POST("/post", WithResponse(201, emptyResp{})),
		PUT("/put", WithResponse(200, emptyResp{})),
		PATCH("/patch", WithResponse(200, emptyResp{})),
		DELETE("/delete", WithResponse(204, emptyResp{})),
		HEAD("/head", WithResponse(200, emptyResp{})),
		OPTIONS("/options", WithResponse(200, emptyResp{})),
		TRACE("/trace", WithResponse(200, emptyResp{})),
	)

	require.NoError(t, err)

	var spec map[string]any
	require.NoError(t, json.Unmarshal(result.JSON, &spec))

	paths, ok := spec["paths"].(map[string]any)
	require.True(t, ok, "paths must exist")

	getPath, ok := paths["/get"].(map[string]any)
	require.True(t, ok, "/get path must exist")
	assert.Contains(t, getPath, "get")

	postPath, ok := paths["/post"].(map[string]any)
	require.True(t, ok, "/post path must exist")
	assert.Contains(t, postPath, "post")

	putPath, ok := paths["/put"].(map[string]any)
	require.True(t, ok, "/put path must exist")
	assert.Contains(t, putPath, "put")

	patchPath, ok := paths["/patch"].(map[string]any)
	require.True(t, ok, "/patch path must exist")
	assert.Contains(t, patchPath, "patch")

	deletePath, ok := paths["/delete"].(map[string]any)
	require.True(t, ok, "/delete path must exist")
	assert.Contains(t, deletePath, "delete")

	headPath, ok := paths["/head"].(map[string]any)
	require.True(t, ok, "/head path must exist")
	assert.Contains(t, headPath, "head")

	optionsPath, ok := paths["/options"].(map[string]any)
	require.True(t, ok, "/options path must exist")
	assert.Contains(t, optionsPath, "options")

	tracePath, ok := paths["/trace"].(map[string]any)
	require.True(t, ok, "/trace path must exist")
	assert.Contains(t, tracePath, "trace")
}

func TestGenerate_OperationMetadata(t *testing.T) {
	type emptyResp struct {
		Body struct{} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
		WithBearerAuth("bearerAuth", "JWT"),
	)

	result, err := api.Generate(context.Background(),
		GET("/test",
			WithSummary("Summary"),
			WithDescription("Description"),
			WithOperationID("testOp"),
			WithDeprecated(),
			WithSecurity("bearerAuth", "read", "write"),
			WithOperationExtension("x-custom", true),
			WithOperationExtension("x-rate-limit", 100),
			WithTags("users", "admin"),
			WithResponse(200, emptyResp{}),
		),
	)

	require.NoError(t, err)

	var spec map[string]any
	require.NoError(t, json.Unmarshal(result.JSON, &spec))

	op := getOperation(t, spec, "get")

	assert.Equal(t, "Summary", op["summary"])
	assert.Equal(t, "Description", op["description"])
	assert.Equal(t, "testOp", op["operationId"])

	deprecated, ok := op["deprecated"].(bool)
	require.True(t, ok, "deprecated must be a bool")
	assert.True(t, deprecated)

	assert.Equal(t, true, op["x-custom"])
	assert.Equal(t, float64(100), op["x-rate-limit"])

	tags, ok := op["tags"].([]any)
	require.True(t, ok, "tags must be an array")
	assert.Equal(t, []any{"users", "admin"}, tags)

	sec, ok := op["security"].([]any)
	require.True(t, ok, "security must be an array")
	require.Len(t, sec, 1)

	secReq, ok := sec[0].(map[string]any)
	require.True(t, ok, "security requirement must be a map")
	assert.Equal(t, []any{"read", "write"}, secReq["bearerAuth"])
}

func TestGenerate_WithOptions(t *testing.T) {
	type emptyResp struct {
		Body struct{} `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/test",
			WithSummary("Original"),
			WithDescription("Original description"),
			WithOptions(
				WithSummary("Overridden"),
				WithOperationID("composedOp"),
			),
			WithResponse(200, emptyResp{}),
		),
	)

	require.NoError(t, err)

	var spec map[string]any
	require.NoError(t, json.Unmarshal(result.JSON, &spec))

	op := getOperation(t, spec, "get")

	assert.Equal(t, "Overridden", op["summary"], "WithOptions should override previous summary")
	assert.Equal(t, "Original description", op["description"], "Description should remain from earlier option")
	assert.Equal(t, "composedOp", op["operationId"])
}

func TestGenerate_RequestExamples(t *testing.T) {
	type Body struct {
		X string `json:"x"`
	}
	type CreateRequest struct {
		Body Body `body:"structured"`
	}
	type Response struct {
		Body Body `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		POST("/test",
			WithRequest(
				CreateRequest{},
				example.New("example1", Body{X: "value1"}),
				example.New("example2", Body{X: "value2"}),
			),
			WithResponse(201, Response{}),
		),
	)

	require.NoError(t, err)

	var spec map[string]any
	require.NoError(t, json.Unmarshal(result.JSON, &spec))

	op := getOperation(t, spec, "post")

	reqBody, ok := op["requestBody"].(map[string]any)
	require.True(t, ok, "requestBody must be a map")

	content, ok := reqBody["content"].(map[string]any)
	require.True(t, ok, "content must be a map")

	jsonContent, ok := content["application/json"].(map[string]any)
	require.True(t, ok, "application/json content must be a map")

	examples, ok := jsonContent["examples"].(map[string]any)
	require.True(t, ok, "examples must be a map")

	assert.Contains(t, examples, "example1")
	assert.Contains(t, examples, "example2")
}

func TestGenerate_ResponseExamples(t *testing.T) {
	type Body struct {
		X string `json:"x"`
	}
	type Response struct {
		Body Body `body:"structured"`
	}

	api := NewAPI(
		WithInfoTitle("Test"),
		WithInfoVersion("1.0.0"),
		WithVersion("3.1.2"),
	)

	result, err := api.Generate(context.Background(),
		GET("/test",
			WithResponse(
				200,
				Response{},
				example.New("success", Body{X: "ok"}),
				example.New("cached", Body{X: "cached_value"}),
			),
		),
	)

	require.NoError(t, err)

	var spec map[string]any
	require.NoError(t, json.Unmarshal(result.JSON, &spec))

	op := getOperation(t, spec, "get")

	responses, ok := op["responses"].(map[string]any)
	require.True(t, ok, "responses must be a map")

	resp, ok := responses["200"].(map[string]any)
	require.True(t, ok, "200 response must be a map")

	content, ok := resp["content"].(map[string]any)
	require.True(t, ok, "content must be a map")

	jsonContent, ok := content["application/json"].(map[string]any)
	require.True(t, ok, "application/json content must be a map")

	examples, ok := jsonContent["examples"].(map[string]any)
	require.True(t, ok, "examples must be a map")

	assert.Contains(t, examples, "success")
	assert.Contains(t, examples, "cached")
}
