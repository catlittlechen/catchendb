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

func Init() bool {
	return true
}
