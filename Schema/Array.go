package Schema

import "encoding/json"

type Array struct {
	elementType interface{}
	_internal   _internal_array
}

type _internal_array struct {
	Description string           `json:"description,omitempty"`
	Type        PrimitiveType    `json:"type,omitempty"`
	Items       *json.RawMessage `json:"items,omitempty"`
}

func (r *Array) ElementType() interface{} {
	return r.elementType
}

func (r *Array) Description() string {
	return r._internal.Description
}

func (r *Array) Resolve(schema *Schema) *Array {

	var elementType interface{}

	if ref, isRef := ReadRef(r._internal.Items); isRef {
		elementType = ref.Resolve(schema)
	} else if primitive, isPrimitive := ReadPrimitive(r._internal.Items); isPrimitive {
		elementType = primitive
	} else {
		panic("not implemented")
	}

	r.elementType = elementType

	return r
}

func ReadArray(data *json.RawMessage) (*Array, bool) {

	var _internal _internal_array

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if _internal.Type == Arr {
		return &Array{elementType: nil, _internal: _internal}, true
	}

	return nil, false
}
