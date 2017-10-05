package main

import "encoding/json"
import "io/ioutil"
import "fmt"
import "os"
import "reflect"
import "bufio"
import "flag"
import "regexp"
import "strings"

import (
	"github.com/DinoInc/DomainInflator/Schema"
	"github.com/DinoInc/DomainInflator/Thrift"
	"github.com/DinoInc/DomainInflator/Utils"
)

var _ = reflect.TypeOf

func handlePrimitiveEnum(thriftStructure Thrift.ThriftStructure, property string, propertyPrimitive *Schema.PropertyPrimitive) {

	enumName := Utils.UpperConcat("Enum", Utils.RemoveUnderscore(thriftStructure.Identifier), property)

	enum := Thrift.ThriftEnum{}
	for _, value := range propertyPrimitive.Enum {
		enum.Items = append(enum.Items, value)
	}

	enums[enumName] = &enum
	thriftStructure.Properties[property] = enumName

}

func handleArray(thriftStructure Thrift.ThriftStructure, property string, propertyArray *Schema.PropertyArray) {

	if ref, isRef := Schema.ReadRef(propertyArray.Items); isRef {
		thriftStructure.Properties[property] = "list<" + ref.Name() + ">"
		resolve(ref)
	} else if primitive, isPrimitive := Schema.ReadPrimitive(propertyArray.Items); isPrimitive {
		thriftStructure.Properties[property] = "list<" + string(primitive.Type) + ">"
	} else {
		panic("not implemented")
	}

}

func handlePrimitive(thriftStructure Thrift.ThriftStructure, property string, propertyPrimitive *Schema.PropertyPrimitive) {

	if propertyPrimitive.Enum != nil {

		handlePrimitiveEnum(thriftStructure, property, propertyPrimitive)

	} else {

		switch propertyPrimitive.Type {
		case Schema.Number, Schema.Boolean, Schema.Str:
			thriftStructure.Properties[property] = string(propertyPrimitive.Type)
		default:
			panic("not implemented")
		}

	}

}

func handleAllOfRef(thriftStructure Thrift.ThriftStructure, ref *Schema.Ref) {
	structureName := ref.Name()
	resolve(ref)

	for property, propertyType := range resolved[structureName].Properties {
		thriftStructure.Properties[property] = propertyType
	}
}

func handleAllOfStructure(thriftStructure Thrift.ThriftStructure, meta *json.RawMessage) {
	var structure Schema.SchemaStructure

	err := json.Unmarshal(*meta, &structure)
	if err != nil {
		panic("not implemented")
	}

	for property, propertyJSON := range structure.Properties {

		//fmt.Println(property)

		// skip _
		if property[0] == '_' {
			continue
		}

		if primitive, isPrimitive := Schema.ReadPrimitive(propertyJSON); isPrimitive {
			handlePrimitive(thriftStructure, property, primitive)
		} else if array, isArray := Schema.ReadArray(propertyJSON); isArray {
			handleArray(thriftStructure, property, array)
		} else if ref, isRef := Schema.ReadRef(propertyJSON); isRef {
			thriftStructure.Properties[property] = ref.Name()
			resolve(ref)
		}

	}
}

func handleAllOf(identifier string, definition Schema.SchemaDefinition) {
	//fmt.Println("--------- " + identifier)

	thriftStructure := Thrift.ThriftStructure{}
	thriftStructure.Identifier = identifier
	thriftStructure.Properties = make(map[string]string)

	for _, meta := range definition.AllOf {

		if ref, isRef := Schema.ReadRef(meta); isRef {
			handleAllOfRef(thriftStructure, ref)
		} else {
			handleAllOfStructure(thriftStructure, meta)
		}

	}

	resolved[identifier] = &thriftStructure
}

// [v] allOf
// [ ] anyOf
// [v] primitive type
// [v] array -> ref
// [v] array -> primitive
// [ ]

var baseDir string
var namespace string
var currentFile = "Patient.schema.json"

func resolve(ref *Schema.Ref) {

	structureName := ref.Name()
	if _, isExists := resolved[structureName]; isExists {
		return
	}

	var structureFile = currentFile
	if !ref.IsSelf() {
		structureFile = ref.File()
	}

	content, e := ioutil.ReadFile(baseDir + structureFile)
	if e != nil {

		content, e = ioutil.ReadFile(baseDir + structureFile + ".schema.json")

		if e != nil {
			fmt.Printf("File error: %v\n", e)
			os.Exit(1)
		}

	}

	var _currentFile = currentFile
	currentFile = structureFile

	var s Schema.Schema
	json.Unmarshal(content, &s)
	definition := s.Definitions[structureName]

	if definition.AllOf != nil {
		handleAllOf(structureName, definition)
	} else {
		panic("not implemented")
	}

	currentFile = _currentFile

	resolvedOrder = append(resolvedOrder, structureName)
}

