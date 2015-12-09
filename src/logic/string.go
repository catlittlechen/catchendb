package logic

import (
	"catchendb/src/data"
	"catchendb/src/node"
	"catchendb/src/util"
	"strconv"
	"time"
)

import lgd "catchendb/src/log"

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
		lgd.Errorf("set fail! key[%s] value[%s]", req.Key, req.Value)
		rsp.C = errCmdSet
	}
	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(req.Key)
		d := new(data.Data)
		d.Key = req.Key
		d.Value = req.Value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(insertType, d)
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
		lgd.Errorf("del fail! key[%s]", key)
		rsp.C = errCmdDel
	}
	if tranObj != nil && tranObj.isBegin() {
		d := new(data.Data)
		d.Key = key
		tranObj.push(deleteType, d)
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
		lgd.Errorf("set fail! key[%s] value[%s]", key, value)
		rsp.C = errCmdSet
	}

	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(key)
		d := new(data.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(insertType, d)
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
		lgd.Errorf("delay fail! key[%s] startTime[%d]", key, req.StartTime)
		rsp.C = errCmdDelAY
	}
	if tranObj != nil && tranObj.isBegin() {
		value, _, end := node.Get(key)
		d := new(data.Data)
		d.Key = key
		d.Value = value
		d.StartTime = req.StartTime
		d.EndTime = end
		tranObj.push(updateType, d)
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
		lgd.Errorf("delay fail! key[%s] endTime[%d]", key, req.EndTime)
		rsp.C = errCmdExpire
	}
	if tranObj != nil && tranObj.isBegin() {
		value, start, _ := node.Get(key)
		d := new(data.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = req.EndTime
		tranObj.push(updateType, d)
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
	registerCMD(cmdSet, 3, handleSet, typeW)
	registerCMD(cmdGet, 2, handleGet, typeR)
	registerCMD(cmdDel, 2, handleDel, typeW)
	registerCMD(cmdSetEX, 5, handleSetEx, typeW)
	registerCMD(cmdDelAY, 3, handleDelay, typeW)
	registerCMD(cmdExpire, 3, handleExpire, typeW)
	registerCMD(cmdTTL, 2, handleTTL, typeR)
	registerCMD(cmdTTS, 2, handleTTS, typeR)
}
