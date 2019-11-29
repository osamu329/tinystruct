package syntax
import (
    "fmt"
)

type Type interface {
    Name() string
}

type Struct struct {
	name   string
	Fields []*Field
}

func (t *Struct) Name() string {
    return t.name
}


func NewStruct(name string, fields []*Field) *Struct {
    return &Struct{name:name, Fields:fields}
}

type Field struct {
	Name string
	Type Type
}


type Typedef struct {
	name string
	typ Type
}

func NewTypedef(name string, t Type) *Typedef {
    return &Typedef{name:name, typ:t}
}

func (t *Typedef) Name() string {
    return t.name
}

func (t *Typedef) BaseType() Type {
    return t.typ
}

type ArrayType struct {
	Base Type
	Len  int
}

func (t *ArrayType) Name() string {
    return fmt.Sprintf("[%d]%s", t.Len, t.Base.Name())
}

type PrimitiveType struct {
	name string
	size int
}

func (t *PrimitiveType) Name() string {
    return t.name
}

var types = map[string]Type{
	"int":      &PrimitiveType{name: "int", size: 1},
	"char":     &PrimitiveType{name: "byte", size: 1},
	"uint8_t":  &PrimitiveType{name: "uint8", size: 1},
	"int8_t":   &PrimitiveType{name: "int8", size: 1},
	"uint16_t": &PrimitiveType{name: "uint16", size: 2},
	"int16_t":  &PrimitiveType{name: "int16", size: 2},
	"uint32_t": &PrimitiveType{name: "uint32", size: 4},
	"int32_t":  &PrimitiveType{name: "int32", size: 4},
	"int64_t":  &PrimitiveType{name: "int64", size: 8},
	"uint64_t": &PrimitiveType{name: "uint64", size: 8},
}

