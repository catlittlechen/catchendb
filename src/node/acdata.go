package node

import (
	"sync"
	"time"
)

import lgd "code.google.com/p/log4go"

type acNodeData struct {
	Mutex     *sync.Mutex
	size      int
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
	memory    []byte
}

func createAcData(key, value string, start, end int64, mutex bool) (data *acNodeData) {
	size := len(key) + len(value) + 1
	data = globalDynamicList.acNodeData(size)
	data.init()
	if !data.setStartTime(start) || !data.setEndTime(end) {
		data.free()
		return nil
	}

	data.setKeyValue(key, value)
	if mutex {
		data.lock()
	}
	return
}

func (nd *acNodeData) init() {
	nd.Mutex = new(sync.Mutex)
}

func (nd *acNodeData) free() {
	nd.memory = nil
}

func (nd *acNodeData) getStartTime() int64 {
	return nd.startTime
}

func (nd *acNodeData) setStartTime(start int64) bool {
	nd.Mutex.Lock()
	defer nd.Mutex.Unlock()
	nd.startTime = start
	return true
}

func (nd *acNodeData) isStart() bool {
	if nd.startTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *acNodeData) getEndTime() int64 {
	return nd.endTime
}

func (nd *acNodeData) setEndTime(end int64) bool {
	nd.Mutex.Lock()
	defer nd.Mutex.Unlock()
	if end != 0 && time.Now().Unix() > end {
		return false
	}
	nd.endTime = end
	return true
}

func (nd *acNodeData) isEnd() bool {
	if nd.endTime != 0 && nd.endTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *acNodeData) key() (key []byte) {
	key = make([]byte, nd.keySize)
	if len(nd.memory) < nd.keySize {
		lgd.Info("Debug! key %d memory %d", nd.keySize, len(nd.memory))
	}
	copy(key, nd.memory[:nd.keySize])
	return
}

func (nd *acNodeData) value() (value []byte) {
	value = make([]byte, nd.valueSize)
	copy(value, nd.memory[nd.keySize:nd.keySize+nd.valueSize])
	return
}

func (nd *acNodeData) setKeyValue(key, value string) bool {
	nd.Mutex.Lock()
	defer nd.Mutex.Unlock()
	nd.keySize = len(key)
	nd.valueSize = len(value)
	copy(nd.memory[:nd.keySize], []byte(key))
	copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *acNodeData) setKey(key string) bool {
	nd.Mutex.Lock()
	defer nd.Mutex.Unlock()
	value := nd.value()
	nd.keySize = len(key)
	size := nd.keySize + nd.valueSize
	if size > nd.size {
		return false
	}
	copy(nd.memory[:nd.keySize], []byte(key))
	copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *acNodeData) setValue(value string) bool {
	nd.Mutex.Lock()
	defer nd.Mutex.Unlock()
	size := nd.keySize + len(value)
	if size < nd.size {
		nd.valueSize = len(value)
		copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
		return true
	}
	return false
}

func (nd *acNodeData) lock() {
	nd.Mutex.Lock()
}

func (nd *acNodeData) unLock() {
	nd.Mutex.Unlock()
}
