package node

import (
	"bytes"
	"catchendb/src/config"
	"hash/fnv"
)

//import lgd "code.google.com/p/log4go"

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
			index++
			if index == hr.size {
				break
			}
		}
		channel <- datastr
	}
	channel <- sign
	return
}

func (hr *hashRoot) outputData(channel chan Data) {
	hrchan := make(chan Data, 1000)
	for _, rbtree := range hr.rbNode {
		go rbtree.outputData(hrchan)
	}
	index := 0
	d := Data{}
	for {
		d = <-hrchan
		if len(d.Key) == 0 {
			index++
			if index == hr.size {
				break
			}
		}
		channel <- d
	}
	channel <- Data{}
	return
}

func (hr *hashRoot) input(line []byte) bool {
	d := Data{}
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

	for index := 0; index < hr.size; index++ {
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
