package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swaggest/openapi-go/openapi3"
)

func TestExecute(t *testing.T) {
	c := require.New(t)

	path := "samples/PetStore-staging-oas30-apigateway.json"
	expectedOutFile := "samples/PetStore-staging-oas30-apigateway-expected.yaml"

	expectedSpec := openapi3.Spec{}
	outSpec := openapi3.Spec{}

	b, err := os.ReadFile(expectedOutFile)
	c.NoError(err)
	c.NoError(expectedSpec.UnmarshalYAML(b))

	b, err = Execute(path)
	c.NoError(err)
	c.NoError(outSpec.UnmarshalYAML(b))

	c.Equal(outSpec, expectedSpec)
}
