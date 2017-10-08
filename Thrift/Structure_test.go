package Thrift

import "fmt"
import "testing"

var _ = fmt.Println

var _structInputInvalid = []string{
	`xstruct someStruct {`,
	`}`,
}

var _structInputValid = []string{
	`struct someStruct {`,
	`	1: optional string id`,
	`	2: bool property1`,
	`	3:bool property1`,
	`}`,
}

func TestStructInvalid(t *testing.T) {
	_, isStructure := ReadStructure(_structInputInvalid)
	if isStructure {
		t.Error("ReadStructure on Invalid return Structure")
	}
}

func TestStructValid(t *testing.T) {
	structure, isStructure := ReadStructure(_structInputValid)
	if !isStructure {
		t.Error("ReadStructure on Valid not return Structure")
	}

	structure.Resolve()

	if len(structure.Properties()) != 3 && structure.Identifier() != "someStruct" {
		t.Error("ReadStructure on Valid not return expected structure")
	}
}
