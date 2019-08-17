package main

import (
	"fmt"
	"time"
	"github.com/dungw3b/dwlog"
)

func main() {
	app := &dwlog.Server {
		Subject: "test",
		Servers: []string{"nats://127.0.0.1:4222"},
		Timeout: 10 * time.Second,
	}

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}