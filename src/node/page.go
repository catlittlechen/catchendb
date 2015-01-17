package node

import (
	"unsafe"
)

type page struct {
	id    uint64
	flag  uint8
	count uint16
	ptr   uintptr
}

const (
	nodePageFlag = 0x01
	listPageFlag = 0x02
)

const (
	pageHeaderSize = int(unsafe.Offsetof(((*page)(nil)).ptr))
)

func (p *page) nodePageElem() (n *nodePageElem) {
	n = (*nodePageElem)(unsafe.Pointer(&p.ptr))
	return
}

func (p *page) listPageElem() (l *listPageElem) {
	l = (*listPageElem)(unsafe.Pointer(&p.ptr))
	return
}
