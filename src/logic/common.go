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
	CMD_SET = "set"
	CMD_GET = "get"
	CMD_DEL = "del"
)

const (
	ERR_URL_PARSE  = 0x100
	ERR_CMD_MISS   = 0x101
	ERR_PARSE_MISS = 0x102
	ERR_CMD_SET    = 0x103
	ERR_CMD_DEL    = 0x104
)

const (
	URL_CMD   = "cmd"
	URL_KEY   = "key"
	URL_VALUE = "value"
)
