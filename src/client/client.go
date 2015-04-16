package main

import (
	"bufio"
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
	bp := ""
	data := make([]byte, 10240)
	var urlData url.Values
	var count int
	var err error
	in := bufio.NewReader(os.Stdin)
	out := os.Stdout
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
		bp, err = in.ReadString('\n')
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}
		bp = strings.Trim(bp, "\n")
		password = &bp
	}

	urlData = url.Values{}
	urlData.Add(handle.URL_CMD, handle.CMD_AUT)
	urlData.Add(handle.URL_USER, *username)
	urlData.Add(handle.URL_PASS, *password)
	_, err = conn.Write([]byte(urlData.Encode()))
	if err != nil {
		out.WriteString("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}
	count, err = conn.Read(data)
	if err != nil {
		out.WriteString("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}

	var rsp Rsp
	err = json.Unmarshal(data[:count], &rsp)
	if err != nil {
		out.WriteString("ccdb>Fatal Error " + err.Error() + "\n")
		return
	}

	if rsp.C != 0 {
		out.WriteString(fmt.Sprintf("ccdb>ERROR %d Access denied for user '%s'@'%s' (using password: YES)\n", rsp.C, *username, serverHost))
		return
	}

	for {
		out.WriteString("ccdb>")
		bp, err = in.ReadString('\n')
		if err != nil {
			out.WriteString("Fatal Error " + err.Error() + "\n")
			return
		}
		bp = strings.Trim(bp, "\n")
		if len(bp) == 0 {
			continue
		}
		if "exit" == bp {
			break
		}

		code := (strings.Split(bp, string(' ')))[0]
		fun, ok := handle.GetHandle(code)
		if !ok {
			out.WriteString(fmt.Sprintf("wrong command[%s]\n", code))
			continue
		}
		data = fun([]byte(bp))
		if data == nil {
			out.WriteString("wrong argv\n")
			continue
		}
		_, err = conn.Write(data)
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
