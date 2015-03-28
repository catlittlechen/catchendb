package handle

import (
	"catchendb/src/util"
	"net/url"
)

func HandleUserAdd(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 4 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_UADD)
	urlData.Add(URL_USER, argv[1])
	urlData.Add(URL_PASS, argv[2])
	urlData.Add(URL_PRIV, argv[3])
	return []byte(urlData.Encode())
}

func handleUserDel(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_UDEL)
	urlData.Add(URL_USER, argv[1])
	return []byte(urlData.Encode())
}

func handleUserPas(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_UPAS)
	urlData.Add(URL_PASS, argv[1])
	return []byte(urlData.Encode())
}

func handleUserPri(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 2 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_UPRI)
	urlData.Add(URL_PRIV, argv[1])
	return []byte(urlData.Encode())
}

func init() {
	registerHandle(CMD_UADD, HandleUserAdd)
	registerHandle(CMD_UDEL, handleUserDel)
	registerHandle(CMD_UPAS, handleUserPas)
	registerHandle(CMD_UPRI, handleUserPri)
}
