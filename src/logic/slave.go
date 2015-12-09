package logic

import (
	"catchendb/src/config"
	"catchendb/src/data"
	"catchendb/src/node"
	"catchendb/src/util"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

import lgd "catchendb/src/log"

var (
	replicationChannel map[string]*chan Req
	channelMutex       *sync.Mutex
)

func addReplicationChannel(name string, chans *chan Req) {
	channelMutex.Lock()
	replicationChannel[name] = chans
	channelMutex.Unlock()
}

func deleteReplicationChannel(name string) {
	channelMutex.Lock()
	delete(replicationChannel, name)
	channelMutex.Unlock()
}

func replicationData(req Req) {
	channelMutex.Lock()
	for _, chans := range replicationChannel {
		go func(channel chan Req) {
			channel <- req
		}(*chans)
	}
	channelMutex.Unlock()
}

func translateData(conn *net.TCPConn, req []byte) (ok bool) {
	var err error
	var rsp Rsp
	var count int
	data := make([]byte, 1024)
	if _, err = conn.Write(util.JSONOut(req)); err != nil {
		lgd.Errorf("Sync Error %s", err)
		return
	}

	if count, err = conn.Read(data); err != nil {
		lgd.Errorf("Sync Fatal Error %s", err)
		return
	}

	if err = json.Unmarshal(data[:count], &rsp); err != nil {
		lgd.Errorf("Sync Fatal Error %s", err)
		return
	}

	if rsp.C != 0 {
		lgd.Errorf("Sync Fatal Error %s", err)
		return
	}
	return true
}

func replicationMaster(name string, conn *net.TCPConn) {
	channelReplication := make(chan Req, 1000)
	addReplicationChannel(name, &channelReplication)
	defer func() {
		deleteReplicationChannel(name)
		close(channelReplication)
	}()
	channel := make(chan data.Data, 1000)
	go node.OutPutData(channel)
	var req Req
	var ok bool
	for {
		d := <-channel
		if len(d.Key) == 0 {
			break
		}
		req = Req{
			C:         CMD_SETEX,
			Key:       d.Key,
			Value:     d.Value,
			StartTime: d.StartTime,
			EndTime:   d.EndTime,
		}
		if ok = translateData(conn, util.JSONOut(req)); !ok {
			return
		}
	}
	close(channel)
	for {
		req = <-channelReplication
		if ok = translateData(conn, util.JSONOut(req)); !ok {
			return
		}
	}
}

func replicationSlave() bool {

	data := make([]byte, 10240)

	serverHost := fmt.Sprintf("%s:%d", config.GlobalConf.MasterSlave.IP, config.GlobalConf.MasterSlave.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverHost)
	if err != nil {
		lgd.Errorf("Fatal Error %s", err)
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		lgd.Errorf("Fatal Error %s", err)
		return false
	}
	req := Req{
		C:        CMD_AUT,
		UserName: config.GlobalConf.MasterSlave.UserName,
		PassWord: config.GlobalConf.MasterSlave.PassWord,
	}
	_, err = conn.Write(util.JSONOut(req))
	if err != nil {
		lgd.Errorf("Fatal Error %s", err)
		return false
	}
	count, err := conn.Read(data)
	if err != nil {
		lgd.Errorf("Fatal Error %s", err)
		return false
	}
	var rsp Rsp
	err = json.Unmarshal(data[:count], &rsp)
	if err != nil {
		lgd.Errorf("Fatal Error %s", err)
		return false
	}

	if rsp.C != 0 {
		lgd.Errorf("ccdb>ERROR %d Access denied for user '%s'@'%s' (using password: YES)", rsp.C, config.GlobalConf.MasterSlave.UserName, serverHost)
		return false
	}

	go syncData(conn)
	return true
}

func syncData(conn *net.TCPConn) {
	var err error
	var req Req
	data := make([]byte, 10240)
	count := 0
	for {
		count, err = conn.Read(data)
		if err != nil {
			lgd.Errorf("Fatal Error %s", err)
			return
		}
		req = Req{}
		err = json.Unmarshal(data[:count], &req)
		if err != nil {
			lgd.Errorf("Fatal Error %s", err)
			return
		}
		_, err = conn.Write(mapAction(req, 2, true, nil))
		if err != nil {
			lgd.Errorf("Sync Fatal Error %s", err)
			return
		}
	}
}

func init() {
	replicationChannel = make(map[string]*chan Req)
	channelMutex = new(sync.Mutex)
}
