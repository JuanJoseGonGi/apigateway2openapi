package main

import (
	"encoding/json"
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

	b := x.Properties

	err := responseOrRef.UnmarshalJSON(b)
	if err != nil {
		// TODO: delete this
		log.Println("Using fallback for hack response")

		var response Response
		if errNew := json.Unmarshal(x.Properties, &response); errNew != nil {
			return err
		}

		b = response.getBytes()
	}

	err = responseOrRef.UnmarshalJSON(b)
	if err != nil {
		// end delete
		return err
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
