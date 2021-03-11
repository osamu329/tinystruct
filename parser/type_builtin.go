package parser

import (
    "encoding/binary"
    "fmt"
    "io"
)

func irange(min, max int64) func(interface{})bool {
    return func(v interface{}) bool {
        i, ok := v.(int64); return ok && min <= i && i<= max
    }
}

func urange(max uint64) func(interface{})bool {
    return func(v interface{}) bool {
        i, ok := v.(uint64); return ok && i<= max
    }
}

func FromTypeInfo(tinfo *TypeInfo) *Type {
    return &Type{Name:tinfo.name, Extra:tinfo}
}

func charTypeCheck(v interface{}) bool {
    _, ok := v.(rune)
    return ok
}

type PackFunc func(io.Writer, ByteOrder, interface{}) error

func packHelper_(conv func(interface{}) interface{}) PackFunc {
    return func (w io.Writer, order ByteOrder, v interface{}) error {
        if err := binary.Write(w, order, conv(v)); err != nil {
            return err
        }
        return nil
    }
}

func packHelper(conv func(int) interface{}, min, max int) PackFunc {
    return func (w io.Writer, order ByteOrder, value interface{}) error {
        i, ok := value.(int);
        if !ok {
            return fmt.Errorf("value %d shoud be int", value)
        }
        if min <= i && i <= max {
        } else {
            return fmt.Errorf("value %d shoud be in range[%d,%d]", value, min, max)
        }
        if err := binary.Write(w, order, conv(i)); err != nil {
            return err
        }
        return nil
    }
}

var (
    Char = FromTypeInfo(&TypeInfo{name:"char", size:1, typecheck: charTypeCheck})

    Uint8 = FromTypeInfo(&TypeInfo{
        size:1, name:"uint8_t", typecheck: urange(0xff),
        pack: packHelper(func(v int) interface{} { return uint8(v)}, 0, 0xff),})

    Uint16 = FromTypeInfo(&TypeInfo{ size:2, name:"uint16_t", typecheck: urange(0xffff),
        pack: packHelper(func(v int) interface{} { return uint16(v)}, 0, 0xffff)},)
    Uint32 = FromTypeInfo(&TypeInfo{ size:4, name:"uint32_t", typecheck: urange(0xffffffff),
        pack: packHelper(func(v int) interface{} { return uint32(v)}, 0, 0xffffff)},)
    Uint64 = FromTypeInfo(&TypeInfo{ size:8, name:"uint64_t", typecheck: urange(0xffffffffffffffff),
        pack: func (w io.Writer, order ByteOrder, value interface{}) error {
            u64, ok := value.(uint64);
            if !ok {
                i, ok := value.(int)
                if !ok {
                    return fmt.Errorf("value %d shoud be uint64 or int, but %T", value, value)
                } else {
                    return binary.Write(w, order, uint64(i))
                }
            }
            return binary.Write(w, order, u64)
        },
    })
    Int8  = FromTypeInfo(&TypeInfo{ size:1, name:"int8_t", typecheck: irange(-0x7f-1, 0x7f),
        pack: packHelper(func(v int) interface{} { return int8(v)}, -0x7f-1, 0xff)},)
    Int16  = FromTypeInfo(&TypeInfo{ size:2, name:"int16_t", typecheck: irange(-0x7fff-1, 0x7fff),
        pack: packHelper(func(v int) interface{} { return int16(v)}, -0x7fff-1, 0xff)},)
    Int32  = FromTypeInfo(&TypeInfo{ size:4, name:"int32_t", typecheck: irange(-0x7fffffff-1, 0x7fffffff),
        pack: packHelper(func(v int) interface{} { return int32(v)}, -0x7fffffff-1, 0xff)},)
    Int64 = FromTypeInfo(&TypeInfo{ size:8, name:"int64_t", typecheck: irange(-0x7fffffffffffffff-1, 0x7fffffffffffffff),
        pack: packHelper(func(v int) interface{} { return int64(v)}, -0x7fffffffffffffff-1, 0x7fffffffffffffff)},)
)

var primitiveTypes []*Type = []*Type{
    Char,Uint8, Uint16, Uint32, Uint64, Int8, Int16, Int32, Int64,
}

