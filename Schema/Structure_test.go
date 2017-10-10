package Schema

import "fmt"
import "testing"
import "strconv"
import "encoding/json"

var _ = fmt.Println

var _structure_input = []string{
	`{
    "description": "deskripsi",
    "properties": {
      "property1": {
        "description": "deskripsi1",
        "type": "string"
      }
    }
  }`,
	`{
    "description": "deskripsi",
    "properties": {
      "property1": {
        "description": "deskripsi1",
        "type": "string",
        "enum": ["a", "b", "c"]
      }
    }
  }`,
	`{
    "description": "deskripsi",
    "properties": {
      "property1": {
        "description": "deskripsi1",
        "type": "array",
        "items": { "type": "string" }
      }
    }
  }`,
	`{
    "description": "deskripsi",
    "properties": {
      "property1": {
        "description": "deskripsi1",
        "type": "null"
      }
    }
  }`,
}

func TestStructureNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isStructure := ReadStructure(&data)

	if isStructure != false {
		t.Error("ReadStructure on NonJSON")
	}
}

func TestStructureNonStructure(t *testing.T) {
	data := json.RawMessage(`{"type": "array"}`)
	_, isStructure := ReadStructure(&data)

	if isStructure != false {
		t.Error("ReadStructure on NonRef")
	}
}

func TestStructureInput(t *testing.T) {
	for i := 0; i < len(_structure_input); i++ {

		data := json.RawMessage(_structure_input[i])
		structure, isStructure := ReadStructure(&data)

		if isStructure != true {
			t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not return Structure")
		}

		receivePanic := false
		var recoveredPanic interface{}
		defer func() {
			recoveredPanic = recover()
			receivePanic = (recoveredPanic != nil)
		}()

		structure.Resolve(nil)

		if len(structure.Properties) != 1 {
			t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not return expected Structure")
		}

		switch i {
		case 0:
			property1, ok := structure.Properties["property1"].(*Primitive)
			if !ok || property1.Type() != Str {
				t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}
		case 1:
			property1, ok := structure.Properties["property1"].(*Enum)
			if !ok || len(property1.Members()) != 3 {
				t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}
		case 2:
			property1, ok := structure.Properties["property1"].(*Array)
			elementType := property1.ElementType()
			if !ok || elementType.(*Primitive).Type() != Str {
				t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}
		case 3:
			if !receivePanic {
				t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] not receive panic")
				receivePanic = false
			}
		}

		if receivePanic {
			t.Error("ReadStructure on Input[" + strconv.Itoa(i) + "] receive panic" + recoveredPanic.(string))
		}

	}
}
