package Schema

import "encoding/json"

type Structure struct {
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}

func ReadStructure(data *json.RawMessage) (*Structure, bool) {
	var structure Structure

	if json.Unmarshal(*data, &structure) != nil {
		return nil, false
	}

	if len(structure.Properties) > 0 {
		return &structure, true
	}

	return nil, false
}
