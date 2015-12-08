package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
)

func JSONOut(obj interface{}) []byte {
	js, _ := json.Marshal(obj)
	return js
}

func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)

	return hex.EncodeToString(h.Sum(nil))
}

func SplitString(data string) (data2 []string) {
	argvs := strings.Split(data, string(' '))
	for _, argv := range argvs {
		if len(argv) != 0 {
			data2 = append(data2, argv)
		}
	}
	return
}
