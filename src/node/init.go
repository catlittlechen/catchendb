package node

import (
	"os"
)

var (
	pageSize int
)

const (
	maxAlloacSize = 0xFFFFFFF
)

func init() {
	pageSize = os.Getpagesize()
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

func Init() bool {
	return true
}
