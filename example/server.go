package main

import (
	"os"
	"fmt"
	"flag"
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
	if err := App.Run(); err != nil {
		fmt.Println(err)
	}
}