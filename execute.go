package apigateway2openapi

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
	spec           *openapi3.Spec
}

func NewSpecHandler(path string) (SpecHandler, error) {
	handler := SpecHandler{
		spec: &openapi3.Spec{},
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

	return handler, json.Unmarshal(b, &handler.specApigateway)
}

func Execute(path string) ([]byte, error) {
	h, err := NewSpecHandler(path)
	if err != nil {
		return nil, err
	}

	// Validate input structure
	if h.specApigateway.AmazonApigatewayDocumentation.DocumentationParts == nil {
		return nil, fmt.Errorf("invalid API Gateway export: no documentation parts found")
	}

	if len(h.specApigateway.AmazonApigatewayDocumentation.DocumentationParts) == 0 {
		log.Println("Warning: API Gateway export contains no documentation parts")
	}

	for _, x := range h.specApigateway.AmazonApigatewayDocumentation.DocumentationParts {
		if x.Location.Method == "*" {
			continue
		}

		if x.Location.Path == "" {
			x.Location.Path = "/"
		}

		switch x.Location.Type {
		case "API":
			// TODO: remove
			err = h.processAPIPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: API")
				return nil, err
			}
		case "METHOD":
			err := h.processMethodPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: METHOD", err.Error())
				return nil, err
			}
		case "MODEL":
			schema := jsonschema.SchemaOrBool{}
			err := schema.UnmarshalJSON(x.Properties)
			if err != nil {
				log.Println("Failed processing documentation part of type: MODEL for:", x.Location.Name, "with error:", err.Error())
				return nil, err
			}

			schemaOrRef := openapi3.SchemaOrRef{}
			schemaOrRef.FromJSONSchema(schema)
			h.spec.ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(x.Location.Name, schemaOrRef)
		case "PATH_PARAMETER":
			err := h.processPathParameterPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: PATH_PARAMETER", err.Error())
				return nil, err
			}
		case "QUERY_PARAMETER":
			err := h.processQueryParameterPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: QUERY_PARAMETER", err.Error())
				return nil, err
			}
		case "REQUEST_BODY":
			// if enableHack == "yes" {
			// 	continue
			// }

			err := h.processRequestBodyPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: REQUEST_BODY", err.Error())
				// return nil, err
			}
		case "RESOURCE":
			m := Method{}

			err := json.Unmarshal(x.Properties, &m)
			if err != nil {
				log.Println("Found error processing documentation part: RESOURCE")
				return nil, err
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[strings.ToLower(x.Location.Path)]; ok && path.Summary == nil && m.Summary != "" {
				h.spec.Paths.WithMapOfPathItemValuesItem(strings.ToLower(x.Location.Path), *path.WithSummary(m.Summary))
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[strings.ToLower(x.Location.Path)]; ok && path.Description == nil && m.Description != "" {
				h.spec.Paths.WithMapOfPathItemValuesItem(strings.ToLower(x.Location.Path), *path.WithSummary(m.Description))
			}
		case "RESPONSE":
			err := h.processResponsePart(x)
			if err != nil {
				log.Println("Found error processing documentation part: RESPONSE", err.Error())
				return nil, err
			}
		}
	}

	return h.spec.MarshalYAML()
}
