package store

import (
	"catchendb/src/config"
	"time"
)

import lgd "code.google.com/p/log4go"

func LoadData() bool {

	files := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName
	lgd.Trace("start loaddata[%s] at the time[%d]", files, time.Now().Unix)

	//TODO

	lgd.Trace("finish loaddata[%s] at the time[%d]", files, time.Now().Unix)
	return true
}

func saveData() bool {

	files := config.GlobalConf.Data.DataPath + config.GlobalConf.Data.DataName
	lgd.Trace("start saveData[%s] at the time[%d]", files, time.Now().Unix)

	//TODO

	lgd.Trace("finish saveData[%s] at the time[%d]", files, time.Now().Unix)
	return true
}

func AutoSaveData() (ret bool) {
	ret = saveData()
	if ret {
		go func() {
			c := time.Tick(config.GlobalConf.Data.DataTime * time.Second)
			for _ = range c {
				saveData()
			}
		}()
	}

	return
}
