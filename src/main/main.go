package main

import (
	"catchendb/src/config"
	"catchendb/src/logic"
	"catchendb/src/node"
	"catchendb/src/util"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

import lgd "code.google.com/p/log4go"

var configFile = flag.String("config", "../etc/config.xml", "configFile")
var displayHelp = flag.Bool("help", false, "display HelpMessage")

func handleServer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			lgd.Error("recover %s", r)
			lgd.Error("stack %s", debug.Stack())
		}
	}()

	res := logic.LYW(r)
	bodysize := len(res)
	w.Header().Add(http.CanonicalHeaderKey("Content-Length"), strconv.Itoa(bodysize))
	w.Header().Add(http.CanonicalHeaderKey("Content-Type"), "application/json")
	w.Write(res)
	return
}

func mainloop() {

	http.HandleFunc(config.GlobalConf.Server.Path, handleServer)
	err := http.ListenAndServe(config.GlobalConf.Server.BindAddr, nil)
	if err != nil {
		lgd.Error("ListenAndServer[server] error[%s]", err)
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
