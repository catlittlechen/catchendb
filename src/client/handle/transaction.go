package handle

import (
	"catchendb/src/util"
)

func HandleBegin(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}

	req := Req{
		C: CMD_BEGIN,
	}
	return util.JSONOut(req)
}

func HandleRollback(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CMD_ROLLBACK,
	}
	return util.JSONOut(req)
}

func HandleCommit(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CMD_COMMIT,
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CMD_BEGIN, HandleBegin)
	registerHandle(CMD_ROLLBACK, HandleRollback)
	registerHandle(CMD_COMMIT, HandleCommit)
}
