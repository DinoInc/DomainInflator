package Thrift

import "strings"
import "bytes"
import "sort"
import "fmt"
import "github.com/DinoInc/DomainInflator/Utils"

type Enum struct {
	identifier string
	items      map[int]string
	index      map[string]int

	_content   []string
	_lastIndex int
}

func NewEnum(identifier string) *Enum {
	return &Enum{identifier: identifier, _lastIndex: 0, items: make(map[int]string), index: make(map[string]int)}
}

func ReadEnum(content []string) (*Enum, bool) {

	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "enum" {
		_content := content[1 : len(content)-1]
		identifier := header["identifier"]

		return &Enum{_content: _content, identifier: identifier, items: make(map[int]string), index: make(map[string]int)}, true
	}

	return nil, false

}

func (r *Enum) Identifier() string {
	return r.identifier
}

func (r *Enum) Items() map[int]string {
	return r.items
}

func (r *Enum) AddItem(member string) {
	r._lastIndex++
	r.items[r._lastIndex] = member
	r.index[member] = r._lastIndex
}

func (r *Enum) Resolve() *Enum {
	for i, item := range r._content {
		identifier := strings.TrimSpace(item)

		r.items[i+1] = identifier
		r.index[identifier] = i + 1
		r._lastIndex = i + 1
	}

	return r
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
		fmt.Fprintf(&buffer, "  %s\n", items[key])
	}

	fmt.Fprintf(&buffer, "}\n\n")

	return buffer.String()
}
