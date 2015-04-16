package logic

import (
	"bytes"
	"catchendb/src/config"
	"catchendb/src/node"
	"catchendb/src/store"
	"catchendb/src/user"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	lengthUser int
)

const (
	magicKey   = "1A089524CB555F689E5E8F72CFFC54C7"
	dataKeyLen = 10
	userKeyLen = 10

//幻数写于文件的最前面
//datakeylen为存储数据的基本长度
//userkeylen为用户数据的基本长度
)

type userInfo struct {
	Username  string `json:"usmd"`
	Password  string `json:"psmd"`
	Privilege int    `json:"psmd"`
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
	var line []byte

	//magicKey
	l := make([]byte, len(magicKey))
	lens, err := fp.Read(l)
	if err != nil {
		if err != io.EOF {
			lgd.Error("file[%s] read error[%s]", filename, err)
			return false
		} else {
			return true
		}
	} else if string(l) != magicKey {
		lgd.Error("file[%s] magicKey[%s]", filename, l)
		return false
	}

	//userData
	l = make([]byte, userKeyLen)
	lens, err = fp.Read(l)
	if err != nil {
		lgd.Error("faile[%s] read error[%s]", filename, err)
		return false
	} else if lens != userKeyLen {
		lgd.Error("file[%s] read error[len is illegal]", filename)
		return false
	} else {
		lengthUser, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Error("file[%s] length fail! data[%s]", filename, l)
			return false
		}
	}
	//count 代表这用户的数量
	//lengthUser 代表这用户的信息长度的存储长度
	count := lengthUser / int(math.Pow10(userKeyLen/2))
	lengthUser = lengthUser % int(math.Pow10(userKeyLen/2))
	length := lengthUser

	for i := 0; i < count; i++ {
		l = make([]byte, length)
		lens, err = fp.Read(l)
		if err != nil {
			lgd.Error("file[%s] read error[%s]", filename, err)
			return false
		}
		length, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Error("file[%s] length fail! data[%s]", filename, l)
			return false
		}
		l = make([]byte, length)
		lens, err = fp.Read(l)
		if err != nil {
			lgd.Error("file[%s] read error[%s]", filename, err)
			return false
		}
		line, err = store.Decode(l)
		if err != nil {
			lgd.Error("data[%s] is illegal", l)
			return false
		}
		u := userInfo{}
		err = json.Unmarshal(line, &u)
		if err != nil {
			lgd.Error("data[%s] is illegal", line)
			return false
		}
		go user.AddUser(u.Username, u.Password, u.Privilege)
		length = lengthUser
	}

	//data
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

	length = lengthData
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
			if err != nil {
				lgd.Error("data[%s] is illegal", l)
				return false
			}
			if node.InPut(line) {
				lgd.Error("data[%s] is illegal", line)
				return false
			}
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

	//magicKey
	_, err = fp.Write([]byte(magicKey))
	if err != nil {
		lgd.Error("file[%s] write fail! err[%s]", filename, err)
		return false
	}
	count := 0
	var datastrSum, datastr, datastr2 []byte
	channel := make(chan []byte, 1000)
	go user.OutPut(channel, outPutSign)
	printsign := "%0" + fmt.Sprintf("%d", lengthUser) + "d"
	for {
		datastr = <-channel
		if bytes.Equal(datastr, outPutSign) {
			break
		}
		datastr = store.Encode(datastr)
		datastr2 = []byte(fmt.Sprintf(printsign, len(datastr)))
		datastrSum = append(datastrSum, datastr2...)
		datastrSum = append(datastrSum, datastr...)
		count += 1
	}

	printsign = "%0" + fmt.Sprintf("%d", userKeyLen/2) + "d"
	datastr = []byte(fmt.Sprintf(printsign+printsign, count, lengthUser))
	datastr = append(datastr, datastrSum...)
	_, err = fp.Write(datastr)
	if err != nil {
		lgd.Error("file[%s] write fail! err[%s]", filename, err)
		return false
	}

	go node.OutPut(channel, outPutSign)
	datastr = []byte(fmt.Sprintf("%0"+fmt.Sprintf("%d", dataKeyLen)+"d", lengthData))
	_, err = fp.Write(datastr)
	if err != nil {
		lgd.Error("file[%s] write fail! err[%s]", filename, err)
		return false
	}
	printsign = "%0" + fmt.Sprintf("%d", lengthData) + "d"
	for {
		datastr = <-channel
		if bytes.Equal(datastr, outPutSign) {
			break
		}
		datastr = store.Encode(datastr)
		datastr = append([]byte(fmt.Sprintf(printsign, len(datastr))), datastr...)
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

func autoSaveData() (ret bool) {
	c := time.Tick(config.GlobalConf.Data.DataTime * time.Second)
	for _ = range c {
		saveData()
	}

	return
}

func init() {
	lengthData = 4
	lengthUser = 4
}
