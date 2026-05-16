package docs

import _ "embed"

// OpenAPIYAML is the OpenAPI 3 description for /api/v1.
//
//go:embed openapi.yaml
var OpenAPIYAML []byte

// SwaggerUIPage is the Swagger UI shell (requests ./openapi.yaml under /swagger).
//
//go:embed swagger-ui.html
var SwaggerUIPage []byte
