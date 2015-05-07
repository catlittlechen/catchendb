package logic

import (
	"catchendb/src/config"
	"catchendb/src/user"
	"catchendb/src/util"
	"net"
	"net/url"
	"sync"
)

import lgd "code.google.com/p/log4go"

var (
	userConnection map[string]int
	userMutex      *sync.Mutex
)

func ReplicationLogic(conn *net.TCPConn) {
	data := make([]byte, 1024)
	count, err := conn.Read(data)
	if err != nil {
		lgd.Error("read error[%s]", err)
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
		lgd.Error("read error[%s]", err)
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
	errRes := util.JsonOut(Rsp{
		C: ERR_URL_PARSE,
	})
	tranObj := new(transaction)
	var urlStr string
	var keyword url.Values
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

		urlStr = string(data[:count])
		//lgd.Trace("%s", urlStr)
		keyword, err = url.ParseQuery(urlStr)
		if err != nil {
			lgd.Warn("ParseQuery fail with the url %s", urlStr)
			_, err = conn.Write(errRes)
			if err != nil {
				lgd.Warn("write error[%s]", err)
				return
			}
			continue
		}

		ok, res = clientTransactionLogic(keyword, tranObj)
		if !ok {
			_, err = conn.Write(res)
			if err != nil {
				lgd.Warn("write error[%s]", err)
				return
			}
			continue
		}
		privilege = user.GetPrivilege(name)
		_, err = conn.Write(mapAction(keyword, privilege, false, tranObj))
		if err != nil {
			lgd.Warn("write error[%s]", err)
			return
		}
	}
}

func clientTransactionLogic(keyword url.Values, tranObj *transaction) (normal bool, res []byte) {
	rsp := Rsp{}
	switch keyword.Get(URL_CMD) {
	case CMD_BEGIN:
		if tranObj.isBegin() {
			rsp.C = ERR_TRA_BEGIN
		} else {
			tranObj.init()
		}
		res = util.JsonOut(rsp)
		return
	case CMD_ROLLBACK:
		if tranObj.isBegin() {
			rsp.C = tranObj.rollback()
		} else {
			rsp.C = ERR_TRA_NO_BEGIN
		}
		res = util.JsonOut(rsp)
		return
	case CMD_COMMIT:
		if tranObj.isBegin() {
			rsp.C = tranObj.commit()
		} else {
			rsp.C = ERR_TRA_NO_BEGIN
		}
		res = util.JsonOut(rsp)
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
		userConnection[name] -= 1
	}
	return
}

func aut(data []byte) (ok bool, name string, r []byte) {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)

	keyword, err := url.ParseQuery(urlStr)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		r = util.JsonOut(rsp)
		return
	}
	ok, name = handleUserAut(keyword, nil)
	if !ok {
		rsp.C = ERR_ACCESS_DENIED
		r = util.JsonOut(rsp)
		return
	}
	userMutex.Lock()
	if _, ok2 := userConnection[name]; ok2 {
		if userConnection[name] >= config.GlobalConf.MaxOnlyUserConnection {
			rsp.C = ERR_USER_MAX_ONLY
			r = util.JsonOut(rsp)
			userMutex.Unlock()
			return
		}
		userConnection[name] += 1
	} else {
		if len(userConnection) >= config.GlobalConf.MaxUserConnection {
			rsp.C = ERR_USER_MAX_USER
			r = util.JsonOut(rsp)
			userMutex.Unlock()
			return
		}
		userConnection[name] = 1
	}
	userMutex.Unlock()
	r = util.JsonOut(rsp)
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
