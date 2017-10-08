package Thrift

import "strings"
import "bytes"
import "sort"
import "fmt"
import "github.com/DinoInc/DomainInflator/Utils"

type Enum struct {
	items     map[int]string
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

func (r *Enum) Items() map[int]string {
	return r.items
}

func ReadEnum(content []string) (*Enum, bool) {

	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "enum" {
		var _internal _internal_enum
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Enum{_internal: _internal, items: make(map[int]string), index: make(map[string]int)}, true
	}

	return nil, false

}

func (r *Enum) Resolve() *Enum {
	for i, item := range r._internal._content {
		identifier := strings.TrimSpace(item)

		r.items[i+1] = identifier
		r.index[identifier] = i + 1
		r._internal._lastIndex = i + 1
	}

	return r
}

func NewEnum(identifier string) *Enum {
	var _internal _internal_enum
	_internal.identifier = identifier
	_internal._lastIndex = 0

	return &Enum{_internal: _internal, items: make(map[int]string), index: make(map[string]int)}
}

func (r *Enum) AddMember(member string) {
	r._internal._lastIndex++
	r.items[r._internal._lastIndex] = member
	r.index[member] = r._internal._lastIndex
}

func (r *Enum) IndexOf(identifier string) (int, bool) {
	index, isExists := r.index[identifier]
	return index, isExists
}

func (r *Enum) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "enum %s {\n", r.Identifier())
	var items = r.Items()

	var keys []int
	for k := range items {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Fprintf(&buffer, "\t%s\n", items[key])
	}

	fmt.Fprintf(&buffer, "}\n\n")

	return buffer.String()
}
