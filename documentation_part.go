package apigateway2openapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/swaggest/openapi-go/openapi3"
)

func (sh *SpecHandler) processAPIPart(x DocumentationPart) error {
	var out map[string]json.RawMessage

	err := json.Unmarshal(x.Properties, &out)
	if err != nil {
		return err
	}

	if _, ok := out["info"]; ok {
		return nil
	}

	err = sh.spec.Info.UnmarshalJSON(x.Properties)
	if err != nil {
		log.Println("Failed processing documentation part of type: API")
		return err
	}

	return nil
}

func (sh *SpecHandler) processResponsePart(x DocumentationPart) error {
	operationKey := strings.ToLower(x.Location.Method)
	path := sh.spec.Paths.MapOfPathItemValues[x.Location.Path]
	operation := path.MapOfOperationValues[operationKey]
	responses := operation.Responses
	responseOrRef := openapi3.ResponseOrRef{}

	err := responseOrRef.UnmarshalJSON(x.Properties)
	if err != nil {
		// Try to handle custom Response format
		var response Response
		if errNew := json.Unmarshal(x.Properties, &response); errNew != nil {
			log.Println("Failed processing documentation part of type: RESPONSE for status:", x.Location.StatusCode)
			return err
		}

		// Convert custom Response to OpenAPI format
		responseData := map[string]any{
			"description": response.Description,
		}

		if len(response.ResponseExamples) > 0 {
			content := make(map[string]map[string]any)
			for contentType, examples := range response.ResponseExamples {
				exampleMap := make(map[string]any)
				for key, example := range examples {
					var value any = example.Value

					// Try to parse JSON if content type suggests it
					if contentType == "application/json" {
						if err := json.Unmarshal([]byte(example.Value), &value); err != nil {
							// If JSON parsing fails, use raw value
							value = example.Value
						}
					}

					exampleMap[key] = map[string]any{"value": value}
				}
				content[contentType] = map[string]any{"examples": exampleMap}
			}
			responseData["content"] = content
		}

		responseBytes, err := json.Marshal(responseData)
		if err != nil {
			return fmt.Errorf("failed to marshal response data: %w", err)
		}

		err = responseOrRef.UnmarshalJSON(responseBytes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal converted response: %w", err)
		}
	}

	sh.spec.Paths.WithMapOfPathItemValuesItem(
		x.Location.Path,
		*path.WithMapOfOperationValuesItem(operationKey, *operation.WithResponses(
			*responses.WithMapOfResponseOrRefValuesItem(
				x.Location.StatusCode,
				responseOrRef,
			),
		)))

	return nil
}

func (sh *SpecHandler) processRequestBodyPart(x DocumentationPart) error {
	b := x.Properties

	operationKey := strings.ToLower(x.Location.Method)
	path := sh.spec.Paths.MapOfPathItemValues[x.Location.Path]
	operation := path.MapOfOperationValues[operationKey]
	requestBodyOrRef := operation.RequestBodyEns()

	err := requestBodyOrRef.UnmarshalJSON(b)
	if err != nil {
		// TODO: delete this
		log.Println("Using fallback for REQUEST_BODY")

		requestBody := RequestBody{}

		newErr := json.Unmarshal(x.Properties, &requestBody)
		if newErr != nil {
			return err
		}

		if requestBody.Content == "" {
			return nil
		}

		b = []byte(`{"content":` + requestBody.Content + `}`)
	}

	err = requestBodyOrRef.UnmarshalJSON(b)
	if err != nil {
		// end here
		return err
	}

	sh.spec.Paths.WithMapOfPathItemValuesItem(
		x.Location.Path,
		*path.WithMapOfOperationValuesItem(operationKey, *operation.WithRequestBody(*requestBodyOrRef)))

	return nil
}

func (sh *SpecHandler) processPathParameterPart(x DocumentationPart) error {
	if x.Location.Path == "" || x.Location.Name == "" {
		return nil
	}

	pathItem, exists := sh.spec.Paths.MapOfPathItemValues[x.Location.Path]
	if !exists {
		pathItem = openapi3.PathItem{}
		sh.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, pathItem)
	}

	parameter := openapi3.ParameterOrRef{
		Parameter: &openapi3.Parameter{
			Name: x.Location.Name,
			In:   "path",
		},
	}

	err := parameter.Parameter.UnmarshalJSON(x.Properties)
	if err != nil {
		log.Println("Failed processing documentation part of type: PATH_PARAMETER for parameter:", x.Location.Name)
		return err
	}

	parameter.Parameter.Name = x.Location.Name
	parameter.Parameter.In = "path"
	requiredVal := true
	parameter.Parameter.Required = &requiredVal

	pathItem.Parameters = append(pathItem.Parameters, parameter)
	sh.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, pathItem)

	return nil
}

func (sh *SpecHandler) processQueryParameterPart(x DocumentationPart) error {
	if x.Location.Path == "" || x.Location.Method == "" || x.Location.Name == "" {
		return nil
	}

	operationKey := strings.ToLower(x.Location.Method)
	pathItem, exists := sh.spec.Paths.MapOfPathItemValues[x.Location.Path]
	if !exists {
		return nil
	}

	operation, exists := pathItem.MapOfOperationValues[operationKey]
	if !exists {
		return nil
	}

	parameter := openapi3.ParameterOrRef{
		Parameter: &openapi3.Parameter{
			Name: x.Location.Name,
			In:   "query",
		},
	}

	err := parameter.Parameter.UnmarshalJSON(x.Properties)
	if err != nil {
		log.Println("Failed processing documentation part of type: QUERY_PARAMETER for parameter:", x.Location.Name)
		return err
	}

	parameter.Parameter.Name = x.Location.Name
	parameter.Parameter.In = "query"

	operation.Parameters = append(operation.Parameters, parameter)
	pathItem.WithMapOfOperationValuesItem(operationKey, operation)
	sh.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, pathItem)

	return nil
}

func (sh *SpecHandler) processMethodPart(x DocumentationPart) error {
	if x.Location.Path == "" || x.Location.Method == "" {
		return nil
	}

	operationKey := strings.ToLower(x.Location.Method)
	pathItem, exists := sh.spec.Paths.MapOfPathItemValues[x.Location.Path]
	if !exists {
		return nil
	}

	operation, exists := pathItem.MapOfOperationValues[operationKey]
	if !exists {
		return nil
	}

	var method Method
	err := json.Unmarshal(x.Properties, &method)
	if err != nil {
		log.Println("Failed processing documentation part of type: METHOD for:", x.Location.Method)
		return err
	}

	if method.Summary != "" && operation.Summary == nil {
		operation.Summary = &method.Summary
	}
	if method.Description != "" && operation.Description == nil {
		operation.Description = &method.Description
	}
	if len(method.Tags) > 0 && len(operation.Tags) == 0 {
		operation.Tags = method.Tags
	}

	pathItem.WithMapOfOperationValuesItem(operationKey, operation)
	sh.spec.Paths.WithMapOfPathItemValuesItem(x.Location.Path, pathItem)

	return nil
}
