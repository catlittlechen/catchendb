package logic

import (
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

func LYW(data []byte) []byte {
	return lyw(data)
}

func AUT(data []byte) (bool, []byte) {
	return aut(data)
}

func lyw(data []byte) []byte {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)
	urlTmp, err := url.Parse(urlStr)
	if err != nil {
		lgd.Error("Parse RequestURL fail with the url %s! ", urlStr)
		rsp.C = ERR_URL_PARSE
		return util.JsonOut(rsp)
	}

	keyword, err := url.ParseQuery(urlTmp.RawQuery)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		return util.JsonOut(rsp)
	}
	return mapAction(keyword)
}

func aut(data []byte) (ok bool, r []byte) {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)
	urlTmp, err := url.Parse(urlStr)
	if err != nil {
		lgd.Error("Parse RequestURL fail with the url %s! ", urlStr)
		rsp.C = ERR_URL_PARSE
		r = util.JsonOut(rsp)
		return
	}

	keyword, err := url.ParseQuery(urlTmp.RawQuery)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		r = util.JsonOut(rsp)
		return
	}
	if !userAut(keyword) {
		rsp.C = ERR_ACCESS_DENIED
		r = util.JsonOut(rsp)
		return
	}
	ok = true
	return ok, util.JsonOut(rsp)
}

func Init() {
	initString()
	go autoSaveData()
}

func init() {

}
