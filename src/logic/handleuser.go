package logic

import (
	"catchendb/src/user"
	"catchendb/src/util"
)

//import lgd "code.google.com/p/log4go"

func handleUserAut(req Req, tranObj *transaction) (ok bool, username string) {
	if req.C != CMDAUT {
		return
	}
	ok = user.VerifyPassword(req.UserName, req.PassWord)
	username = req.UserName
	return
}

func handleUserAdd(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.AddUser(req.UserName, req.PassWord, req.Privilege) {
		rsp.C = ERRUSERDUPLICATE
	}
	return util.JSONOut(rsp)
}

func handleUserDelete(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.DeleteUser(req.UserName) {
		rsp.C = ERRUSERNOTEXIST
	}
	return util.JSONOut(rsp)
}

func handleUserPass(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, req.PassWord, -1) {
		rsp.C = ERRUSERNOTEXIST
	}
	return util.JSONOut(rsp)
}

func handleUserPriv(req Req, tranObj *transaction) []byte {
	rsp := Rsp{}
	if !user.MotifyUserInfo(req.UserName, "", req.Privilege) {
		rsp.C = ERRUSERPRIVILEGE
	}
	return util.JSONOut(rsp)
}

func initUser() {
	registerCMD(CMDUADD, 4, handleUserAdd, TypeX)
	registerCMD(CMDUDEL, 2, handleUserDelete, TypeX)
	registerCMD(CMDUPAS, 3, handleUserPass, TypeX)
	registerCMD(CMDUPRI, 3, handleUserPriv, TypeX)
}
