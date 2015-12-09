package node

import (
	"catchendb/src/data"
	"runtime/debug"
	"sync"
)

import lgd "catchendb/src/log"

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

func (ac *acNodeRoot) insertNode(key, value string, start, end int64, tranID int) bool {
	defer func() {
		if re := recover(); re != nil {
			lgd.Errorf("recover %s", re)
			lgd.Errorf("stack %s", debug.Stack())
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
			if ret := child.transaction(tranID); ret == 1 {
				child.setValue("")
				child.setStartTime(0)
				child.setEndTime(0)
			} else if ret == 3 {
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
			if ret := child3.transaction(tranID); ret == 1 {
				child3.setValue("")
				child3.setStartTime(0)
				child3.setEndTime(0)
			} else if ret == 3 {
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
			if ret := child.transaction(tranID); ret == 1 {
				node.unlock(k)
				return true
			} else if ret == 3 {
				node.unlock(k)
				return false
			}
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
			if ret := child2.transaction(tranID); ret == 1 {
				child2.setValue("")
				child2.setStartTime(0)
				child2.setEndTime(0)
			} else if ret == 3 {
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
			if child.isEnd() {
				go ac.deleteNode(key, 0)
				return nil
			}
			if len(child.value()) == 0 {
				return nil
			}
			return child
		}
		parent.unlock(k)
		key = key[index:]
		parent = child
	}
}

func (ac *acNodeRoot) searchNode(key string) (value string, start, end int64) {
	node := ac.search(key)
	if node != nil {
		value = string(node.value())
		return value, node.getStartTime(), node.getEndTime()
	}
	return
}

func (ac *acNodeRoot) setStartTime(key string, start int64, tranID int) bool {
	var node *acNodePageElem
	if node = ac.search(key); node == nil {
		return false
	}
	if ret := node.transaction(tranID); ret == 1 {
		return true
	} else if ret == 3 {
		return false
	}
	return node.setStartTime(start)
}

func (ac *acNodeRoot) setEndTime(key string, end int64, tranID int) bool {
	if node := ac.search(key); node != nil {
		if ret := node.transaction(tranID); ret == 1 {
			return true
		} else if ret == 3 {
			return false
		}
		return node.setEndTime(end)
	}
	return false
}

//DoSomething
func (ac *acNodeRoot) deleteNode(key string, tranID int) bool {
	node := ac.search(key)
	ret := 0
	if ret = node.transaction(tranID); ret == 1 {
		return true
	} else if ret == 3 {
		return false
	}
	node.setValue("")
	node.setStartTime(0)
	node.setEndTime(0)

	for node.getChildNum() == 0 {
		if ret = node.transaction(tranID); ret != 2 {
			return true
		}
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

func (ac *acNodeRoot) outputData(chanData chan data.Data) {
	ac.preorder(chanData, "", ac.node)
	d := new(data.Data)
	chanData <- *d
}

func (ac *acNodeRoot) output(chans chan []byte, sign []byte) {
	chanData := make(chan data.Data, 1000)
	var datastr []byte
	go ac.outputData(chanData)
	for {
		d := <-chanData
		if len(d.Key) == 0 {
			break
		}
		datastr, _ = d.Encode()
		chans <- datastr
	}
	chans <- sign
}

func (ac *acNodeRoot) preorder(chans chan data.Data, key string, node *acNodePageElem) {
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
	d := data.Data{}

	if !d.Decode(line) {
		return false
	}
	go func() {
		if !ac.insertNode(d.Key, d.Value, d.StartTime, d.EndTime, 0) {
			lgd.Errorf("insert node fail! --> data %+v", d)
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
	data      *data.AcNodeData

	transactionID int
}

//@ret 1 success 2 go to the next action 3 timeout
func (ac *acNodePageElem) transaction(tranID int) (ret int) {
	ret = 1
	if tranID < 0 {
		ac.transactionID = 0
		if tranID == -2 {
			//回滚
			return
		}
	} else if ac.transactionID == 0 {
		if tranID != 0 {
			ac.transactionID = tranID
			//抢占锁
			return
		}
	} else {
		if ac.transactionID != tranID {
			//不同个事务
			ret = 3
		}
		return
	}
	ret = 2
	return
}

func (ac *acNodePageElem) setData(key, value string, start, end int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()

	data := data.CreateAcData(key, value, start, end)
	if data == nil {
		return false
	}
	data2 := ac.data
	ac.data = data
	if data2 != nil {
		data2.Free()
	}
	return true
}

func (ac *acNodePageElem) getData(key string) (d *data.Data) {
	if ac.isEnd() {
		return nil
	}

	d = new(data.Data)
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
		lgd.Errorf("%+v", ac.nodeChildMutex)
		lgd.Errorf("%+d", uint8(key))
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
		ac.childNum++
	}
	ac.child[child] = node
}

func (ac *acNodePageElem) delChild(child byte) {
	ac.childMutex.Lock()
	defer ac.childMutex.Unlock()

	if ac.child[child] != nil {
		ac.childNum--
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
	return ac.data.IsStart()
}

func (ac *acNodePageElem) isEnd() bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.IsEnd()
}

func (ac *acNodePageElem) getStartTime() int64 {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.GetStartTime()
}

func (ac *acNodePageElem) getEndTime() int64 {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.GetEndTime()
}

func (ac *acNodePageElem) setTime(startTime, endTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetStartTime(startTime) && ac.data.SetEndTime(endTime)
}

func (ac *acNodePageElem) setStartTime(startTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetStartTime(startTime)
}

func (ac *acNodePageElem) setEndTime(endTime int64) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetEndTime(endTime)
}

func (ac *acNodePageElem) key() (key []byte) {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.Key()
}

func (ac *acNodePageElem) value() (value []byte) {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.Value()
}

func (ac *acNodePageElem) setKeyValue(key, value string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetKeyValue(key, value)
}

func (ac *acNodePageElem) setKey(key string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetKey(key)
}

func (ac *acNodePageElem) setValue(value string) bool {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	return ac.data.SetValue(value)
}

func (ac *acNodePageElem) free() {
	ac.dataMutex.Lock()
	defer ac.dataMutex.Unlock()
	ac.data.Free()
	ac.data = nil
	ac.child = nil
	ac.parent = nil
	ac.childMutex = nil
}

func acInit() bool {
	iRoot = new(acNodeRoot)
	return iRoot.init()
}
