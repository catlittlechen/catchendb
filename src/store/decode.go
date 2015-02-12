package store

import (
	"bytes"
	"compress/zlib"
)

import lgd "code.google.com/p/log4go"

func decode(buf []byte) ([]byte, error) {
	var obj bytes.Buffer
	b := bytes.NewBuffer(buf)
	r, err := zlib.NewReader(b)
	if err != nil {
		lgd.Error("zlib fail! error[%s] buf[%s]", err, buf)
		return nil, err
	}
	_, err = obj.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return obj.Bytes(), nil
}
