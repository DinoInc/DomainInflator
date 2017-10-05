package Schema

import "encoding/json"

type PrimitiveType string

const (
	Null    PrimitiveType = "null"
	Boolean PrimitiveType = "bool"
	Object  PrimitiveType = "object"
	Arr     PrimitiveType = "array"
	Number  PrimitiveType = "i32"
	Str     PrimitiveType = "string"
)

type Primitive struct {
	_internal _internal_primitive
}

func (r *Primitive) Type() PrimitiveType {
	return r._internal.Type
}

func (r *Primitive) Description() string {
	return r._internal.Description
}

func (r *Primitive) Resolve(schema *Schema) *Primitive {
	// do nothing
	return r
}

type _internal_primitive struct {
	Description string        `json:"description,omitempty"`
	Type        PrimitiveType `json:"type,omitempty"`
}

func ReadPrimitive(data *json.RawMessage) (*Primitive, bool) {
	var _internal _internal_primitive

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if _internal.Type == Str || _internal.Type == Number || _internal.Type == Boolean {
		return &Primitive{_internal: _internal}, true
	}

	return nil, false
}
