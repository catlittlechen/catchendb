package logic

import (
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

const (
	TYPE_R = 4
	TYPE_W = 2
	TYPE_X = 1
)

var functionAction map[string]func(url.Values) []byte
var functionArgv map[string]int
var functionType map[string]int

func mapAction(keyword url.Values, privilege int) []byte {
	rsp := Rsp{
		C: ERR_CMD_MISS,
	}
	key := keyword.Get(URL_CMD)
	if function, ok := functionAction[key]; ok {
		typ := functionType[key]
		switch typ {
		case TYPE_R:
			if privilege < 4 || privilege > 7 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
		case TYPE_W:
			if privilege/2 != 1 && privilege/2 != 3 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
		case TYPE_X:
			if privilege < 0 || privilege > 7 && privilege%2 != 1 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
		}
		if len(keyword) != functionArgv[key] {
			lgd.Error("argv[%+v] is illegal", keyword)
			rsp.C = ERR_PARSE_MISS
			return util.JsonOut(rsp)
		}
		return function(keyword)
	}

	return util.JsonOut(rsp)
}

func registerCMD(key string, argvcount int, function func(url.Values) []byte, typ int) {
	if _, ok := functionAction[key]; ok {
		lgd.Error("duplicate key %s", key)
		return
	} else {
		lgd.Info("reister cmd %s", key)
		functionAction[key] = function
		functionArgv[key] = argvcount
		functionType[key] = typ
	}
}

func init() {
	functionAction = make(map[string]func(url.Values) []byte)
	functionArgv = make(map[string]int)
	functionType = make(map[string]int)
}
