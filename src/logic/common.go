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
	ERR_URL_PARSE  = 0x100
	ERR_CMD_MISS   = 0x101
	ERR_PARSE_MISS = 0x102
	ERR_CMD_SET    = 0x103
	ERR_CMD_DEL    = 0x104
	ERR_START_TIME = 0x105
	ERR_END_TIME   = 0x106
	ERR_CMD_DELAY  = 0x107
	ERR_CMD_EXPIRE = 0x108
)

const (
	URL_CMD   = "cmd"
	URL_KEY   = "key"
	URL_VALUE = "value"
	URL_START = "start"
	URL_END   = "end"
)
