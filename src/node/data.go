package node

import (
	"encoding/json"
)

type data struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	StartTime int64  `json:"start"`
	EndTime   int64  `json:"end"`
}

func (d *data) decode(line []byte) bool {
	err := json.Unmarshal(line, d)
	if err != nil {
		return false
	}
	return true
}

func (d *data) encode() (line []byte, ok bool) {
	var err error
	line, err = json.Marshal(d)
	if err != nil {
		return
	}
	ok = true
	return
}
