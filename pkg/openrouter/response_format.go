package openrouter

import (
	"github.com/kklipsch/billy-bot/pkg/jsonschema"
)

// ResponseFormatEnabled represents the response_format field in OpenRouter API requests
type ResponseFormatEnabled struct {
	Type       string             `json:"type"`
	JSONSchema *jsonschema.Schema `json:"json_schema"`
}

// NewResponseFormatEnabled creates a new ResponseFormatEnabled with the given type and JSON schema
func NewResponseFormatEnabled(schema *jsonschema.Schema) ResponseFormatEnabled {
	return ResponseFormatEnabled{
		Type:       "json_schema",
		JSONSchema: schema,
	}
}
