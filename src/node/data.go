package node

import (
	"encoding/json"
	"time"
	"unsafe"
)

type data struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	StartTime int64  `json:"start"`
	EndTime   int64  `json:"end"`
}

func (d *data) decode(line []byte) bool {
	err := json.Unmarshal(line, d)
	if err != nil {
		return false
	}
	return true
}

func (d *data) encode() (line []byte, ok bool) {
	var err error
	line, err = json.Marshal(d)
	if err != nil {
		return
	}
	ok = true
	return
}

const (
	nodeDataSize = int(unsafe.Sizeof(nodeData{}))
)

type nodeData struct {
	ptr            *page
	lastModifyTime int64
	startTime      int64
	endTime        int64
	keySize        int
	valueSize      int
}

func createData(key, value string, start, end int64) (data *nodeData) {
	size := pageHeaderSize + nodeDataSize + len(key) + len(value)
	page := globalPageList.allocate(int(size/pageSize) + 1)
	if page != nil {
		data = page.nodeData()
		data.init()
		if !data.setStartTime(start) || !data.setEndTime(end) {
			data.free()
			return nil
		}
		data.setKeyValue(key, value)
		return
	}
	return nil
}

func (nd *nodeData) init() {
	nd.lastModifyTime = time.Now().Unix()
}

func (nd *nodeData) free() {
	globalPageList.free(nd.ptr)
}

func (nd *nodeData) getLastModifyTime() int64 {
	return nd.lastModifyTime
}

func (nd *nodeData) getStartTime() int64 {
	return nd.startTime
}

func (nd *nodeData) setStartTime(start int64) bool {
	nd.lastModifyTime = time.Now().Unix()
	nd.startTime = start
	return true
}

func (nd *nodeData) isStart() bool {
	if nd.startTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *nodeData) getEndTime() int64 {
	return nd.endTime
}

func (nd *nodeData) setEndTime(end int64) bool {
	if end != 0 && time.Now().Unix() > end {
		return false
	}
	nd.lastModifyTime = time.Now().Unix()
	nd.endTime = end
	return true
}

func (nd *nodeData) isEnd() bool {
	if nd.endTime != 0 && nd.endTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *nodeData) key() (key []byte) {
	key = make([]byte, nd.keySize)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(key, buf[nodeDataSize:nodeDataSize+nd.keySize])
	return
}

func (nd *nodeData) value() (value []byte) {
	value = make([]byte, nd.keySize)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(value, buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize])
	return
}

func (nd *nodeData) setKeyValue(key, value string) bool {
	nd.lastModifyTime = time.Now().Unix()
	nd.keySize = len(key)
	nd.valueSize = len(value)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(buf[nodeDataSize:nodeDataSize+nd.keySize], []byte(key))
	copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *nodeData) setKey(key string) bool {
	value := nd.value()
	nd.keySize = len(key)
	size := pageHeaderSize + nodeDataSize + nd.keySize + nd.valueSize
	if size > int(nd.ptr.count)*pageSize {
		return false
	}
	nd.lastModifyTime = time.Now().Unix()
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(buf[nodeDataSize:nodeDataSize+nd.keySize], []byte(key))
	copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *nodeData) setValue(value string) bool {
	size := pageHeaderSize + nodeDataSize + nd.keySize + len(value)
	if size < int(nd.ptr.count)*pageSize {
		nd.lastModifyTime = time.Now().Unix()
		buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
		nd.valueSize = len(value)
		copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
		return true
	}
	return false
}
