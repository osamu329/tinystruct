package syntax
import (
    "fmt"
)

type Type interface {
    Name() string
	IsArray() bool
}

type Struct struct {
	name   string
	Fields []*Field
}

func (t *Struct) Name() string {
    return t.name
}

func (t *Struct) IsArray() bool {
	return false
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

func (t *Typedef) IsArray() bool {
	return false
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

func (t *ArrayType) IsString() bool {
	return t.Base.Name() == "char"
}

func (t *ArrayType) IsArray() bool {
	return true
}

type PrimitiveType struct {
	name string
	size int
}

func (t *PrimitiveType) Name() string {
    return t.name
}

func (t *PrimitiveType) IsArray() bool {
	return false
}

var types = map[string]Type{
	"int":      &PrimitiveType{name: "int", size: 1},
	"char":     &PrimitiveType{name: "char", size: 1},
	"uint8_t":  &PrimitiveType{name: "uint8_t", size: 1},
	"int8_t":   &PrimitiveType{name: "int8_t", size: 1},
	"uint16_t": &PrimitiveType{name: "uint16_t", size: 2},
	"int16_t":  &PrimitiveType{name: "int16_t", size: 2},
	"uint32_t": &PrimitiveType{name: "uint32_t", size: 4},
	"int32_t":  &PrimitiveType{name: "int32_t", size: 4},
	"int64_t":  &PrimitiveType{name: "int64_t", size: 8},
	"uint64_t": &PrimitiveType{name: "uint64_t", size: 8},
}

