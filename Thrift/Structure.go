package Thrift

import "fmt"
import "strconv"
import "strings"
import "regexp"
import "github.com/DinoInc/DomainInflator/Utils"

var _ = fmt.Println

type Structure struct {
	properties map[int]*Property
	index      map[string]int
	_internal  _internal_structure
}

type Property struct {
	req          string
	propertyType string
	identifier   string
}

type _internal_structure struct {
	identifier string
	_lastIndex int
	_content   []string
}

var __reProperty = regexp.MustCompile(`(?P<comment>//)?(?P<index>[1-9][0-9]*):(\s*(?P<req>[A-z][A-z0-9]*))?(\s*(?P<type>[A-z][A-z0-9]*))(\s*(?P<identifier>[A-z][A-z0-9]*))`)

func (r *Structure) Identifier() string {
	return r._internal.identifier
}

func (r *Structure) Resolve() *Structure {

	for _, unparsedProperty := range r._internal._content {

		unparsedProperty = strings.TrimSpace(unparsedProperty)
		match := Utils.ReSubMatchMap(__reProperty, unparsedProperty)

		var index int
		var err error
		index, err = strconv.Atoi(match["index"])

		if err != nil {
			panic(err)
		}

		var req = match["req"]
		var identifier = match["identifier"]
		var propertyType = match["type"]

		r.properties[index] = &Property{req: req, propertyType: propertyType, identifier: identifier}
		r.index[identifier] = index

		r._internal._lastIndex = Utils.Max(r._internal._lastIndex, index)
	}

	return r
}

func ReadStructure(content []string) (*Structure, bool) {
	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "struct" {
		var _internal _internal_structure
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Structure{_internal: _internal, properties: make(map[int]*Property), index: make(map[string]int)}, true
	}

	return nil, false
}

func NewStructure(identifier string) *Structure {
	var _internal _internal_structure
	_internal.identifier = identifier

	return &Structure{_internal: _internal, properties: make(map[int]*Property), index: make(map[string]int)}
}

func (r *Structure) AddProperty(identifier string, propertyType string) {
	var property Property
	property.identifier = identifier
	property.propertyType = propertyType

	r._internal._lastIndex++
	r.properties[r._internal._lastIndex] = &property
	r.index[identifier] = r._internal._lastIndex
}

func (r *Structure) IndexOf(identifier string) (int, bool) {
	index, isExists := r.index[identifier]
	return index, isExists
}
