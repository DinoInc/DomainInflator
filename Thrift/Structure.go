package Thrift

import "fmt"
import "strconv"
import "strings"
import "regexp"
import "github.com/DinoInc/DomainInflator/Utils"

var _ = fmt.Println

type Structure struct {
	properties map[int]*Property
	order      map[string]int
	_internal  _internal_structure
}

type Property struct {
	req          string
	propertyType string
	identifier   string
}

type _internal_structure struct {
	identifier string
	_content   []string
}

var __reProperty = regexp.MustCompile(`(?P<order>[1-9][0-9]*):(\s*(?P<req>[A-z][A-z0-9]*))?(\s*(?P<type>[A-z][A-z0-9]*))(\s*(?P<identifier>[A-z][A-z0-9]*))`)

func (r *Structure) Resolve() *Structure {
	for _, unparsedProperty := range r._internal._content {

		unparsedProperty = strings.TrimSpace(unparsedProperty)
		match := Utils.ReSubMatchMap(__reProperty, unparsedProperty)

		var order int
		var err error
		order, err = strconv.Atoi(match["order"])

		if err != nil {
			panic(err)
		}

		var req = match["req"]
		var identifier = match["identifier"]
		var propertyType = match["type"]

		r.properties[order] = &Property{req: req, propertyType: propertyType, identifier: identifier}

	}

	return nil
}

func ReadStructure(content []string) (*Structure, bool) {
	header := Utils.ReSubMatchMap(__reHeader, content[0])

	if header["type"] == "struct" {
		var _internal _internal_structure
		_internal._content = content[1 : len(content)-1]
		_internal.identifier = header["identifier"]

		return &Structure{_internal: _internal, properties: make(map[int]*Property)}, true
	}

	return nil, false
}
