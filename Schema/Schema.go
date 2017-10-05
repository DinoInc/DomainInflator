package Schema

import "encoding/json"

type SchemaDefinition struct {
	AllOf []*json.RawMessage `json:"allOf,omitempty"`
	AnyOf []*json.RawMessage `json:"anyOf,omitempty"`
}

type Schema struct {
	Definitions map[string]SchemaDefinition `json:"definitions,omitempty"`
}
