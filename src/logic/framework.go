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
		C: errCmdMiss,
	}
	if function, ok := functionAction[req.C]; ok {
		typ := functionType[req.C]
		switch typ {
		case typeR:
			if privilege < 4 || privilege > 7 {
				rsp.C = errAccessDenied
				return util.JSONOut(rsp)
			}
		case typeW:
			if (privilege/2 != 1 && privilege/2 != 3) || (!config.GlobalConf.MasterSlave.IsMaster && !replication) {
				rsp.C = errAccessDenied
				return util.JSONOut(rsp)
			}
		case typeX:
			if privilege != 1 && privilege != 3 && privilege != 7 {
				rsp.C = errAccessDenied
				return util.JSONOut(rsp)
			}
			if tranObj.isBegin() {
				rsp.C = errTraUser
				return util.JSONOut(rsp)
			}
		}
		if len(req.Key) == 0 && len(req.UserName) == 0 {
			lgd.Errorf("argv[%+v] is illegal", req)
			rsp.C = errParseMiss
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
