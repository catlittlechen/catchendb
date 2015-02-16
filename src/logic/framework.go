package logic

import (
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

var functionAction map[string]func(url.Values) []byte

func mapAction(keyword url.Values) []byte {
	key := keyword.Get(URL_CMD)
	if function, ok := functionAction[key]; ok {
		return function(keyword)
	}
	rsp := Rsp{
		C: ERR_CMD_MISS,
	}
	return util.JsonOut(rsp)
}

func resgiterCMD(key string, function func(url.Values) []byte) {
	if _, ok := functionAction[key]; ok {
		lgd.Error("duplicate key %s", key)
		return
	} else {
		lgd.Info("reister cmd %s", key)
		functionAction[key] = function
	}
}

func init() {
	functionAction = make(map[string]func(url.Values) []byte)
}
