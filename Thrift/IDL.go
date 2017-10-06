package Thrift

import "strings"
import "regexp"
import "fmt"

var __reHeader = regexp.MustCompile(`^(?P<type>[A-z]+) (?P<identifier>[A-z][A-z0-9]*) \{`)

type _enumState uint32

const (
	none _enumState = iota
	begin
	content
	end
)

type IDL struct {
	enum      []*Enum
	structure []*Structure
	_content  []byte

	_state _enumState
}

func (r *IDL) Enum() []*Enum {
	return r.enum
}

func (r *IDL) Structure() []*Enum {
	return r.enum
}

func (r *IDL) Resolve() {
	var _content string = string(r._content)
	var _lines []string = strings.Split(_content, "\n")

	var _start int
	var _end int
	var i int = 0
	for i < len(_lines) {
		r.processLine(_lines[i])

		if r._state == begin {
			_start = i
		}

		i++

		if r._state == end {
			_end = i

			fmt.Println(i, _lines[_start:_end])

			if enum, isEnum := ReadEnum(_lines[_start:_end]); isEnum {
				r.enum = append(r.enum, enum)
			} else if structure, isStructure := ReadStructure(_lines[_start:_end]); isStructure {
				r.structure = append(r.structure, structure)
			}
		}
	}
}

func (r *IDL) processLine(line string) {
	if strings.Contains(line, "{") && (r._state == end || r._state == none) {
		r._state = begin
	} else if strings.Contains(line, "}") && (r._state == begin || r._state == content) {
		r._state = end
	} else if r._state == begin {
		r._state = content
	} else if r._state == end {
		r._state = none
	}
}

func ReadIDL(content []byte) *IDL {
	return &IDL{_content: content, _state: none}
}
