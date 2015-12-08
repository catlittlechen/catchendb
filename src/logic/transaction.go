package logic

import (
	"catchendb/src/data"
	"catchendb/src/node"
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

func getTransactionID() (id int) {
	transactionMutex.Lock()
	transactionID++
	id = transactionID
	transactionMutex.Unlock()
	return
}

type transaction struct {
	ID         int
	ChangeLog  []transactionLog
	ChangeData map[string]*data.Data
}

func (t *transaction) init() {
	t.ID = getTransactionID()
	t.ChangeLog = []transactionLog{}
	t.ChangeData = make(map[string]*data.Data)
}

func (t *transaction) getID() int {
	return t.ID
}

func (t *transaction) push(typ int, newData *data.Data) {
	tl := transactionLog{
		typ:     typ,
		newData: newData,
	}
	t.ChangeLog = append(t.ChangeLog, tl)
	t.ChangeData[newData.Key] = newData
}

func (t *transaction) getData(key string) *data.Data {
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
	var req Req
	for _, tl := range t.ChangeLog {
		req = Req{}
		switch tl.typ {
		case INSERT_TYPE, UPDATE_TYPE:
			node.Put(tl.newData.Key, tl.newData.Value, tl.newData.StartTime, tl.newData.EndTime, t.ID)
			req.C = CMD_SETEX
			req.Key = tl.newData.Key
			req.Value = tl.newData.Value
			req.StartTime = tl.newData.StartTime
			req.EndTime = tl.newData.EndTime
		case DELETE_TYPE:
			node.Del(tl.newData.Key, t.ID)
			req.C = CMD_DEL
			req.Key = tl.newData.Key
		}

		go replicationData(req)
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
	newData *data.Data
}

func init() {
	transactionID = 0
	transactionMutex = new(sync.Mutex)
}
