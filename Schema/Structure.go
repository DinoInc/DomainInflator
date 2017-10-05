package Schema

import "encoding/json"

type SchemaStructure struct {
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}

func ReadSchemaStructure(data *json.RawMessage) (*SchemaStructure, bool) {
	var structure SchemaStructure

	if json.Unmarshal(*data, &structure) != nil {
		return nil, false
	}

	if len(structure.Properties) > 0 {
		return &structure, true
	}

	return nil, false
}
