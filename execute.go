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

var enableHack = os.Getenv("HACK")

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
			if enableHack != "yes" {
				continue
			}

			// TODO: Remove this logic, once aws_resource are migrated to the module
			if strings.ToLower(x.Location.Method) != "options" {
				continue
			}

			if path, ok := h.spec.Paths.MapOfPathItemValues[x.Location.Path]; ok && path.Summary == nil {
				m := Method{}

				err := json.Unmarshal(x.Properties, &m)
				if err != nil {
					log.Println("Found error processing documentation part: METHOD")
					return nil, err
				}

				h.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, *path.WithSummary(m.Summary).WithDescription(m.Description))
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
			// TODO
		case "QUERY_PARAMETER":
			// TODO
		case "REQUEST_BODY":
			if enableHack == "yes" {
				continue
			}

			err := h.processRequestBodyPart(x)
			if err != nil {
				log.Println("Found error processing documentation part: REQUEST_BODY")
				return nil, err
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
				log.Println("Found error processing documentation part: RESPONSE")
				return nil, err
			}
		}
	}

	return h.spec.MarshalYAML()
}
