package store

import (
	"bytes"
	"compress/zlib"
)

func encode(obj []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(obj)
	w.Close()
	return buf.Bytes()
}
