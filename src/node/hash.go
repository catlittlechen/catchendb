package node

import (
	"bytes"
	"catchendb/src/config"
	"hash/fnv"
)

func hash(s string, count int) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) % count
}

type hashRoot struct {
	size    int
	rbNode  []nodeRoot
	channel chan []byte
}

func (hr *hashRoot) insertNode(key, value string, start, end int64) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].insertNode(key, value, start, end)
}

func (hr *hashRoot) searchNode(key string) (string, int64, int64) {
	index := hash(key, hr.size)
	return hr.rbNode[index].searchNode(key)
}

func (hr *hashRoot) setStartTime(key string, start int64) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].setStartTime(key, start)
}

func (hr *hashRoot) setEndTime(key string, end int64) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].setEndTime(key, end)
}

func (hr *hashRoot) deleteNode(key string) bool {
	index := hash(key, hr.size)
	return hr.rbNode[index].deleteNode(key)
}

func (hr *hashRoot) output(channel chan []byte, sign []byte) {
	for _, rbtree := range hr.rbNode {
		go rbtree.output(hr.channel, sign)
	}
	index := 0
	datastr := []byte("")
	for {
		datastr = <-hr.channel
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

func (hr *hashRoot) input(line []byte) bool {
	d := data{}
	if !d.decode(line) {
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
	hr.channel = make(chan []byte)

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
