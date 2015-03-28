package logic

import (
	"catchendb/src/user"
	"catchendb/src/util"
	"net/url"
	"strconv"
)

//import lgd "code.google.com/p/log4go"

func handleUserAut(keyword url.Values) (ok bool, privilege int) {
	code := keyword.Get(URL_CMD)
	username := keyword.Get(URL_USER)
	password := keyword.Get(URL_PASS)
	if code != CMD_AUT {
		return
	}
	if ok = user.VerifyPassword(username, password); ok {
		privilege = user.GetPrivilege(username)
	}
	return
}

func handleUserAdd(keyword url.Values) []byte {
	rsp := Rsp{}
	username := keyword.Get(URL_USER)
	password := keyword.Get(URL_PASS)
	privilege, err := strconv.Atoi(keyword.Get(URL_PRIV))
	if err != nil {
		rsp.C = ERR_USER_PRIVILEGE
		return util.JsonOut(rsp)
	}
	if !user.AddUser(username, password, privilege) {
		rsp.C = ERR_USER_DUPLICATE
	}
	return util.JsonOut(rsp)
}

func handleUserDelete(keyword url.Values) []byte {
	rsp := Rsp{}
	username := keyword.Get(URL_USER)
	if !user.DeleteUser(username) {
		rsp.C = ERR_USER_NOT_EXIST
	}
	return util.JsonOut(rsp)
}

func handleUserPass(keyword url.Values) []byte {
	rsp := Rsp{}
	username := keyword.Get(URL_USER)
	password := keyword.Get(URL_PASS)
	if !user.MotifyUserInfo(username, password, -1) {
		rsp.C = ERR_USER_NOT_EXIST
	}
	return util.JsonOut(rsp)
}

func handleUserPriv(keyword url.Values) []byte {
	rsp := Rsp{}
	username := keyword.Get(URL_USER)
	privilege, err := strconv.Atoi(keyword.Get(URL_PRIV))
	if err != nil {
		rsp.C = ERR_PARSE_MISS
		return util.JsonOut(rsp)
	}
	if !user.MotifyUserInfo(username, "", privilege) {
		rsp.C = ERR_USER_PRIVILEGE
	}
	return util.JsonOut(rsp)
}

func initUser() {
	registerCMD(CMD_UADD, 3, handleUserAdd, TYPE_X)
	registerCMD(CMD_UDEL, 2, handleUserDelete, TYPE_X)
	registerCMD(CMD_UPAS, 3, handleUserPass, TYPE_X)
	registerCMD(CMD_UPRI, 3, handleUserPriv, TYPE_X)
}
