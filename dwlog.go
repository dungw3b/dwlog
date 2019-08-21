/*
* @project  dwlog
* @file     dwlog.go
* @author   dungw3b
* @date     2019-08-20
*/
package dwlog

import (
	"os"
	"fmt"
	"time"
	"context"
	"google.golang.org/grpc"
	"github.com/hectane/go-nonblockingchan"
)

type DWLog struct {
	Name string
	Host string
	Level Level
	Tformat string
	Server string
	Timeout time.Duration

	conn *grpc.ClientConn
	client LogServiceClient
	nbchan *nbc.NonBlockingChan
}

func (log *DWLog) Connect() error {
	if len(log.Server) == 0 {
		return fmt.Errorf("log.Server is required")
	}
	if log.Timeout == 0 {
		log.Timeout = time.Duration(20*time.Second)
	}
	if len(log.Host) == 0 {
		log.Host, _ = os.Hostname()
	}
	if len(log.Tformat) == 0 {
		log.Tformat = "2006-01-02 15:04:05"
	}
	if log.Level == 0 {
		log.Level = ERROR
	}

	var err error
	log.conn, err = grpc.Dial(log.Server, grpc.WithInsecure(), grpc.WithTimeout(log.Timeout))
	if err != nil {
		return err
	}
	log.client = NewLogServiceClient(log.conn)
	log.nbchan = nbc.New()
	go log.run()

	return nil
}

func (log *DWLog) run() {
	for {
		select {
		case item := <- log.nbchan.Recv:
			message := item.(*MessageRequest)
			if message == nil {
				return
			}
			// push message
			_, err := log.client.Log(context.Background(), message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (log *DWLog) send(level Level, str string) {
	message := &MessageRequest {
		Name: log.Name,
		Host: log.Host,
		Time: time.Now().Format(log.Tformat),
		Level: level.Val(),
		Message: str,
	}
	log.nbchan.Send <- message
}

func (log *DWLog) Error(str ...string) {
	if log.Level > ERROR {
		return
	}
	var msg string
	for _,s := range str {
		msg = msg + s
	}
	if len(msg) == 0 {
		return
	}
	log.send(ERROR, msg)
}

func (log *DWLog) Info(str ...string) {
	if log.Level > INFO {
		return
	}
	var msg string
	for _,s := range str {
		msg = msg + s
	}
	if len(msg) == 0 {
		return
	}
	log.send(INFO, msg)
}

func (log *DWLog) Debug(str ...string) {
	if log.Level > DEBUG {
		return
	}
	var msg string
	for _,s := range str {
		msg = msg + s
	}
	if len(msg) == 0 {
		return
	}
	log.send(DEBUG, msg)
}

func (log *DWLog) Close() {
	if log.conn != nil {
		log.conn.Close()
	}
	log.nbchan.Send <- nil
	fmt.Println("log closed")
}