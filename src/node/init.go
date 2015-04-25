package node

import (
	"catchendb/src/config"
)

var (
	pageSize int
)

const (
	maxAlloacSize = 0xFFFFFFF
)

func Put(key, value string, startTime, endTime int64) bool {
	return acRoot.insertNode(key, value, startTime, endTime)
}

func Get(key string) (string, int64, int64) {
	return acRoot.searchNode(key)
}

func Set(key string, start, end int64) bool {
	if start != 0 && !acRoot.setStartTime(key, start) {
		return false
	}
	if end != 0 && !acRoot.setEndTime(key, end) {
		return false
	}
	return true
}

func Del(key string) bool {
	return acRoot.deleteNode(key)
}

func OutPut(channel chan []byte, sign []byte) {
	acRoot.output(channel, sign)
}

func InPut(line []byte) bool {
	return acRoot.input(line)
}

func Init() bool {
	pageSize = config.GlobalConf.PageSize
	return acInit()
}
