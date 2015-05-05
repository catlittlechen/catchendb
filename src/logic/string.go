package logic

import (
	"catchendb/src/node"
	"catchendb/src/util"
	"net/url"
	"strconv"
	"time"
)

import lgd "code.google.com/p/log4go"

var (
	nowTime int64
)

func handleSet(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	value := keyword.Get(URL_VALUE)
	rsp := Rsp{
		C: 0,
	}
	if !node.Put(key, value, 0, 0, tranObj.getID()) {
		lgd.Error("set fail! key[%s] value[%s]", key, value)
		rsp.C = ERR_CMD_SET
	}
	go replicationData(keyword)
	if tranObj.isBegin() {
		_, start, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(INSERT_TYPE, d)
	}
	return util.JsonOut(rsp)
}

func handleGet(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	value, startTime, _ := node.Get(key)
	nowTime = time.Now().Unix()
	if nowTime > startTime {
		rsp.D = value
	}
	return util.JsonOut(rsp)
}

func handleDel(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	if !node.Del(key, tranObj.getID()) {
		lgd.Error("del fail! key[%s]", key)
		rsp.C = ERR_CMD_DEL
	}
	go replicationData(keyword)
	if tranObj.isBegin() {
		d := new(node.Data)
		d.Key = key
		tranObj.push(DELETE_TYPE, d)
	}
	return util.JsonOut(rsp)
}

func handleSetEx(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	value := keyword.Get(URL_VALUE)
	startTime, err := strconv.Atoi(keyword.Get(URL_START))
	if err != nil {
		rsp.C = ERR_START_TIME
		return util.JsonOut(rsp)
	}
	endTime, err := strconv.Atoi(keyword.Get(URL_END))
	if err != nil {
		rsp.C = ERR_END_TIME
		return util.JsonOut(rsp)
	}

	if !node.Put(key, value, int64(startTime), int64(endTime), tranObj.getID()) {
		lgd.Error("set fail! key[%s] value[%s]", key, value)
		rsp.C = ERR_CMD_SET
	}
	go replicationData(keyword)
	if tranObj.isBegin() {
		_, start, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = end
		tranObj.push(INSERT_TYPE, d)
	}
	return util.JsonOut(rsp)
}

func handleDelay(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	startTime, err := strconv.Atoi(keyword.Get(URL_START))
	if err != nil {
		rsp.C = ERR_START_TIME
		return util.JsonOut(rsp)
	}
	if !node.Set(key, int64(startTime), 0, tranObj.getID()) {
		lgd.Error("delay fail! key[%s] startTime[%d]", key, startTime)
		rsp.C = ERR_CMD_DELAY
	}
	go replicationData(keyword)
	if tranObj.isBegin() {
		value, _, end := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = int64(startTime)
		d.EndTime = end
		tranObj.push(UPDATE_TYPE, d)
	}
	return util.JsonOut(rsp)
}

func handleExpire(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	endTime, err := strconv.Atoi(keyword.Get(URL_END))
	if err != nil {
		rsp.C = ERR_END_TIME
		return util.JsonOut(rsp)
	}
	if !node.Set(key, 0, int64(endTime), tranObj.getID()) {
		lgd.Error("delay fail! key[%s] endTime[%d]", key, endTime)
		rsp.C = ERR_CMD_EXPIRE
	}
	go replicationData(keyword)
	if tranObj.isBegin() {
		value, start, _ := node.Get(key)
		d := new(node.Data)
		d.Key = key
		d.Value = value
		d.StartTime = start
		d.EndTime = int64(endTime)
		tranObj.push(UPDATE_TYPE, d)
	}
	return util.JsonOut(rsp)
}

func handleTTL(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	_, _, endTime := node.Get(key)
	nowTime = time.Now().Unix()
	if endTime-nowTime < 0 {
		rsp.D = "-1"
	} else {
		rsp.D = strconv.Itoa(int(endTime - nowTime))
	}
	return util.JsonOut(rsp)
}

func handleTTS(keyword url.Values, tranObj *transaction) []byte {
	key := keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	_, startTime, _ := node.Get(key)
	nowTime = time.Now().Unix()
	if startTime-nowTime < 0 {
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
