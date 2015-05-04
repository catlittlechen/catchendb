package main

import ()

func handleBegin(keyword url.Value) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	return util.JsonOut(rsp)
}

func handleCommit(keyword url.Value) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	return util.JsonOut(rsp)
}

func handleRollBack(keyword url.Value) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	return util.JsonOut(rsp)
}

func initTransaction() {
	registerCMD(CMD_BEGIN, 1, handleBegin, TYPE_W)
	registerCMD(CMD_COMMIT, 1, handleCommit, TYPE_W)
	registerCMD(CMD_ROLLBACK, 1, handleRollBack, TYPE_W)
}
