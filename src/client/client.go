package main

import (
	"catchendb/src/client/handle"
	"catchendb/src/util"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"
)

type Rsp struct {
	C int    `json:"c"`
	M string `json:"m"`
	D string `json:"d"`
}

var displayHelp = flag.Bool("help", false, "displayHelpMessage")
var username = flag.String("u", "root", "username")
var password = flag.String("p", "", "password")
var host = flag.String("h", "127.0.0.1", "host")
var port = flag.Int("P", 13570, "port")

func Init() bool {
	flag.Parse()
	fmt.Printf("start time %s\n", util.FormalTime(time.Now().Unix()))
	fmt.Printf("start version %s\n", VERSION)
	fmt.Printf("username----> %s\nhost----> %s\t\tport---> %d\n", *username, *host, *port)
	if *displayHelp {
		flag.PrintDefaults()
		return false
	}
	return true
}

func mainloop() {
	bp := make([]byte, 10240)
	data := make([]byte, 10240)
	data2 := make([]byte, 10240)
	var urlData url.Values
	var count int
	var err error
	in := os.Stdin
	out := os.Stdout
	defer in.Close()
	defer out.Close()

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

	if *password == "" {
		out.WriteString("ccdb>")
		count, err = in.Read(bp)
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}
		data2 = bp[:count-1]
		password = (*string)(unsafe.Pointer(&data2))
	}

	urlData = url.Values{}
	urlData.Add(handle.URL_CMD, handle.CMD_AUT)
	urlData.Add(handle.URL_USER, *username)
	urlData.Add(handle.URL_PASS, *password)
	out.WriteString("ccdb>")
	_, err = conn.Write([]byte(urlData.Encode()))
	if err != nil {
		out.WriteString("Fatal Error " + err.Error() + "\n")
		return
	}
	count, err = conn.Read(data)
	if err != nil {
		out.WriteString("Fatal Error " + err.Error() + "\n")
		return
	}

	var rsp Rsp
	err = json.Unmarshal(data[:count], &rsp)
	if err != nil {
		out.WriteString("Fatal Error " + err.Error() + "\n")
		return
	}

	if rsp.C != 0 {
		out.WriteString(fmt.Sprintf("ERROR %d Access denied for user '%s'@'%s' (using password: NO)\n", rsp.C, *username, serverHost))
		return
	}

	for {
		count, err = in.Read(bp)
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}
		fmt.Println(count)
		if count == 0 || count == 1 {
			out.WriteString("ccdb>")
			continue
		}
		data2 = bp[:count-1]
		if "exit" == string(data2) {
			break
		}

		code := (strings.Split(string(data2), string(' ')))[0]
		fun, ok := handle.GetHandle(code)
		if !ok {
			out.WriteString(fmt.Sprintf("wrong command[%s]\n", code))
			out.WriteString("ccdb>")
			continue
		}
		bp = fun(data2)
		if bp == nil {
			out.WriteString(fmt.Sprintf("wrong argv\n"))
			out.WriteString("ccdb>")
			continue
		}
		_, err = conn.Write(bp)
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}
		count, err = conn.Read(data)
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}

		var rsp Rsp
		err = json.Unmarshal(data[:count], &rsp)
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}

		if rsp.C != 0 {
			out.WriteString(fmt.Sprintf("ERROR %d \n", rsp.C))
			out.WriteString("ccdb>")
			continue
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
