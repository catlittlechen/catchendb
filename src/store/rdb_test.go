package store

import (
	"reflect"
	"testing"
)

func TestFunc(t *testing.T) {
	testfunc([]byte("asd"), t)
}
func testfunc(obj interface{}, t *testing.T) {
	b := encode(obj)
	if b == nil {
		t.Fatal("obj is null")
	}
	obj2, err := decode(b)
	if err != nil {
		t.Fatal(err)
		return
	}
	if !reflect.DeepEqual(obj, obj2) {
		t.Fatal("translate bytes error")
	}
}
