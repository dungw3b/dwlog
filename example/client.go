/*
* @project  dwlog
* @file     client.go
* @author   dungw3b
* @date     2019-08-20
*/
package main

import (
	"fmt"
	"github.com/dungw3b/dwlog"
)

func main() {
	log := dwlog.DWLog {
		Name: "test",
		Server: "127.0.0.1:2701",
	}
	err := log.Connect()
	if err != nil {
		fmt.Println(err)
	}
	defer log.Close()

	log.Error("error1 ", "error2 ", "error3 ")
	log.Info("info1 ", "info2 ", "info3 ")
	log.Debug("debug1 ", "debug2 ", "debug3 ")
	
	select{}
}