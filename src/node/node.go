package node

import (
	"bytes"
	"time"
	"unsafe"
)

import lgd "code.google.com/p/log4go"

const (
	nodePageSize = int(unsafe.Sizeof(nodePageElem{}))

	NODE_KEY_SMALL = -1
	NODE_KEY_EQUAL = 0
	NODE_KEY_LARGE = 1
)

var (
	nowTime  int64
	treeRoot *nodeRoot
	channel  chan []byte
)

type nodeRoot struct {
	node *nodePageElem
}

func (nr *nodeRoot) input(line []byte) bool {
	d := data{}
	if !d.decode(line) {
		return false
	}
	go nr.insertNode(d.Key, d.Value, d.StartTime, d.EndTime)
	return true
}

func (nr *nodeRoot) output(channe chan []byte, sign []byte) {
	channel = channe
	nr.preorder(nr.node)
	channel <- sign
}

func (nr *nodeRoot) preorder(node *nodePageElem) {

	if node != nil {
		d := new(data)
		d.Key = string(node.key())
		d.Value = string(node.value())
		d.StartTime = node.getStartTime()
		d.EndTime = node.getEndTime()
		datastr, _ := d.encode()
		channel <- datastr
		nr.preorder(node.lChild)
		nr.preorder(node.rChild)
	}
}

func (nr *nodeRoot) first(nodeIndex *nodePageElem) (node *nodePageElem) {
	node = nodeIndex

	if node != nil {
		for node.lChild != nil {
			node = node.lChild
		}
	}
	return
}

func (nr *nodeRoot) last(nodeIndex *nodePageElem) (node *nodePageElem) {
	node = nodeIndex

	if node != nil {
		for node.rChild != nil {
			node = node.rChild
		}
	}
	return
}

func (nr *nodeRoot) next(nodeIndex *nodePageElem) (node *nodePageElem) {

	node = nodeIndex
	if node.rChild != nil {
		node = nr.first(node.rChild)
		return
	}

	nodeTmp := node.parent
	for nodeTmp != nil && node == nodeTmp.rChild {
		node = nodeTmp
		nodeTmp = nodeTmp.parent
	}
	node = nodeTmp
	return
}

func (nr *nodeRoot) pre(nodeIndex *nodePageElem) (node *nodePageElem) {

	node = nodeIndex
	if node.lChild != nil {
		node = (nr).last(node.lChild)
		return
	}

	nodeTmp := node.parent
	for nodeTmp != nil && node == nodeTmp.lChild {
		node = nodeTmp
		nodeTmp = nodeTmp.parent
	}
	node = nodeTmp
	return
}

func (nr *nodeRoot) search(key string) (node *nodePageElem) {
	node = nr.node

	for node != nil {
		switch node.compareKey(key) {
		case NODE_KEY_EQUAL:
			if node.isEnd() {
				nr.delet(node)
				return nil
			} else {
				return node
			}
		case NODE_KEY_SMALL:
			node = node.rChild
		case NODE_KEY_LARGE:
			node = node.lChild
		}
	}
	return nil
}

func (nr *nodeRoot) searchNode(key string) (value string, start, end int64) {
	node := nr.search(key)
	if node != nil {
		return string(node.value()), node.getStartTime(), node.getEndTime()
	}
	return
}

func (nr *nodeRoot) setStartTime(key string, start int64) bool {
	if node := nr.search(key); node != nil {
		return node.setStartTime(start)
	}
	return false
}

func (nr *nodeRoot) setEndTime(key string, end int64) bool {
	if node := nr.search(key); node != nil {
		return node.setEndTime(end)
	}
	return false
}

func (nr *nodeRoot) leftRotate(node *nodePageElem) {

	nodeTmp := node.rChild
	node.rChild = nodeTmp.lChild
	if nodeTmp.lChild != nil {
		nodeTmp.lChild.parent = node
	}
	nodeTmp.parent = node.parent
	if node.parent == nil {
		nr.node = nodeTmp
	} else {
		if node.parent.lChild == node {
			node.parent.lChild = nodeTmp
		} else {
			node.parent.rChild = nodeTmp
		}
	}

	nodeTmp.lChild = node
	node.parent = nodeTmp
}

func (nr *nodeRoot) rightRotate(node *nodePageElem) {

	nodeTmp := node.lChild
	node.lChild = nodeTmp.rChild
	if nodeTmp.rChild != nil {
		nodeTmp.rChild.parent = node
	}
	nodeTmp.parent = node.parent
	if node.parent == nil {
		nr.node = nodeTmp
	} else {
		if node == node.parent.rChild {
			node.parent.rChild = nodeTmp
		} else {
			node.parent.lChild = nodeTmp
		}
	}

	nodeTmp.rChild = node
	node.parent = nodeTmp
}

