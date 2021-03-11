package parser

import (
    "testing"
    _ "log"
    "bytes"
    "github.com/stretchr/testify/assert"
    "encoding/binary"
)

func TestParseTypedef(t *testing.T) {
    _, err := ParseTypedef("testdata/typedef.yaml")
    if err != nil {
        t.Fatal(err)
    }
    //log.Printf("%#v", root)
}

func TestPack(t *testing.T) {
    testcases := []struct {
        name string
        typ *Type
        value interface{}
        expected []byte
    } {
        {name: "uint8(0)", typ:Uint8, value:0, expected:[]byte{0}},
        {name: "uint8(0xff)", typ:Uint8, value:0xff, expected:[]byte{0xff}},
        {name: "int8(0xff)", typ:Int8, value:-1, expected:[]byte{0xff}},
        {name: "uint64(0xff)", typ:Int64, value:0x7fff_ffff_ffff_ffff,
            expected:[]byte{0xff,0xff,0xff,0xff, 0xff,0xff,0xff,0x7f,}},
        {name: "uint64(0)", typ:Uint64, value:0,
            expected:[]byte{0,0,0,0, 0,0,0,0,}},
        {name: "uint64(0xff)", typ:Uint64, value:uint64(0xffff_ffff_ffff_ffff),
            expected:[]byte{0xff,0xff,0xff,0xff, 0xff,0xff,0xff,0xff,}},
    }
    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            w := &bytes.Buffer{}
            assert.Nil(t, tc.typ.Pack(w, binary.LittleEndian, tc.value))
            assert.Equal(t, tc.expected, w.Bytes())
        })
    }
}
