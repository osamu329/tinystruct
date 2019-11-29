package syntax

type Typedef struct {
	Name string
	Type Type
}

type ArrayType struct {
	Base Type
	Len  int
}

type PrimitiveType struct {
	Name string
	Size int
}

var types = map[string]Type{
	"int":      &PrimitiveType{Name: "int", Size: 1},
	"char":     &PrimitiveType{Name: "byte", Size: 1},
	"uint8_t":  &PrimitiveType{Name: "uint8", Size: 1},
	"int8_t":   &PrimitiveType{Name: "int8", Size: 1},
	"uint16_t": &PrimitiveType{Name: "uint16", Size: 2},
	"int16_t":  &PrimitiveType{Name: "int16", Size: 2},
	"uint32_t": &PrimitiveType{Name: "uint32", Size: 4},
	"int32_t":  &PrimitiveType{Name: "int32", Size: 4},
	"int64_t":  &PrimitiveType{Name: "int64", Size: 8},
	"uint64_t": &PrimitiveType{Name: "uint64", Size: 8},
}

type Struct struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name string
	Type Type
}

type Type interface {
}
