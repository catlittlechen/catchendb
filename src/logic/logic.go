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

func ClientLogic(conn *net.TCPConn) {
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

	privilege := 0
	errRes := util.JSONOut(Rsp{
		C: ERR_URL_PARSE,
	})
	tranObj := new(transaction)
	var req Req
	for {
		defer func() {
			if tranObj.isBegin() {
				tranObj.rollback()
			}
		}()
		count, err = conn.Read(data)
		if err != nil {
			lgd.Warn("read error[%s]", err)
			return
		}

		req = Req{}
		err = json.Unmarshal(data[:count], &req)
		if err != nil {
			lgd.Warn("ParseQuery fail with the data %s", string(data[:count]))
			_, err = conn.Write(errRes)
			if err != nil {
				lgd.Warn("write error[%s]", err)
				return
			}
			continue
		}

		ok, res = clientTransactionLogic(req, tranObj)
		if !ok {
			_, err = conn.Write(res)
			if err != nil {
				lgd.Warn("write error[%s]", err)
				return
			}
			continue
		}
		privilege = user.GetPrivilege(name)
		_, err = conn.Write(mapAction(req, privilege, false, tranObj))
		if err != nil {
			lgd.Warn("write error[%s]", err)
			return
		}
	}
}

func clientTransactionLogic(req Req, tranObj *transaction) (normal bool, res []byte) {
	rsp := Rsp{}
	switch req.C {
	case CMD_BEGIN:
		if tranObj.isBegin() {
			rsp.C = ERR_TRA_BEGIN
		} else {
			tranObj.init()
		}
		res = util.JSONOut(rsp)
		return
	case CMD_ROLLBACK:
		if tranObj.isBegin() {
			rsp.C = tranObj.rollback()
		} else {
			rsp.C = ERR_TRA_NO_BEGIN
		}
		res = util.JSONOut(rsp)
		return
	case CMD_COMMIT:
		if tranObj.isBegin() {
			rsp.C = tranObj.commit()
		} else {
			rsp.C = ERR_TRA_NO_BEGIN
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
		rsp.C = ERR_URL_PARSE
		r = util.JSONOut(rsp)
		return
	}
	ok, name = handleUserAut(req, nil)
	if !ok {
		rsp.C = ERR_ACCESS_DENIED
		r = util.JSONOut(rsp)
		return
	}
	userMutex.Lock()
	if _, ok2 := userConnection[name]; ok2 {
		if userConnection[name] >= config.GlobalConf.MaxOnlyUserConnection {
			rsp.C = ERR_USER_MAX_ONLY
			r = util.JSONOut(rsp)
			userMutex.Unlock()
			return
		}
		userConnection[name]++
	} else {
		if len(userConnection) >= config.GlobalConf.MaxUserConnection {
			rsp.C = ERR_USER_MAX_USER
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
