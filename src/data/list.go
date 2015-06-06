package data

import (
	"sync"
	"unsafe"
)

import lgd "code.google.com/p/log4go"

var pageListSize int
var globalPageList pageList

type pageNode struct {
	next      *pageNode
	pagebyte  *[mmapBranch]byte
	freelist  *freeNode
	freecount int
	fid       pid
	count     uint64
}

func (pn *pageNode) allocate(n int) *page {
	if pn.freecount < n {
		lgd.Trace("freecount[%d/%d] limit. ", pn.freecount, pn.count)
		return nil
	}

	freenode := pn.freelist
	for freenode != nil {
		if freenode.rBound-freenode.lBound > pid(n) {
			id := freenode.lBound
			freenode.lBound += pid(n)
			if freenode.lBound == freenode.rBound {
				freenode.pre.next = freenode.next
			}
			p := (*page)(unsafe.Pointer(&pn.pagebyte[int(id-pn.fid)*pageSize]))
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
	for {
		if freenode.lBound > id {
			if freenode.lBound == id+count {
				if freenode.pre != nil && freenode.pre.rBound == id {
					freenode.pre.rBound = freenode.rBound
					freenode.pre.next = freenode.next
					return
				} else {
					freenode.lBound = id
					return
				}
			} else if freenode.lBound > id+count {
				if freenode.pre != nil {
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
					fn := new(freeNode)
					fn.lBound = id
					fn.lBound = id + count
					fn.next = freenode
					pn.freelist = fn
					return
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
	pb, err := mmap(newPageListSize)
	if err != nil {
		lgd.Error("mmap error %s, pageListSize %d", err, newPageListSize)
		return false
	} else {
		pn.pagebyte = (*[mmapBranch]byte)(unsafe.Pointer(&pb[0]))
		pageListSize = newPageListSize
	}
	pn.count = uint64(newPageListSize) / uint64(pageSize)
	pn.freelist = new(freeNode)
	pn.freelist.lBound = pid(count)
	pn.freelist.rBound = pid(count + pn.count)
	pn.fid = pid(count)
	pn.freecount = int(pn.count)
	return true
}

type pageList struct {
	head  *pageNode
	tail  *pageNode
	mutex *sync.Mutex
}

func (pl *pageList) allocate(n int) *page {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
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
		return p
	}
	return nil
}

func (pl *pageList) free(p *page) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
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
	globalPageList = pageList{
		mutex: new(sync.Mutex),
	}
	pageListSize = 0
}
