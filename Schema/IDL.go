package Schema

type PropertyEnum struct {
	Identifier string
	Items      []string
}

type PropertyList struct {
	ElementType string
}

type Structure struct {
	Identifier  string
	Description string
	Properties  map[string]interface{}
}
