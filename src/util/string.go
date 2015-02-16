package util

import (
	"encoding/json"
)

func JsonOut(obj interface{}) []byte {
	js, _ := json.Marshal(obj)
	return js
}
