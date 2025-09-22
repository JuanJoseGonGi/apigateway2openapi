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
	Description string   `json:"description"`
	Summary     string   `json:"summary"`
	Tags        []string `json:"tags"`
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

