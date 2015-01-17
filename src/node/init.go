package node

import (
	"os"
)

var (
	pageSize      int
	maxAlloacSize = 0xFFFFFFF
)

func init() {
	pageSize = os.Getpagesize()
}

func Init() bool {
	return true
}
