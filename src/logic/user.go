package logic

import (
	"catchendb/src/util"
	"net/url"
)

import lgd "code.google.com/p/log4go"

var (
	mapUser map[string]user
)

const (
	preUserSig = "zjm(@)*@*XX!@)@"
)

type user struct {
	Username  string `json:"usmd"`
	Password  string `json:"psmd"`
	Privilege int    `json:"psmd"`
}

func (u *user) init(user, pass string, pri int) {
	u.Username = user
	u.Password = pass
	u.Privilege = pri
}

func (u *user) verifyPassword(password string) bool {
	password = preUserSig + util.Md5(password)
	if u.Password == util.Md5(password) {
		return true
	}
	return false
}

func (u *user) setPassword(password string) {
	password = preUserSig + util.Md5(password)
	u.Password = util.Md5(password)
}

func (u *user) getPrivilege() int {
	return u.Privilege
}

func (u *user) setPrivilege(pri int) {
	u.Privilege = pri
}

func verifyPassword(username, password string) bool {
	u, ok := mapUser[username]
	if !ok {
		lgd.Warn("No username[%s] exists")
		return false
	}
	return u.verifyPassword(password)
}

func addUser(name, pass string, pri int) (ret int) {
	u := new(user)
	if _, ok := mapUser[name]; ok {
		ret = ERR_USER_DUPLICATE
		return
	}
	u.init(name, pass, pri)
	mapUser[name] = *u
	return
}

func deleteUser(name string) (ret int) {
	if _, ok := mapUser[name]; !ok {
		ret = ERR_USER_NOT_EXIST
		return
	} else {
		delete(mapUser, name)
	}
	return
}

func motifyUserInfo(name, pass string, pri int) (ret int) {
	if u, ok := mapUser[name]; !ok {
		ret = ERR_USER_NOT_EXIST
		return
	} else {
		if pass != "" {
			u.setPassword(pass)
		}
		if pri != -1 {
			u.setPrivilege(pri)
		}
	}
	return
}

func userAut(keyword url.Values) bool {
	return false
}

func init() {
	mapUser = make(map[string]user)
	u := new(user)
	u.init("root", "root", 7)
	mapUser["root"] = *u
}
