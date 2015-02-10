package node

import (
	"unsafe"
)

import lgd "code.google.com/p/log4go"

var pageListSize int
var globalPageList pageList

type pageNode struct {
	next      *pageNode
	pagebyte  []byte
	freelist  *freeNode
	freecount int
	count     uint64
}

func (pn *pageNode) allocate(n int) *page {
	if pn.freecount < n {
		return nil
	}

	freenode := pn.freelist
	for freenode != nil {
		if freenode.rBound-freenode.lBound > pid(n) {
			id := freenode.lBound
			freenode.lBound -= pid(n)
			if freenode.lBound == freenode.rBound {
				freenode.pre.next = freenode.next
			}
			p := (*page)(unsafe.Pointer(&pn.pagebyte[int(id)*pageSize]))
			p.id = id
			p.count = uint16(n)
			pn.freecount -= n
			return p
		} else {
			freenode = freenode.next
		}
	}
	return nil
}

func (pn *pageNode) free(p *page) {
	id := p.id
	count := pid(p.count)
	pn.freecount += int(count)
	freenode := pn.freelist
	for true {
		if freenode.lBound > id {
			if freenode.lBound == id+count {
				if freenode.pre.rBound == id {
					freenode.pre.rBound = freenode.rBound
					freenode.pre.next = freenode.next
					return
				} else {
					freenode.lBound = id
					return
				}
			} else if freenode.lBound > id+count {
				if freenode.pre.rBound == id {
					freenode.pre.rBound = id + count
					return
				} else if freenode.pre.rBound > id {
					fn := new(freeNode)
					fn.lBound = id
					fn.lBound = id + count
					fn.next = freenode
					fn.pre = freenode.pre
					freenode.pre = fn
					return
				} else {
					lgd.Error("Bug!")
				}
			} else {
				lgd.Error("Bug!")
			}
		} else {
			if freenode.next != nil {
				freenode = freenode.next
			} else if freenode.rBound == id {
				freenode.rBound = id + count
				return
			} else {
				fn := new(freeNode)
				fn.lBound = id
				fn.rBound = id + count
				freenode.next = fn
				fn.pre = freenode
				return
			}
		}
	}
	return
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
	pn.freelist = new(freeNode)
	pn.freelist.lBound = pid(count)
	pn.freelist.lBound = pid(count + pn.count)
	pn.freecount = int(pn.count)
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

type freeNode struct {
	lBound pid
	rBound pid
	next   *freeNode
	pre    *freeNode
}

func init() {
	globalPageList = pageList{}
	pageListSize = 0
}
