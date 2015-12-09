package node

import (
	"catchendb/src/config"
	"catchendb/src/data"
)

func Put(key, value string, startTime, endTime int64, tranID int) bool {
	return iRoot.insertNode(key, value, startTime, endTime, tranID)
}

func Get(key string) (string, int64, int64) {
	return iRoot.searchNode(key)
}

func Set(key string, start, end int64, tranID int) bool {
	if start != 0 && !iRoot.setStartTime(key, start, tranID) {
		return false
	}
	if end != 0 && !iRoot.setEndTime(key, end, tranID) {
		return false
	}
	return true
}

func Del(key string, tranID int) bool {
	return iRoot.deleteNode(key, tranID)
}

func OutPut(channel chan []byte, sign []byte) {
	iRoot.output(channel, sign)
}

func OutPutData(channel chan data.Data) {
	iRoot.outputData(channel)
}

func InPut(line []byte) bool {
	return iRoot.input(line)
}

func Init() bool {
	if config.GlobalConf.MasterSlave.IsMaster {
		return acInit()
	}
	return hashInit()
}
