package logic

import (
	"catchendb/src/user"
	"catchendb/src/util"
)

//import lgd "catchendb/src/log"

func handleUserAut(req Req, tranObj *transaction) (ok bool, username string) {
	if req.C != CMD_AUT {
		return
	}
	ok = user.VerifyPassword(req.UserName, req.PassWord)
	username = req.UserName
	return
}

func handleUserAdd(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.AddUser(req.UserName, req.PassWord, req.Privilege) {
		rsp.C = ERR_USER_DUPLICATE
	}
	return util.JsonOut(rsp)
}

func handleUserDelete(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.DeleteUser(req.UserName) {
		rsp.C = ERR_USER_NOT_EXIST
	}
	return util.JsonOut(rsp)
}

func handleUserPass(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, req.PassWord, -1) {
		rsp.C = ERR_USER_NOT_EXIST
	}
	return util.JsonOut(rsp)
}

func handleUserPriv(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, "", req.Privilege) {
		rsp.C = ERR_USER_PRIVILEGE
	}
	return util.JsonOut(rsp)
}

func initUser() {
	registerCMD(CMD_UADD, 4, handleUserAdd, TYPE_X)
	registerCMD(CMD_UDEL, 2, handleUserDelete, TYPE_X)
	registerCMD(CMD_UPAS, 3, handleUserPass, TYPE_X)
	registerCMD(CMD_UPRI, 3, handleUserPriv, TYPE_X)
}
