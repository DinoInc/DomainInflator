package Engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

import (
	"github.com/DinoInc/DomainInflator/Schema"
	"github.com/DinoInc/DomainInflator/Thrift"
)

type SchemaIDL struct {
	structure map[string]*Schema.Structure
}

type ThriftIDL struct {
	enums     map[string]*Thrift.Enum
	structure map[string]*Thrift.Structure
}

type Engine struct {
	currentFile   string
	baseDir       string
	schemaIDL     SchemaIDL
	thriftIDL     ThriftIDL
	deviations    []*Deviation
	resolvedOrder []string
}

func NewEngine(baseDir string) *Engine {
	e := Engine{}

	e.schemaIDL.structure = make(map[string]*Schema.Structure)

	e.thriftIDL.enums = make(map[string]*Thrift.Enum)
	e.thriftIDL.structure = make(map[string]*Thrift.Structure)

	e.resolvedOrder = make([]string, 0)
	e.deviations = make([]*Deviation, 0)

	e.baseDir = baseDir

	return &e
}

func (e *Engine) SetCurrentFile(currentFile string) {
	e.currentFile = currentFile
}

func (e *Engine) handlePrimitiveEnum(schemaStructure Schema.Structure, name string, propertyPrimitive *Schema.PropertyPrimitive) {

	enum := Schema.PropertyEnum{Identifier: name}
	for _, value := range propertyPrimitive.Enum {
		enum.Items = append(enum.Items, value)
	}

	schemaStructure.Properties[name] = enum

}

func (e *Engine) handlePrimitive(schemaStructure Schema.Structure, name string, propertyPrimitive *Schema.PropertyPrimitive) {

	if propertyPrimitive.Enum != nil {

		e.handlePrimitiveEnum(schemaStructure, name, propertyPrimitive)

	} else {

		switch propertyPrimitive.Type {
		case Schema.Number, Schema.Boolean, Schema.Str:
			schemaStructure.Properties[name] = propertyPrimitive.Type
		default:
			panic("not implemented")
		}

	}

}

func (e *Engine) handleArray(schemaStructure Schema.Structure, name string, propertyArray *Schema.PropertyArray) {

	if ref, isRef := Schema.ReadRef(propertyArray.Items); isRef {
		schemaStructure.Properties[name] = Schema.PropertyList{ElementType: ref.Name()}
		e.Resolve(ref)
	} else if primitive, isPrimitive := Schema.ReadPrimitive(propertyArray.Items); isPrimitive {
		schemaStructure.Properties[name] = Schema.PropertyList{ElementType: string(primitive.Type)}
	} else {
		panic("not implemented")
	}

}

func (e *Engine) handleAllOfRef(schemaStructure Schema.Structure, ref *Schema.Ref) {
	structureName := ref.Name()
	e.Resolve(ref)

	for property, propertyType := range e.schemaIDL.structure[structureName].Properties {
		schemaStructure.Properties[property] = propertyType
	}
}

func (e *Engine) handleAllOfStructure(schemaStructure Schema.Structure, meta *json.RawMessage) {

	structure, isStruct := Schema.ReadSchemaStructure(meta)
	if !isStruct {
		panic("not implemented")
	}

	for property, propertyJSON := range structure.Properties {

		//fmt.Println(property)

		// skip _
		if property[0] == '_' {
			continue
		}

		if primitive, isPrimitive := Schema.ReadPrimitive(propertyJSON); isPrimitive {
			e.handlePrimitive(schemaStructure, property, primitive)
		} else if array, isArray := Schema.ReadArray(propertyJSON); isArray {
			e.handleArray(schemaStructure, property, array)
		} else if ref, isRef := Schema.ReadRef(propertyJSON); isRef {
			schemaStructure.Properties[property] = ref.Name()
			e.Resolve(ref)
		}

	}
}

func (e *Engine) handleAllOf(identifier string, definition Schema.SchemaDefinition) {
	//fmt.Println("--------- " + identifier)

	schemaStructure := Schema.Structure{}
	schemaStructure.Identifier = identifier
	schemaStructure.Properties = make(map[string]interface{})

	for _, meta := range definition.AllOf {

		if ref, isRef := Schema.ReadRef(meta); isRef {
			e.handleAllOfRef(schemaStructure, ref)
		} else {
			e.handleAllOfStructure(schemaStructure, meta)
		}

	}

	e.schemaIDL.structure[identifier] = &schemaStructure
}

func (e *Engine) Resolve(ref *Schema.Ref) {

	structureName := ref.Name()
	if _, isExists := e.schemaIDL.structure[structureName]; isExists {
		return
	}

	var structureFile = e.currentFile
	if !ref.IsSelf() {
		structureFile = ref.File()
	}

	content, err := ioutil.ReadFile(e.baseDir + structureFile)
	if err != nil {

		content, err = ioutil.ReadFile(e.baseDir + structureFile + ".schema.json")

		if err != nil {
			fmt.Printf("File error: %v\n", e)
			os.Exit(1)
		}

	}

	var _currentFile = e.currentFile
	e.currentFile = structureFile

	jsonContent := json.RawMessage(content)
	s, isSchema := Schema.ReadSchema(&jsonContent)
	if !isSchema {
		panic("not implemented")
	}

	definition := s.Definitions[structureName]

	if definition.AllOf != nil {
		e.handleAllOf(structureName, definition)
	} else {
		panic("not implemented")
	}

	e.currentFile = _currentFile

	e.resolvedOrder = append(e.resolvedOrder, structureName)
}

func (e *Engine) Print() {
	for _, x := range e.schemaIDL.structure {
		fmt.Println(x)
	}
}
