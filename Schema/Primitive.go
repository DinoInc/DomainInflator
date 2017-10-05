package Schema

import "encoding/json"

type PropertyPrimitive struct {
	Description string      `json:"description,omitempty"`
	Type        elementType `json:"type,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

func ReadPrimitive(data *json.RawMessage) (*PropertyPrimitive, bool) {
	var property PropertyPrimitive

	if json.Unmarshal(*data, &property) != nil {
		return nil, false
	}

	if property.Type == Str || property.Type == Number || property.Type == Boolean {
		return &property, true
	}

	return nil, false
}
