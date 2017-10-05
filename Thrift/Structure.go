package Thrift

type ThriftEnum struct {
	Items []string
}

type ThriftStructure struct {
	Identifier  string
	Description string            `json:"description,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}
