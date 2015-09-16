package logic

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

const (
	CMDAUT      = "aut"
	CMDUADD     = "useradd"
	CMDUDEL     = "userdelete"
	CMDUPAS     = "userpass"
	CMDUPRI     = "userprivilege"
	CMDSET      = "set"
	CMDGET      = "get"
	CMDDEL      = "del"
	CMDSETEX    = "setex"
	CMDDELAY    = "delay"
	CMDEXPIRE   = "expire"
	CMDTTL      = "ttl"
	CMDTTS      = "tts"
	CMDBEGIN    = "begin"
	CMDCOMMIT   = "commit"
	CMDROLLBACK = "rollback"
)

const (
	ERRURLPARSE     = 100
	ERRCMDMISS      = 101
	ERRPARSEMISS    = 102
	ERRACCESSDENIED = 103
	ERRSYSTEMBUSY   = 104
	ERRSYSTEMEROR   = 105

	ERRUSERPASS      = 200
	ERRUSERDUPLICATE = 201
	ERRUSERNOTEXIST  = 202
	ERRUSERPRIVILEGE = 203
	ERRUSERMAXONLY   = 204
	ERRUSERMAXUSER   = 205

	ERRCMDSET    = 301
	ERRCMDDEL    = 302
	ERRSTARTTIME = 303
	ERRENDTIME   = 304
	ERRCMDDELAY  = 305
	ERRCMDEXPIRE = 306

	ERRTRABEGIN   = 401
	ERRTRANOBEGIN = 402
	ERRTRAUSER    = 403
)

const (
	URLCMD   = "cmd"
	URLKEY   = "key"
	URLVALUE = "value"
	URLSTART = "start"
	URLEND   = "end"
	URLUSER  = "user"
	URLPASS  = "pass"
	URLPRIV  = "priv"
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
