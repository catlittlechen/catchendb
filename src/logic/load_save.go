package logic

import (
	"bytes"
	"catchendb/src/config"
	"catchendb/src/node"
	"catchendb/src/store"
	"catchendb/src/user"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

import lgd "catchendb/src/log"

var (
	outPutSign   = []byte("quit")
	outPutMutex  *sync.Mutex
	outPutBoolen bool
)

var (
	lengthData     int
	lengthUser     int
	lastModifyTime int64
)

const (
	magicKey   = "1A089524CB555F689E5E8F72CFFC54C7"
	timeKeyLen = 20
	dataKeyLen = 10
	userKeyLen = 10

//幻数写于文件的最前面
//datakeylen为存储数据的基本长度
//userkeylen为用户数据的基本长度
)

func LoadData() bool {

	if !config.GlobalConf.MasterSlave.IsMaster {
		return true
	}

	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName
	lgd.Trace("start loaddata[%s] at the time[%d]", filename, time.Now().Unix())

	fp, err := os.Open(filename)
	if err != nil {
		lgd.Errorf("file[%s] open fail! error[%s]", filename, err)
		return false
	}
	defer fp.Close()

	var line []byte

	//magicKey
	l := make([]byte, len(magicKey))
	_, err = fp.Read(l)
	if err != nil {
		if err != io.EOF {
			lgd.Errorf("file[%s] read error[%s]", filename, err)
		}
		return err == io.EOF

	}
	if string(l) != magicKey {
		lgd.Errorf("file[%s] magicKey[%s]", filename, l)
		return false
	}

	//lastModifyTime
	l = make([]byte, timeKeyLen)
	lens, err := fp.Read(l)
	if err != nil {
		lgd.Errorf("faile[%s] read error[%s]", filename, err)
		return false
	} else if lens != timeKeyLen {
		lgd.Errorf("file[%s] read error[len is illegal]", filename)
		return false
	} else {
		lastModifyTime, err = strconv.ParseInt(string(l), 10, 64)
		if err != nil {
			lgd.Errorf("file[%s] length fail! data[%s]", filename, l)
			return false
		}
	}

	//userData
	l = make([]byte, userKeyLen)
	lens, err = fp.Read(l)
	if err != nil {
		lgd.Errorf("faile[%s] read error[%s]", filename, err)
		return false
	} else if lens != userKeyLen {
		lgd.Errorf("file[%s] read error[len is illegal]", filename)
		return false
	} else {
		lengthUser, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Errorf("file[%s] length fail! data[%s]", filename, l)
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
			lgd.Errorf("file[%s] read error[%s]", filename, err)
			return false
		}
		length, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Errorf("file[%s] length fail! data[%s]", filename, l)
			return false
		}
		l = make([]byte, length)
		lens, err = fp.Read(l)
		if err != nil {
			lgd.Errorf("file[%s] read error[%s]", filename, err)
			return false
		}
		line, err = store.Decode(l)
		if err != nil {
			lgd.Errorf("data[%s] is illegal", l)
			return false
		}
		if !user.InPut(line) {
			lgd.Errorf("data[%s] is illegal", line)
			return false
		}
		length = lengthUser
	}

	//data
	l = make([]byte, dataKeyLen)
	lens, err = fp.Read(l)
	if err != nil {
		lgd.Errorf("file[%s] read error[%s]", filename, err)
		return false
	} else if lens != dataKeyLen {
		lgd.Errorf("file[%s] read error[len is illegal]", filename)
		return false
	} else {
		lengthData, err = strconv.Atoi(string(l))
		if err != nil {
			lgd.Errorf("file[%s] length fail! data[%s]", filename, l)
			return false
		}
	}

	length = lengthData
	lengthBool := true
	for {
		l = make([]byte, length)
		_, err = fp.Read(l)
		if err != nil {
			if err == io.EOF && lengthBool {
				break
			} else {
				lgd.Errorf("file[%s] read error[%s]", filename, err)
				return false
			}
		}
		if lengthBool {
			length, err = strconv.Atoi(string(l))
			if err != nil {
				lgd.Errorf("file[%s] length fail! data[%s]", filename, l)
				return false
			}
		} else {
			length = lengthData
			line, err = store.Decode(l)
			if err != nil {
				lgd.Errorf("data[%s] is illegal", l)
				return false
			}
			if !node.InPut(line) {
				lgd.Errorf("data[%s] is illegal", line)
				return false
			}
		}
		lengthBool = !lengthBool
	}

	lgd.Trace("finish loaddata[%s] at the time[%d]", filename, time.Now().Unix())
	return true
}

