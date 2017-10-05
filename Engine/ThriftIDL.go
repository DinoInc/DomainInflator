package Engine

import (
	"fmt"
)

import (
	"github.com/DinoInc/DomainInflator/Schema"
	"github.com/DinoInc/DomainInflator/Thrift"
)

var _ = fmt.Println

type ThriftIDL struct {
	// state
	resolvedOrder []string

	Enum      map[string]*Thrift.Enum
	Structure map[string]*Thrift.Structure
}

func (e *ThriftIDL) ResolveType(propertyMeta interface{}, se *Schema.SchemaIDL) string {
	switch propertyMeta.(type) {

	case Schema.SchemaEnum:
		return propertyMeta.(Schema.SchemaEnum).Identifier

	case Schema.SchemaPrimitive:
		return string(propertyMeta.(Schema.SchemaPrimitive).Type)

	case Schema.SchemaList:
		elementType := e.ResolveType(propertyMeta.(Schema.SchemaList).ElementType, se)
		return "list<" + elementType + ">"

	case Schema.SchemaStructure:
		identifier := propertyMeta.(Schema.SchemaStructure).Identifier
		e.Resolve(identifier, se)
		return identifier

	default:
		fmt.Println(propertyMeta)
		panic("not implemented")
	}

	return ""
}

func (e *ThriftIDL) Resolve(structureName string, se *Schema.SchemaIDL) {

	if _, isExists := e.Structure[structureName]; isExists {
		return
	}

	schemaStructure := se.Structure[structureName]

	thriftStructure := &Thrift.Structure{
		Identifier:  schemaStructure.Identifier,
		Description: schemaStructure.Description,
		Properties:  make(map[string]string),
	}

	for propertyName, propertyMeta := range se.Structure[structureName].Properties {
		thriftStructure.Properties[propertyName] = e.ResolveType(propertyMeta, se)
	}

	e.Structure[structureName] = thriftStructure

}
