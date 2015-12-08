package main

import (
	"catchendb/src/client/handle"
	"catchendb/src/util"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

type Req struct {
	C string `json:"c"`

	UserName  string `json:"usr"`
	PassWord  string `json:"pas"`
	Privilege int    `json:"pri"`

	Key       string `json:"key"`
	Value     string `json:"val"`
	StartTime int64  `json:"sta"`
	EndTime   int64  `json:"end"`
}

var displayHelp = flag.Bool("help", false, "displayHelpMessage")
var username = flag.String("u", "root", "username")
var password = flag.String("p", "root", "password")
var host = flag.String("h", "127.0.0.1", "host")
var port = flag.Int("P", 13570, "port")
var capt = flag.Int("c", 13570, "capt")
var onlyget = flag.Bool("o", false, "onlyget")
var onlyset = flag.Bool("s", false, "onlyget")
var gorout = flag.Int("g", 1, "goroutince")

var wg sync.WaitGroup

func Init() bool {
	flag.Parse()
	fmt.Printf("start time %s\n", util.FormalTime(time.Now().Unix()))
	if *displayHelp {
		flag.PrintDefaults()
		return false
	}
	return true
}

func mainloop(capts, begin int) {
	defer wg.Done()
	bp := ""
	data := make([]byte, 10240)
	data2 := make([]byte, 10240)
	var req Req
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

	req = Req{
		C:        handle.CMD_AUT,
		UserName: *username,
		PassWord: *password,
	}
	_, err = conn.Write(util.JsonOut(req))
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

	countap := begin - 1
	end := begin + capts
	var code string
	ok := false
	var fun func([]byte) []byte
	for {
		countap += 1
		if countap == end {
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
				fmt.Println("wrong argv")
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
			} else if len(rsp.D) != 0 {
				fmt.Println(rsp.D + "\n")
			}
		}
		if !*onlyset {
			bp = fmt.Sprintf("get %d", countap)

			code = (strings.Split(bp, string(' ')))[0]
			fun, ok = handle.GetHandle(code)
			if !ok {
				fmt.Printf("wrong command[%s]\n", code)
				continue
			}
			data2 = fun([]byte(bp))
			if data2 == nil {
				fmt.Println("wrong argv")
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
				fmt.Printf("%d\n", countap)
			}

		}
	}
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	fmt.Println(time.Now().UnixNano())
	capts := *capt
	for i := 0; i < *gorout; i++ {
		wg.Add(1)
		go mainloop(capts, capts*i)
	}
	wg.Wait()
	fmt.Println(time.Now().UnixNano())
	return
}
