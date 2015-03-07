package logic

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

const (
	TYPE_STRING = "string"
)

const (
	CMD_AUT    = "aut"
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
	ERR_URL_PARSE     = 0x100
	ERR_CMD_MISS      = 0x101
	ERR_PARSE_MISS    = 0x102
	ERR_ACCESS_DENIED = 0x103

	ERR_USER_PASS      = 0x200
	ERR_USER_DUPLICATE = 0x201
	ERR_USER_NOT_EXIST = 0x202
	ERR_USER_PRIVILEGE = 0x203

	ERR_CMD_SET    = 0x301
	ERR_CMD_DEL    = 0x302
	ERR_START_TIME = 0x303
	ERR_END_TIME   = 0x304
	ERR_CMD_DELAY  = 0x305
	ERR_CMD_EXPIRE = 0x306
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
