/*
* @project  dwlog
* @file     server.go
* @author   dungw3b
* @date     2019-08-20
*/
package main

import (
	"os"
	"os/signal"
	"fmt"
	"flag"
	"syscall"
	"io/ioutil"
	"github.com/olebedev/config"
	"github.com/dungw3b/dwlog"
)

var App *dwlog.Server

func init() {
	path := flag.String("c", "", "full path to config file Ex. conf/dwlog.json")
	flag.Parse()
	if len(*path) == 0 {
		fmt.Println("\nUsage of dwlog:");
		flag.PrintDefaults();
	}
	data, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Println("Can not read config file "+ *path)
		os.Exit(1)
	}
	cfg, err := config.ParseJson(string(data))
	if err != nil {
		fmt.Println("Can not parse JSON config file "+ *path)
		os.Exit(1)
	}
	
	App = &dwlog.Server {
		Listen: cfg.UString("service.listen"),
		Port: uint32(cfg.UInt("service.port")),
		FileCount: uint32(cfg.UInt("service.filecount")),
		Data: cfg.UString("service.data"),
	}
}

func main() {
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan bool)
	signal.Notify(signalChan)
	signal.Ignore(syscall.SIGHUP)

	go App.Run(exitChan)

	// Signal manager
	go func() {
		<- signalChan
		App.Close()
		exitChan <- true
	}()
	<- exitChan
}