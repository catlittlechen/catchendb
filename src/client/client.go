package main

import (
	"catchendb/src/logic"
	"catchendb/src/util"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
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
	fmt.Printf("username----> %s\n host----> %s\t\tport---> %d\n", *username, *host, *port)
	if *displayHelp {
		flag.PrintDefaults()
		return false
	}
	return true
}

func mainloop() {
	bp := make([]byte, 1024)
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
		bp = bp[:count-1]
		password = (*string)(unsafe.Pointer(&bp))

	}

	urlData = url.Values{}
	urlData.Add(logic.URL_CMD, logic.CMD_AUT)
	urlData.Add(logic.URL_USER, *username)
	urlData.Add(logic.URL_PASS, *password)
	out.WriteString("ccdb>")
	_, err = conn.Write([]byte(urlData.Encode()))
	if err != nil {
		out.WriteString("Fatal Error " + err.Error() + "\n")
		return
	}
	data := make([]byte, 1024)
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
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	mainloop()
	return
}
