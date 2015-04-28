package node

import (
	"sync"
)

import lgd "code.google.com/p/log4go"

var (
	channelAC chan []byte
	acRoot    *acNodeRoot
)

type acNodeRoot struct {
	node *acNodePageElem
}

func (ac *acNodeRoot) init() bool {
	ac.node = ac.createNode("", "", 0, 0, nil)
	if ac.node == nil {
		return false
	}
	return true
}

func (ac *acNodeRoot) createNode(key, value string, start, end int64, parent *acNodePageElem) (node *acNodePageElem) {
	node = new(acNodePageElem)
	node.init()
	node.data = createAcData(key, value, start, end)
	if node.data == nil {
		return nil
	}
	node.setParent(parent)
	return
}

func (ac *acNodeRoot) createData(key, value string, start, end int64, node *acNodePageElem) bool {
	data := createAcData(key, value, start, end)
	if data == nil {
		return false
	}
	data2 := node.data
	node.data = data
	data2.free()
	return true
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
			node.lock()
			if node.isChanging() {
				status = true
				node.unlock()
				continue
			}
			if node.getChild(key[0]) != child {
				node.unlock()
				continue
			}
			child = ac.createNode(key, value, start, end, node)
			if child == nil {
				node.unlock()
				return false
			}
			node.setChild(key[0], child)
			node.unlock()
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
				node.unlock()
				continue
			}
			node.changeStatus()
			node.unlock()
		}
		if !ok {
			acKey := child.key()
			child2 := ac.createNode(key[:index], "", 0, 0, node)
			child3 := ac.createNode(key[index:], value, start, end, child2)
			if child2 == nil || child3 == nil {
				node.lock()
				node.changeStatus()
				node.openBlock()
				node.unlock()
				return false
			}
			node.lock()
			child.setKey(string(acKey[index:]))
			node.setChild(key[0], child2)
			child2.setChild(key[index], child3)
			child2.setChild(acKey[index], child)
			child.setParent(child2)
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
			if child.setValue(value) {
				child.setStartTime(start)
				child.setEndTime(end)
			} else if !ac.createData(key, value, start, end, child) {
				node.lock()
				node.changeStatus()
				node.openBlock()
				node.unlock()
				return false
			}
			node.lock()
			node.changeStatus()
			node.openBlock()
			node.unlock()
			return true
		case 1:
			acKey := child.key()
			child2 := ac.createNode(key, value, start, end, node)
			if child2 == nil {
				node.lock()
				node.changeStatus()
				node.openBlock()
				node.unlock()
				return false
			}
			node.lock()
			node.setChild(key[0], child2)
			child2.setChild(acKey[index], child)
			child.setKey(string(acKey[index:]))
			child.setParent(child2)
			node.changeStatus()
			node.openBlock()
			node.unlock()
			return true
		}

	}
	return false
}

func (ac *acNodeRoot) search(key string) (node *acNodePageElem) {
	node = ac.node
	index := 0
	ok := false
	lenc := 0
	status := false
	child := ac.node
	for node != nil {
		if status {
			<-node.channel
			status = false
		}
		child = ac.node.getChild(key[0])
		if child == nil {
			return
		}
		ok, lenc, index = node.compareKey(key)
		if !ok || lenc == 1 {
			return nil
		}
		if lenc == 0 {
			return node
		}
		key = key[index:]
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
	node := ac.search(key)
	node.setValue("")

	for node.getChildNum() == 0 {
		parent := node.getParent()
		if parent == nil {
			break
		}
		parent.lock()
		if parent.isChanging() {
			parent.unlock()
			continue
		}
		if parent.getChild((node.key())[0]) != node {
			parent.unlock()
			continue
		}
		parent.delChild((node.key())[0])
		node.free()
		node = parent
		parent.unlock()
	}
	return true
}

func (ac *acNodeRoot) output(chans chan []byte, sign []byte) {
	channelAC = chans
	ac.preorder("", ac.node)
	channelAC <- sign
}

func (ac *acNodeRoot) preorder(key string, node *acNodePageElem) {
	d := node.getData(key)
	if d != nil {
		datastr, _ := d.encode()
		channelAC <- datastr
	}
	child := node.getAllChild()
	key += string(node.key())
	for _, v := range child {
		if v != nil {
			ac.preorder(key, v)
		}
	}
}

func (ac *acNodeRoot) input(line []byte) bool {
	d := data{}

	if !d.decode(line) {
		return false
	}
	ac.insertNode(d.Key, d.Value, d.StartTime, d.EndTime)
	return true
}

type acNodePageElem struct {
	parent    *acNodePageElem
	child     map[byte]*acNodePageElem
	childNum  int
	nodeMutex *sync.Mutex
	channel   chan bool
	status    bool

	data *acNodeData
}

func (ac *acNodePageElem) getData(key string) (d *data) {
	if ac.isEnd() {
		return nil
	}

	d = new(data)
	d.Value = string(ac.value())
	if len(d.Value) == 0 {
		return nil
	}
	d.StartTime = ac.getStartTime()
	d.EndTime = ac.getEndTime()
	d.Key = key + string(ac.key())
	return
}

func (ac *acNodePageElem) init() {
	ac.channel = make(chan bool)
	ac.nodeMutex = new(sync.Mutex)
	ac.child = make(map[byte]*acNodePageElem)
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

func (ac *acNodePageElem) getAllChild() map[byte]*acNodePageElem {
	return ac.child
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
	return ac.data.isStart()
}

func (ac *acNodePageElem) isEnd() bool {
	return ac.data.isEnd()
}

func (ac *acNodePageElem) getStartTime() int64 {
	return ac.data.getStartTime()
}

func (ac *acNodePageElem) getEndTime() int64 {
	return ac.data.getEndTime()
}

func (ac *acNodePageElem) setTime(startTime, endTime int64) bool {
	return ac.data.setStartTime(startTime) && ac.data.setEndTime(endTime)
}

func (ac *acNodePageElem) setStartTime(startTime int64) bool {
	return ac.data.setStartTime(startTime)
}

func (ac *acNodePageElem) setEndTime(endTime int64) bool {
	return ac.data.setEndTime(endTime)
}

func (ac *acNodePageElem) key() (key []byte) {
	if ac.data == nil {
		lgd.Error("Debug!")
	}
	return ac.data.key()
}

func (ac *acNodePageElem) value() (value []byte) {
	return ac.data.value()
}

func (ac *acNodePageElem) setKeyValue(key, value string) bool {
	return ac.data.setKeyValue(key, value)
}

func (ac *acNodePageElem) setKey(key string) bool {
	return ac.data.setKey(key)
}

func (ac *acNodePageElem) setValue(value string) bool {
	return ac.data.setValue(value)
}

func (ac *acNodePageElem) free() {
	ac.data.free()
	ac.data = nil
}

func acInit() bool {
	iRoot = new(acNodeRoot)
	return iRoot.init()
}
