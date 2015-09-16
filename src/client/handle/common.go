package handle

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
	CMDROLLBACK = "rollback"
	CMDCOMMIT   = "commit"
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
