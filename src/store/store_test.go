package store

import (
	"bytes"
	"testing"
)

func TestFunc(t *testing.T) {
	var input = []byte("catchen")
	testfunc(input, t)
}
func testfunc(obj []byte, t *testing.T) {
	buf := encode(obj)
	t.Log(buf)
	obj2, err := decode(buf)
	if err != nil {
		t.Fatal(err)
		return
	} else {
		t.Log(obj2)
	}
	if !bytes.Equal(obj2, obj) {
		t.Fatal(string(obj2))
	}
	return
}
