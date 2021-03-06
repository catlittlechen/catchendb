package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type xmlServer struct {
	BindAddr        string `xml:"bindaddr"`
	ReplicationAddr string `xml:"replicationaddr"`
}

type xmlData struct {
	DataPath string        `xml:"datapath"`
	DataName string        `xml:"dataname"`
	DataTime time.Duration `xml:"datatimes"`
}

type xmlMasterSlave struct {
	IsMaster bool `xml:"ismaster"`

	HashSize    int    `xml:"hashsize"`
	ReFlushTime int    `xml:"reflushs"`
	IP          string `xml:"masterip"`
	Port        int    `xml:"masterport"`
	UserName    string `xml:"masteruser"`
	PassWord    string `xml:"masterpass"`
}

type xmlConfig struct {
	XMLName               xml.Name       `xml:"config"`
	Server                xmlServer      `xml:"server"`
	Log                   string         `xml:"log"`
	Data                  xmlData        `xml:"data"`
	MasterSlave           xmlMasterSlave `xml:"masterslave"`
	PageSize              int            `xml:"pagesize"`
	MaxOnlyUserConnection int            `xml:"maxonlyuserconnection"`
	MaxUserConnection     int            `xml:"maxuserconnection"`
	MaxTransactionTime    time.Duration  `xml:"maxtransactiontime"`
}

var GlobalConf xmlConfig

func (conf *xmlConfig) LoadConfig(filename string) bool {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfig[EROR] Can‘t open file[%s] for reading, error[%s]", filename, err)
		return false
	}

	content, err := ioutil.ReadAll(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfig[EROR] Can‘t read file[%s], error[%s]", filename, err)
		return false
	}

	if err = xml.Unmarshal(content, conf); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfig[EROR] Can‘t parse file[%s] into XML, error[%s]", filename, err)
		return false
	}

	fmt.Fprintf(os.Stdout, "GlobalConf:\n %+v\n", *conf)
	err = fp.Close()
	if err != nil {
		return false
	}
	return true
}

func LoadConf(filename string) bool {
	if !GlobalConf.LoadConfig(filename) {
		return false
	}

	return true
}
