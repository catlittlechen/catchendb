package logic

import (
	"catchendb/src/node"
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

func handleSet(keyword url.Values) []byte {
	key := TYPE_STRING + "&" + keyword.Get(URL_KEY)
	value := keyword.Get(URL_VALUE)
	rsp := Rsp{
		C: 0,
	}
	if !node.Put(key, value) {
		lgd.Error("set fail! key[%s] value[%s]", key, value)
		rsp.C = ERR_CMD_SET
	}
	return util.JsonOut(rsp)
}

func handleGet(keyword url.Values) []byte {
	key := TYPE_STRING + "&" + keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	rsp.D = node.Get(key)
	return util.JsonOut(rsp)
}

func handleDel(keyword url.Values) []byte {
	key := TYPE_STRING + "&" + keyword.Get(URL_KEY)
	rsp := Rsp{
		C: 0,
	}
	if !node.Del(key) {
		lgd.Error("del fail! key[%s]", key)
		rsp.C = ERR_CMD_DEL
	}
	return util.JsonOut(rsp)
}
