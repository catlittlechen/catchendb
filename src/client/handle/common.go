package handle

const (
	CMD_AUT      = "aut"
	CMD_UADD     = "useradd"
	CMD_UDEL     = "userdelete"
	CMD_UPAS     = "userpass"
	CMD_UPRI     = "userprivilege"
	CMD_SET      = "set"
	CMD_GET      = "get"
	CMD_DEL      = "del"
	CMD_SETEX    = "setex"
	CMD_DELAY    = "delay"
	CMD_EXPIRE   = "expire"
	CMD_TTL      = "ttl"
	CMD_TTS      = "tts"
	CMD_BEGIN    = "begin"
	CMD_ROLLBACK = "rollback"
	CMD_COMMIT   = "commit"
)

const (
	URL_CMD   = "cmd"
	URL_KEY   = "key"
	URL_VALUE = "value"
	URL_START = "start"
	URL_END   = "end"
	URL_USER  = "user"
	URL_PASS  = "pass"
	URL_PRIV  = "priv"
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
