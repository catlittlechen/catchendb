package main

import (
	"catchendb/src/config"
	"catchendb/src/logic"
	"catchendb/src/node"
	"catchendb/src/util"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"runtime/debug"
	"syscall"
	"time"
)

import lgd "code.google.com/p/log4go"

var configFile = flag.String("config", "../etc/config.xml", "configFile")
var displayHelp = flag.Bool("help", false, "display HelpMessage")

func handleServer(conn *net.TCPConn) {
	defer func() {
		if re := recover(); re != nil {
			lgd.Error("recover %s", re)
			lgd.Error("stack %s", debug.Stack())
		}
	}()
	defer conn.Close()

	data := make([]byte, 1024)
	count, err := conn.Read(data)
	if err != nil {
		lgd.Error("read error[%s]", err)
		return
	}
	ok, res := logic.AUT(data[:count])
	conn.Write(res)
	if !ok {
		return
	}
	for {
		count, err = conn.Read(data)
		if err != nil {
			lgd.Error("read error[%s]", err)
			return
		}
		res := logic.LYW(data[:count])
		conn.Write(res)
	}
}

func mainloop() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.GlobalConf.Server.Path)
	if err != nil {
		lgd.Error("ResolveTCPAddr[server] error[%s]", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		lgd.Error("listenTCP[%s] error[%s]", tcpAddr, err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			lgd.Error("listener Accept error[%s]", err)
			return
		}

		go handleServer(conn)
	}
}

func Init() bool {
	flag.Parse()
	fmt.Printf("start time %s\n", util.FormalTime(time.Now().Unix()))
	fmt.Printf("start version %s\n", VERSION)
	fmt.Printf("help %t, config:%s\n", *displayHelp, *configFile)
	if *displayHelp || *configFile == "" {
		flag.PrintDefaults()
		return false
	}
	syscall.Umask(0)
	os.Chdir(path.Dir(os.Args[0]))
	if config.LoadConf(*configFile) {
		lgd.LoadConfiguration(config.GlobalConf.Log)
		os.Chmod(config.GlobalConf.Log, 0666)
	} else {
		return false
	}
	return node.Init() && logic.LoadData() && true
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	lgd.Info("start")
	logic.Init()
	mainloop()
	return
}
