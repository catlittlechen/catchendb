package node

import (
	"sort"
	"unsafe"
)

import lgd "code.google.com/p/log4go"

var pageListSize int

type pageNode struct {
	next     *pageNode
	pagebyte []byte
	freelist pids
	count    int
}

func (pn *pageNode) allocate(n int) *page {
	if len(pn.freelist) < n {
		return nil
	}

	for i, id := range pn.freelist {
		if id+pid(n) == pn.freelist[i+n] {
			copy(pn.freelist[i:], pn.freelist[i+n:])
			pn.freelist = pn.freelist[:len(pn.freelist)-n]
			//TODO
			var p *page
			p.id = id
			p.count = uint16(n)
			return p
		}
	}
	return nil
}

func (pn *pageNode) free(p *page) {
	id := p.id
	count := p.count
	var i uint16
	for i = 0; i < count; i = i + 1 {
		pn.freelist = append(pn.freelist, id+pid(i))
	}
	sort.Sort(pn.freelist)
}

func (pn *pageNode) mmap() bool {
	newPageSize := mmapSize(pageListSize)
	var err error
	pn.pagebyte, err = mmap(newPageSize)
	if err != nil {
		lgd.Error("mmap error %s", err)
		return false
	} else {
		pageListSize = newPageSize
	}
	pn.count = int(unsafe.Sizeof(pn.pagebyte)) / pageSize
	//TODO
	return false
}

type pageList struct {
	head *pageNode
	tail *pageNode
}

func (pl *pageList) allocate(n int) *page {
	for node := pl.head; node != nil; node = node.next {
		p := node.allocate(n)
		if p != nil {
			return p
		}
	}
	node := pageNode{}
	if node.mmap() {
		pl.tail.next = &node
		pl.tail = &node
		p := pl.allocate(n)
		if p != nil {
			return p
		}
	}
	return nil
}

func init() {
	pageListSize = 0
}
