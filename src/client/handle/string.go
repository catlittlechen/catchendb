package handle

import (
	"catchendb/src/util"
	"net/url"
)

func HandleSet(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_SET)
	urlData.Add(URL_KEY, argv[1])
	urlData.Add(URL_VALUE, argv[2])
	return []byte(urlData.Encode())
}

func HandleGet(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_GET)
	urlData.Add(URL_KEY, argv[1])
	return []byte(urlData.Encode())
}

func HandleDel(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_DEL)
	urlData.Add(URL_KEY, argv[1])
	return []byte(urlData.Encode())
}

func HandleSetEx(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 5 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_SETEX)
	urlData.Add(URL_KEY, argv[1])
	urlData.Add(URL_VALUE, argv[2])
	urlData.Add(URL_START, argv[3])
	urlData.Add(URL_END, argv[4])
	return []byte(urlData.Encode())
}

func HandleDelay(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_DELAY)
	urlData.Add(URL_KEY, argv[1])
	urlData.Add(URL_START, argv[2])
	return []byte(urlData.Encode())
}

func HandleExpire(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 3 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_EXPIRE)
	urlData.Add(URL_KEY, argv[1])
	urlData.Add(URL_END, argv[2])
	return []byte(urlData.Encode())
}

func HandleTTL(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_TTL)
	urlData.Add(URL_KEY, argv[1])
	return []byte(urlData.Encode())
}

func HandleTTS(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_TTS)
	urlData.Add(URL_KEY, argv[1])
	return []byte(urlData.Encode())
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