func (nr *nodeRoot) insertFixTree(node *nodePageElem) {
	var parent, grandparent, uncle *nodePageElem

	for node.parent != nil && node.parent.isRed() {
		parent = node.parent
		grandparent = parent.parent

		if parent == grandparent.lChild {
			uncle = grandparent.rChild
			if uncle != nil && uncle.isRed() {
				uncle.setBlack()
				parent.setBlack()
				grandparent.setRed()
				node = grandparent
				continue
			}
			if parent.rChild == node {
				nr.leftRotate(parent)
				parent, node = node, parent
			}
			parent.setBlack()
			grandparent.setRed()
			nr.rightRotate(grandparent)
		} else {
			uncle = grandparent.lChild
			if uncle != nil && uncle.isRed() {
				uncle.setBlack()
				parent.setBlack()
				grandparent.setRed()
				node = grandparent
				continue
			}
			if parent.lChild == node {
				nr.rightRotate(parent)
				parent, node = node, parent
			}
			parent.setBlack()
			grandparent.setRed()
			nr.leftRotate(grandparent)
		}
	}

	nr.node.setBlack()
}

func (nr *nodeRoot) insert(node *nodePageElem) {
	var nodeY *nodePageElem
	nodeX := nr.node

	for nodeX != nil {
		nodeY = nodeX
		if node.compare(nodeX) {
			nodeX = nodeX.lChild
		} else {
			nodeX = nodeX.rChild
		}
	}
	node.parent = nodeY

	if nodeY != nil {
		if node.compare(nodeY) {
			nodeY.lChild = node
		} else {
			nodeY.rChild = node
		}
	} else {
		nr.node = node
	}

	node.setRed()
	nr.insertFixTree(node)
}

func (nr *nodeRoot) createNode(key, value string, startTime, endTime int64, parent, lChild, rChild *nodePageElem) (node *nodePageElem) {
	size := pageHeaderSize + nodePageSize + len(key) + len(value)
	lgd.Trace("size[%d] pagecount[%d]", size, int(size/pageSize)+1)
	page := globalPageList.allocate(int(size/pageSize) + 1)
	if page != nil {
		node = page.nodePageElem()
		node.lChild = lChild
		node.rChild = rChild
		node.parent = parent
		if !node.setTime(startTime, endTime) {
			node.free()
			return nil
		}
		node.keySize = len(key)
		node.valueSize = len(value)
		node.setKeyValue(key, value)
		return node
	}

	return nil
}

func (nr *nodeRoot) insertNode(key, value string, startTime, endTime int64) bool {
	nowTime = time.Now().Unix()
	if endTime != 0 && endTime < nowTime {
		return true
	}

	node := nr.search(key)
	if node != nil {
		if node.setValue(value) {
			if !node.setTime(startTime, endTime) {
				return false
			}
			return true
		}
		nodeTmp := nr.createNode(key, value, startTime, endTime, node.parent, node.lChild, node.rChild)
		if node == nil {
			lgd.Error("reset value fail!")
			return false
		} else {
			node.lChild.parent = nodeTmp
			node.rChild.parent = nodeTmp
			if node.parent.lChild == node {
				node.parent.lChild = nodeTmp
			} else {
				node.parent.rChild = nodeTmp
			}
			if node.isRed() {
				nodeTmp.setRed()
			}
			node.free()
		}
	}
	if node = nr.createNode(key, value, startTime, endTime, nil, nil, nil); node != nil {
		nr.insert(node)
		return true
	} else {
		lgd.Error("createNode fail!")
		return false
	}
	return false
}

func (nr *nodeRoot) deleteFixTree(node, parent *nodePageElem) {
	for (node == nil || node.isBlack()) && node != nr.node {
		if parent.lChild == node {
			other := parent.rChild
			if other.isRed() {
				other.setBlack()
				parent.setRed()
				nr.leftRotate(parent)
				other = parent.rChild
			}
			if (other.lChild == nil || other.lChild.isBlack()) && (other.rChild == nil || other.rChild.isBlack()) {
				other.setRed()
				node = parent
				parent = node.parent
			} else {
				if other.rChild == nil || other.rChild.isBlack() {
					other.lChild.setBlack()
					other.setRed()
					nr.rightRotate(other)
					other = parent.rChild
				}
				if parent.isBlack() {
					other.setBlack()
				} else {
					other.setRed()
				}
				parent.setBlack()
				other.rChild.setBlack()
				nr.leftRotate(parent)
				node = nr.node
				break
			}
		} else {
			other := parent.lChild
			if other.isRed() {
				other.setBlack()
				parent.setRed()
				nr.rightRotate(parent)
				other = parent.lChild
			}
			if (other.lChild == nil || other.lChild.isBlack()) && (other.rChild == nil || other.rChild.isBlack()) {
				other.setRed()
				node = parent
				parent = node.parent
			} else {
				if other.lChild == nil || other.lChild.isBlack() {
					other.rChild.setBlack()
					other.setRed()
					nr.leftRotate(other)
					other = parent.lChild
				}
				if parent.isRed() {
					other.setRed()
				} else {
					other.setBlack()
				}
				parent.setBlack()
				other.lChild.setBlack()
				nr.rightRotate(parent)
				node = nr.node
				break
			}
		}
	}
	if node != nil {
		node.setBlack()
	}
}

