package store

import (
	"bytes"
	"compress/zlib"
)

import lgd "catchendb/src/log"

func decode(buf []byte) ([]byte, error) {
	var obj bytes.Buffer
	b := bytes.NewBuffer(buf)
	r, err := zlib.NewReader(b)
	if err != nil {
		lgd.Errorf("zlib fail! error[%s] buf[%s]", err, buf)
		return nil, err
	}
	_, err = obj.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return obj.Bytes(), nil
}
