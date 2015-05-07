package handle

import (
	"catchendb/src/util"
	"net/url"
)

func HandleBegin(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_BEGIN)
	return []byte(urlData.Encode())
}

func HandleRollback(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_ROLLBACK)
	return []byte(urlData.Encode())
}

func HandleCommit(data []byte) []byte {
	argv := util.SplitString(string(data))
	if len(argv) != 1 {
		return nil
	}
	urlData := url.Values{}
	urlData.Add(URL_CMD, CMD_COMMIT)
	return []byte(urlData.Encode())
}

func init() {
	registerHandle(CMD_BEGIN, HandleBegin)
	registerHandle(CMD_ROLLBACK, HandleRollback)
	registerHandle(CMD_COMMIT, HandleCommit)
}
