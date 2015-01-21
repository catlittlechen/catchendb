package node

import (
	"unsafe"
)

type page struct {
	id    pid
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

func (p *page) pagelist() (l *pagelist) {
	l = (*pagelist)(unsafe.Pointer(&p.ptr))
	return
}

type pid uint64

type pids []pid

func (p pids) Len() int {
	return len(p)
}

func (p pids) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p pids) Less(i, j int) bool {
	return p[i] < p[j]
}
