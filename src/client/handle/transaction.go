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
		C: CMDBEGIN,
	}
	return util.JSONOut(req)
}

func HandleRollback(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CMDROLLBACK,
	}
	return util.JSONOut(req)
}

func HandleCommit(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CMDCOMMIT,
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CMDBEGIN, HandleBegin)
	registerHandle(CMDROLLBACK, HandleRollback)
	registerHandle(CMDCOMMIT, HandleCommit)
}
