package Schema

import "fmt"
import "testing"
import "encoding/json"

var _ = fmt.Println

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
