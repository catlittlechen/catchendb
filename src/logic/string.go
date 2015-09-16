package logic

import (
	"catchendb/src/node"
	"catchendb/src/util"
	"strconv"
	"time"
)

import lgd "code.google.com/p/log4go"

var (
	nowTime int64
)

func handleSet(req Req, tranObj *transaction) []byte {
	rsp := Rsp{
		C: 0,
	}
	tid := 0
	if tranObj != nil {
		tid = tranObj.getID()
	}
	if !node.Put(req.Key, req.Value, 0, 0, tid) {
		lgd.Error("set fail! key[%s] value[%s]", req.Key, req.Value)
		rsp.C = ERRCMDSET
	}
	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(req.Key)
		d := new(node.Data)
		d.Key = req.Key
		d.Value = req.Value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(InsertType, d)
	} else {
		go replicationData(req)
	}
	return util.JSONOut(rsp)
}

func handleGet(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	nowTime = time.Now().Unix()
	if tranObj != nil && tranObj.isBegin() {
		if data := tranObj.getData(key); data != nil {
			if nowTime > data.StartTime {
				rsp.D = data.Value
			}
			return util.JSONOut(rsp)
		}
	}
	value, startTime, _ := node.Get(key)
	if nowTime > startTime {
		rsp.D = value
	}
	return util.JSONOut(rsp)
}

func handleDel(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	tid := 0
	if tranObj != nil {
		tid = tranObj.getID()
	}
	if !node.Del(key, tid) {
		lgd.Error("del fail! key[%s]", key)
		rsp.C = ERRCMDDEL
	}
	if tranObj != nil && tranObj.isBegin() {
		d := new(node.Data)
		d.Key = key
		tranObj.push(DeleteType, d)
	} else {
		go replicationData(req)
	}
	return util.JSONOut(rsp)
}

func handleSetEx(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	value := req.Value
	tid := 0
	if tranObj != nil {
		tid = tranObj.getID()
	}
	if !node.Put(key, value, req.StartTime, req.EndTime, tid) {
		lgd.Error("set fail! key[%s] value[%s]", key, value)
		rsp.C = ERRCMDSET
	}

	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(InsertType, d)
	} else {
		go replicationData(req)
	}
	return util.JSONOut(rsp)
}

func handleDelay(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	tid := 0
	if tranObj != nil {
		tid = tranObj.getID()
	}
	if !node.Set(key, req.StartTime, 0, tid) {
		lgd.Error("delay fail! key[%s] startTime[%d]", key, req.StartTime)
		rsp.C = ERRCMDDELAY
	}
	if tranObj != nil && tranObj.isBegin() {
		value, _, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = req.StartTime
		d.EndTime = end
		tranObj.push(UpdateType, d)
	} else {
		go replicationData(req)
	}
	return util.JSONOut(rsp)
}

func handleExpire(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	tid := 0
	if tranObj != nil {
		tid = tranObj.getID()
	}
	if !node.Set(key, 0, req.EndTime, tid) {
		lgd.Error("delay fail! key[%s] endTime[%d]", key, req.EndTime)
		rsp.C = ERRCMDEXPIRE
	}
	if tranObj != nil && tranObj.isBegin() {
		value, start, _ := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = req.EndTime
		tranObj.push(UpdateType, d)
	} else {
		go replicationData(req)
	}
	return util.JSONOut(rsp)
}

func handleTTL(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	nowTime = time.Now().Unix()
	if tranObj != nil && tranObj.isBegin() {
		if data := tranObj.getData(key); data != nil {
			if nowTime < data.EndTime {
				rsp.D = strconv.Itoa(int(data.EndTime - nowTime))
			} else {
				rsp.D = "-1"
			}
			return util.JSONOut(rsp)
		}
	}
	_, _, endTime := node.Get(key)
	if endTime < nowTime {
		rsp.D = "-1"
	} else {
		rsp.D = strconv.Itoa(int(endTime - nowTime))
	}
	return util.JSONOut(rsp)
}

func handleTTS(req Req, tranObj *transaction) []byte {
	key := req.Key
	rsp := Rsp{
		C: 0,
	}
	nowTime = time.Now().Unix()
	if tranObj != nil && tranObj.isBegin() {
		if data := tranObj.getData(key); data != nil {
			if nowTime > data.StartTime {
				rsp.D = strconv.Itoa(int(nowTime - data.StartTime))
			} else {
				rsp.D = "-1"
			}
			return util.JSONOut(rsp)
		}
	}
	_, startTime, _ := node.Get(key)
	if startTime < nowTime {
		rsp.D = "-1"
	} else {
		rsp.D = strconv.Itoa(int(startTime - nowTime))
	}
	return util.JSONOut(rsp)

}

func initString() {
	registerCMD(CMDSET, 3, handleSet, TypeW)
	registerCMD(CMDGET, 2, handleGet, TypeR)
	registerCMD(CMDDEL, 2, handleDel, TypeW)
	registerCMD(CMDSETEX, 5, handleSetEx, TypeW)
	registerCMD(CMDDELAY, 3, handleDelay, TypeW)
	registerCMD(CMDEXPIRE, 3, handleExpire, TypeW)
	registerCMD(CMDTTL, 2, handleTTL, TypeR)
	registerCMD(CMDTTS, 2, handleTTS, TypeR)
}
