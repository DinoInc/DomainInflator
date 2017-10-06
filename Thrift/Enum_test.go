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

	if len(enum.Items) != 3 {
		t.Error("ReadEnum on Valid not return expected enum")
	}
}
