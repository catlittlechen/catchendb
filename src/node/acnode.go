package node

import (
	"runtime/debug"
	"sync"
)

import lgd "code.google.com/p/log4go"

var (
	acRoot *acNodeRoot
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
	if !node.setData(key, value, start, end) {
		return nil
	}
	node.setParent(parent)
	return
}

func (ac *acNodeRoot) createData(key, value string, start, end int64, node *acNodePageElem) bool {
	return node.setData(key, value, start, end)
}

func (ac *acNodeRoot) insertNode(key, value string, start, end int64) bool {
	defer func() {
		if re := recover(); re != nil {
			lgd.Error("recover %s", re)
			lgd.Error("stack %s", debug.Stack())
			lgd.Info("key %s", key)
		}
	}()

	node := ac.node
	ok := false
	lenc := 0
	index := 0
	var acKey []byte
	var k byte
	for {
		k = key[0]
		node.lock(k)
		child := node.getChild(k)
		if child == nil {
			child = ac.createNode(key, value, start, end, node)
			if child == nil {
				node.unlock(k)
				return false
			}
			node.setChild(k, child)
			node.unlock(k)
			return true
		}
		ok, lenc, index = child.compareKey(key)
		acKey = child.key()
		if !ok {
			child2 := ac.createNode(key[:index], "", 0, 0, node)
			child3 := ac.createNode(key[index:], value, start, end, child2)
			if child2 == nil || child3 == nil {
				node.unlock(k)
				return false
			}
			child.setKey(string(acKey[index:]))
			node.setChild(k, child2)
			child2.setChild(key[index], child3)
			child2.setChild(acKey[index], child)
			child.setParent(child2)
			node.unlock(k)
			return true
		}
		switch lenc {
		case -1:
			node.unlock(k)
			key = key[index:]
			node = child
		case 0:
			if child.setValue(value) {
				child.setStartTime(start)
				child.setEndTime(end)
			} else if !ac.createData(key, value, start, end, child) {
				node.unlock(k)
				return false
			}
			node.unlock(k)
			return true
		case 1:
			child2 := ac.createNode(key, value, start, end, node)
			if child2 == nil {
				node.unlock(k)
				return false
			}
			node.setChild(k, child2)
			child2.setChild(acKey[index], child)
			child.setKey(string(acKey[index:]))
			child.setParent(child2)
			node.unlock(k)
			return true
		}

	}
	return false
}

