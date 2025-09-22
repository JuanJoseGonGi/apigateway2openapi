# AWS API Gateway to OpenAPI 3.0 Converter

![Code lines](https://sloc.xyz/github/mauriciocm9/apigateway2openapi/?category=code)
![Comments](https://sloc.xyz/github/mauriciocm9/apigateway2openapi/?category=comments)

A comprehensive Go tool that converts AWS API Gateway documentation parts to valid OpenAPI 3.0 specifications with **complete coverage** of all AWS documentation part types.

## 🎯 Why Use This Tool

Transform your AWS API Gateway exports into industry-standard OpenAPI specifications for:
- **API Documentation**: Generate comprehensive API documentation for development teams
- **Code Generation**: Use with OpenAPI generators to create client SDKs and server stubs
- **API Testing**: Integrate with testing frameworks that support OpenAPI specifications
- **DevOps Integration**: Automate API documentation as part of your CI/CD pipeline

## ✨ Features

- **Complete AWS API Gateway Support**: Supports all AWS documentation part types for comprehensive API documentation
- **Professional CLI**: Intuitive command-line interface with configurable options and detailed help
- **Flexible Output**: Custom output filenames, verbose logging, and configurable processing options
- **Production Ready**: Enterprise-grade reliability with robust error handling and input validation
- **Well Tested**: Comprehensive test suite ensuring consistent and accurate conversions

## 📋 Supported Documentation Parts

- ✅ **API**: General API information and metadata
- ✅ **METHOD**: Complete HTTP method definitions (GET, POST, PUT, DELETE, PATCH, etc.)
- ✅ **MODEL**: Schema definitions converted to OpenAPI components
- ✅ **PATH_PARAMETER**: Path parameter documentation (e.g., `/pets/{petId}`)
- ✅ **QUERY_PARAMETER**: Query parameter documentation
- ✅ **REQUEST_BODY**: Request body schemas and examples
- ✅ **RESOURCE**: Resource-level documentation
- ✅ **RESPONSE**: Response definitions with proper error handling

## 🚀 Installation

### Build from source
```bash
git clone https://github.com/mauriciocm9/apigateway2openapi.git
cd apigateway2openapi
go build -o apigateway2openapi cmd/apigateway2openapi/main.go
```

### Using go install
```bash
go install github.com/mauriciocm9/apigateway2openapi/cmd/apigateway2openapi@latest
```

## 💻 Usage

### Basic Usage
```bash
# Convert with default output (out.yaml)
./apigateway2openapi input-file.json

# Custom output file
./apigateway2openapi input-file.json -o my-api-spec.yaml

# Verbose mode for detailed logging
./apigateway2openapi input-file.json -v
```

### CLI Options
```bash
# Show help and all available options
./apigateway2openapi --help

# All options example
./apigateway2openapi samples/PetStore-staging-oas30-apigateway.json \
  --output custom-spec.yaml \
  --verbose
```

### Using go run
```bash
# Run directly with go
go run cmd/apigateway2openapi/main.go samples/PetStore-staging-oas30-apigateway.json
```

## 📖 Examples

### Convert Pet Store API
```bash
# Using the included sample
./apigateway2openapi samples/PetStore-staging-oas30-apigateway.json -o petstore-openapi.yaml -v
```

Expected output:
```
Processing API Gateway file: samples/PetStore-staging-oas30-apigateway.json
Output will be written to: petstore-openapi.yaml
Successfully converted API Gateway documentation to OpenAPI spec: petstore-openapi.yaml
```

## 🧪 Testing

```bash
# Run all tests
go test -v

# Run specific test
go test -run TestSpecHandler_processMethodPart -v

# Run tests with coverage
go test -v -cover
```

## 🛠 Development

### Project Structure
```
├── cmd/apigateway2openapi/    # CLI application entry point
├── samples/                   # Sample API Gateway exports for testing
├── models.go                  # Data structures for AWS and OpenAPI formats
├── execute.go                 # Main conversion engine with input validation
├── documentation_part.go      # Type-specific processors for all documentation parts
├── *_test.go                 # Comprehensive test suite
└── CLAUDE.md                 # Detailed architecture documentation
```

### Build and Test
```bash
# Install dependencies
go get .

# Build
go build -v ./...

# Test
go test -v

# Build CLI binary
go build -o apigateway2openapi cmd/apigateway2openapi/main.go
```

## 🎯 Key Features

### Production Ready
- **Enterprise Grade**: Robust error handling and input validation
- **Performance Optimized**: Efficient processing of large API Gateway exports
- **Reliable**: Comprehensive test suite ensuring consistent output quality

### Developer Experience
- **Modern CLI**: Intuitive command-line interface with helpful flags and options
- **Flexible Output**: Configurable output files and detailed logging options
- **Easy Integration**: Simple to incorporate into CI/CD pipelines and automation workflows

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

- 📫 Create an issue for bug reports or feature requests
- 📚 Check [CLAUDE.md](./CLAUDE.md) for detailed architecture documentation
- 🔍 Review the test files for usage examples
