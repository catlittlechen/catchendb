package handle

import (
	"catchendb/src/util"
	"strconv"
)

func HandleSet(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	req := Req{
		C:     CmdSet,
		Key:   argv[1],
		Value: argv[2],
	}
	return util.JSONOut(req)
}

func HandleGet(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CmdGet,
		Key: argv[1],
	}
	return util.JSONOut(req)
}

func HandleDel(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CmdDel,
		Key: argv[1],
	}
	return util.JSONOut(req)
}

func HandleSetEx(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 5 {
		return nil
	}
	start, err := strconv.ParseInt(argv[3], 10, 64)
	if err != nil {
		return nil
	}
	end, err := strconv.ParseInt(argv[4], 10, 64)
	if err != nil {
		return nil
	}
	req := Req{
		C:         CmdSetEX,
		Key:       argv[1],
		Value:     argv[2],
		StartTime: start,
		EndTime:   end,
	}
	return util.JSONOut(req)
}

func HandleDelay(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	start, err := strconv.ParseInt(argv[2], 10, 64)
	if err != nil {
		return nil
	}

	req := Req{
		C:         CmdDelAY,
		Key:       argv[1],
		StartTime: start,
	}
	return util.JSONOut(req)
}

func HandleExpire(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	end, err := strconv.ParseInt(argv[2], 10, 64)
	if err != nil {
		return nil
	}

	req := Req{
		C:       CmdExpire,
		Key:     argv[1],
		EndTime: end,
	}
	return util.JSONOut(req)
}

func HandleTTL(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CmdTTL,
		Key: argv[1],
	}
	return util.JSONOut(req)
}

func HandleTTS(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CmdTTS,
		Key: argv[1],
	}
	return util.JSONOut(req)
}

func init() {
	registerHandle(CmdSet, HandleSet)
	registerHandle(CmdGet, HandleGet)
	registerHandle(CmdDel, HandleDel)
	registerHandle(CmdSetEX, HandleSetEx)
	registerHandle(CmdDelAY, HandleDelay)
	registerHandle(CmdExpire, HandleExpire)
	registerHandle(CmdTTL, HandleTTL)
	registerHandle(CmdTTS, HandleTTS)
}
