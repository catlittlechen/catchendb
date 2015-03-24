package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type xmlServer struct {
	BindAddr string `xml:"bindaddr"`
	Path     string `xml:"path"`
	TempPath string `xml:"temppath"`
}

type xmlData struct {
	DataPath string        `xml:"datapath"`
	DataName string        `xml:"dataname"`
	DataTime time.Duration `xml:"datatimes"`
}

type xmlConfig struct {
	XMLName  xml.Name  `xml:"config"`
	Server   xmlServer `xml:"server"`
	Log      string    `xml:"log"`
	Data     xmlData   `xml:"data"`
	PageSize int       `xml:"pagesize"`
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
	fp.Close()
	return true
}

func LoadConf(filename string) bool {
	if !GlobalConf.LoadConfig(filename) {
		return false
	}

	return true
}
