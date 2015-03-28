package handle

var handle2data map[string]func([]byte) []byte

func GetHandle(code string) (fun func([]byte) []byte, ok bool) {
	fun, ok = handle2data[code]
	return
}

func registerHandle(code string, fun func([]byte) []byte) {
	handle2data[code] = fun
}

func init() {
	handle2data = make(map[string]func([]byte) []byte)
}
