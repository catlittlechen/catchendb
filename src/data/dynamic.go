package data

var globalDynamicList dynamicList

type dynamicList struct {
}

func (dl *dynamicList) allocateSize(size int) (sizePage int) {
	sizePage = pageSize
	if size < sizePage {
		for size < sizePage {
			sizePage /= 2
		}
		sizePage *= 2
	} else {
		for size > sizePage {
			sizePage *= 2
		}
	}
	return
}

func (dl *dynamicList) acNodeData(size int) (and *AcNodeData) {
	size = dl.allocateSize(size)
	and = new(AcNodeData)
	and.size = size
	and.memory = make([]byte, size)
	return
}

func init() {
	globalDynamicList = dynamicList{}
}
