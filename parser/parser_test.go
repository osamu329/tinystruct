package parser

import (
	"testing"
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
