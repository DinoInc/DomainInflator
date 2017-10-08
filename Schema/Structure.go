package Schema

import "encoding/json"

type Structure struct {
	Properties map[string]interface{}
	_internal  _internal_structure
}

type _internal_structure struct {
	Identifier  string
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}

func (r *Structure) Identifier() string {
	return r._internal.Identifier
}

func (r *Structure) SetIdentifier(identifier string) {
	r._internal.Identifier = identifier
}

func (r *Structure) Resolve(schema *Schema) *Structure {
	r.Properties = make(map[string]interface{})

	for propertyName, propertyMeta := range r._internal.Properties {

		if propertyName[0] == '_' {
			continue
		}

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
