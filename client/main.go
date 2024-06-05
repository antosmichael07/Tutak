package main

import (
	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
)

var logger = lgr.NewLogger("Tutak")

func main() {
	client := tcp.NewClient("localhost:8080")
	client.Connect()

	client.On("test", func(data string) {
		logger.Log(lgr.Info, "wow")
	})

	client.On("test2", func(data string) {
		logger.Log(lgr.Info, "bomba")
	})

	client.SendData("execute", "Hello from client")

	client.Listen()
}
