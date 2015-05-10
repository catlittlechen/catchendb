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
		rsp.C = ERR_CMD_SET
	}
	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(req.Key)
		d := new(node.Data)
		d.Key = req.Key
		d.Value = req.Value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(INSERT_TYPE, d)
	} else {
		go replicationData(req)
	}
	return util.JsonOut(rsp)
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
			return util.JsonOut(rsp)
		}
	}
	value, startTime, _ := node.Get(key)
	if nowTime > startTime {
		rsp.D = value
	}
	return util.JsonOut(rsp)
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
		rsp.C = ERR_CMD_DEL
	}
	if tranObj != nil && tranObj.isBegin() {
		d := new(node.Data)
		d.Key = key
		tranObj.push(DELETE_TYPE, d)
	} else {
		go replicationData(req)
	}
	return util.JsonOut(rsp)
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
		rsp.C = ERR_CMD_SET
	}

	if tranObj != nil && tranObj.isBegin() {
		_, start, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(INSERT_TYPE, d)
	} else {
		go replicationData(req)
	}
	return util.JsonOut(rsp)
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
		rsp.C = ERR_CMD_DELAY
	}
	if tranObj != nil && tranObj.isBegin() {
		value, _, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = req.StartTime
		d.EndTime = end
		tranObj.push(UPDATE_TYPE, d)
	} else {
		go replicationData(req)
	}
	return util.JsonOut(rsp)
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
		rsp.C = ERR_CMD_EXPIRE
	}
	if tranObj != nil && tranObj.isBegin() {
		value, start, _ := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = req.EndTime
		tranObj.push(UPDATE_TYPE, d)
	} else {
		go replicationData(req)
	}
	return util.JsonOut(rsp)
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
			return util.JsonOut(rsp)
		}
	}
	_, _, endTime := node.Get(key)
	if endTime < nowTime {
		rsp.D = "-1"
	} else {
		rsp.D = strconv.Itoa(int(endTime - nowTime))
	}
	return util.JsonOut(rsp)
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
			return util.JsonOut(rsp)
		}
	}
	_, startTime, _ := node.Get(key)
	if startTime < nowTime {
		rsp.D = "-1"
	} else {
		rsp.D = strconv.Itoa(int(startTime - nowTime))
	}
	return util.JsonOut(rsp)

}

func initString() {
	registerCMD(CMD_SET, 3, handleSet, TYPE_W)
	registerCMD(CMD_GET, 2, handleGet, TYPE_R)
	registerCMD(CMD_DEL, 2, handleDel, TYPE_W)
	registerCMD(CMD_SETEX, 5, handleSetEx, TYPE_W)
	registerCMD(CMD_DELAY, 3, handleDelay, TYPE_W)
	registerCMD(CMD_EXPIRE, 3, handleExpire, TYPE_W)
	registerCMD(CMD_TTL, 2, handleTTL, TYPE_R)
	registerCMD(CMD_TTS, 2, handleTTS, TYPE_R)
}
