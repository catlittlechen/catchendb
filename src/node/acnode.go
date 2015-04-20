package node

import (
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
		node.parent = parent
		if !node.setTime(start, end) {
			node.free()
			return nil
		}
		node.keySize = len(key)
		node.valueSize = len(value)
		node.setKeyValue(key, value)
		return
	}
	return
}

func (ac *acNodeRoot) insertNode(key, value string, start, end int64) bool {
	node := ac.node.getChild(key[0])

	return false
}

func (ac *acNodeRoot) search(key string) (node *acNodePageElem) {
	node = ac.node.getChild(key[0])
	index := 0
	for node != nil {
		index = node.compareKey(key)
		if index < 0 {
			return nil
		}
		if index == 0 {
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
		return string(node.value()), node.getStartTime(), node.getEndTime()
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
	return false
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
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
}

func (ac *acNodePageElem) compareKey(key string) (index int) {
	acKey := ac.key()
	acBool := true
	count := len(acKey)
	if count > len(key) {
		count = len(key)
		acBool = false
	}
	for index = 0; index < count; index++ {
		if acKey[index] == key[index] {
			continue
		}
		if acBool {
			index = -index
		}
		return
	}
	index = 0
	return
}

func (ac *acNodePageElem) getChild(child byte) (node *acNodePageElem) {
	return ac.child[child]
}

func (ac *acNodePageElem) setChild(child byte, node *acNodePageElem) {
	ac.child[child] = node
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
