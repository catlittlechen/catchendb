package node

import (
	"sort"
	"unsafe"
)

import lgd "code.google.com/p/log4go"

//TODO freelist需要改变策略

var pageListSize int
var globalPageList pageList

type pageNode struct {
	next     *pageNode
	pagebyte []byte
	freelist pids
	count    uint64
}

func (pn *pageNode) allocate(n int) *page {
	if len(pn.freelist) < n {
		return nil
	}

	for i, id := range pn.freelist {
		if id+pid(n) == pn.freelist[i+n] {
			copy(pn.freelist[i:], pn.freelist[i+n:])
			pn.freelist = pn.freelist[:len(pn.freelist)-n]
			p := (*page)(unsafe.Pointer(&pn.pagebyte[int(id)*pageSize]))
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

func (pn *pageNode) mmap(count uint64) bool {
	newPageListSize := mmapSize(pageListSize)
	var err error
	pn.pagebyte, err = mmap(newPageListSize)
	if err != nil {
		lgd.Error("mmap error %s", err)
		return false
	} else {
		pageListSize = newPageListSize
	}
	pn.count = uint64(unsafe.Sizeof(pn.pagebyte)) / uint64(pageSize)
	var i uint64
	for i = 0; i < pn.count; i++ {
		pn.freelist = append(pn.freelist, pid(i+count))
	}
	return true
}

type pageList struct {
	head *pageNode
	tail *pageNode
}

func (pl *pageList) allocate(n int) *page {
	var count uint64
	node := pl.head
	for node != nil {
		p := node.allocate(n)
		if p != nil {
			return p
		} else {
			count += node.count
		}
		node = node.next
	}

	node = &pageNode{}
	if node.mmap(count) {
		if pl.head == nil {
			pl.head = node
		}
		if pl.tail != nil {
			pl.tail.next = node
		}
		pl.tail = node
		p := node.allocate(n)
		if p != nil {
			return p
		}
	}
	return nil
}

func (pl *pageList) free(p *page) {
	var count uint64
	for node := pl.head; node != nil; node = node.next {
		count += node.count
		if p.id < pid(count) {
			node.free(p)
		}
	}
}

func init() {
	globalPageList = pageList{}
	pageListSize = 0
}
