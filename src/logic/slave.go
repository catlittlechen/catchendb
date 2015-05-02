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
	replicationChannel map[string]*chan url.Values
	channelMutex       *sync.Mutex
)

func addReplicationChannel(name string, chans *chan url.Values) {
	channelMutex.Lock()
	replicationChannel[name] = chans
	channelMutex.Unlock()
}

func deleteReplicationChannel(name string) {
	channelMutex.Lock()
	delete(replicationChannel, name)
	channelMutex.Unlock()
}

func replicationData(data url.Values) {
	channelMutex.Lock()
	for _, chans := range replicationChannel {
		go func() {
			(*chans) <- data
		}()
	}
	channelMutex.Unlock()
}

func Replication(name string, conn *net.TCPConn) {
	replicationMaster(name, conn)
	return
}

func replicationMaster(name string, conn *net.TCPConn) {
	channelReplication := make(chan url.Values, 1000)
	addReplicationChannel(name, &channelReplication)
	defer func() {
		deleteReplicationChannel(name)
		close(channelReplication)
	}()
	channel := make(chan node.Data, 1000)
	go node.OutPutData(channel)
	urlData := url.Values{}
	var err error
	var count int
	var rsp Rsp
	data := make([]byte, 1024)
	for {
		d := <-channel
		if len(d.Key) == 0 {
			break
		}
		urlData = url.Values{}
		urlData.Add(URL_CMD, CMD_SETEX)
		urlData.Add(URL_KEY, d.Key)
		urlData.Add(URL_VALUE, d.Value)
		urlData.Add(URL_START, fmt.Sprintf("%d", d.StartTime))
		urlData.Add(URL_END, fmt.Sprintf("%d", d.EndTime))
		_, err = conn.Write([]byte(urlData.Encode()))
		if err != nil {
			lgd.Error("Sync Error %s", err)
			return
		}
		count, err = conn.Read(data)
		if err != nil {
			lgd.Error("Sync Fatal Error %s", err)
			return
		}

		err = json.Unmarshal(data[:count], &rsp)
		if err != nil {
			lgd.Error("Sync Fatal Error %s", err)
			return

		}
		if rsp.C != 0 {
			lgd.Error("Sync Fatal Error %s", err)
			return
		}
	}
	close(channel)
	for {
		urlData := <-channelReplication
		_, err = conn.Write([]byte(urlData.Encode()))
		if err != nil {
			lgd.Error("Sync Error %s", err)
			return
		}
		count, err = conn.Read(data)
		if err != nil {
			lgd.Error("Sync Fatal Error %s", err)
			return
		}

		err = json.Unmarshal(data[:count], &rsp)
		if err != nil {
			lgd.Error("Sync Fatal Error %s", err)
			return

		}
		if rsp.C != 0 {
			lgd.Error("Sync Fatal Error %s", err)
			return
		}
	}
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
	replicationChannel = make(map[string]*chan url.Values)
	channelMutex = new(sync.Mutex)
}
