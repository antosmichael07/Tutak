package main

import (
	custom_logger "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
)

var logger = custom_logger.NewLogger()

func main() {
	client := tcp.NewClient("localhost:8080")
	client.Connect()

	client.On("test", func(data string) {
		logger.Log("Data received: %s", data)
	})

	client.SendData("test", "Hello from client")

	client.Listen()
}
