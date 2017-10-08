package Converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/DinoInc/DomainInflator/Schema"
	"github.com/DinoInc/DomainInflator/Thrift"
	"github.com/DinoInc/DomainInflator/Utils"
)

var _primitiveMapping = map[Schema.PrimitiveType]string{
	Schema.Str:     "string",
	Schema.Number:  "i32",
	Schema.Boolean: "bool",
}

type Converter struct {
	_structureName map[*Schema.Structure]string
	_structure     *Schema.Structure
	thriftIDL      *Thrift.IDL
	resolveOrder   []string
}

func NewConverter(baseDir string) *Converter {
	Schema.SetBaseDir(baseDir)
	return &Converter{_structureName: make(map[*Schema.Structure]string)}
}

func (c *Converter) NewIDL() {
	c.thriftIDL = Thrift.NewIDL()
}

func (c *Converter) ReadIDL(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c.thriftIDL = Thrift.ReadIDL(data)
	c.thriftIDL.Resolve()
}

func (c *Converter) IDL() *Thrift.IDL {
	return c.thriftIDL
}

func (c *Converter) ResolveDefinitionOf(identifier string) {

	data := json.RawMessage(`{ "$ref":"` + identifier + `#/definition/` + identifier + `"}`)
	ref, _ := Schema.ReadRef(&data)

	c._structure = ref.Resolve(nil)
}

func (c *Converter) _TypeOf(context string, s interface{}) string {
	switch s.(type) {
	case *Schema.Array:
		return "list<" + c._TypeOf(context, s.(*Schema.Array).ElementType()) + ">"
	case *Schema.Structure:
		return c._ConvertStructure(s.(*Schema.Structure))
	case *Schema.Enum:
		return c._ConvertEnum(context, s.(*Schema.Enum))
	case *Schema.Primitive:
		return _primitiveMapping[s.(*Schema.Primitive).Type()]
	default:
		panic("not implemented")
	}
}

func (c *Converter) _ConvertStructure(s *Schema.Structure) string {

	//remove duplicates
	if name, isFound := c._structureName[s]; isFound {
		return name
	}

	thriftStructure := Thrift.NewStructure(s.Identifier())

	for propertyName, propertyMeta := range s.Properties {
		context := Utils.UpperConcat(Utils.RemoveUnderscore(s.Identifier()), propertyName)
		thriftStructure.AddProperty(propertyName, c._TypeOf(context, propertyMeta))
	}

	c.thriftIDL.AddStructure(thriftStructure)

	c.resolveOrder = append(c.resolveOrder, thriftStructure.Identifier())
	c._structureName[s] = thriftStructure.Identifier()

	return thriftStructure.Identifier()

}

func (c *Converter) _ConvertEnum(context string, e *Schema.Enum) string {
	thriftEnum := Thrift.NewEnum(context)

	for _, member := range e.Members() {
		thriftEnum.AddMember(member)
	}

	c.thriftIDL.AddEnum(thriftEnum)

	return thriftEnum.Identifier()
}

func (c *Converter) Convert() {
	c._ConvertStructure(c._structure)
}

func (c *Converter) Thrift() string {
	var buffer bytes.Buffer

	enum := c.thriftIDL.Enum()

	var enumId []string
	for id, _ := range enum {
		enumId = append(enumId, id)
	}
	sort.Strings(enumId)

	for _, id := range enumId {
		fmt.Fprintf(&buffer, "%s", enum[id].String())
	}

	structure := c.thriftIDL.Structure()

	for _, id := range c.resolveOrder {
		fmt.Fprintf(&buffer, "%s", structure[id].String())
	}

	return buffer.String()
}
