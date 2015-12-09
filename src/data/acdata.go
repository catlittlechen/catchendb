package data

import (
	"time"
)

import lgd "catchendb/src/log"

type AcNodeData struct {
	size      int
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
	memory    []byte
}

func CreateAcData(key, value string, start, end int64) (data *AcNodeData) {
	size := len(key) + len(value) + 1
	data = globalDynamicList.acNodeData(size)
	data.Init()
	if !data.SetStartTime(start) || !data.SetEndTime(end) {
		data.Free()
		return nil
	}

	data.SetKeyValue(key, value)
	return
}

func (nd *AcNodeData) Init() {
}

func (nd *AcNodeData) Free() {
	nd.memory = nil
}

func (nd *AcNodeData) GetStartTime() int64 {
	return nd.startTime
}

func (nd *AcNodeData) SetStartTime(start int64) bool {
	nd.startTime = start
	return true
}

func (nd *AcNodeData) IsStart() bool {
	if nd.startTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *AcNodeData) GetEndTime() int64 {
	return nd.endTime
}

func (nd *AcNodeData) SetEndTime(end int64) bool {
	if end != 0 && time.Now().Unix() > end {
		return false
	}
	nd.endTime = end
	return true
}

func (nd *AcNodeData) IsEnd() bool {
	if nd.endTime != 0 && nd.endTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *AcNodeData) Key() (key []byte) {
	key = make([]byte, nd.keySize)
	if len(nd.memory) < nd.keySize {
		lgd.Info("Debug! key %d memory %d", nd.keySize, len(nd.memory))
	}
	copy(key, nd.memory[:nd.keySize])
	return
}

func (nd *AcNodeData) Value() (value []byte) {
	value = make([]byte, nd.valueSize)
	copy(value, nd.memory[nd.keySize:nd.keySize+nd.valueSize])
	return
}

func (nd *AcNodeData) SetKeyValue(key, value string) bool {
	nd.keySize = len(key)
	nd.valueSize = len(value)
	copy(nd.memory[:nd.keySize], []byte(key))
	copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *AcNodeData) SetKey(key string) bool {
	value := nd.Value()
	nd.keySize = len(key)
	size := nd.keySize + nd.valueSize
	if size > nd.size {
		return false
	}
	copy(nd.memory[:nd.keySize], []byte(key))
	copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *AcNodeData) SetValue(value string) bool {
	size := nd.keySize + len(value)
	if size < nd.size {
		nd.valueSize = len(value)
		copy(nd.memory[nd.keySize:nd.keySize+nd.valueSize], []byte(value))
		return true
	}
	return false
}
