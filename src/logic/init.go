package logic

import (
	"catchendb/src/util"
	"fmt"
	"net/url"
	"sync"
)

import lgd "code.google.com/p/log4go"

var (
	userPrivilege map[string]int
	getNameMutex  *sync.Mutex
	nameInt       int
)

func LYW(data []byte, name string) []byte {
	return lyw(data, name)
}

func AUT(data []byte) (bool, string, []byte) {
	return aut(data)
}

func DisConnection(name string) {
	delete(userPrivilege, name)
}

func lyw(data []byte, name string) []byte {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)

	keyword, err := url.ParseQuery(urlStr)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		return util.JsonOut(rsp)
	}
	privilege, ok := userPrivilege[name]
	if !ok {
		lgd.Error("System Error For Get Privilege")
		rsp.C = ERR_SYSTEM_EROR
		return util.JsonOut(rsp)
	}
	return mapAction(keyword, privilege)
}

func aut(data []byte) (ok bool, name string, r []byte) {
	lgd.Info("Request %s", string(data))

	rsp := Rsp{}
	urlStr := string(data)

	keyword, err := url.ParseQuery(urlStr)
	if err != nil {
		lgd.Error("ParseQuery fail with the url %s", urlStr)
		rsp.C = ERR_URL_PARSE
		r = util.JsonOut(rsp)
		return
	}
	privilege := 0
	ok, privilege = handleUserAut(keyword)
	if !ok {
		rsp.C = ERR_ACCESS_DENIED
		r = util.JsonOut(rsp)
		return
	}
	name = getName()
	if _, ok2 := userPrivilege[name]; !ok2 {
		userPrivilege[name] = privilege
	} else {
		ok = false
		rsp.C = ERR_SYSTEM_EROR
	}
	r = util.JsonOut(rsp)
	return
}

func getName() string {
	getNameMutex.Lock()
	defer getNameMutex.Unlock()
	nameInt += 1
	return fmt.Sprintf("%d", nameInt)
}

func Init() {
	initString()
	initUser()
	go autoSaveData()
}

func init() {
	userPrivilege = make(map[string]int)
	nameInt = 0
	getNameMutex = new(sync.Mutex)
}
