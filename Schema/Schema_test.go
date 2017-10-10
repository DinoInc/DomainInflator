package Schema

import "fmt"
import "testing"
import "strconv"
import "encoding/json"

var _ = fmt.Println

var _schema_input = []string{
	`
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
}
  `,
	`
{
  "definitions": {
    "Structure3": {
      "allOf": [
        {
          "$ref": "#/definitions/Structure4"
        },
        {
          "description": "Some Structure called Structure3",
          "properties": {
            "primitive1": {
              "description": "Some primitive property of structure3",
              "type": "string"
            }
          }
        }
      ]
    },
    "Structure4": {
      "allOf": [
        {
          "description": "Some Structure called Structure4",
          "properties": {
            "primitive2": {
              "description": "Some primitive property of structure4",
              "type": "string"
            }
          }
        }
      ]
    }
  }
}
  `,
	`
{
  "definitions": {
    "Structure5": {
      "allOf": [
        {
          "description": "Some Structure called Structure5",
          "properties": {
            "array1": {
              "description": "Some primitive property of structure5",
              "type": "array",
              "items": {
                "$ref": "#/definitions/Structure6"
              }
            }
          }
        }
      ]
    },
    "Structure6": {
      "allOf": [
        {
          "description": "Some Structure called Structure6",
          "properties": {
            "primitive2": {
              "description": "Some primitive property of structure6",
              "type": "string"
            }
          }
        }
      ]
    }
  }
}
  `,
	`
{
  "definitions": {
    "Structure7": {
      "anyOf": [
        {
          "description": "Some Structure called Structure7",
          "properties": {
            "array1": {
              "description": "Some primitive property of structure7",
              "type": "array",
              "items": {
                "$ref": "#/definitions/Structure6"
              }
            }
          }
        },
        {
          "description": "Some Structure called Structure8",
          "properties": {
            "primitive2": {
              "description": "Some primitive property of structure8",
              "type": "string"
            }
          }
        }
      ]
    }
  }
}
  `,
}

func TestSchemaNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isSchema := ReadSchema(&data)

	if isSchema != false {
		t.Error("ReadSchema on NonJSON")
	}
}

func TestSchemaNonSchema(t *testing.T) {
	data := json.RawMessage(`{"type": "array"}`)
	_, isSchema := ReadSchema(&data)

	if isSchema != false {
		t.Error("ReadSchema on NonSchema")
	}
}

func TestSchemaInput(t *testing.T) {
	for i := 0; i < len(_schema_input); i++ {

		__resolved = make(map[string]*Structure)

		data := json.RawMessage(_schema_input[i])
		schema, isSchema := ReadSchema(&data)

		receivePanic := false
		var recoveredPanic interface{}
		defer func() {
			recoveredPanic = recover()
			receivePanic = (recoveredPanic != nil)
		}()

		if isSchema != true {
			t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return Structure")
		}

		switch i {

		case 0:
			s1 := schema.Resolve("Structure1")

			if len(s1.Properties) != 2 {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}

			// multiple resolve

			s2 := schema.Resolve("Structure1")

			if len(s1.Properties) != len(s2.Properties) {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}

		case 1:
			s1 := schema.Resolve("Structure3")

			if len(s1.Properties) != 2 {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}

		case 2:
			s1 := schema.Resolve("Structure5")

			if len(s1.Properties) != 1 {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}

			array := s1.Properties["array1"].(*Array)
			structure := array.ElementType().(*Structure)
			primitive := structure.Properties["primitive2"].(*Primitive)

			if primitive.Type() != Str {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not return expected Structure")
			}

		case 3:
			_ = schema.Resolve("Structure7")

			if !receivePanic {
				t.Error("ReadSchema on Input[" + strconv.Itoa(i) + "] not receive panic ")

				receivePanic = false
			}
		}

	}
}