func saveData() bool {
	outPutMutex.Lock()
	if outPutBoolen {
		outPutMutex.Unlock()
		return true
	}
	outPutBoolen = true
	outPutMutex.Unlock()

	defer func() {
		outPutBoolen = false
	}()

	lastModifyTime = time.Now().Unix()
	filename := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName + ".tmp"
	lgd.Trace("start saveData[%s] at the time[%d]", filename, time.Now().Unix())

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		lgd.Errorf("file[%s] open fail! err[%s]", filename, err)
		return false
	}

	//magicKey
	_, err = fp.Write([]byte(magicKey))
	if err != nil {
		lgd.Errorf("file[%s] write fail! err[%s]", filename, err)
		return false
	}

	printsign := "%0" + strconv.Itoa(timeKeyLen) + "d"
	l := fmt.Sprintf(printsign, lastModifyTime)
	_, err = fp.Write([]byte(l))
	if err != nil {
		lgd.Errorf("file[%s] write fail! err[%s]", filename, err)
		return false
	}

	//user
	count := 0
	var datastrSum, datastr, datastr2 []byte
	channel := make(chan []byte, 1000)
	go user.OutPut(channel, outPutSign)
	printsign = "%0" + strconv.Itoa(lengthUser) + "d"
	for {
		datastr = <-channel
		if bytes.Equal(datastr, outPutSign) {
			break
		}
		datastr = store.Encode(datastr)
		datastr2 = []byte(fmt.Sprintf(printsign, len(datastr)))
		datastrSum = append(datastrSum, datastr2...)
		datastrSum = append(datastrSum, datastr...)
		count++
	}

	printsign = "%0" + strconv.Itoa(userKeyLen/2) + "d"
	datastr = []byte(fmt.Sprintf(printsign+printsign, count, lengthUser))
	datastr = append(datastr, datastrSum...)
	_, err = fp.Write(datastr)
	if err != nil {
		lgd.Errorf("file[%s] write fail! err[%s]", filename, err)
		return false
	}

	go node.OutPut(channel, outPutSign)

	datastr = []byte(fmt.Sprintf("%0"+strconv.Itoa(dataKeyLen)+"d", lengthData))
	_, err = fp.Write(datastr)
	if err != nil {
		lgd.Errorf("file[%s] write fail! err[%s]", filename, err)
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
		_, err = fp.Write(datastr)
		if err != nil {
			lgd.Errorf("file[%s] write fail! err[%s]", filename, err)
			return false
		}
	}
	fp.Close()
	err = os.Rename(filename, config.GlobalConf.Data.DataPath+config.GlobalConf.Data.DataName)
	if err != nil {
		lgd.Errorf("file[%s] rename fail! err[%s]", filename, err)
		return false
	}

	lgd.Trace("finish saveData[%s] at the time[%d]", filename, time.Now().Unix())
	return true
}

func autoSaveData() (ret bool) {
	for {
		time.Sleep(config.GlobalConf.Data.DataTime * time.Second)
		saveData()
	}
}

func init() {
	lengthData = 4
	lengthUser = 4
	outPutMutex = new(sync.Mutex)
	outPutBoolen = false
}
