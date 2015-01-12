package store

import (
	"bytes"
	"github.com/cupcake/rdb"
)

import lgd "code.google.com/p/log4go"

func encode(obj interface{}) []byte {
	var buf bytes.Buffer
	robj := rdb.NewEncoder(&buf)

	switch v := obj.(type) {
	case []byte:
		robj.EncodeType(rdb.TypeString)
		robj.EncodeString(v)
	default:
		lgd.Error("invalid type %T", obj)
		return nil
	}

	robj.EncodeDumpFooter()
	return buf.Bytes()
}
