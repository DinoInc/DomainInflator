package Schema

import "encoding/json"

type Enum struct {
	_internal _internal_enum
}

func (r *Enum) Resolve(schema *Schema) *Enum {
	// do nothing
	return r
}

type _internal_enum struct {
	Identifier  string
	Description string        `json:"description,omitempty"`
	Type        PrimitiveType `json:"type,omitempty"`
	Enum        []string      `json:"enum,omitempty"`
}

func (r *Enum) Members() []string {
	return r._internal.Enum
}

func (r *Enum) Description() string {
	return r._internal.Description
}

func (r *Enum) Identifier() string {
	return r._internal.Identifier
}

func (r *Enum) SetIdentifier(identifier string) {
	r._internal.Identifier = identifier
}

func ReadEnum(data *json.RawMessage) (*Enum, bool) {
	var _internal _internal_enum

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if _internal.Type == Str && len(_internal.Enum) > 0 {
		return &Enum{_internal: _internal}, true
	}

	return nil, false
}
