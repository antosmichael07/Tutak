package main

import (
	custom_logger "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
)

var logger = custom_logger.NewLogger()

func main() {
	server := tcp.NewServer("localhost:8080")

	server.SendData("test", "Hello from server")

	server.Start()
}
