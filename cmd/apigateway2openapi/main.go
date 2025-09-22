package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mauriciocm9/apigateway2openapi"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "apigateway2openapi [input-file]",
	Short: "Convert AWS API Gateway documentation to OpenAPI 3.0 specification",
	Long: `apigateway2openapi converts AWS API Gateway documentation parts to valid OpenAPI 3.0 specifications.

Supported Documentation parts:
- API: General API information
- METHOD: HTTP method definitions
- MODEL: Schema definitions
- PATH_PARAMETER: Path parameter documentation
- QUERY_PARAMETER: Query parameter documentation
- REQUEST_BODY: Request body schemas
- RESOURCE: Resource-level documentation
- RESPONSE: Response definitions`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]

		if verbose {
			log.Printf("Processing API Gateway file: %s", inputFile)
			log.Printf("Output will be written to: %s", outputFile)
		}

		b, err := apigateway2openapi.Execute(inputFile)
		if err != nil {
			return fmt.Errorf("error processing API Gateway file: %w", err)
		}

		err = os.WriteFile(outputFile, b, 0644)
		if err != nil {
			return fmt.Errorf("error writing output file %s: %w", outputFile, err)
		}

		if verbose {
			log.Printf("Successfully converted API Gateway documentation to OpenAPI spec: %s", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "out.yaml", "Output file for the OpenAPI specification")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
