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
		C:     CMD_SET,
		Key:   argv[1],
		Value: argv[2],
	}
	return util.JsonOut(req)
}

func HandleGet(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CMD_GET,
		Key: argv[1],
	}
	return util.JsonOut(req)
}

func HandleDel(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CMD_DEL,
		Key: argv[1],
	}
	return util.JsonOut(req)
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
		C:         CMD_SETEX,
		Key:       argv[1],
		Value:     argv[2],
		StartTime: start,
		EndTime:   end,
	}
	return util.JsonOut(req)
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
		C:         CMD_DELAY,
		Key:       argv[1],
		StartTime: start,
	}
	return util.JsonOut(req)
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
		C:       CMD_EXPIRE,
		Key:     argv[1],
		EndTime: end,
	}
	return util.JsonOut(req)
}

func HandleTTL(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CMD_TTL,
		Key: argv[1],
	}
	return util.JsonOut(req)
}

func HandleTTS(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	req := Req{
		C:   CMD_TTS,
		Key: argv[1],
	}
	return util.JsonOut(req)
}

func init() {
	registerHandle(CMD_SET, HandleSet)
	registerHandle(CMD_GET, HandleGet)
	registerHandle(CMD_DEL, HandleDel)
	registerHandle(CMD_SETEX, HandleSetEx)
	registerHandle(CMD_DELAY, HandleDelay)
	registerHandle(CMD_EXPIRE, HandleExpire)
	registerHandle(CMD_TTL, HandleTTL)
	registerHandle(CMD_TTS, HandleTTS)
}
