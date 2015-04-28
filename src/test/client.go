package main

import (
	"catchendb/src/client/handle"
	"catchendb/src/util"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

var displayHelp = flag.Bool("help", false, "displayHelpMessage")
var username = flag.String("u", "root", "username")
var password = flag.String("p", "root", "password")
var host = flag.String("h", "127.0.0.1", "host")
var port = flag.Int("P", 13570, "port")
var capt = flag.Int("c", 13570, "capt")
var onlyget = flag.Bool("o", false, "onlyget")

func Init() bool {
	flag.Parse()
	fmt.Printf("start time %s\n", util.FormalTime(time.Now().Unix()))
	if *displayHelp {
		flag.PrintDefaults()
		return false
	}
	return true
}

func mainloop() {
	bp := ""
	data := make([]byte, 10240)
	data2 := make([]byte, 10240)
	var urlData url.Values
	var count int
	var err error

	serverHost := fmt.Sprintf("%s:%d", *host, *port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverHost)
	if err != nil {
		fmt.Printf("Fatal Error %s\n", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("Fatal Error %s\n", err)
		return
	}

	urlData = url.Values{}
	urlData.Add(handle.URL_CMD, handle.CMD_AUT)
	urlData.Add(handle.URL_USER, *username)
	urlData.Add(handle.URL_PASS, *password)
	_, err = conn.Write([]byte(urlData.Encode()))
	if err != nil {
		fmt.Println("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}
	count, err = conn.Read(data)
	if err != nil {
		fmt.Println("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}

	var rsp Rsp
	err = json.Unmarshal(data[:count], &rsp)
	if err != nil {
		fmt.Println("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}

	if rsp.C != 0 {
		fmt.Printf("ccdb>ERROR %d Access denied for user '%s'@'%s' (using password: YES)\n", rsp.C, *username, serverHost)
		return
	}

	countap := 1
	fmt.Println(time.Now().Unix())
	var code string
	ok := false
	var fun func([]byte) []byte
	for {
		countap += 1
		if countap == *capt {
			fmt.Println(time.Now().Unix())
			break
		}
		if !*onlyget {
			bp = fmt.Sprintf("set %d %d", countap, countap)

			code = (strings.Split(bp, string(' ')))[0]
			fun, ok = handle.GetHandle(code)
			if !ok {
				fmt.Printf("wrong command[%s]\n", code)
				continue
			}
			data2 = fun([]byte(bp))
			if data2 == nil {
				fmt.Println("wrong argv\n")
				continue
			}
			_, err = conn.Write(data2)
			if err != nil {
				fmt.Println("Fatal Error " + err.Error() + "\n")
				return
			}
			count, err = conn.Read(data)
			if err != nil {
				fmt.Println("Fatal Error " + err.Error() + "\n")
				return
			}

			var rsp Rsp
			err = json.Unmarshal(data[:count], &rsp)
			if err != nil {
				fmt.Println("Fatal Data " + string(data[:count]) + "Error " + err.Error() + "\n")
				return
			}

			if rsp.C != 0 {
				fmt.Printf("ERROR %d \n", rsp.C)
				continue
			} else if len(rsp.D) != 0 {
				fmt.Println(rsp.D + "\n")
			}
		}
		bp = fmt.Sprintf("get %d", countap)

		code = (strings.Split(bp, string(' ')))[0]
		fun, ok = handle.GetHandle(code)
		if !ok {
			fmt.Printf("wrong command[%s]\n", code)
			continue
		}
		data2 = fun([]byte(bp))
		if data2 == nil {
			fmt.Println("wrong argv\n")
			continue
		}
		_, err = conn.Write(data2)
		if err != nil {
			fmt.Println("Fatal Error " + err.Error() + "\n")
			return
		}
		count, err = conn.Read(data)
		if err != nil {
			fmt.Println("Fatal Error " + err.Error() + "\n")
			return
		}

		err = json.Unmarshal(data[:count], &rsp)
		if err != nil {
			fmt.Println("Fatal Data " + string(data[:count]) + "Error " + err.Error() + "\n")
			return
		}

		if rsp.C != 0 {
			fmt.Printf("ERROR %d \n", rsp.C)
			continue
		} else if rsp.D != fmt.Sprintf("%d", countap) {
			fmt.Println(rsp.D + "\n")
		}
	}
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	mainloop()
	return
}
