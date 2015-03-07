package user

import ()

const (
	preUserSig = "zjm(@)*@*XX!@)@"
)

func VerifyPassword(username, password string) bool {
	return verifyPassword(username, password)
}

func AddUser(name, pass string, pri int) bool {
	return addUser(name, pass, pri)
}

func DeleteUser(name string) bool {
	return deleteUser(name)
}

func MotifyUserInfo(name, pass string, pri int) bool {
	return motifyUserInfo(name, pass, pri)
}

func InPut(line []byte) bool {
	return input(line)
}

func OutPut(channel chan []byte, sign []byte) {
	output(channel, sign)
}

func init() {

}
