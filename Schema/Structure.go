package Schema

import "encoding/json"

type Structure struct {
	Properties map[string]interface{}
	_internal  _internal_structure
}

type _internal_structure struct {
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}

func (r *Structure) Resolve(schema *Schema) *Structure {
	r.Properties = make(map[string]interface{})
	for propertyName, propertyMeta := range r._internal.Properties {

		var property interface{}
		if enum, isEnum := ReadEnum(propertyMeta); isEnum {
			property = enum.Resolve(schema)
		} else if array, isArray := ReadArray(propertyMeta); isArray {
			property = array.Resolve(schema)
		} else if primitive, isPrimitive := ReadPrimitive(propertyMeta); isPrimitive {
			property = primitive.Resolve(schema)
		} else if ref, isRef := ReadRef(propertyMeta); isRef {
			property = ref.Resolve(schema)
		} else {
			panic("not implemented")
		}

		r.Properties[propertyName] = property

	}

	return r
}

func ReadStructure(data *json.RawMessage) (*Structure, bool) {
	var _internal _internal_structure

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if len(_internal.Properties) > 0 {
		return &Structure{_internal: _internal}, true
	}

	return nil, false
}
