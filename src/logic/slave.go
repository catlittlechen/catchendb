package logic

import (
	"catchendb/src/config"
	"catchendb/src/node"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"sync"
)

import lgd "code.google.com/p/log4go"

var (
	replicationChannel map[string]chan node.Data
	channelMutex       *sync.Mutex
)

func Replication(conn *net.TCPConn) {
	replicationMaster(conn)
	return
}

func replicationMaster(conn *net.TCPConn) {

}

func replicationSlave() bool {

	data := make([]byte, 10240)

	serverHost := fmt.Sprintf("%s:%d", config.GlobalConf.MasterSlave.IP, config.GlobalConf.MasterSlave.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverHost)
	if err != nil {
		lgd.Error("Fatal Error %s", err)
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		lgd.Error("Fatal Error %s", err)
		return false
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_AUT)
	urlData.Add(URL_USER, config.GlobalConf.MasterSlave.UserName)
	urlData.Add(URL_PASS, config.GlobalConf.MasterSlave.PassWord)
	_, err = conn.Write([]byte(urlData.Encode()))
	if err != nil {
		lgd.Error("Fatal Error %s", err)
		return false
	}
	count, err := conn.Read(data)
	if err != nil {
		lgd.Error("Fatal Error %s", err)
		return false
	}
	var rsp Rsp
	err = json.Unmarshal(data[:count], &rsp)
	if err != nil {
		lgd.Error("Fatal Error %s", err)
		return false
	}

	if rsp.C != 0 {
		lgd.Error("ccdb>ERROR %d Access denied for user '%s'@'%s' (using password: YES)", rsp.C, config.GlobalConf.MasterSlave.UserName, serverHost)
		return false
	}

	go syncData(conn)
	return true
}

func syncData(conn *net.TCPConn) {
	var err error
	data := make([]byte, 10240)
	count := 0
	for {
		count, err = conn.Read(data)
		if err != nil {
			lgd.Error("Fatal Error %s", err)
			return
		}
		keyword, err := url.ParseQuery(string(data[:count]))
		if err != nil {
			lgd.Error("Fatal Error %s", err)
			return
		}
		_, err = conn.Write(mapAction(keyword, 2, true))
		if err != nil {
			lgd.Error("Sync Fatal Error %s", err)
			return
		}
	}
}

func init() {
	replicationChannel = make(map[string]chan node.Data)
	channelMutex = new(sync.Mutex)
}
