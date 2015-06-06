package data

import (
	"catchendb/src/config"
)

var (
	pageSize int
)

const (
	maxAlloacSize = 0xFFFFFFF
)

func Init() bool {
	pageSize = config.GlobalConf.PageSize
	return true
}