func (nr *nodeRoot) delet(node *nodePageElem) {
	var child, parent, replace *nodePageElem
	var color bool

	if node.lChild != nil && node.rChild != nil {
		replace = node.rChild
		for replace.lChild != nil {
			replace = replace.lChild
		}
		if node.parent != nil {
			if node.parent.lChild == node {
				node.parent.lChild = replace
			} else {
				node.parent.rChild = replace
			}
		} else {
			nr.node = replace
		}

		child = replace.rChild
		parent = replace.parent
		color = replace.isRed()

		if parent == node {
			parent = replace
		} else {
			if child != nil {
				child.parent = parent
			}
			parent.lChild = child
			replace.rChild = node.rChild
			node.rChild.parent = replace
		}
		replace.parent = node.parent
		replace.colorType = node.colorType
		replace.lChild = node.lChild
		node.lChild.parent = replace

		if color == false {
			nr.deleteFixTree(child, parent)
		}
		node.free()
		return
	}
	if node.lChild != nil {
		child = node.lChild
	} else {
		child = node.rChild
	}

	parent = node.parent
	color = node.isRed()

	if child != nil {
		child.parent = parent
	}

	if parent != nil {
		if parent.lChild == node {
			parent.lChild = child
		} else {
			parent.rChild = child
		}
	} else {
		nr.node = child
	}
	if color == false {
		nr.deleteFixTree(child, parent)
	}
	node.free()
}

func (nr *nodeRoot) deleteNode(key string) bool {
	if node := nr.search(key); node != nil {
		nr.delet(node)
	}
	return true
}

type nodePageElem struct {
	ptr       *page
	colorType bool
	lChild    *nodePageElem
	rChild    *nodePageElem
	parent    *nodePageElem
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
}

func (n *nodePageElem) isRed() bool {
	return n.colorType
}

func (n *nodePageElem) isBlack() bool {
	return !n.colorType
}

func (n *nodePageElem) setRed() {
	n.colorType = true
}

func (n *nodePageElem) setBlack() {
	n.colorType = false
}

func (n *nodePageElem) isStart() bool {
	nowTime = time.Now().Unix()
	if n.startTime == 0 || n.startTime < nowTime {
		return true
	}
	return false
}

func (n *nodePageElem) isEnd() bool {
	nowTime = time.Now().Unix()
	if n.endTime != 0 && n.endTime < nowTime {
		return true
	}
	return false
}

func (n *nodePageElem) getStartTime() int64 {
	return n.startTime
}

func (n *nodePageElem) getEndTime() int64 {
	return n.endTime
}

func (n *nodePageElem) setTime(startTime, endTime int64) bool {
	nowTime = time.Now().Unix()
	if (startTime != 0 && nowTime > startTime) || (endTime != 0 && nowTime > endTime) {
		return false
	}
	n.startTime = startTime
	n.endTime = endTime
	return true
}

func (n *nodePageElem) setStartTime(startTime int64) bool {
	nowTime = time.Now().Unix()
	if nowTime > startTime {
		return false
	}
	n.startTime = startTime
	return true
}

func (n *nodePageElem) setEndTime(endTime int64) bool {
	nowTime = time.Now().Unix()
	if nowTime > endTime {
		return false
	}
	n.endTime = endTime
	return true
}

func (n *nodePageElem) compare(node *nodePageElem) bool {
	return bytes.Compare(n.key(), node.key()) < 0
}

func (n *nodePageElem) compareKey(key string) int {
	ok := bytes.Compare(n.key(), []byte(key))
	if ok < 0 {
		return NODE_KEY_SMALL
	} else if ok > 0 {
		return NODE_KEY_LARGE
	}
	return NODE_KEY_EQUAL
}

func (n *nodePageElem) key() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize : nodePageSize+n.keySize]
}

func (n *nodePageElem) value() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize+n.keySize : nodePageSize+n.keySize+n.valueSize]
}

func (n *nodePageElem) setKeyValue(key, value string) bool {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	copy(buf[nodePageSize:nodePageSize+n.keySize], []byte(key))
	copy(buf[nodePageSize+n.keySize:nodePageSize+n.keySize+n.valueSize], []byte(value))
	return true
}

func (n *nodePageElem) setValue(value string) bool {
	size := pageHeaderSize + nodePageSize + n.keySize + len(value)
	if size < int(n.ptr.count)*pageSize {
		buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
		n.valueSize = len(value)
		copy(buf[nodePageSize+n.keySize:nodePageSize+n.keySize+n.valueSize], []byte(value))
		return true
	}
	return false
}

func (n *nodePageElem) free() {
	globalPageList.free(n.ptr)
}

func init() {
	treeRoot = new(nodeRoot)
}
