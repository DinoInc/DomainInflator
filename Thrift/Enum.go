package Thrift

import "strings"
import "github.com/DinoInc/DomainInflator/Utils"

type Enum struct {
	Items     map[int]string
	index     map[string]int
	_internal _internal_enum
}

type _internal_enum struct {
	identifier string
	_content   []string
	_lastIndex int
}

func (r *Enum) Identifier() string {
	return r._internal.identifier
}

func ReadEnum(content []string) (*Enum, bool) {

	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "enum" {
		var _internal _internal_enum
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Enum{_internal: _internal, Items: make(map[int]string), index: make(map[string]int)}, true
	}

	return nil, false

}

func (r *Enum) Resolve() *Enum {
	for i, item := range r._internal._content {
		identifier := strings.TrimSpace(item)

		r.Items[i+1] = identifier
		r.index[identifier] = i + 1
		r._internal._lastIndex = i + 1
	}

	return r
}

func NewEnum(identifier string) *Enum {
	var _internal _internal_enum
	_internal.identifier = identifier
	_internal._lastIndex = 0

	return &Enum{_internal: _internal, Items: make(map[int]string), index: make(map[string]int)}
}

func (r *Enum) AddMember(member string) {
	r._internal._lastIndex++
	r.Items[r._internal._lastIndex] = member
	r.index[member] = r._internal._lastIndex
}

func (r *Enum) IndexOf(identifier string) (int, bool) {
	index, isExists := r.index[identifier]
	return index, isExists
}
