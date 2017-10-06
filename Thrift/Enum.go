package Thrift

import "strings"
import "github.com/DinoInc/DomainInflator/Utils"

type Enum struct {
	Items     []string
	_internal _internal_enum
}

type _internal_enum struct {
	identifier string
	_content   []string
}

func ReadEnum(content []string) (*Enum, bool) {

	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "enum" {
		var _internal _internal_enum
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Enum{_internal: _internal}, true
	}

	return nil, false

}

func (r *Enum) Resolve() {
	for _, item := range r._internal._content {
		r.Items = append(r.Items, strings.TrimSpace(item))
	}
}
