package apigateway2openapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swaggest/openapi-go/openapi3"
)

func TestSpecHandler_processPathParameterPart(t *testing.T) {
	c := require.New(t)

	// Create a basic spec with a path
	spec := &openapi3.Spec{}
	spec.Paths.WithMapOfPathItemValuesItem("/pets/{petId}", openapi3.PathItem{})

	handler := SpecHandler{
		spec: spec,
	}

	// Create a path parameter documentation part
	docPart := DocumentationPart{
		Location: Location{
			Type: "PATH_PARAMETER",
			Path: "/pets/{petId}",
			Name: "petId",
		},
		Properties: json.RawMessage(`{"description": "The id of the pet to retrieve", "required": true, "schema": {"type": "string"}}`),
	}

	err := handler.processPathParameterPart(docPart)
	c.NoError(err)

	// Verify the parameter was added
	pathItem := spec.Paths.MapOfPathItemValues["/pets/{petId}"]
	c.Len(pathItem.Parameters, 1)
	c.Equal("petId", pathItem.Parameters[0].Parameter.Name)
	c.Equal("path", pathItem.Parameters[0].Parameter.In)
	c.NotNil(pathItem.Parameters[0].Parameter.Required)
	c.True(*pathItem.Parameters[0].Parameter.Required)
}

func TestSpecHandler_processQueryParameterPart(t *testing.T) {
	c := require.New(t)

	// Create a basic spec with a path and operation
	spec := &openapi3.Spec{}
	operation := openapi3.Operation{}
	pathItem := openapi3.PathItem{}
	pathItem.WithMapOfOperationValuesItem("get", operation)
	spec.Paths.WithMapOfPathItemValuesItem("/pets", pathItem)

	handler := SpecHandler{
		spec: spec,
	}

	// Create a query parameter documentation part
	docPart := DocumentationPart{
		Location: Location{
			Type:   "QUERY_PARAMETER",
			Path:   "/pets",
			Method: "GET",
			Name:   "limit",
		},
		Properties: json.RawMessage(`{"description": "Maximum number of pets to return", "schema": {"type": "integer"}}`),
	}

	err := handler.processQueryParameterPart(docPart)
	c.NoError(err)

	// Verify the parameter was added
	pathItem = spec.Paths.MapOfPathItemValues["/pets"]
	operation = pathItem.MapOfOperationValues["get"]
	c.Len(operation.Parameters, 1)
	c.Equal("limit", operation.Parameters[0].Parameter.Name)
	c.Equal("query", operation.Parameters[0].Parameter.In)
}

func TestSpecHandler_processMethodPart(t *testing.T) {
	c := require.New(t)

	// Create a basic spec with a path and operation
	spec := &openapi3.Spec{}
	operation := openapi3.Operation{}
	pathItem := openapi3.PathItem{}
	pathItem.WithMapOfOperationValuesItem("get", operation)
	spec.Paths.WithMapOfPathItemValuesItem("/pets", pathItem)

	handler := SpecHandler{
		spec: spec,
	}

	// Create a method documentation part
	docPart := DocumentationPart{
		Location: Location{
			Type:   "METHOD",
			Path:   "/pets",
			Method: "GET",
		},
		Properties: json.RawMessage(`{"summary": "List all pets", "description": "Returns a list of all pets", "tags": ["pets"]}`),
	}

	err := handler.processMethodPart(docPart)
	c.NoError(err)

	// Verify the method information was added
	pathItem = spec.Paths.MapOfPathItemValues["/pets"]
	operation = pathItem.MapOfOperationValues["get"]
	c.NotNil(operation.Summary)
	c.Equal("List all pets", *operation.Summary)
	c.NotNil(operation.Description)
	c.Equal("Returns a list of all pets", *operation.Description)
	c.Len(operation.Tags, 1)
	c.Equal("pets", operation.Tags[0])
}