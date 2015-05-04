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
	conn.Write(res)
	if !ok {
		return

	}
	replicationMaster(name, conn)
	disConnection(name)
}

func ClientLogic(conn *net.TCPConn) {
	data := make([]byte, 1024)
	count, err := conn.Read(data)
	if err != nil {
		lgd.Error("read error[%s]", err)
		return

	}
	ok, name, res := aut(data[:count])
	conn.Write(res)
	if !ok {
		return
	}
	for {
		count, err = conn.Read(data)
		if err != nil {
			lgd.Warn("read error[%s]", err)
			disConnection(name)
			return
		}
		res := lyw(data[:count], name, false)
		conn.Write(res)
	}
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

func lyw(data []byte, name string, replication bool) []byte {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)

	keyword, err := url.ParseQuery(urlStr)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		return util.JsonOut(rsp)
	}
	privilege := user.GetPrivilege(name)
	return mapAction(keyword, privilege, replication)
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
	ok, name = handleUserAut(keyword)
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
