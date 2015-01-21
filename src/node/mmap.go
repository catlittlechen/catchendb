package node

import (
	"sync"

	"syscall"
)

const (
	minMmapSize = 1 << 22
	mmapBranch  = 1 << 30
)

var (
	mmapLock sync.Mutex
)

func mmapSize(size int) int {
	if size <= minMmapSize {
		size = minMmapSize
	} else if size > mmapBranch {
		size = mmapBranch
	}

	if size%pageSize != 0 {
		size = ((size / pageSize) + 1) * pageSize
	}

	return size
}

func mmap(size int) ([]byte, error) {
	mmapLock.Lock()
	defer mmapLock.Unlock()
	dataMap, err := syscall.Mmap(-1, 0, size, 0, 0)
	return dataMap, err
}

func munmap(dataMap []byte) error {
	err := syscall.Munmap(dataMap)
	return err
}

func init() {
}
