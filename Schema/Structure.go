package Schema

import "encoding/json"

type SchemaStructure struct {
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}
