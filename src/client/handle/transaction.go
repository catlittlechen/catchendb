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
		C: CmdBegin,
	}
	return util.JSONOut(req)
}

func HandleRollback(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CmdRollback,
	}
	return util.JSONOut(req)
}

func HandleCommit(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	req := Req{
		C: CmdCommit,
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CmdBegin, HandleBegin)
	registerHandle(CmdRollback, HandleRollback)
	registerHandle(CmdCommit, HandleCommit)
}
