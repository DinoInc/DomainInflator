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
	_internal  _internal_structure
}

type _internal_structure struct {
	identifier string
	_lastTag   int
	_content   []string
}

var __reProperty = regexp.MustCompile(`(?P<comment>//)?(?P<tag>[1-9][0-9]*):(\s*(?P<req>[A-z][A-z0-9]*))?(\s*(?P<type>[A-z][A-z0-9]*))(\s*(?P<identifier>[A-z][A-z0-9]*))`)

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

func (r *Structure) Identifier() string {
	return r._internal.identifier
}

func (r *Structure) Properties() map[int]*Property {
	return r.properties
}

func (r *Structure) TagOf(identifier string) (int, bool) {
	tag, isExists := r.tag[identifier]
	return tag, isExists
}

func (r *Structure) Resolve() *Structure {

	for _, unparsedProperty := range r._internal._content {

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

		r._internal._lastTag = Utils.Max(r._internal._lastTag, tag)
	}

	return r
}

func ReadStructure(content []string) (*Structure, bool) {
	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "struct" {
		var _internal _internal_structure
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Structure{_internal: _internal, properties: make(map[int]*Property), tag: make(map[string]int)}, true
	}

	return nil, false
}

func NewStructure(identifier string) *Structure {
	var _internal _internal_structure
	_internal.identifier = identifier

	return &Structure{_internal: _internal, properties: make(map[int]*Property), tag: make(map[string]int)}
}

func (r *Structure) AddProperty(identifier string, propertyType string) {
	var property Property
	property.req = "optional"
	property.identifier = identifier
	property.propertyType = propertyType

	r._internal._lastTag++
	r.properties[r._internal._lastTag] = &property
	r.tag[identifier] = r._internal._lastTag
}

func (r *Structure) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "struct %s {\n", r.Identifier())
	var properties = r.Properties()

	var keys []int
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		property := properties[key]
		fmt.Fprintf(&buffer, "\t%d: %s %s %s\n", key, property.Req(), property.Type(), property.Identifier())
	}
	fmt.Fprintf(&buffer, "}\n\n")

	return buffer.String()
}
