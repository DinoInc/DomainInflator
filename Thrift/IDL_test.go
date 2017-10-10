package Thrift

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

var _idl_valid = []byte(`
namespace go some_idl

enum someEnum {
	a
	b
	c
}

struct someStruct {
	1: optional string someProperty1
	2: optional integer someProperty2
}`)

func TestIDLValid(t *testing.T) {
	idl := ReadIDL(_idl_valid)

	idl.Resolve()

	if len(idl.Enum()) != 1 {
		t.Error("ReadIDL on Valid not return expected enum")
	}

	if len(idl.Structure()) != 1 {
		t.Error("ReadIDL on Valid not return expected structure")
	}
}

func TestNewIDL(t *testing.T) {
	idl := NewIDL()

	if len(idl.Enum()) != 0 {
		t.Error("New IDL not return Empty Enum")
	}
	if len(idl.Structure()) != 0 {
		t.Error("New IDL not return Empty Structure")
	}
}

func TestNewIDLAddStructure(t *testing.T) {
	idl := NewIDL()

	structure := NewStructure("someStruct")
	structure.AddProperty("someProperty1", "string")
	structure.AddProperty("someProperty2", "bool")
	idl.AddStructure(structure)

	if _, isExist := idl.FindStructure("someStruct"); !isExist {
		t.Error("New IDL on AppendStructure not return expected structure")
	}
}

func TestNewIDLAddEnum(t *testing.T) {
	idl := NewIDL()

	enum := NewEnum("someEnum")
	enum.AddItem("a")
	enum.AddItem("b")
	enum.AddItem("c")
	idl.AddEnum(enum)

	if _, isExist := idl.FindEnum("someEnum"); !isExist {
		t.Error("Existing IDL on AppendEnum not return expected enum")
	}
}

func TestIDLAddStructure(t *testing.T) {
	idl := ReadIDL(_idl_valid)
	idl.Resolve()

	structure2 := NewStructure("someStruct2")
	structure2.AddProperty("someProperty1", "string")
	structure2.AddProperty("someProperty2", "bool")
	idl.AddStructure(structure2)

	if _, isExist := idl.FindStructure("someStruct2"); !isExist {
		t.Error("Existing IDL on AppendStructure not return expected structure")
	}
}

func TestIDLAddEnum(t *testing.T) {
	idl := ReadIDL(_idl_valid)
	idl.Resolve()

	enum2 := NewEnum("someEnum2")
	enum2.AddItem("a")
	enum2.AddItem("b")
	enum2.AddItem("c")
	idl.AddEnum(enum2)

	if _, isExist := idl.FindEnum("someEnum2"); !isExist {
		t.Error("Existing IDL on AppendEnum not return expected enum")
	}
}

func TestIDLAddStructureProperty(t *testing.T) {
	idl := ReadIDL(_idl_valid)
	idl.Resolve()

	structure, _ := idl.FindStructure("someStruct")
	structure.AddProperty("someProperty3", "someEnum")

	s, _ := idl.FindStructure("someStruct")
	if tag, isExist := s.TagOf("someProperty3"); !isExist || tag != 3 {
		t.Error("Existing IDL on AppendProperty not return expected property")
	}

}

func TestIDLAddEnumMember(t *testing.T) {
	idl := ReadIDL(_idl_valid)
	idl.Resolve()

	enum, _ := idl.FindEnum("someEnum")
	enum.AddItem("d")

	e, _ := idl.FindEnum("someEnum")
	if index, isExist := e.IndexOf("d"); !isExist || index != 4 {
		t.Error("Existing IDL on AppendEnum not return expected property")
	}
}
