package logic

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

const (
	CMD_AUT    = "aut"
	CMD_UADD   = "useradd"
	CMD_UDEL   = "userdelete"
	CMD_UPAS   = "userpass"
	CMD_UPRI   = "userprivilege"
	CMD_SET    = "set"
	CMD_GET    = "get"
	CMD_DEL    = "del"
	CMD_SETEX  = "setex"
	CMD_DELAY  = "delay"
	CMD_EXPIRE = "expire"
	CMD_TTL    = "ttl"
	CMD_TTS    = "tts"
)

const (
	ERR_URL_PARSE     = 100
	ERR_CMD_MISS      = 101
	ERR_PARSE_MISS    = 102
	ERR_ACCESS_DENIED = 103
	ERR_SYSTEM_BUSY   = 104
	ERR_SYSTEM_EROR   = 105

	ERR_USER_PASS      = 200
	ERR_USER_DUPLICATE = 201
	ERR_USER_NOT_EXIST = 202
	ERR_USER_PRIVILEGE = 203
	ERR_USER_MAX_ONLY  = 204
	ERR_USER_MAX_USER  = 205

	ERR_CMD_SET    = 301
	ERR_CMD_DEL    = 302
	ERR_START_TIME = 303
	ERR_END_TIME   = 304
	ERR_CMD_DELAY  = 305
	ERR_CMD_EXPIRE = 306
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
