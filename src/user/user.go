package user

import (
	"catchendb/src/util"
	"encoding/json"
)

import lgd "code.google.com/p/log4go"

type userInfo struct {
	Username  string `json:"usmd"`
	Password  string `json:"pamd"`
	Privilege int    `json:"pri"`
}

func (u *userInfo) init(user, pass string, pri int) {
	u.Username = user
	u.setPassword(pass)
	u.setPrivilege(pri)
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

func (u *userInfo) encode() (line []byte, ok bool) {
	var err error
	line, err = json.Marshal(*u)
	if err != nil {
		return
	}
	ok = true
	return
}

func (u *userInfo) decode(line []byte) bool {
	err := json.Unmarshal(line, u)
	if err != nil {
		return false
	}
	return true
}

var (
	mapUser map[string]userInfo
)

func verifyPassword(username, password string) bool {
	u, ok := mapUser[username]
	if !ok {
		lgd.Warn("No username[%s] exists", username)
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

func deleteUser(name string) (ret bool) {
	if _, ok := mapUser[name]; !ok {
		ret = false
	} else {
		delete(mapUser, name)
		ret = true
	}
	return
}

func motifyUserInfo(name, pass string, pri int) (ret bool) {
	ret = true
	if u, ok := mapUser[name]; !ok {
		ret = false
	} else {
		if pass != "" {
			u.setPassword(pass)
		} else if pri != -1 {
			u.setPrivilege(pri)
		} else {
			ret = false
		}
	}
	return
}

func getPrivilege(name string) (ret int) {
	if u, ok := mapUser[name]; !ok {
		ret = 0
	} else {
		ret = u.getPrivilege()
	}
	return ret
}

func input(line []byte) bool {
	u := userInfo{}
	if !u.decode(line) {
		return false
	}

	mapUser[u.Username] = u
	return true
}

func output(channel chan []byte, outPutSign []byte) {
	for k := range mapUser {
		u := mapUser[k]
		line, _ := u.encode()
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
