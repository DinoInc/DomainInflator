package Thrift

import "strings"
import "regexp"

var __reHeader = regexp.MustCompile(`^(?P<type>[A-z]+) (?P<identifier>[A-z][A-z0-9]*) \{`)

type _enumState uint32

const (
	none _enumState = iota
	begin
	content
	end
)

type IDL struct {
	enum         map[string]*Enum
	structure    map[string]*Structure
	resolveOrder []string
	_content     []byte

	_state _enumState
}

func ReadIDL(content []byte) *IDL {
	return &IDL{_content: content, _state: none, structure: make(map[string]*Structure), enum: make(map[string]*Enum)}
}

func NewIDL() *IDL {
	return &IDL{_content: nil, _state: none, structure: make(map[string]*Structure), enum: make(map[string]*Enum)}
}

func (r *IDL) Enum() map[string]*Enum {
	return r.enum
}

func (r *IDL) AddEnum(e *Enum) {
	r.enum[e.Identifier()] = e
}

func (r *IDL) FindEnum(identifier string) (*Enum, bool) {
	enum, isExist := r.enum[identifier]
	return enum, isExist
}

func (r *IDL) Structure() map[string]*Structure {
	return r.structure
}

func (r *IDL) AddStructure(s *Structure) {
	r.structure[s.Identifier()] = s
}

func (r *IDL) FindStructure(identifier string) (*Structure, bool) {
	structure, isExist := r.structure[identifier]
	return structure, isExist
}

func (r *IDL) ResolveOrder() []string {
	return r.resolveOrder
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

			if enum, isEnum := ReadEnum(_lines[_start:_end]); isEnum {
				r.enum[enum.Identifier()] = enum.Resolve()
			} else if structure, isStructure := ReadStructure(_lines[_start:_end]); isStructure {
				r.structure[structure.Identifier()] = structure.Resolve()
				r.resolveOrder = append(r.resolveOrder, structure.Identifier())
			}
		}
	}
}
