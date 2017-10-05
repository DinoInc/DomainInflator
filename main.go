package main

import "encoding/json"
import "io/ioutil"
import "fmt"
import "regexp"
import "os"
import "reflect"
import "strings"
import "bufio"
import "flag"

var _ = reflect.TypeOf

type elementType string

const (
	null    elementType = "null"
	boolean elementType = "bool"
	object  elementType = "object"
	array   elementType = "array"
	number  elementType = "i32"
	str     elementType = "string"
)

type Ref struct {
	Ref string `json:"$ref"`
}

type PropertyPrimitive struct {
	Description string      `json:"description,omitempty"`
	Type        elementType `json:"type,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

type PropertyArray struct {
	Description string           `json:"description,omitempty"`
	Type        elementType      `json:"type,omitempty"`
	Items       *json.RawMessage `json:"items,omitempty"`
}

type SchemaStructure struct {
	Description string                      `json:"description,omitempty"`
	Properties  map[string]*json.RawMessage `json:"properties,omitempty"`
}

type ThriftStructure struct {
	Identifier  string
	Description string            `json:"description,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}

type ThriftEnum struct {
	Items []string
}

type SchemaDefinition struct {
	AllOf []*json.RawMessage `json:"allOf,omitempty"`
	AnyOf []*json.RawMessage `json:"anyOf,omitempty"`
}

type Schema struct {
	Definitions map[string]SchemaDefinition `json:"definitions,omitempty"`
}

var reservedSet map[string]bool
var resolvedList = []string{"BEGIN", "END", "__CLASS__", "__DIR__", "__FILE__", "__FUNCTION__", "__LINE__", "__METHOD__", "__NAMESPACE__", "abstract", "alias", "and", "args", "as", "assert", "begin", "break", "case", "catch", "class", "clone", "continue", "declare", "def", "default", "del", "delete", "do", "dynamic", "elif", "else", "elseif", "elsif", "end", "enddeclare", "endfor", "endforeach", "endif", "endswitch", "endwhile", "ensure", "except", "exec", "finally", "float", "for", "foreach", "from", "function", "global", "goto", "if", "implements", "import", "in", "inline", "instanceof", "interface", "is", "lambda", "module", "native", "new", "next", "nil", "not", "or", "package", "pass", "public", "print", "private", "protected", "raise", "redo", "rescue", "retry", "register", "return", "self", "sizeof", "static", "super", "switch", "synchronized", "then", "this", "throw", "transient", "try", "undef", "unless", "unsigned", "until", "use", "var", "virtual", "volatile", "when", "while", "with", "xor", "yield"}

func IsReservedWord(word string) bool {

	if reservedSet == nil {
		reservedSet = make(map[string]bool)
		for _, reservedWord := range resolvedList {
			reservedSet[reservedWord] = true
		}
	}

	_, isInReservedWord := reservedSet[word]
	return isInReservedWord
}

func readPrimitive(data *json.RawMessage) (*PropertyPrimitive, bool) {
	var property PropertyPrimitive

	if json.Unmarshal(*data, &property) != nil {
		return nil, false
	}

	if property.Type != str && property.Type != number && property.Type != boolean {
		return nil, false
	}

	return &property, true
}

func readArray(data *json.RawMessage) (*PropertyArray, bool) {
	var property PropertyArray

	if json.Unmarshal(*data, &property) != nil {
		return nil, false
	}

	if property.Type != array {
		return nil, false
	}

	return &property, true
}

func readRef(data *json.RawMessage) (*Ref, bool) {
	var ref Ref

	if json.Unmarshal(*data, &ref) != nil {
		return nil, false
	}

	if ref.Ref == "" {
		return nil, false
	}

	return &ref, true
}

func getStructureName(refString string) string {
	return regexp.MustCompile(`[a-zA-z]*$`).FindString(refString)
}

func handlePrimitiveEnum(thriftStructure ThriftStructure, property string, propertyPrimitive *PropertyPrimitive) {

	enumName := upperConcat("Enum", removeUnderscore(thriftStructure.Identifier), property)

	enum := ThriftEnum{}
	for _, value := range propertyPrimitive.Enum {
		enum.Items = append(enum.Items, value)
	}

	enums[enumName] = &enum
	thriftStructure.Properties[property] = enumName

}

func handleArray(thriftStructure ThriftStructure, property string, propertyArray *PropertyArray) {

	if ref, isRef := readRef(propertyArray.Items); isRef {
		thriftStructure.Properties[property] = "list<" + getStructureName(ref.Ref) + ">"
		resolve(*ref)
	} else if primitive, isPrimitive := readPrimitive(propertyArray.Items); isPrimitive {
		thriftStructure.Properties[property] = "list<" + string(primitive.Type) + ">"
	} else {
		panic("not implemented")
	}

}

func handlePrimitive(thriftStructure ThriftStructure, property string, propertyPrimitive *PropertyPrimitive) {

	if propertyPrimitive.Enum != nil {

		handlePrimitiveEnum(thriftStructure, property, propertyPrimitive)

	} else {

		switch propertyPrimitive.Type {
		case number, boolean, str:
			thriftStructure.Properties[property] = string(propertyPrimitive.Type)
		default:
			panic("not implemented")
		}

	}

}

func handleAllOf(identifier string, definition SchemaDefinition) {
	//fmt.Println("--------- " + identifier)

	thriftStructure := ThriftStructure{}
	thriftStructure.Identifier = identifier
	thriftStructure.Properties = make(map[string]string)

	for _, meta := range definition.AllOf {

		if ref, isRef := readRef(meta); isRef {

			structureName := getStructureName(ref.Ref)
			resolve(*ref)

			for property, propertyType := range resolved[structureName].Properties {
				thriftStructure.Properties[property] = propertyType
			}

		} else {

			var structure SchemaStructure
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

				if primitive, isPrimitive := readPrimitive(propertyJSON); isPrimitive {
					handlePrimitive(thriftStructure, property, primitive)
				} else if array, isArray := readArray(propertyJSON); isArray {
					handleArray(thriftStructure, property, array)
				} else if ref, isRef := readRef(propertyJSON); isRef {
					thriftStructure.Properties[property] = getStructureName(ref.Ref)
					resolve(*ref)
				}

			}
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

func getStructureFile(refString string) string {
	return regexp.MustCompile(`^[a-zA-z.]*`).FindString(refString)
}

func resolve(ref Ref) {

	structureName := getStructureName(ref.Ref)
	if _, isExists := resolved[structureName]; isExists {
		return
	}

	var structureFile = currentFile
	if ref.Ref[0] != '#' {
		structureFile = getStructureFile(ref.Ref)
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

	var s Schema
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

var enums map[string]*ThriftEnum
var resolved map[string]*ThriftStructure
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

func upperConcat(s ...string) string {
	result := strings.ToLower(s[0])
	for i := 1; i < len(s); i++ {
		result += strings.Title(s[i])
	}
	return result
}

func removeUnderscore(s string) string {
	re := regexp.MustCompile(`_+`)
	replaced := re.ReplaceAllString(s, "")
	return replaced
}

func isValidIdentifier(s string) bool {
	isMatch, _ := regexp.MatchString(`^[A-z_][A-z0-9._]*$`, s)
	return isMatch
}

func main() {

	enums = make(map[string]*ThriftEnum)
	resolved = make(map[string]*ThriftStructure)
	resolved = make(map[string]*ThriftStructure)
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

		var ref Ref
		if json.Unmarshal(content, &ref) != nil {
			panic("not implemented")
		}

		if ref.Ref == "" {
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

			if IsReservedWord(propertyName) {
				newPropertyName := upperConcat(structName, propertyName)

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
