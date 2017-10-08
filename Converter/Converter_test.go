package Converter

import "fmt"
import "testing"
import "os"
import "path/filepath"

import "github.com/DinoInc/DomainInflator/Thrift"

var _ = fmt.Println

var _converter_schema_input = `
{
  "definitions": {
    "Structure1": {
      "allOf": [
        {
          "description": "Some Structure called Structure1",
          "properties": {
            "primitive1": {
              "description": "Some primitive property of structure1",
              "type": "string",
              "enum": ["a", "b", "c"]
            },
            "reference2": {
              "$ref": "#/definitions/Structure2"
            }
          }
        }
      ]
    },
    "Structure2": {
      "allOf": [
        {
          "description": "Some Structure called Structure2",
          "properties": {
            "primitive2": {
              "description": "Some primitive property of structure2",
              "type": "string"
            }
          }
        }
      ]
    }
  }
}`

var _converter_thrift_input = `
enum Structure1Primitive1 {
  a
  b
  c
}

struct Structure2 {
  1: optional string primitive2
}

struct Structure1 {
  1: optional Structure1Primitive1 primitive1
  2: optional Structure2 reference2
}
`

func __PutToTMP(filename string, content string) {
	tmpfile := filepath.Join(os.TempDir(), filename)

	file, err := os.Create(tmpfile) // For read access.
	if err != nil {
		panic(err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		panic(err)
	}
}

func TestNewThrift(t *testing.T) {

	__PutToTMP("Structure1.schema.json", _converter_schema_input)
	defer os.Remove("Structure1.schema.json")

	e := NewConverter("/tmp")
	e.NewIDL()
	e.ResolveDefinitionOf("Structure1")
	e.Convert()

	thriftString := e.Thrift()
	thrift := Thrift.ReadIDL([]byte(thriftString))
	thrift.Resolve()

	if _, isExist := thrift.FindEnum("Structure1Primitive1"); !isExist {
		t.Error("Thrift not generate expected enum")
	}

	if _, isExist := thrift.FindStructure("Structure1"); !isExist {
		t.Error("Thrift not generate expected structure {Structure1}")
	}

	if _, isExist := thrift.FindStructure("Structure2"); !isExist {
		t.Error("Thrift not generate expected structure {Structure2}")
	}
}

func TestReadThrift(t *testing.T) {

	__PutToTMP("Structure1.schema.json", _converter_schema_input)
	defer os.Remove("Structure1.schema.json")
	__PutToTMP("Structure1.thrift", _converter_thrift_input)
	defer os.Remove("Structure1.thrift")

	e := NewConverter("/tmp")
	e.ReadIDL("/tmp/Structure1.thrift")
	e.ResolveDefinitionOf("Structure1")

	// e.ResolveDefinitionOf("Structure1")
	// fmt.Println(e.Thrift())
}
