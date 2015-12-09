package logic

import (
	"catchendb/src/config"
	"catchendb/src/user"
	"catchendb/src/util"
	"encoding/json"
	"net"
	"sync"
)

import lgd "catchendb/src/log"

var (
	userConnection map[string]int
	userMutex      *sync.Mutex
)

func ReplicationLogic(conn *net.TCPConn) {
	data := make([]byte, 1024)
	count, err := conn.Read(data)
	if err != nil {
		lgd.Errorf("read error[%s]", err)
		return
	}
	ok, name, res := aut(data[:count])
	_, err = conn.Write(res)
	if !ok {
		return
	}
	defer disConnection(name)
	if err != nil {
		lgd.Warn("write error[%s]", err)
		return
	}
	replicationMaster(name, conn)
}

func clientAuth(conn *net.TCPConn) (name string, res []byte, ok bool) {
	data := make([]byte, 1024)
	count, err := conn.Read(data)
	if err != nil {
		lgd.Errorf("read error[%s]", err)
		return
	}

	ok, name, res = aut(data[:count])
	if !ok {
		return
	}
	return
}

func ClientLogic(conn *net.TCPConn) {
	name, res, ok := clientAuth(conn)
	if !ok {
		return
	}
	defer disConnection(name)

	_, err := conn.Write(res)
	if err != nil {
		lgd.Warn("write error[%s]", err)
		return
	}

	privilege := user.GetPrivilege(name)
	errRes := util.JSONOut(Rsp{
		C: errUrlParse,
	})

	tranObj := new(transaction)
	defer func() {
		if tranObj.isBegin() {
			tranObj.rollback()
		}
	}()

	data := make([]byte, 1024)
	rdata := make([]byte, 1024)
	var req Req
	var count int
	for {
		count, err = conn.Read(data)
		if err != nil {
			lgd.Warn("read error[%s]", err)
			return
		}

		if err = json.Unmarshal(data[:count], &req); err == nil {
			if ok, res = clientTransactionLogic(req, tranObj); ok {
				rdata = mapAction(req, privilege, false, tranObj)
			} else {
				rdata = res
			}
		} else {
			rdata = errRes
		}

		if _, err = conn.Write(rdata); err != nil {
			lgd.Warn("write error[%s]", err)
			return
		}

	}
}

func clientTransactionLogic(req Req, tranObj *transaction) (normal bool, res []byte) {
	rsp := Rsp{}
	switch req.C {
	case cmdBegin:
		if tranObj.isBegin() {
			rsp.C = errTraBegin
		} else {
			tranObj.init()
		}
		res = util.JSONOut(rsp)
		return
	case cmdRollback:
		if tranObj.isBegin() {
			rsp.C = tranObj.rollback()
		} else {
			rsp.C = errTraNoBegin
		}
		res = util.JSONOut(rsp)
		return
	case cmdCommit:
		if tranObj.isBegin() {
			rsp.C = tranObj.commit()
		} else {
			rsp.C = errTraNoBegin
		}
		res = util.JSONOut(rsp)
		return
	default:
		normal = true
	}
	return
}

func disConnection(name string) {
	userMutex.Lock()
	defer userMutex.Unlock()
	if userConnection[name] == 1 {
		delete(userConnection, name)
	} else {
		userConnection[name]--
	}
	return
}

func aut(data []byte) (ok bool, name string, r []byte) {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	req := Req{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		lgd.Errorf("ParseQuery fail with the data %s", string(data))
		rsp.C = errUrlParse
		r = util.JSONOut(rsp)
		return
	}
	ok, name = handleUserAut(req, nil)
	if !ok {
		rsp.C = errAccessDenied
		r = util.JSONOut(rsp)
		return
	}
	userMutex.Lock()
	if _, ok2 := userConnection[name]; ok2 {
		if userConnection[name] >= config.GlobalConf.MaxOnlyUserConnection {
			rsp.C = errUserMaxOnly
			r = util.JSONOut(rsp)
			userMutex.Unlock()
			return
		}
		userConnection[name]++
	} else {
		if len(userConnection) >= config.GlobalConf.MaxUserConnection {
			rsp.C = errUserMaxUser
			r = util.JSONOut(rsp)
			userMutex.Unlock()
			return
		}
		userConnection[name] = 1
	}
	userMutex.Unlock()
	r = util.JSONOut(rsp)
	return
}

func Init() bool {
	initString()
	initUser()
	if !config.GlobalConf.MasterSlave.IsMaster && !replicationSlave() {
		return false
	}
	go autoSaveData()
	return true
}

func init() {
	userConnection = make(map[string]int)
	userMutex = new(sync.Mutex)
}
