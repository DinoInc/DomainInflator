package Thrift

import "fmt"
import "testing"

var _ = fmt.Println

var _enumInputInvalid = []string{
	`xenum someEnum {`,
	`}`,
}

var _enumInputValid = []string{
	`enum someEnum {`,
	`	a`,
	`	b`,
	`	c`,
	`}`,
}

func TestEnumInvalid(t *testing.T) {
	_, isEnum := ReadEnum(_enumInputInvalid)
	if isEnum {
		t.Error("ReadEnum on Invalid return Enum")
	}
}

func TestEnumValid(t *testing.T) {
	enum, isEnum := ReadEnum(_enumInputValid)
	if !isEnum {
		t.Error("ReadEnum on Valid not return Enum")
	}

	enum.Resolve()

	if enum.Identifier() != "someEnum" {
		t.Error("ReadEnum on Valid not return expected enum")
	}

	if index, isExist := enum.IndexOf("c"); !isExist || index != 3 {
		t.Error("ReadEnum on Valid not return expected enum")
	}

}

func TestEnumNew(t *testing.T) {

	enum2 := NewEnum("someEnum2")
	enum2.AddMember("a")
	enum2.AddMember("b")
	enum2.AddMember("c")

	if index, isExist := enum2.IndexOf("c"); !isExist || index != 3 {
		t.Error("NewEnum not return expected enum")
	}

}
