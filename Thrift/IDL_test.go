package Thrift

import (
	"testing"
)

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
