package Schema

import "io/ioutil"
import "encoding/json"
import "regexp"

type Ref struct {
	_resolved *Structure
	_internal _internal_ref
}

type _internal_ref struct {
	URI string `json:"$ref"`
}

func (r *Ref) Name() string {
	return regexp.MustCompile(`[a-zA-z][a-zA-z0-9]*$`).FindString(r._internal.URI)
}

func haveJSONSchemaExtension(s string) bool {
	haveJSONSchema, _ := regexp.MatchString(`\.schema\.json$`, s)
	return haveJSONSchema
}

func (r *Ref) File() string {

	file := regexp.MustCompile(`^[a-zA-z][a-zA-z0-9.]*`).FindString(r._internal.URI)
	if haveJSONSchemaExtension(file) {
		return file
	} else {
		return file + ".schema.json"
	}
}

func (r *Ref) IsSelf() bool {
	return (r._internal.URI[0] == '#')
}

func (r *Ref) Resolve(_schema *Schema) *Structure {
	if r._resolved != nil {
		return r._resolved
	}

	if !r.IsSelf() {

		data, err := ioutil.ReadFile(__baseDir + r.File())
		if err != nil {
			panic(err)
		}

		_dataRawMessage := json.RawMessage(data)
		schema, isSchema := ReadSchema(&_dataRawMessage)

		if !isSchema {
			panic("reference to non .schema.json file")
		}

		_schema = schema
	}

	r._resolved = _schema.Resolve(r.Name())
	r._resolved.SetIdentifier(r.Name())

	return r._resolved
}

func ReadRef(data *json.RawMessage) (*Ref, bool) {
	var _internal _internal_ref

	if json.Unmarshal(*data, &_internal) != nil {
		return nil, false
	}

	if _internal.URI != "" {
		return &Ref{_internal: _internal}, true
	}

	return nil, false
}
