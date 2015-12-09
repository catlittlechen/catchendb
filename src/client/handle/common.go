package handle

const (
	CmdAut      = "aut"
	CmdUAdd     = "useradd"
	CmdUDel     = "userdelete"
	CmdUPas     = "userpass"
	CmdUPri     = "userprivilege"
	CmdSet      = "set"
	CmdGet      = "get"
	CmdDel      = "del"
	CmdSetEX    = "setex"
	CmdDelAY    = "delay"
	CmdExpire   = "expire"
	CmdTTL      = "ttl"
	CmdTTS      = "tts"
	CmdBegin    = "begin"
	CmdRollback = "rollback"
	CmdCommit   = "commit"
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
