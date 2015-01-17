package node

import (
	"catchendb/src/config"
	"os"
	"syscall"
)

import lgd "code.google.com/p/log4go"

var (
	pageSize int
	tempfile *os.File
	dataMap []byte
)

const (
	minMmapSize = 1 << 22
	mmapBranch  = 1 << 30
)

func init() {
	pageSize = os.Getpagesize()
}

func fileLock(f *os.File) error {
	for {
		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return nil

		} else if err != syscall.EWOULDBLOCK {
			return err

		}

	}
	return nil

}

func fileUnlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

}

func Init() bool {
	var err error
	tempfile, err = os.OpenFile(config.GlobalConf.Server.TempPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		lgd.Error("tempfile open fail with error %s", err)
		return false
	}

	if err := fileLock(tempfile); err != nil {
		lgd.Error("tempfile filelock error %s", err)
		return false
	}
	return true
}
