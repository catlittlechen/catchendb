package logic

import (
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

var functionAction map[string]func(url.Values) []byte
var functionArgv map[string]int

func mapAction(keyword url.Values) []byte {
	rsp := Rsp{
		C: ERR_CMD_MISS,
	}
	key := keyword.Get(URL_CMD)
	if function, ok := functionAction[key]; ok {
		if len(keyword) != functionArgv[key] {
			lgd.Error("argv[%+v] is illegal", keyword)
			rsp.C = ERR_PARSE_MISS
			return util.JsonOut(rsp)
		}
		return function(keyword)
	}

	return util.JsonOut(rsp)
}

func registerCMD(key string, argvcount int, function func(url.Values) []byte) {
	if _, ok := functionAction[key]; ok {
		lgd.Error("duplicate key %s", key)
		return
	} else {
		lgd.Info("reister cmd %s", key)
		functionAction[key] = function
		functionArgv[key] = argvcount
	}
}

func init() {
	functionAction = make(map[string]func(url.Values) []byte)
	functionArgv = make(map[string]int)
}
