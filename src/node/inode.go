package node

var (
	iRoot iNodeRoot
)

type iNodeRoot interface {
	insertNode(string, string, int64, int64) bool

	searchNode(string) (string, int64, int64)

	setStartTime(string, int64) bool

	setEndTime(string, int64) bool

	deleteNode(string) bool

	output(chan []byte, []byte)

	input([]byte) bool

	init() bool

	outputData(chan Data)
}
