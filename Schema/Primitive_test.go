package Schema

import "fmt"
import "testing"
import "encoding/json"

var _ = fmt.Println

func TestPrimitiveNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isPrimitive := ReadPrimitive(&data)

	if isPrimitive != false {
		t.Error("ReadPrimitive on NonJSON")
	}
}

func TestPrimitiveNonPrimitive(t *testing.T) {
	data := json.RawMessage(`{"type": "array"}`)
	_, isPrimitive := ReadPrimitive(&data)

	if isPrimitive != false {
		t.Error("ReadPrimitive on NonPrimitive")
	}
}

func TestPrimitiveDescription(t *testing.T) {
	data := json.RawMessage(`{"description":"deskripsi", "type":"string"}`)
	primitive, _ := ReadPrimitive(&data)

	primitive = primitive.Resolve(nil)

	if primitive.Description() != "deskripsi" {
		t.Error("ReadPrimitive description not match")
	}
}

func TestPrimitiveBoolean(t *testing.T) {
	data := json.RawMessage(`{"description":"deskripsi", "type":"bool"}`)
	primitive, _ := ReadPrimitive(&data)

	primitive = primitive.Resolve(nil)

	if primitive.Type() != Boolean {
		t.Error("ReadPrimitive type not match")
	}
}
