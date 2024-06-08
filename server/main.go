package main

import (
	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

func main() {
	rl.SetTraceLogLevel(rl.LogError)
	server := tcp.NewServer("localhost:8080")

	players := []Player{}
	bounding_boxes := []rl.BoundingBox{}
	trigger_boxes := []TriggerBox{}
	interractable_boxes := []InteractableBox{}

	init_player_events(&server, &players, bounding_boxes, trigger_boxes, interractable_boxes)

	server.Logger.Level = lgr.Warning
	server.Start()
}
