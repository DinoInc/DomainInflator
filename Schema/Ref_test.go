package Schema

import "fmt"
import "testing"
import "encoding/json"
import "os"
import "path/filepath"
import "github.com/DinoInc/DomainInflator/Utils"

var _ = fmt.Println

var _ref_schema_input = `
{
  "definitions": {
    "Structure1": {
      "allOf": [
        {
          "description": "Some Structure called Structure1",
          "properties": {
            "primitive1": {
              "description": "Some primitive property of structure1",
              "type": "string"
            },
            "reference2": {
              "$ref": "#/definitions/Structure2"
            }
          }
        }
      ]
    }
  }
}`

func TestRefNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isRef := ReadRef(&data)

	if isRef != false {
		t.Error("ReadRef on NonJSON")
	}
}

func TestRefNonRef(t *testing.T) {
	data := json.RawMessage(`{"type": "array"}`)
	_, isRef := ReadRef(&data)

	if isRef != false {
		t.Error("ReadRef on NonRef")
	}
}

func TestRefSelfBackboneElement(t *testing.T) {
	data := json.RawMessage(`{"$ref": "#/definitions/BackboneElement"}`)
	ref, isRef := ReadRef(&data)

	if isRef != true {
		t.Error("ReadRef on SelfBackboneElement not return Ref")
	}

	if ref.IsSelf() != true {
		t.Error("ReadRef on SelfBackboneElement, IsSelf not return true")
	}

	if ref.Name() != "BackboneElement" {
		t.Error("ReadRef on SelfBackboneElement, Name not return BackboneElement")
	}
}

func TestRefHumanName(t *testing.T) {
	data := json.RawMessage(`{"$ref": "HumanName.schema.json#/definitions/HumanName"}`)
	ref, isRef := ReadRef(&data)

	if isRef != true {
		t.Error("ReadRef on HumanName not return Ref")
	}

	if ref.IsSelf() != false {
		t.Error("ReadRef on HumanName, IsSelf not return false")
	}

	if ref.Name() != "HumanName" {
		t.Error("ReadRef on HumanName, Name not return HumanName")
	}

	if ref.File() != "HumanName.schema.json" {
		t.Error("ReadRef on HumanName, File not return HumanName.schema.json")
	}
}

func TestRefWithoutExtension(t *testing.T) {
	data := json.RawMessage(`{"$ref": "HumanName#/definitions/HumanName"}`)
	ref, isRef := ReadRef(&data)

	if isRef != true {
		t.Error("ReadRef on WithoutExtension not return Ref")
	}

	if ref.IsSelf() != false {
		t.Error("ReadRef on WithoutExtension, IsSelf not return false")
	}

	if ref.Name() != "HumanName" {
		t.Error("ReadRef on WithoutExtension, Name not return HumanName")
	}

	if ref.File() != "HumanName.schema.json" {
		t.Error("ReadRef on WithoutExtension, File not return HumanName.schema.json")
	}
}

func TestRefMultipleResolve(t *testing.T) {

	__resolved = make(map[string]*Structure)

	data := json.RawMessage(_ref_schema_input)
	schema, _ := ReadSchema(&data)

	data = json.RawMessage(`{"$ref": "#/definitions/Structure1"}`)
	ref, _ := ReadRef(&data)

	s3 := ref.Resolve(schema)
	s4 := ref.Resolve(schema)

	if len(s3.Properties) != len(s4.Properties) {
		t.Error("ReadRef on MultipleResolve not return expected Structure")
	}
}

func TestRefNonSelf(t *testing.T) {

	tmpfile := Utils.TempFileName("x", ".schema.json")

	file, err := os.Create(tmpfile) // For read access.
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile)

	if _, err := file.Write([]byte(_schema_input[0])); err != nil {
		panic(err)
	}

	dir, filename := filepath.Split(tmpfile)
	data := json.RawMessage(`{ "$ref":"` + filename + `#/definition/Structure1"}`)

	SetBaseDir(dir)
	ref, _ := ReadRef(&data)

	structure := ref.Resolve(nil)

	if len(structure.Properties) != 2 {
		t.Error("ReadSchema on NonSelfRef not return expected Structure")
	}
}
