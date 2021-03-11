package parser

import (
	"testing"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "github.com/stretchr/testify/assert"
    "encoding/binary"
    "bytes"
)

func TestParse(t *testing.T) {
	fname := "testdata/source.h"
	f, err := ParseFile(fname)
	if err != nil {
		t.Fatalf("Parse %s failed, err:%s", fname, err)
	}
	if f == nil {
		t.Fatalf("f is nil")
	}
}

func TestYaml(t *testing.T) {
    // uint64 でしか表現できない範囲だと、uint64
    // それ以外では int になる
    yaml_data, err := ioutil.ReadFile("testdata/test.yaml")
    assert.Equal(t, nil, err)
    dict := make(map[string]interface{})
    err = yaml.Unmarshal(yaml_data, &dict)
    assert.Equal(t, nil, err)
    t.Logf("%v a:%T b:%T", dict, dict["a"].(uint64), dict["b"].(int))
}

func TestBinary(t *testing.T) {
    w := &bytes.Buffer{}
    var v interface{} = byte(0)
    assert.Nil(t, binary.Write(w, binary.LittleEndian, v))
    assert.Equal(t, 1, len(w.Bytes()))
}
