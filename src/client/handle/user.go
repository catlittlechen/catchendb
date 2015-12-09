package handle

import (
	"catchendb/src/util"
	"strconv"
)

func HandleUserAdd(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 4 {
		return nil
	}
	pri, err := strconv.Atoi(argv[3])
	if err != nil {
		return nil
	}
	req := Req{
		C:         CmdUAdd,
		UserName:  argv[1],
		PassWord:  argv[2],
		Privilege: pri,
	}
	return util.JSONOut(req)
}

func handleUserDel(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:        CmdUDel,
		UserName: argv[1],
	}
	return util.JSONOut(req)
}

func handleUserPas(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:        CmdUPas,
		UserName: argv[1],
	}
	return util.JSONOut(req)
}

func handleUserPri(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	pri, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil
	}
	req := Req{
		C:         CmdUPri,
		Privilege: pri,
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CmdUAdd, HandleUserAdd)
	registerHandle(CmdUDel, handleUserDel)
	registerHandle(CmdUPas, handleUserPas)
	registerHandle(CmdUPri, handleUserPri)
}
