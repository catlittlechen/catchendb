package logic

import (
	"catchendb/src/config"
	"catchendb/src/util"
)

import lgd "code.google.com/p/log4go"

const (
	TYPE_R = 4
	TYPE_W = 2
	TYPE_X = 1
)

var functionAction map[string]func(Req, *transaction) []byte
var functionArgv map[string]int
var functionType map[string]int

func mapAction(req Req, privilege int, replication bool, tranObj *transaction) []byte {
	rsp := Rsp{
		C: ERR_CMD_MISS,
	}
	if function, ok := functionAction[req.C]; ok {
		typ := functionType[req.C]
		switch typ {
		case TYPE_R:
			if privilege < 4 || privilege > 7 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
		case TYPE_W:
			if (privilege/2 != 1 && privilege/2 != 3) || (!config.GlobalConf.MasterSlave.IsMaster && !replication) {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
		case TYPE_X:
			if privilege != 1 && privilege != 3 && privilege != 7 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JsonOut(rsp)
			}
			if tranObj.isBegin() {
				rsp.C = ERR_TRA_USER
				return util.JsonOut(rsp)
			}
		}
		if len(req.Key) == 0 && len(req.UserName) == 0 {
			lgd.Error("argv[%+v] is illegal", req)
			rsp.C = ERR_PARSE_MISS
			return util.JsonOut(rsp)
		}
		return function(req, tranObj)
	}

	return util.JsonOut(rsp)
}

func registerCMD(key string, argvcount int, function func(Req, *transaction) []byte, typ int) {
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
	functionAction = make(map[string]func(Req, *transaction) []byte)
	functionArgv = make(map[string]int)
	functionType = make(map[string]int)
}
