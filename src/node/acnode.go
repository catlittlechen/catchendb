package node

import (
	"sync"
	"time"
	"unsafe"
)

var (
	acRoot    *acNodeRoot
	channelAC chan []byte
)

const (
	acNodePageSize = int(unsafe.Sizeof(acNodePageElem{}))
)

type acNodeRoot struct {
	node *acNodePageElem
}

func (ac *acNodeRoot) init() bool {
	size := pageHeaderSize + acNodePageSize
	page := globalPageList.allocate(int(size/pageSize) + 1)
	if page != nil {
		ac.node = page.acNodePageElem()
		return true
	}
	return false
}

func (ac *acNodeRoot) createNode(key, value string, start, end int64, parent *acNodePageElem) (node *acNodePageElem) {
	size := pageHeaderSize + acNodePageSize + len(key) + len(value)
	page := globalPageList.allocate(int(size/pageSize) + 1)
	if page != nil {
		node = page.acNodePageElem()
		node.setParent(parent)
		if !node.setTime(start, end) {
			node.free()
			return nil
		}
		node.channel = make(chan bool)
		node.keySize = len(key)
		node.valueSize = len(value)
		node.setKeyValue(key, value)
		return
	}
	return
}

func (ac *acNodeRoot) insertNode(key, value string, start, end int64) bool {
	node := ac.node
	ok := false
	lenc := 0
	index := 0
	status := false
	for {
		if status {
			<-node.channel
			status = false
		}
		child := node.getChild(key[0])
		if child == nil {
			child = ac.createNode(key, value, start, end, node)
			return true
		}
		ok, lenc, index = child.compareKey(key)
		if !ok || lenc == 1 {
			node.lock()
			if node.isChanging() {
				status = true
				node.unlock()
				continue
			}
			if node.getChild(key[0]) != child {
				status = true
				node.unlock()
				continue
			}
			node.changeStatus()
			node.unlock()
		}
		if !ok {
			acKey := child.key()
			child.setKey(string(acKey[index:]))
			child2 := ac.createNode(key[:index], "", 0, 0, node)
			child3 := ac.createNode(key[index:], value, start, end, child2)
			node.setChild(key[0], child2)
			child2.setChild(key[index], child3)
			child2.setChild(acKey[index], child)
			child.setParent(child2)
			node.lock()
			node.changeStatus()
			node.openBlock()
			node.unlock()
			return true
		}
		switch lenc {
		case -1:
			key = key[index:]
			node = child
		case 0:
			child.setValue(value)
			child.setStartTime(start)
			child.setEndTime(end)
			node.lock()
			node.changeStatus()
			node.openBlock()
			node.unlock()
			return true
		case 1:
			acKey := child.key()
			child2 := ac.createNode(key, value, start, end, node)
			node.setChild(key[0], child2)
			child2.setChild(acKey[index], child)
			child.setParent(child2)
			node.lock()
			node.changeStatus()
			node.openBlock()
			node.unlock()
			return true
		}

	}
	return false
}

func (ac *acNodeRoot) search(key string) (node *acNodePageElem) {
	node = ac.node.getChild(key[0])
	index := 0
	ok := false
	lenc := 0
	for node != nil {
		ok, lenc, index = node.compareKey(key)
		if !ok || lenc == 1 {
			return nil
		}
		if lenc == 0 {
			return node
		}
		key = key[:index]
		node = node.getChild(key[0])
	}
	return
}

func (ac *acNodeRoot) searchNode(key string) (value string, start, end int64) {
	node := ac.search(key)
	if node != nil {
		value = string(node.value())
		if value == "" {
			return
		}
		return value, node.getStartTime(), node.getEndTime()
	}
	return
}

func (ac *acNodeRoot) setStartTime(key string, start int64) bool {
	if node := ac.search(key); node != nil {
		return node.setStartTime(start)
	}
	return false
}

func (ac *acNodeRoot) setEndTime(key string, end int64) bool {
	if node := ac.search(key); node != nil {
		return node.setEndTime(end)
	}
	return false
}

func (ac *acNodeRoot) deleteNode(key string) bool {
	//TODO
	node := ac.search(key)
	node.lock()
	if node.getChildNum() == 1 {
		if node.isChanging() {
			node.unlock()
			<-node.channel
		}
	} else {
		node.unlock()
	}
	parent := node.getParent()
	parent.lock()
	if parent.getChildNum() == 1 && len(parent.value()) == 0 {

	}
	node.setValue("")

	return true
}

func (ac *acNodeRoot) outPut(chans chan []byte, sign []byte) {
	channelAC = chans
	ac.preorder()
	channelAC <- sign
}

func (ac *acNodeRoot) preorder() {

}

func (ac *acNodeRoot) inPut(line []byte) bool {
	d := data{}

	if !d.decode(line) {
		return false
	}
	go ac.insertNode(d.Key, d.Value, d.StartTime, d.EndTime)
	return true
}

