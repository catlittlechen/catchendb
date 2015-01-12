package store

import (
	"fmt"
	"github.com/cupcake/rdb"
	"github.com/cupcake/rdb/nopdecoder"
)

func decode(obj []byte) (interface{}, error) {
	d := &decoder{}
	if err := rdb.DecodeDump(obj, 0, nil, 0, d); err != nil {
		return nil, err
	}
	return d.obj, d.err
}

type decoder struct {
	nopdecoder.NopDecoder
	obj interface{}
	err error
}

func (d *decoder) Set(key, value []byte, expire int64) {
	d.initObject([]byte(value))
}

func (d *decoder) initObject(obj interface{}) {
	if d.err != nil {
		return
	}

	if d.obj != nil {
		d.err = fmt.Errorf("invalid object, init again")
	} else {
		d.obj = obj
	}
}
