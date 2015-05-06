package logic

import (
	"catchendb/src/node"
	"fmt"
	"net/url"
	"sync"
)

var (
	transactionID    int
	transactionMutex *sync.Mutex
)

const (
	INSERT_TYPE = 1
	DELETE_TYPE = 2
	UPDATE_TYPE = 3
)

func getTransactionId() (id int) {
	transactionMutex.Lock()
	transactionID += 1
	id = transactionID
	transactionMutex.Unlock()
	return
}

type transaction struct {
	ID         int
	ChangeLog  []transactionLog
	ChangeData map[string]*node.Data
}

func (t *transaction) init() {
	t.ID = getTransactionId()
	t.ChangeLog = []transactionLog{}
	t.ChangeData = make(map[string]*node.Data)
}

func (t *transaction) getID() int {
	return t.ID
}

func (t *transaction) push(typ int, newData *node.Data) {
	tl := transactionLog{
		typ:     typ,
		newData: newData,
	}
	t.ChangeLog = append(t.ChangeLog, tl)
	t.ChangeData[newData.Key] = newData
}

func (t *transaction) getData(key string) *node.Data {
	return t.ChangeData[key]
}

func (t *transaction) isBegin() bool {
	return t.ID != 0
}

func (t *transaction) unInit() {
	t.ID = 0
	t.ChangeLog = []transactionLog{}
	t.ChangeData = nil
}

func (t *transaction) commit() (res int) {
	t.ID = -1
	for _, tl := range t.ChangeLog {
		keyword := url.Values{}
		switch tl.typ {
		case INSERT_TYPE, UPDATE_TYPE:
			node.Put(tl.newData.Key, tl.newData.Value, tl.newData.StartTime, tl.newData.EndTime, t.ID)
			keyword.Add(URL_CMD, CMD_SETEX)
			keyword.Add(URL_KEY, tl.newData.Key)
			keyword.Add(URL_VALUE, tl.newData.Value)
			keyword.Add(URL_START, fmt.Sprintf("%d", tl.newData.StartTime))
			keyword.Add(URL_END, fmt.Sprintf("%d", tl.newData.EndTime))
		case DELETE_TYPE:
			node.Del(tl.newData.Key, t.ID)
			keyword.Add(URL_CMD, CMD_DEL)
			keyword.Add(URL_KEY, tl.newData.Key)
		}

		go replicationData(keyword)
	}
	t.unInit()
	return
}

func (t *transaction) rollback() (res int) {
	t.ID = -2
	for _, tl := range t.ChangeLog {
		switch tl.typ {
		case INSERT_TYPE, UPDATE_TYPE:
			node.Put(tl.newData.Key, "", 0, 0, t.ID)
		case DELETE_TYPE:
			node.Del(tl.newData.Key, t.ID)
		}
	}
	t.unInit()
	return
}

type transactionLog struct {
	typ     int
	newData *node.Data
}

func init() {
	transactionID = 0
	transactionMutex = new(sync.Mutex)
}
