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

func init() {
}

func Put(key, value string) bool {
	return treeRoot.insertNode(key, value)
}

func Get(key string) string {
	return treeRoot.searchNode(key)
}

func Del(key string) bool {
	return treeRoot.deleteNode(key)
}

func OutPut(channel chan []byte, sign []byte) {
	treeRoot.output(channel, sign)
}

func Init() bool {
	pageSize = config.GlobalConf.PageSize
	return true
}
