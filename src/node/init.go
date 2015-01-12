package node

import (
	"os"
)

var (
	pageSize int
)

func init() {

	pageSize = os.Getpagesize()

}
