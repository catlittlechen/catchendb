package main

import (
	"catchendb/src/config"
	"catchendb/src/logic"
	"catchendb/src/util"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime/debug"
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

	return
}

func mainloop() {

	logic.AutoSaveData()

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
	return config.LoadConf(*configFile) && logic.LoadData() && true
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}

	lgd.Info("start")
	mainloop()
	return
}
