package dwlog

import (
	_"fmt"
	"time"
	_"sync"
	"github.com/nats-io/nats.go"
)

type DWLog struct {
	Name string
	Host string
	Servers []string
	Timeout time.Duration
	Level Level

	tformat string
	conn *nats.Conn
	encode *nats.EncodedConn
}

/*
func (app *DWLog) Log1111(level Level, msg string) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	var (
		line []byte
		message *Message
	)
	message = &Message {
		Time: time.Now().Format(app.tformat),
		Level: level,
		Message: msg,
	}

	line = message.CSVFormat(app)
	line = append(line, '\n')
	_, err := app.writer.Write(line)
	if err != nil {
		fmt.Println(err)
	}
}
*/