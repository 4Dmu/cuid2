package cuid2

import (
	"testing"
)

func TestCuid2(t *testing.T) {
	gen := New(GenOpts{})
	id1 := gen.Cuid2()
	id2 := gen.Cuid2()

	if !IsCuid(id1, gen.length, bigLength) {
		t.Fatalf("(%s) is not a valid cuid", id1)
	}

	if !IsCuid(id2, gen.length, bigLength) {
		t.Fatalf("(%s) is not a valid cuid", id2)
	}

	if id1 == id2 {
		t.Fatal("ids are the same")
	}
}
