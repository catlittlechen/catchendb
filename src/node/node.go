package node

import (
	"unsafe"
)

const (
	nodePageSize = uint32(unsafe.Sizeof(nodePageElem{}))
)

type nodePageElem struct {
	lChird    *nodePageElem
	rChird    *nodePageElem
	keySize   uint32
	valueSize uint32
}

func (n *nodePageElem) key() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize : nodePageSize+n.keySize]

}

func (n *nodePageElem) value() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize+n.keySize : nodePageSize+n.keySize+n.valueSize]

}
