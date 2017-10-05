package Engine

import (
	"github.com/DinoInc/DomainInflator/Schema"
	"github.com/DinoInc/DomainInflator/Thrift"
)

type Engine struct {
	SchemaIDL  *Schema.SchemaIDL
	ThriftIDL  ThriftIDL
	deviations []*Deviation
}

func NewEngine(baseDir string) *Engine {
	e := Engine{}

	e.SchemaIDL = Schema.NewSchemaIDL(baseDir)

	e.ThriftIDL.Enum = make(map[string]*Thrift.Enum)
	e.ThriftIDL.Structure = make(map[string]*Thrift.Structure)

	e.deviations = make([]*Deviation, 0)

	return &e
}
