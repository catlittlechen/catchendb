package logic

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

const (
	cmdAut      = "aut"
	cmdUadd     = "useradd"
	cmdUdel     = "userdelete"
	cmdUpas     = "userpass"
	cmdUpri     = "userprivilege"
	cmdSet      = "set"
	cmdGet      = "get"
	cmdDel      = "del"
	cmdSetEX    = "setex"
	cmdDelAY    = "delay"
	cmdExpire   = "expire"
	cmdTTL      = "ttl"
	cmdTTS      = "tts"
	cmdBegin    = "begin"
	cmdCommit   = "commit"
	cmdRollback = "rollback"
)

const (
	errURLParse     = 100
	errCmdMiss      = 101
	errParseMiss    = 102
	errAccessDenied = 103
	errSystemBusy   = 104
	errSystemEror   = 105

	errUserPass      = 200
	errUserDuplicate = 201
	errUserNotExist  = 202
	errUserPrivilege = 203
	errUserMaxOnly   = 204
	errUserMaxUser   = 205

	errCmdSet    = 301
	errCmdDel    = 302
	errStartTime = 303
	errEndTime   = 304
	errCmdDelAY  = 305
	errCmdExpire = 306

	errTraBegin   = 401
	errTraNoBegin = 402
	errTraUser    = 403
)

const (
	urlCmd   = "cmd"
	urlKey   = "key"
	urlValue = "value"
	urlStart = "start"
	urlEnd   = "end"
	urlUser  = "user"
	urlPass  = "pass"
	urlPriv  = "priv"
)

type Req struct {
	C string `json:"c"`

	UserName  string `json:"usr"`
	PassWord  string `json:"pas"`
	Privilege int    `json:"pri"`

	Key       string `json:"key"`
	Value     string `json:"val"`
	StartTime int64  `json:"sta"`
	EndTime   int64  `json:"end"`
}
