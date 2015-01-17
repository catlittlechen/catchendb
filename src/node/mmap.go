package node

import (
	"sync"

	"syscall"
)

var (
	mmapLock sync.Mutex
)

func mmapSize(size int) int {
	if size <= minMmapSize {
		return minMmapSize
	} else if size < mmapBranch {
		size *= 2
	} else {
		size += mmapBranch
	}

	if size%pageSize != 0 {
		size = ((size / pageSize) + 1) * pageSize
	}

	return size
}

func mmap(size int) error {
	mmapLock.Lock()
	defer mmapLock.Unlock()
	var err error
	size = mmapSize(size)
	
	dataMap, err = syscall.Mmap(int(tempfile.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	return nil
}

func init() {
}
