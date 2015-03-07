package user

import (
	"catchendb/src/util"
	"encoding/json"
)

import lgd "code.google.com/p/log4go"

type userInfo struct {
	Username  string `json:"usmd"`
	Password  string `json:"psmd"`
	Privilege int    `json:"psmd"`
}

func (u *userInfo) init(user, pass string, pri int) {
	u.Username = user
	u.Password = pass
	u.Privilege = pri
}

func (u *userInfo) verifyPassword(password string) bool {
	password = preUserSig + util.Md5(password)
	if u.Password == util.Md5(password) {
		return true
	}
	return false
}

func (u *userInfo) setPassword(password string) {
	password = preUserSig + util.Md5(password)
	u.Password = util.Md5(password)
}

func (u *userInfo) getPrivilege() int {
	return u.Privilege
}

func (u *userInfo) setPrivilege(pri int) {
	u.Privilege = pri
}

var (
	mapUser map[string]userInfo
)

func verifyPassword(username, password string) bool {
	u, ok := mapUser[username]
	if !ok {
		lgd.Warn("No username[%s] exists")
		return false
	}
	return u.verifyPassword(password)
}

func addUser(name, pass string, pri int) bool {
	u := new(userInfo)
	if _, ok := mapUser[name]; ok {
		return false
	}
	u.init(name, pass, pri)
	mapUser[name] = *u
	return true
}

func deleteUser(name string) bool {
	if _, ok := mapUser[name]; !ok {
		return false
	} else {
		delete(mapUser, name)
	}
	return true
}

func motifyUserInfo(name, pass string, pri int) bool {
	if u, ok := mapUser[name]; !ok {
		return false
	} else {
		if pass != "" {
			u.setPassword(pass)
		} else if pri != -1 {
			u.setPrivilege(pri)
		} else {
			return false
		}
	}
	return true
}

func input(line []byte) bool {
	u := userInfo{}
	err := json.Unmarshal(line, &u)
	if err != nil {
		return false
	}
	mapUser[u.Username] = u
	return true
}

func output(channel chan []byte, outPutSign []byte) {
	for k, _ := range mapUser {
		u := mapUser[k]
		line, _ := json.Marshal(u)
		channel <- line
	}
	channel <- outPutSign
}

func init() {
	mapUser = make(map[string]userInfo)
	u := new(userInfo)
	u.init("root", "root", 7)
	mapUser["root"] = *u
}
