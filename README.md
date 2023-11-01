# Convert apigateway documentation parts to OpenAPI valid specification

![Code lines](https://sloc.xyz/github/mauriciocm9/apigateway2openapi/?category=code)
![Comments](https://sloc.xyz/github/mauriciocm9/apigateway2openapi/?category=comments)

## Supported Documentation parts

* API
* METHOD
* MODEL
* REQUEST_BODY
* RESOURCE
* RESPONSE

## Example

```bash
go run . API-staging-oas30-apigateway.json
```

## Roadmap

| Feature/Improvement              | Description                               | Status       |
|----------------------------------|-------------------------------------------|--------------|
| Support PATH_PARAMETER           | -                                         | Not Started  |
| Support QUERY_PARAMETER          | -                                         | Not Started  |
| Remove getBytes method           | This is a temp hack                       | Not Started  |
| Don't set summary using options  | This is a temp hack                       | Not Started  |