func (ac *acNodeRoot) search(key string) (node *acNodePageElem) {
	parent := ac.node
	child := ac.node
	index := 0
	ok := false
	lenc := 0
	var k byte
	for {
		k = key[0]
		parent.lock(k)
		child = parent.getChild(k)
		if child == nil {
			parent.unlock(k)
			return nil
		}
		ok, lenc, index = child.compareKey(key)
		if !ok || lenc == 1 {
			parent.unlock(k)
			return nil
		}
		if lenc == 0 {
			parent.unlock(k)
			return child
		}
		parent.unlock(k)
		key = key[index:]
		parent = child
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
	node.setStartTime(0)
	node.setEndTime(0)

	for node.getChildNum() == 0 {
		parent := node.getParent()
		if parent == nil {
			break
		}
		k := (node.key())[0]
		parent.lock(k)
		if parent.getChild(k) != node {
			parent.unlock(k)
			continue
		}
		parent.delChild(k)
		node.free()
		node = parent
		parent.unlock(k)
	}
	return true
}

func (ac *acNodeRoot) outputData(chanData chan Data) {
	ac.preorder(chanData, "", ac.node)
	d := new(Data)
	chanData <- *d
}

func (ac *acNodeRoot) output(chans chan []byte, sign []byte) {
	chanData := make(chan Data, 1000)
	var datastr []byte
	go ac.outputData(chanData)
	for {
		d := <-chanData
		if len(d.Key) == 0 {
			break
		}
		datastr, _ = d.encode()
		chans <- datastr
	}
	chans <- sign
}

func (ac *acNodeRoot) preorder(chans chan Data, key string, node *acNodePageElem) {
	d := node.getData(key)
	if d != nil {
		chans <- *d
	}
	child := node.getAllChild()
	key += string(node.key())
	for _, v := range child {
		if v != nil {
			ac.preorder(chans, key, v)
		}
	}
}

func (ac *acNodeRoot) input(line []byte) bool {
	d := Data{}

	if !d.decode(line) {
		return false
	}
	go func() {
		if !ac.insertNode(d.Key, d.Value, d.StartTime, d.EndTime) {
			lgd.Error("insert node fail! --> data %+v", d)
		}
	}()
	return true
}

type acNodePageElem struct {
	parent *acNodePageElem

	childMutex *sync.Mutex
	child      map[byte]*acNodePageElem
	childNum   int

	nodeChildMutex map[byte]*sync.Mutex
	nodeMutex      *sync.Mutex

	dataMutex *sync.Mutex
	data      *acNodeData
}

func (ac *acNodePageElem) setData(key, value string, start, end int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()

	data := createAcData(key, value, start, end)
	if data == nil {
		return false
	}
	data2 := ac.data
	ac.data = data
	if data2 != nil {
		data2.free()
	}
	return true
}

func (ac *acNodePageElem) getData(key string) (d *Data) {
	if ac.isEnd() {
		return nil
	}

	d = new(Data)
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
	ac.nodeMutex = new(sync.Mutex)
	ac.dataMutex = new(sync.Mutex)
	ac.child = make(map[byte]*acNodePageElem)
	ac.childMutex = new(sync.Mutex)
	ac.nodeChildMutex = make(map[byte]*sync.Mutex)
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

func (ac *acNodePageElem) lock(key byte) {
	ac.nodeMutex.Lock()
	var mutex *sync.Mutex
	ok := false
	if mutex, ok = ac.nodeChildMutex[key]; !ok {
		mutex = new(sync.Mutex)
		ac.nodeChildMutex[key] = mutex
	}
	ac.nodeMutex.Unlock()

	mutex.Lock()
}

func (ac *acNodePageElem) unlock(key byte) {
	ac.nodeMutex.Lock()
	var mutex *sync.Mutex
	ok := false
	if mutex, ok = ac.nodeChildMutex[key]; !ok {
		lgd.Error("%+v", ac.nodeChildMutex)
		lgd.Error(uint8(key))
	}
	ac.nodeMutex.Unlock()

	mutex.Unlock()
}

func (ac *acNodePageElem) getChildNum() int {
	return ac.childNum
}

func (ac *acNodePageElem) getAllChild() (child map[byte]*acNodePageElem) {
	ac.childMutex.Lock()
	defer ac.childMutex.Unlock()

	child = ac.child
	return
}

func (ac *acNodePageElem) getChild(child byte) (node *acNodePageElem) {
	ac.childMutex.Lock()
	defer ac.childMutex.Unlock()

	return ac.child[child]
}

func (ac *acNodePageElem) setChild(child byte, node *acNodePageElem) {
	ac.childMutex.Lock()
	defer ac.childMutex.Unlock()

	if ac.child[child] == nil {
		ac.childNum += 1
	}
	ac.child[child] = node
}

func (ac *acNodePageElem) delChild(child byte) {
	ac.childMutex.Lock()
	defer ac.childMutex.Unlock()

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
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.isStart()
}

func (ac *acNodePageElem) isEnd() bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.isEnd()
}

func (ac *acNodePageElem) getStartTime() int64 {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.getStartTime()
}

func (ac *acNodePageElem) getEndTime() int64 {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.getEndTime()
}

func (ac *acNodePageElem) setTime(startTime, endTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setStartTime(startTime) && ac.data.setEndTime(endTime)
}

func (ac *acNodePageElem) setStartTime(startTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setStartTime(startTime)
}

func (ac *acNodePageElem) setEndTime(endTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setEndTime(endTime)
}

func (ac *acNodePageElem) key() (key []byte) {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.key()
}

func (ac *acNodePageElem) value() (value []byte) {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.value()
}

func (ac *acNodePageElem) setKeyValue(key, value string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setKeyValue(key, value)
}

func (ac *acNodePageElem) setKey(key string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setKey(key)
}

func (ac *acNodePageElem) setValue(value string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.setValue(value)
}

func (ac *acNodePageElem) free() {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	ac.data.free()
	ac.data = nil
	ac.child = nil
	ac.parent = nil
	ac.childMutex = nil
}

func acInit() bool {
	iRoot = new(acNodeRoot)
	return iRoot.init()
}
