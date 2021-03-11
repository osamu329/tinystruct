package parser

import (
    "encoding/binary"
    "io"
)

type ByteOrder = binary.ByteOrder
var LittleEndian = binary.LittleEndian
type Type struct {
    Name string
    Extra TypeBase
}

func (t *Type) TypeCheck(v interface{}) bool {
    return t.Extra.TypeCheck(v)
}

func (t *Type) Size() int {
    return t.Extra.Size()
}

func (t *Type) Pack(w io.Writer, order ByteOrder, v interface{}) error {
    return t.Extra.Pack(w, order, v)
}

type TypeBase interface {
    TypeCheck(v interface{}) bool
    Pack(w io.Writer, order ByteOrder, v interface{}) error
    Size() int
}

type Field struct {
    Name string
    Type *Type
}

type Struct struct {
    name string
    fields []*Field
    field_map map[string]*Field
    size int
}

func NewStruct(name string, fields []*Field) *Struct {
    st := &Struct{name:name, fields:fields, field_map:make(map[string]*Field), size:0}
    for _, f := range fields {
        st.field_map[f.Name] = f
        st.size += f.Type.Size()
    }
    return st
}

func (s *Struct) Name() string {
    return s.name
}

func (s *Struct) Size() int {
    return s.size
}

func (s *Struct) FieldByName(name string) *Field {
    return s.field_map[name]
}

func (s *Struct) TypeCheck(v interface{}) bool {
    dict, ok := v.(map[string]interface{})
    if !ok {
        return false
    }
    for fieldname, fieldvalue := range dict {
        field := s.FieldByName(fieldname)
        if field == nil {
            return false
        }
        if !field.Type.TypeCheck(fieldvalue) {
            return false
        }
    }
    return true
}

func (s *Struct) Pack(w io.Writer, order ByteOrder, v interface{}) error {
    dict, ok := v.(map[string]interface{})
    if !ok {
        _, err := w.Write(make([]byte, s.Size()))
        return err
    }
    for _, field := range s.fields {
        if value, ok := dict[field.Name]; !ok {
            _, err := w.Write(make([]byte, field.Type.Size()))
            return err
        } else {
            return field.Type.Pack(w, order, value)
        }
    }
    return nil
}

type TypeInfo struct {
    name string
    size int
    typecheck func(v interface{}) bool
    pack PackFunc
}

func (t *TypeInfo) TypeCheck(v interface{}) bool {
    return t.typecheck(v)
}

func (t *TypeInfo) Size() int {
    return t.size
}

func (t *TypeInfo) Pack(w io.Writer, order ByteOrder, v interface{}) error {
    return t.pack(w, order, v)
}

