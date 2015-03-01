package logic

import (
	"catchendb/src/util"
	"net/http"
	"net/url"
)

import lgd "code.google.com/p/log4go"

func LYW(r *http.Request) []byte {
	return lyw(r)
}

func lyw(r *http.Request) []byte {
	lgd.Info("Request %+v", r)

	rsp := Rsp{}
	urlStr := r.URL.String()
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

func Init() {
	initString()
	go autoSaveData()
}

func init() {

}
