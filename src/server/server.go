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
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

import lgd "code.google.com/p/log4go"

var configFile = flag.String("config", "../etc/config.xml", "configFile")
var displayHelp = flag.Bool("help", false, "display HelpMessage")

func handleReplicationServer(conn *net.TCPConn) {
	defer func() {
		if re := recover(); re != nil {
			lgd.Error("recover %s", re)
			lgd.Error("stack %s", debug.Stack())
		}
	}()
	defer conn.Close()
	logic.ReplicationLogic(conn)
}

func replicationloop() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.GlobalConf.Server.ReplicationAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ResolveTCPAddr[%s] error[%s]", config.GlobalConf.Server.ReplicationAddr, err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listenTCP[%s] error[%s]", tcpAddr, err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			lgd.Error("listener Accept error[%s]", err)
			return
		}

		go handleReplicationServer(conn)
	}
}

func handleServer(conn *net.TCPConn) {
	defer func() {
		if re := recover(); re != nil {
			lgd.Error("recover %s", re)
			lgd.Error("stack %s", debug.Stack())
		}
	}()
	defer conn.Close()
	logic.ClientLogic(conn)
}

func mainloop() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.GlobalConf.Server.BindAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ResolveTCPAddr[%s] error[%s]", config.GlobalConf.Server.BindAddr, err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listenTCP[%s] error[%s]", tcpAddr, err)
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
	config.LoadConf(*configFile)

	lgd.LoadConfiguration(config.GlobalConf.Log)
	os.Chmod(config.GlobalConf.Log, 0666)
	runtime.GOMAXPROCS(runtime.NumCPU())
	return node.Init() && logic.LoadData() && true
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	lgd.Info("start")
	if !logic.Init() {
		return
	}
	go replicationloop()
	mainloop()
	return
}
