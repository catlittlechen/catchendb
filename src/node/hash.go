package node

import (
	"bytes"
	"catchendb/src/config"
	"catchendb/src/data"
	"hash/fnv"
)

//import lgd "catchendb/src/log"

func hash(s string, count int) (index int) {
	h := fnv.New32a()
	h.Write([]byte(s))
	index = int(h.Sum32()) % count
	if index < 0 {
		index = -index
	}
	return
}

type hashRoot struct {
	size   int
	rbNode []nodeRoot
}

func (hr *hashRoot) insertNode(key, value string, start, end int64, tranID int) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].insertNode(key, value, start, end)
}

func (hr *hashRoot) searchNode(key string) (string, int64, int64) {
	index := hash(key, hr.size)
	return hr.rbNode[index].searchNode(key)
}

func (hr *hashRoot) setStartTime(key string, start int64, tranID int) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].setStartTime(key, start)
}

func (hr *hashRoot) setEndTime(key string, end int64, tranID int) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].setEndTime(key, end)
}

func (hr *hashRoot) deleteNode(key string, tranID int) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].deleteNode(key)
}

func (hr *hashRoot) output(channel chan []byte, sign []byte) {
	hrchan := make(chan []byte, 1000)
	for _, rbtree := range hr.rbNode {
		go rbtree.output(hrchan, sign)
	}
	index := 0
	datastr := []byte("")
	for {
		datastr = <-hrchan
		if bytes.Equal(datastr, sign) {
			index += 1
			if index == hr.size {
				break
			}
		}
		channel <- datastr
	}
	channel <- sign
	return
}

func (hr *hashRoot) outputData(channel chan data.Data) {
	hrchan := make(chan data.Data, 1000)
	for _, rbtree := range hr.rbNode {
		go rbtree.outputData(hrchan)
	}
	index := 0
	d := data.Data{}
	for {
		d = <-hrchan
		if len(d.Key) == 0 {
			index += 1
			if index == hr.size {
				break
			}
		}
		channel <- d
	}
	channel <- data.Data{}
	return
}

func (hr *hashRoot) input(line []byte) bool {
	d := data.Data{}
	if !d.Decode(line) {
		return false
	}
	index := hash(d.Key, hr.size)
	return hr.rbNode[index].insertNode(d.Key, d.Value, d.StartTime, d.EndTime)

}

func (hr *hashRoot) init() bool {
	hr.size = config.GlobalConf.MasterSlave.HashSize
	if hr.size < 1 {
		return false
	}
	hr.rbNode = make([]nodeRoot, hr.size)

	for index := 0; index < hr.size; index += 1 {
		hr.rbNode[index] = *new(nodeRoot)
		if !hr.rbNode[index].init() {
			return false
		}
	}
	return true
}

func hashInit() bool {
	iRoot = new(hashRoot)
	return iRoot.init()
}
