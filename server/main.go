package main

import (
	"time"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
)

var logger = lgr.NewLogger("Tutak")

func main() {
	server := tcp.NewServer("localhost:8080")

	server.On("execute", func(data string) {
		server.SendData("test", "Hello from server")
		time.Sleep(3*time.Second)
		server.SendData("test2", "Hello from server")
	})

	server.Start()
}