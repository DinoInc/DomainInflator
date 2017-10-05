package Schema

import "encoding/json"

type SchemaDefinition struct {
	AllOf []*json.RawMessage `json:"allOf,omitempty"`
	AnyOf []*json.RawMessage `json:"anyOf,omitempty"`
}

type Schema struct {
	Definitions map[string]SchemaDefinition `json:"definitions,omitempty"`
}

func ReadSchema(data *json.RawMessage) (*Schema, bool) {
	var schema Schema

	if json.Unmarshal(*data, &schema) != nil {
		return nil, false
	}

	if len(schema.Definitions) > 0 {
		return &schema, true
	}

	return nil, false
}
