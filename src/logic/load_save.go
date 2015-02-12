package logic

import (
	"bufio"
	"bytes"
	"catchendb/src/config"
	"catchendb/src/node"
	"catchendb/src/store"
	"encoding/json"
	"os"
	"time"
)

import lgd "code.google.com/p/log4go"

var (
	outPutSign = []byte("quit")
)

type data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func LoadData() bool {

	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName
	lgd.Trace("start loaddata[%s] at the time[%d]", filename, time.Now().Unix)

	fp, err := os.Open(filename)
	if err != nil {
		lgd.Error("file[%s] open fail! error[%s]", filename, err)
		return false
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	var line []byte
	var d data
	appends := false
	for {
		l, isPretex, err := reader.ReadLine()
		if err != nil {
			lgd.Error("file[%s] read error[%s]", filename, err)
			return false
		}
		if appends {
			l = append(line, l...)
			appends = false
		}
		if isPretex {
			line = l
		} else {
			line, err = store.Decode(l)
			if err != nil {
				lgd.Error("data[line] is illegal")
				return false
			}
			err = json.Unmarshal(line, &d)
			if err != nil {
				lgd.Error("data[line] is illegal")
				return false
			}
			go node.Put(d.Key, d.Value)
			line = []byte("")
		}
	}

	lgd.Trace("finish loaddata[%s] at the time[%d]", filename, time.Now().Unix)
	return true
}

func saveData() bool {

	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName + ".tmp"
	lgd.Trace("start saveData[%s] at the time[%d]", filename, time.Now().Unix)

	fp, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		lgd.Error("file[%s] open fail! err[%s]", filename, err)
		return false
	}

	channel := make(chan []byte, 1000)
	go node.OutPut(channel, outPutSign)

	var datastr []byte
	var line []byte
	for {
		datastr = <-channel
		if bytes.Equal(datastr, outPutSign) {
			break
		}
		line = append(store.Encode(datastr), '\n')
		_, err = fp.Write(line)
		if err != nil {
			lgd.Error("file[%s] write fail! err[%s]", filename, err)
			return false
		}
	}
	fp.Close()
	err = os.Rename(filename, config.GlobalConf.Data.DataPath+config.GlobalConf.Data.DataName)
	if err != nil {
		lgd.Error("file[%s] rename fail! err[%s]", filename, err)
		return false
	}

	lgd.Trace("finish saveData[%s] at the time[%d]", filename, time.Now().Unix)
	return true
}

func AutoSaveData() (ret bool) {
	c := time.Tick(config.GlobalConf.Data.DataTime * time.Second)
	for _ = range c {
		saveData()
	}

	return
}
