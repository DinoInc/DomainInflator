package Schema

import "fmt"
import "testing"
import "encoding/json"

var _ = fmt.Println

func TestArrayNonJSON(t *testing.T) {
	data := json.RawMessage(``)
	_, isArray := ReadArray(&data)

	if isArray != false {
		t.Error("ReadArray on NonJSON")
	}
}

func TestArrayNonArray(t *testing.T) {
	data := json.RawMessage(`{"type": "string"}`)
	_, isArray := ReadArray(&data)

	if isArray != false {
		t.Error("ReadArray on NonPrimitive")
	}
}

func TestArrayDescription(t *testing.T) {
	data := json.RawMessage(`{"description":"deskripsi", "type": "array", "items": {"type": "string"}}`)
	array, isArray := ReadArray(&data)

	if isArray != true {
		t.Error("ReadArray on Description")
	}

	if array.Description() != "deskripsi" {
		t.Error("ReadArray on Description return description not match")
	}
}

func TestEnumItemsString(t *testing.T) {
	data := json.RawMessage(`{"type": "array", "items": {"type": "string"}}`)
	array, isArray := ReadArray(&data)

	if isArray != true {
		t.Error("ReadEnum on Items=String return not an array")
	}

	array = array.Resolve(nil)
	elementType := array.ElementType().(*Primitive)

	if elementType.Type() != Str {
		t.Error("ReadEnum on Items=String return not an array with items=string")
	}
}

func TestEnumItemsNull(t *testing.T) {
	data := json.RawMessage(`{"type": "array", "items": {"type": "null"}}`)
	array, isArray := ReadArray(&data)

	if isArray != true {
		t.Error("ReadEnum on Items=Null return not an array")
	}

	receivePanic := false
	defer func() {
		r := recover()
		receivePanic = (r != nil)
	}()

	array = array.Resolve(nil)

	if receivePanic {
		t.Error("ReadEnum on Items=Null return not receive panic")
	}
}
