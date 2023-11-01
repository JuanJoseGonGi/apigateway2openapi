package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type SpecHandler struct {
	specApigateway OpenapiGateway
	spec           openapi3.Spec
	reflector      *openapi3.Reflector
}

func NewSpecHandler(path string) (SpecHandler, error) {
	handler := SpecHandler{
		reflector: openapi3.NewReflector(),
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return handler, err
	}

	err = handler.spec.UnmarshalJSON(b)
	if err != nil {
		fmt.Println("Found error in file:", path)
		return handler, err
	}

	handler.reflector.Spec = &handler.spec

	return handler, json.Unmarshal(b, &handler.specApigateway)
}

func execute(path string) error {
	h, err := NewSpecHandler(path)
	if err != nil {
		return err
	}

	for _, x := range h.specApigateway.AmazonApigatewayDocumentation.DocumentationParts {
		switch x.Location.Type {
		case "API":
			err = h.spec.Info.UnmarshalJSON(x.Properties)
			if err != nil {
				log.Println("Failed processing documentation part of type: API")
				return err
			}
		case "METHOD":
			// TODO: Remove this logic, once aws_resource are migrated to the module
			if strings.ToLower(x.Location.Method) != "options" {
				continue
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[x.Location.Path]; ok && path.Summary == nil {
				m := Method{}

				err := json.Unmarshal(x.Properties, &m)
				if err != nil {
					return err
				}

				h.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, *path.WithSummary(m.Summary).WithDescription(m.Description))
			}
		case "MODEL":
			schema := jsonschema.SchemaOrBool{}
			err := schema.UnmarshalJSON(x.Properties)
			if err != nil {
				log.Println("Failed processing documentation part of type: MODEL for:", x.Location.Name, "with error:", err.Error())
				return err
			}

			schemaOrRef := openapi3.SchemaOrRef{}
			schemaOrRef.FromJSONSchema(schema)
			h.spec.ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(x.Location.Name, schemaOrRef)
		case "PATH_PARAMETER":
			// TODO
		case "QUERY_PARAMETER":
			// TODO
		case "REQUEST_BODY":
			requestBody := RequestBody{}

			err := json.Unmarshal(x.Properties, &requestBody)
			if err != nil {
				return err
			}

			if requestBody.Content == "" {
				continue
			}

			operationKey := strings.ToLower(x.Location.Method)
			path := h.spec.Paths.MapOfPathItemValues[x.Location.Path]
			operation := path.MapOfOperationValues[operationKey]
			requestBodyOrRef := operation.RequestBodyEns()

			// TODO: fix this hack. It should be: requestBodyOrRef.UnmarshalJSON(x.Properties)
			err = requestBodyOrRef.UnmarshalJSON([]byte(`{"content":` + requestBody.Content + `}`))
			if err != nil {
				return err
			}

			h.spec.Paths.WithMapOfPathItemValuesItem(
				x.Location.Path,
				*path.WithMapOfOperationValuesItem(operationKey, *operation.WithRequestBody(*requestBodyOrRef)))
		case "RESOURCE":
			m := Method{}

			err := json.Unmarshal(x.Properties, &m)
			if err != nil {
				return err
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[strings.ToLower(x.Location.Path)]; ok && path.Summary == nil && m.Summary != "" {
				h.spec.Paths.WithMapOfPathItemValuesItem(strings.ToLower(x.Location.Path), *path.WithSummary(m.Summary))
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[strings.ToLower(x.Location.Path)]; ok && path.Description == nil && m.Description != "" {
				h.spec.Paths.WithMapOfPathItemValuesItem(strings.ToLower(x.Location.Path), *path.WithSummary(m.Description))
			}
		case "RESPONSE":
			var response Response
			if err := json.Unmarshal([]byte(x.Properties), &response); err != nil {
				return err
			}

			operationKey := strings.ToLower(x.Location.Method)
			path := h.spec.Paths.MapOfPathItemValues[x.Location.Path]
			operation := path.MapOfOperationValues[operationKey]
			responses := operation.Responses
			responseOrRef := openapi3.ResponseOrRef{}

			err := responseOrRef.UnmarshalJSON(response.getBytes())
			if err != nil {
				return err
			}

			h.spec.Paths.WithMapOfPathItemValuesItem(
				x.Location.Path,
				*path.WithMapOfOperationValuesItem(operationKey, *operation.WithResponses(
					*responses.WithMapOfResponseOrRefValuesItem(
						x.Location.StatusCode,
						responseOrRef,
					),
				)))
		}
	}

	b, err := h.spec.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile("out.json", b, 0600)
}

func main() {
	if len(os.Args) != 2 {
		log.Println("Missing openapi path")
		return
	}

	err := execute(os.Args[1])
	if err != nil {
		log.Println("Found error trying to read file:", err.Error())
		return
	}

}
