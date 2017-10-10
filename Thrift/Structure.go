package Thrift

import "fmt"
import "strconv"
import "strings"
import "regexp"
import "sort"
import "bytes"
import "github.com/DinoInc/DomainInflator/Utils"

var _ = fmt.Println

type Structure struct {
	properties map[int]*Property
	tag        map[string]int
	identifier string
	_lastTag   int
	_content   []string
}

var __reProperty = regexp.MustCompile(`(?P<tag>[1-9][0-9]*):(\s*(?P<req>[A-z][A-z0-9]*))?(\s*(?P<type>[A-z][A-z0-9]*(<[A-z][A-z0-9]*>)?))(\s*(?P<identifier>[A-z][A-z0-9]*))`)

func ReadStructure(content []string) (*Structure, bool) {
	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "struct" {
		_content := content[1 : len(content)-1]
		identifier := header["identifier"]

		return &Structure{_content: _content, identifier: identifier, properties: make(map[int]*Property), tag: make(map[string]int)}, true
	}

	return nil, false
}

func NewStructure(identifier string) *Structure {
	return &Structure{identifier: identifier, properties: make(map[int]*Property), tag: make(map[string]int)}
}

type Property struct {
	req          string
	propertyType string
	identifier   string
}

func (r *Property) Req() string {
	return r.req
}

func (r *Property) Identifier() string {
	return r.identifier
}

func (r *Property) Type() string {
	return r.propertyType
}

func (r *Property) SetType(propertyType string) {
	r.propertyType = propertyType
}

func (r *Structure) Identifier() string {
	return r.identifier
}

func (r *Structure) AddProperty(identifier string, propertyType string) {
	var property Property
	property.req = "optional"
	property.identifier = identifier
	property.propertyType = propertyType

	r._lastTag++
	r.properties[r._lastTag] = &property
	r.tag[identifier] = r._lastTag
}

func (r *Structure) FindProperty(identifier string) (*Property, bool) {
	if tag, isExists := r.tag[identifier]; isExists {
		return r.properties[tag], true
	}
	return nil, false
}

func (r *Structure) TagOf(identifier string) (int, bool) {
	tag, isExists := r.tag[identifier]
	return tag, isExists
}

func (r *Structure) Resolve() *Structure {

	for _, unparsedProperty := range r._content {

		unparsedProperty = strings.TrimSpace(unparsedProperty)
		match := Utils.ReSubMatchMap(__reProperty, unparsedProperty)

		var tag int
		var err error
		tag, err = strconv.Atoi(match["tag"])

		if err != nil {
			panic(err)
		}

		var req = match["req"]
		var identifier = match["identifier"]
		var propertyType = match["type"]

		r.properties[tag] = &Property{req: req, propertyType: propertyType, identifier: identifier}
		r.tag[identifier] = tag

		r._lastTag = Utils.Max(r._lastTag, tag)
	}

	return r
}

func (r *Structure) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "struct %s {\n", r.Identifier())

	var keys []int
	for k := range r.properties {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		property := r.properties[key]
		fmt.Fprintf(&buffer, "  %d: %s %s %s\n", key, property.Req(), property.Type(), property.Identifier())
	}
	fmt.Fprintf(&buffer, "}\n\n")

	return buffer.String()
}
