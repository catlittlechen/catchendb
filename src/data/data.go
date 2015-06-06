package data

import (
	"encoding/json"
	"time"
	"unsafe"
)

type Data struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	StartTime int64  `json:"start"`
	EndTime   int64  `json:"end"`
}

func (d *Data) Decode(line []byte) bool {
	err := json.Unmarshal(line, d)
	if err != nil {
		return false
	}
	return true
}

func (d *Data) Encode() (line []byte, ok bool) {
	var err error
	line, err = json.Marshal(d)
	if err != nil {
		return
	}
	ok = true
	return
}

const (
	nodeDataSize = int(unsafe.Sizeof(NodeData{}))
)

type NodeData struct {
	ptr       *page
	startTime int64
	endTime   int64
	keySize   int
	valueSize int
}

func CreateNodeData(key, value string, start, end int64) (data *NodeData) {
	size := pageHeaderSize + nodeDataSize + len(key) + len(value)
	page := globalPageList.allocate(int(size/pageSize) + 1)
	if page != nil {
		data = page.nodeData()
		data.Init()
		if !data.SetStartTime(start) || !data.SetEndTime(end) {
			data.Free()
			return nil
		}
		data.SetKeyValue(key, value)
		return
	}
	return nil
}

func (nd *NodeData) Init() {
}

func (nd *NodeData) Free() {
	globalPageList.free(nd.ptr)
}

func (nd *NodeData) GetStartTime() int64 {
	return nd.startTime
}

func (nd *NodeData) SetStartTime(start int64) bool {
	nd.startTime = start
	return true
}

func (nd *NodeData) IsStart() bool {
	if nd.startTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *NodeData) GetEndTime() int64 {
	return nd.endTime
}

func (nd *NodeData) SetEndTime(end int64) bool {
	if end != 0 && time.Now().Unix() > end {
		return false
	}
	nd.endTime = end
	return true
}

func (nd *NodeData) IsEnd() bool {
	if nd.endTime != 0 && nd.endTime < time.Now().Unix() {
		return true
	}
	return false
}

func (nd *NodeData) Key() (key []byte) {
	key = make([]byte, nd.keySize)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(key, buf[nodeDataSize:nodeDataSize+nd.keySize])
	return
}

func (nd *NodeData) Value() (value []byte) {
	value = make([]byte, nd.keySize)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(value, buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize])
	return
}

func (nd *NodeData) SetKeyValue(key, value string) bool {
	nd.keySize = len(key)
	nd.valueSize = len(value)
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(buf[nodeDataSize:nodeDataSize+nd.keySize], []byte(key))
	copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *NodeData) SetKey(key string) bool {
	value := nd.Value()
	nd.keySize = len(key)
	size := pageHeaderSize + nodeDataSize + nd.keySize + nd.valueSize
	if size > int(nd.ptr.count)*pageSize {
		return false
	}
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
	copy(buf[nodeDataSize:nodeDataSize+nd.keySize], []byte(key))
	copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
	return true
}

func (nd *NodeData) SetValue(value string) bool {
	size := pageHeaderSize + nodeDataSize + nd.keySize + len(value)
	if size < int(nd.ptr.count)*pageSize {
		buf := (*[maxAlloacSize]byte)(unsafe.Pointer(nd))
		nd.valueSize = len(value)
		copy(buf[nodeDataSize+nd.keySize:nodeDataSize+nd.keySize+nd.valueSize], []byte(value))
		return true
	}
	return false
}
