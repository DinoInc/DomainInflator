package Schema

import "encoding/json"
import "regexp"

type Ref struct {
	URI string `json:"$ref"`
}

func (r *Ref) Name() string {
	return regexp.MustCompile(`[a-zA-z]*$`).FindString(r.URI)
}

func (r *Ref) File() string {
	return regexp.MustCompile(`^[a-zA-z.]*`).FindString(r.URI)
}

func (r *Ref) IsSelf() bool {
	return (r.URI[0] == '#')
}

func ReadRef(data *json.RawMessage) (*Ref, bool) {
	var ref Ref

	if json.Unmarshal(*data, &ref) != nil {
		return nil, false
	}

	if ref.URI != "" {
		return &ref, true
	}

	return nil, false
}
