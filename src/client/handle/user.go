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
		C:         CMDUADD,
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
		C:        CMDUDEL,
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
		C:        CMDUPAS,
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
		C:         CMDUPRI,
		Privilege: pri,
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CMDUADD, HandleUserAdd)
	registerHandle(CMDUDEL, handleUserDel)
	registerHandle(CMDUPAS, handleUserPas)
	registerHandle(CMDUPRI, handleUserPri)
}
