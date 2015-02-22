package logic

import (
	"bytes"
	"catchendb/src/config"
	"catchendb/src/node"
	"catchendb/src/store"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

import lgd "code.google.com/p/log4go"

var (
	outPutSign = []byte("quit")
)

var (
	lengthData int
)

const (
	magicKey   = "1A089524CB555F689E5E8F72CFFC54C7"
	dataKeyLen = 10
)

type data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func LoadData() bool {

	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName
	lgd.Trace("start loaddata[%s] at the time[%d]", filename, time.Now().Unix())

	fp, err := os.Open(filename)
	if err != nil {
		lgd.Error("file[%s] open fail! error[%s]", filename, err)
		return false
	}
	defer fp.Close()

	l := make([]byte, len(magicKey))
	lens, err := fp.Read(l)
	if err != nil && err != io.EOF {
		lgd.Error("file[%s] read error[%s]", filename, err)
		return false
	} else if err == io.EOF {
		return true
	} else if string(l) != magicKey {
		lgd.Error("file[%s] magicKey[%s]", filename, l)
		return false
	}

	l = make([]byte, dataKeyLen)
	lens, err = fp.Read(l)
	if err != nil {
		lgd.Error("file[%s] read error[%s]", filename, err)
		return false
	} else if lens != dataKeyLen {
		lgd.Error("file[%s] read error[len is illegal]", filename)
		return false
	} else {
		lengthData, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Error("file[%s] length fail! data[%s]", filename, l)
			return false
		}
	}

	var line []byte
	var d data
	length := lengthData
	lengthBool := true
	for {
		l = make([]byte, length)
		lens, err = fp.Read(l)
		if err != nil {
			if err == io.EOF && lengthBool {
				break
			} else {
				lgd.Error("file[%s] read error[%s]", filename, err)
				return false
			}
		}
		if lengthBool {
			length, err = strconv.Atoi(string(l))
			if err != nil {
				lgd.Error("file[%s] length fail! data[%s]", filename, l)
				return false
			}
		} else {
			length = lengthData
			line, err = store.Decode(l)
			lgd.Debug(line)
			if err != nil {
				lgd.Error("data[%s] is illegal", l)
				return false
			}
			err = json.Unmarshal(line, &d)
			if err != nil {
				lgd.Error("data[%s] is illegal", line)
				return false
			}
			go node.Put(d.Key, d.Value)
		}
		lengthBool = !lengthBool
	}

	lgd.Trace("finish loaddata[%s] at the time[%d]", filename, time.Now().Unix())
	return true
}

func saveData() bool {

	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName + ".tmp"
	lgd.Trace("start saveData[%s] at the time[%d]", filename, time.Now().Unix())

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		lgd.Error("file[%s] open fail! err[%s]", filename, err)
		return false
	}

	channel := make(chan []byte, 1000)
	go node.OutPut(channel, outPutSign)

	var datastr []byte
	var datastr2 string
	datastr2 = magicKey + fmt.Sprintf("%0"+fmt.Sprintf("%d", dataKeyLen)+"d", lengthData)
	_, err = fp.Write([]byte(datastr2))
	if err != nil {
		lgd.Error("file[%s] write fail! err[%s]", filename, err)
		return false
	}
	printsign := "%0" + fmt.Sprintf("%d", lengthData) + "d"
	for {
		datastr = <-channel
		if bytes.Equal(datastr, outPutSign) {
			break
		}
		datastr = store.Encode(datastr)
		datastr2 = fmt.Sprintf(printsign, len(datastr))
		_, err = fp.Write([]byte(datastr2))
		lgd.Debug(datastr2)
		if err != nil {
			lgd.Error("file[%s] write fail! err[%s]", filename, err)
			return false
		}
		_, err = fp.Write(datastr)
		lgd.Debug(string(datastr))
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

	lgd.Trace("finish saveData[%s] at the time[%d]", filename, time.Now().Unix())
	return true
}

func AutoSaveData() (ret bool) {
	c := time.Tick(config.GlobalConf.Data.DataTime * time.Second)
	for _ = range c {
		saveData()
	}

	return
}

func init() {
	lengthData = 4
}
