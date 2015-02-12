package store

func Encode(obj []byte) []byte {
	return encode(obj)
}

func Decode(obj []byte) ([]byte, error) {
	return decode(obj)
}

func init() {
}
