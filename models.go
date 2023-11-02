package apigateway2openapi

import (
	"encoding/json"
)

type Location struct {
	Type       string
	Path       string
	Method     string
	Name       string
	StatusCode string `json:"statusCode"`
}

type DocumentationPart struct {
	Location   Location
	Properties json.RawMessage
}

type AmazonApigatewayDocumentation struct {
	Version            string
	DocumentationParts []DocumentationPart `json:"documentationParts"`
}

type OpenapiGateway struct {
	AmazonApigatewayDocumentation AmazonApigatewayDocumentation `json:"x-amazon-apigateway-documentation"`
}

type Method struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type RequestBody struct {
	Content     string `json:"content"`
	Description string `json:"description"`
}

type Example struct {
	Value string `json:"value"`
}

type Response struct {
	Description      string                        `json:"description"`
	ResponseExamples map[string]map[string]Example `json:"responseExamples"`
	ResponseModels   map[string]string             `json:"responseModels"`
}

// TODO: Delete this method when responses are set correctly in documentation part
func (r *Response) getBytes() []byte {
	out := make(map[string]any)
	respExamples := make(map[string]map[string]any)

	for contentType, examples := range r.ResponseExamples {
		cc := make(map[string]map[string]any)
		for key, example := range examples {
			var obj any

			obj = example.Value
			if contentType == "application/json" {
				if err := json.Unmarshal([]byte(example.Value), &obj); err != nil {
					panic(err)
				}
			}

			if cc[key] == nil {
				cc[key] = make(map[string]any)
			}

			cc[key]["value"] = obj
			if respExamples[contentType] == nil {
				respExamples[contentType] = make(map[string]any)
			}
			respExamples[contentType]["examples"] = cc
		}
	}

	out["description"] = r.Description
	out["content"] = respExamples

	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	return b
}
