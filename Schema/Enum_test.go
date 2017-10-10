package Schema

import "fmt"
import "testing"
import "encoding/json"

var _ = fmt.Println

func TestEnumNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isEnum := ReadEnum(&data)

	if isEnum != false {
		t.Error("ReadEnum on NonJSON")
	}
}

func TestEnumNonEnum(t *testing.T) {
	data := json.RawMessage(`{"type": "string"}`)
	_, isEnum := ReadEnum(&data)

	if isEnum != false {
		t.Error("ReadEnum on NonPrimitive")
	}
}

func TestEnum(t *testing.T) {
	data := json.RawMessage(`{"description": "deskripsi", "type": "string", "enum": ["a", "b", "c"]}`)
	enum, isEnum := ReadEnum(&data)

	if isEnum != true {
		t.Error("ReadEnum on [a, b, c] return not an enum")
	}

	enum = enum.Resolve(nil)

	if len(enum.Members()) != 3 {
		t.Error("ReadEnum on [a, b, c] return member not match")
	}

	if enum.Description() != "deskripsi" {
		t.Error("ReadEnum on [a, b, c]  return description not match")
	}
}
