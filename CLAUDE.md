# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive Go tool that converts AWS API Gateway documentation parts to OpenAPI 3.0 specifications. It reads an exported API Gateway documentation JSON file and transforms it into a valid OpenAPI YAML specification with full support for all AWS documentation part types.

## Common Commands

### Build
```bash
go build -v ./...
```

### Build CLI binary
```bash
go build -o apigateway2openapi cmd/apigateway2openapi/main.go
```

### Test
```bash
go test -v
```

### Run specific tests
```bash
go test -run TestSpecHandler_processMethodPart -v
```

### Run the CLI tool
```bash
# Using default output (out.yaml)
./apigateway2openapi samples/PetStore-staging-oas30-apigateway.json

# With custom output file and verbose mode
./apigateway2openapi samples/PetStore-staging-oas30-apigateway.json -o custom-output.yaml -v

# Show help
./apigateway2openapi --help
```

### Using go run
```bash
go run cmd/apigateway2openapi/main.go samples/PetStore-staging-oas30-apigateway.json
```

### Install dependencies
```bash
go get .
```

## Architecture

### Core Components

- **models.go**: Contains data structures that mirror AWS API Gateway documentation parts and OpenAPI schema structures:
  - `AmazonApigatewayDocumentation`: Root structure for API Gateway exports
  - `DocumentationPart`: Individual documentation pieces with location and properties
  - `OpenapiGateway`: Output OpenAPI structure
  - `Method`: Enhanced with Tags support for complete method documentation
  - `Response`, `RequestBody`, `Example`: Supporting structures

- **execute.go**: Main processing engine containing:
  - `SpecHandler`: Core processor that handles the conversion logic
  - `Execute()`: Main entry point with input validation and comprehensive error handling
  - No global variables - clean architecture

- **documentation_part.go**: Complete type-specific processors for all documentation part types:
  - `processAPIPart()`: Handles API-level documentation
  - `processMethodPart()`: Complete HTTP method processing for all verbs (GET, POST, etc.)
  - `processPathParameterPart()`: **NEW** - Processes path parameters (e.g., `/pets/{petId}`)
  - `processQueryParameterPart()`: **NEW** - Processes query parameters
  - `processResponsePart()`: Enhanced response processing without hack methods
  - `processRequestBodyPart()`: Handles request body schemas

- **cmd/apigateway2openapi/main.go**: Professional CLI built with Cobra:
  - Configurable output filename (`-o/--output`)
  - Verbose mode (`-v/--verbose`)
  - Comprehensive help and usage information
  - Proper error handling and user feedback

### Supported Documentation Part Types

The tool now supports **ALL** AWS API Gateway documentation part types:
- ✅ **API**: General API information
- ✅ **METHOD**: Complete HTTP method definitions (all verbs: GET, POST, PUT, DELETE, PATCH, etc.)
- ✅ **MODEL**: Schema definitions converted to OpenAPI components
- ✅ **PATH_PARAMETER**: Path parameter documentation (e.g., `/pets/{petId}`)
- ✅ **QUERY_PARAMETER**: Query parameter documentation
- ✅ **REQUEST_BODY**: Request body schemas
- ✅ **RESOURCE**: Resource-level documentation
- ✅ **RESPONSE**: Response definitions with proper error handling

### Processing Flow

1. **Input Validation**: Validate JSON structure and required fields
2. **Parse**: Load API Gateway documentation from JSON file
3. **Initialize**: Create OpenAPI spec handler
4. **Process**: Iterate through each documentation part with type-specific processors
5. **Convert**: Transform AWS format to OpenAPI 3.0 specification
6. **Output**: Generate YAML file with configurable filename

### Technical Improvements

**Eliminated Technical Debt:**
- ❌ Removed `getBytes()` hack method - now uses proper OpenAPI marshaling
- ❌ Removed `enableHack` global variable - clean architecture
- ❌ Eliminated panic() calls - comprehensive error handling
- ❌ Removed hardcoded output filename - configurable via CLI

**New Features:**
- ✅ Professional CLI with Cobra framework
- ✅ Input validation and error prevention
- ✅ Comprehensive test suite for all processors
- ✅ Support for all AWS documentation part types
- ✅ Enhanced error messages and logging
- ✅ Configurable output and verbose mode

### Dependencies

Key external libraries:
- `github.com/swaggest/openapi-go`: OpenAPI 3.0 specification generation
- `github.com/swaggest/jsonschema-go`: JSON Schema handling
- `github.com/spf13/cobra`: Professional CLI framework
- `github.com/stretchr/testify`: Testing framework

### Testing

The project includes comprehensive tests:
- `execute_test.go`: Integration tests for the main conversion flow
- `documentation_part_test.go`: Unit tests for all processor methods
- Tests for PATH_PARAMETER, QUERY_PARAMETER, and METHOD processing

### Code Quality

- No global variables or hacks
- Proper error handling throughout
- Input validation prevents runtime crashes
- Clean, maintainable architecture
- Comprehensive documentation and examples