type acNodePageElem struct {
	ptr       *page
	parent    *acNodePageElem
	child     [256]*acNodePageElem
	childNum  int
	nodeMutex *sync.Mutex
	channel   chan bool
	status    bool
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
}

func (ac *acNodePageElem) compareKey(key string) (ok bool, lenc int, index int) {
	acKey := ac.key()
	lenc = -1
	count := len(acKey)
	if count == len(key) {
		lenc = 0
	} else if count > len(key) {
		count = len(key)
		lenc = 1
	}
	for index = 0; index < count; index++ {
		if acKey[index] == key[index] {
			continue
		}
		return
	}
	ok = true
	return
}

func (ac *acNodePageElem) lock() {
	ac.nodeMutex.Lock()
}

func (ac *acNodePageElem) unlock() {
	ac.nodeMutex.Unlock()
}

func (ac *acNodePageElem) isChanging() bool {
	return ac.status
}

func (ac *acNodePageElem) changeStatus() {
	ac.status = !ac.status
}

func (ac *acNodePageElem) openBlock() {
	close(ac.channel)
	ac.channel = make(chan bool)
}

func (ac *acNodePageElem) getChildNum() int {
	return ac.childNum
}

func (ac *acNodePageElem) getChild(child byte) (node *acNodePageElem) {
	return ac.child[child]
}

func (ac *acNodePageElem) setChild(child byte, node *acNodePageElem) {
	if ac.child[child] == nil {
		ac.childNum += 1
	}
	ac.child[child] = node
}

func (ac *acNodePageElem) delChild(child byte) {
	if ac.child[child] != nil {
		ac.childNum -= 1
		ac.child[child] = nil
	}
}

func (ac *acNodePageElem) setParent(parent *acNodePageElem) {
	ac.parent = parent
}

func (ac *acNodePageElem) getParent() *acNodePageElem {
	return ac.parent
}

func (ac *acNodePageElem) isStart() bool {
	nowTime = time.Now().Unix()
	if ac.startTime == 0 || ac.startTime < nowTime {
		return true
	}
	return false
}

func (ac *acNodePageElem) isEnd() bool {
	nowTime = time.Now().Unix()
	if ac.endTime != 0 && ac.endTime < nowTime {
		return true
	}
	return false
}

func (ac *acNodePageElem) getStartTime() int64 {
	return ac.startTime
}

func (ac *acNodePageElem) getEndTime() int64 {
	return ac.endTime
}

func (ac *acNodePageElem) setTime(startTime, endTime int64) bool {
	nowTime = time.Now().Unix()
	if (startTime != 0 && nowTime > startTime) || (endTime != 0 && nowTime > endTime) {
		return false
	}
	ac.startTime = startTime
	ac.endTime = endTime
	return true
}

func (ac *acNodePageElem) setStartTime(startTime int64) bool {
	nowTime = time.Now().Unix()
	if nowTime > startTime {
		return false
	}
	ac.startTime = startTime
	return true
}

func (ac *acNodePageElem) setEndTime(endTime int64) bool {
	nowTime = time.Now().Unix()
	if nowTime > endTime {
		return false
	}
	ac.endTime = endTime
	return true
}

func (ac *acNodePageElem) key() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(ac))
	return buf[nodePageSize : nodePageSize+ac.keySize]
}

func (ac *acNodePageElem) value() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(ac))
	return buf[nodePageSize+ac.keySize : nodePageSize+ac.keySize+ac.valueSize]
}

func (ac *acNodePageElem) setKeyValue(key, value string) bool {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(ac))
	ac.keySize = len(key)
	ac.valueSize = len(value)
	copy(buf[nodePageSize:nodePageSize+ac.keySize], []byte(key))
	copy(buf[nodePageSize+ac.keySize:nodePageSize+ac.keySize+ac.valueSize], []byte(value))
	return true
}

func (ac *acNodePageElem) setKey(key string) bool {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(ac))
	value := ac.value()
	ac.keySize = len(key)
	ac.valueSize = len(value)
	copy(buf[nodePageSize:nodePageSize+ac.keySize], []byte(key))
	copy(buf[nodePageSize+ac.keySize:nodePageSize+ac.keySize+ac.valueSize], []byte(value))
	return true
}

func (ac *acNodePageElem) setValue(value string) bool {
	size := pageHeaderSize + nodePageSize + ac.keySize + len(value)
	if size < int(ac.ptr.count)*pageSize {
		buf := (*[maxAlloacSize]byte)(unsafe.Pointer(ac))
		ac.valueSize = len(value)
		copy(buf[nodePageSize+ac.keySize:nodePageSize+ac.keySize+ac.valueSize], []byte(value))
		return true
	}
	return false
}

func (ac *acNodePageElem) free() {
	globalPageList.free(ac.ptr)
}

func acInit() bool {
	acRoot = new(acNodeRoot)
	return acRoot.init()
}
