package logic

import (
	"catchendb/src/config"
	"catchendb/src/util"
)

import lgd "catchendb/src/log"

const (
	typeR = 4
	typeW = 2
	typeX = 1
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
		case typeR:
			if privilege < 4 || privilege > 7 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JSONOut(rsp)
			}
		case typeW:
			if (privilege/2 != 1 && privilege/2 != 3) || (!config.GlobalConf.MasterSlave.IsMaster && !replication) {
				rsp.C = ERR_ACCESS_DENIED
				return util.JSONOut(rsp)
			}
		case typeX:
			if privilege != 1 && privilege != 3 && privilege != 7 {
				rsp.C = ERR_ACCESS_DENIED
				return util.JSONOut(rsp)
			}
			if tranObj.isBegin() {
				rsp.C = ERR_TRA_USER
				return util.JSONOut(rsp)
			}
		}
		if len(req.Key) == 0 && len(req.UserName) == 0 {
			lgd.Errorf("argv[%+v] is illegal", req)
			rsp.C = ERR_PARSE_MISS
			return util.JSONOut(rsp)
		}
		return function(req, tranObj)
	}

	return util.JSONOut(rsp)
}

func registerCMD(key string, argvcount int, function func(Req, *transaction) []byte, typ int) {
	if _, ok := functionAction[key]; ok {
		lgd.Errorf("duplicate key %s", key)
	} else {
		lgd.Info("reister cmd %s", key)
		functionAction[key] = function
		functionArgv[key] = argvcount
		functionType[key] = typ
	}
	return
}

func init() {
	functionAction = make(map[string]func(Req, *transaction) []byte)
	functionArgv = make(map[string]int)
	functionType = make(map[string]int)
}
