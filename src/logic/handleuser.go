package logic

import (
	"catchendb/src/user"
	"catchendb/src/util"
)

//import lgd "catchendb/src/log"

func handleUserAut(req Req, tranObj *transaction) (ok bool, username string) {
	if req.C != cmdAut {
		return
	}
	ok = user.VerifyPassword(req.UserName, req.PassWord)
	username = req.UserName
	return
}

func handleUserAdd(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.AddUser(req.UserName, req.PassWord, req.Privilege) {
		rsp.C = errUserDuplicate
	}
	return util.JSONOut(rsp)
}

func handleUserDelete(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.DeleteUser(req.UserName) {
		rsp.C = errUserNotExist
	}
	return util.JSONOut(rsp)
}

func handleUserPass(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, req.PassWord, -1) {
		rsp.C = errUserNotExist
	}
	return util.JSONOut(rsp)
}

func handleUserPriv(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, "", req.Privilege) {
		rsp.C = errUserPrivilege
	}
	return util.JSONOut(rsp)
}

func initUser() {
	registerCMD(cmdUadd, 4, handleUserAdd, typeX)
	registerCMD(cmdUdel, 2, handleUserDelete, typeX)
	registerCMD(cmdUpas, 3, handleUserPass, typeX)
	registerCMD(cmdUpri, 3, handleUserPriv, typeX)
}
