package dwlog

import (
	"fmt"
	"time"
	"strings"
	"github.com/nats-io/nats.go"
)

func (app *DWLog) Connect() error {
	if len(app.Name) == 0 {
		app.Name = "dwlog"
	}
	if len(app.Host) == 0 {
		app.Host = "DWLog"
	}
	if len(app.Servers) == 0 {
		return fmt.Errorf("List servers are required")
	}
	if app.Timeout == 0 {
		app.Timeout = time.Duration(60*time.Second)
	}
	if app.Level == 0 {
		app.Level = ERROR
	}
	app.tformat = "2006-01-02 15:04:05"

	var err error
	app.conn, err = nats.Connect(strings.Join(app.Servers, ","), nats.Name(app.Name), nats.Timeout(app.Timeout), nats.NoEcho())
	if err != nil {
		return err
	}
	app.encode, err = nats.NewEncodedConn(app.conn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}

	return nil
}

func (app *DWLog) Close() {
	if app.conn != nil {
		if app.encode != nil {
			app.encode.Close()
		}
		app.conn.Close()
	}
}

func (app *DWLog) Error(msg string) {
	if app.Level < ERROR {
		return
	}
	go app.Log(ERROR, msg)
}

func (app *DWLog) Info(msg string) {
	if app.Level < INFO {
		return
	}
	go app.Log(INFO, msg)
}

func (app *DWLog) Debug(msg string) {
	if app.Level < DEBUG {
		return
	}
	go app.Log(DEBUG, msg)
}

func (app *DWLog) Log(level Level, msg string) {
	message := &Message {
		Time: time.Now().Format(app.tformat),
		Level: level,
		Message: msg,
	}
	err := app.encode.Publish(app.Name, message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message)
}