package Schema

import "encoding/json"

type PropertyArray struct {
	Description string           `json:"description,omitempty"`
	Type        elementType      `json:"type,omitempty"`
	Items       *json.RawMessage `json:"items,omitempty"`
}

func ReadArray(data *json.RawMessage) (*PropertyArray, bool) {
	var property PropertyArray

	if json.Unmarshal(*data, &property) != nil {
		return nil, false
	}

	if property.Type == Array {
		return &property, true
	}

	return nil, false
}
