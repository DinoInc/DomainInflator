package Schema

import "encoding/json"

type Schema struct {
	Definitions map[string]interface{}
	_internal   _internal_schema
}

type _internal_schema struct {
	Definitions map[string]__internal_schema_definition `json:"definitions,omitempty"`
}
type __internal_schema_definition struct {
	AllOf []*json.RawMessage `json:"allOf,omitempty"`
	AnyOf []*json.RawMessage `json:"anyOf,omitempty"`
}

func (r *Schema) Resolve(structureName string) *Structure {

	if structure, exists := __resolved[structureName]; exists {
		return structure
	}

	if len(r._internal.Definitions[structureName].AnyOf) > 0 {
		panic("not implemented")
	}

	var definition Structure
	definition.Properties = make(map[string]interface{})

	for _, s := range r._internal.Definitions[structureName].AllOf {

		var _structure *Structure

		if ref, isRef := ReadRef(s); isRef {
			_structure = ref.Resolve(r)
		} else if structure, isStructure := ReadStructure(s); isStructure {
			_structure = structure.Resolve(r)
		} else {
			panic("not implemented")
		}

		for propertyName, propertyMeta := range _structure.Properties {
			definition.Properties[propertyName] = propertyMeta
		}

	}

	__resolved[structureName] = &definition

	return &definition

}

func ReadSchema(data *json.RawMessage) (*Schema, bool) {
	var _internal _internal_schema

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if len(_internal.Definitions) > 0 {
		return &Schema{_internal: _internal}, true
	}

	return nil, false
}
