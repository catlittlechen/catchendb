package logic

import (
	"catchendb/src/config"
	"catchendb/src/util"
)

import lgd "code.google.com/p/log4go"

const (
	TypeR = 4
	TypeW = 2
	TypeX = 1
)

var functionAction map[string]func(Req, *transaction) []byte
var functionArgv map[string]int
var functionType map[string]int

func mapAction(req Req, privilege int, replication bool, tranObj *transaction) []byte {
	rsp := Rsp{
		C: ERRCMDMISS,
	}
	if function, ok := functionAction[req.C]; ok {
		typ := functionType[req.C]
		switch typ {
		case TypeR:
			if privilege < 4 || privilege > 7 {
				rsp.C = ERRACCESSDENIED
				return util.JSONOut(rsp)
			}
		case TypeW:
			if (privilege/2 != 1 && privilege/2 != 3) || (!config.GlobalConf.MasterSlave.IsMaster && !replication) {
				rsp.C = ERRACCESSDENIED
				return util.JSONOut(rsp)
			}
		case TypeX:
			if privilege != 1 && privilege != 3 && privilege != 7 {
				rsp.C = ERRACCESSDENIED
				return util.JSONOut(rsp)
			}
			if tranObj.isBegin() {
				rsp.C = ERRTRAUSER
				return util.JSONOut(rsp)
			}
		}
		if len(req.Key) == 0 && len(req.UserName) == 0 {
			lgd.Error("argv[%+v] is illegal", req)
			rsp.C = ERRPARSEMISS
			return util.JSONOut(rsp)
		}
		return function(req, tranObj)
	}

	return util.JSONOut(rsp)
}

func registerCMD(key string, argvcount int, function func(Req, *transaction) []byte, typ int) {
	if _, ok := functionAction[key]; ok {
		lgd.Error("duplicate key %s", key)
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