var enums map[string]*Thrift.ThriftEnum
var resolved map[string]*Thrift.ThriftStructure
var resolvedOrder []string

type DeviationContext int64

const (
	PropertyIdentifier DeviationContext = 0
	EnumIdentifier     DeviationContext = 1
)

type ThriftDeviation struct {
	Original    string
	Replacement string
	Context     DeviationContext
}

var deviations []*ThriftDeviation

func isValidIdentifier(s string) bool {
	isMatch, _ := regexp.MatchString(`^[A-z_][A-z0-9._]*$`, s)
	return isMatch
}

func main() {

	enums = make(map[string]*Thrift.ThriftEnum)
	resolved = make(map[string]*Thrift.ThriftStructure)
	resolved = make(map[string]*Thrift.ThriftStructure)
	deviations = make([]*ThriftDeviation, 0)

	pBaseDir := flag.String("schema-dir", "./schemas/", "JSON schema Directory")
	pSchemas := flag.String("schema", "", "schemas to resolve")
	pNamespace := flag.String("namespace", "", "thrift namespace")
	flag.Parse()

	if *pSchemas == "" {
		fmt.Fprintf(os.Stderr, "missing required --schema argument/flag\n")
		os.Exit(1)
	}

	if *pNamespace == "" {
		fmt.Fprintf(os.Stderr, "missing required --namespace argument/flag\n")
		os.Exit(1)
	}

	schemaList := strings.Split(*pSchemas, ",")

	for _, schema := range schemaList {

		baseDir = *pBaseDir
		currentFile = schema
		namespace = *pNamespace

		content, e := ioutil.ReadFile(baseDir + currentFile + ".schema.json")
		if e != nil {
			fmt.Printf("File error: %v\n", e)
			os.Exit(1)
		}

		var ref *Schema.Ref
		if json.Unmarshal(content, &ref) != nil {
			panic("not implemented")
		}

		if ref.Name() == "" {
			panic("not implemented")
		}

		resolve(ref)

	}

	thriftFile, _ := os.Create(namespace + ".thrift")
	defer thriftFile.Close()
	thriftWriter := bufio.NewWriter(thriftFile)

	fmt.Fprintf(thriftWriter, "namespace go %s\n", namespace)
	fmt.Fprintf(thriftWriter, "namespace java %s\n\n", namespace)

	for enumName, enum := range enums {
		fmt.Fprintf(thriftWriter, "enum %s {\n", enumName)
		for _, enumValue := range enum.Items {

			if isValidIdentifier(enumValue) {
				fmt.Fprintf(thriftWriter, "\t%s\n", enumValue)
			} else {
				newEnumIdentifier := strings.Replace(enumValue, "-", "_", -1)

				fmt.Fprintf(thriftWriter, "\t%s\n", newEnumIdentifier)

				deviations = append(deviations, &ThriftDeviation{
					Original:    enumValue,
					Replacement: newEnumIdentifier,
					Context:     EnumIdentifier,
				})

			}

		}
		fmt.Fprintf(thriftWriter, "}\n")

	}

	for _, structName := range resolvedOrder {
		structMeta := resolved[structName]

		fmt.Fprintf(thriftWriter, "struct %s {\n", structName)

		var i = 1
		for propertyName, propertyType := range structMeta.Properties {

			if Thrift.IsReservedWord(propertyName) {
				newPropertyName := Utils.UpperConcat(structName, propertyName)

				deviations = append(deviations, &ThriftDeviation{
					Original:    propertyName,
					Replacement: newPropertyName,
					Context:     PropertyIdentifier,
				})
				propertyName = newPropertyName
			}

			fmt.Fprintf(thriftWriter, "\t%d: optional %s %s\n", i, propertyType, propertyName)
			i = i + 1
		}
		fmt.Fprintf(thriftWriter, "}\n")
	}

	sedFile, _ := os.Create(namespace + ".sed")
	sedWriter := bufio.NewWriter(sedFile)
	defer sedFile.Close()

	for _, deviation := range deviations {
		switch deviation.Context {
		case PropertyIdentifier:
			fmt.Fprintf(sedWriter, "s/json:\"%s,omitempty\"/json:\"%s,omitempty\"/g\n", deviation.Replacement, deviation.Original)
		case EnumIdentifier:
			fmt.Fprintf(sedWriter, "s/\"%s\"/\"%s\"/g\n", deviation.Replacement, deviation.Original)
		}
	}

	thriftWriter.Flush()
	sedWriter.Flush()

}
