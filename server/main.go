package main

import (
	"time"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

func main() {
	rl.SetTraceLogLevel(rl.LogError)
	server := tcp.NewServer("localhost:8080")
	server.Logger.Output.File = true

	players := []Player{}
	bounding_boxes := []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
	}
	trigger_boxes := []TriggerBox{}
	interractable_boxes := []InteractableBox{}

	init_player_events(&server, &players, bounding_boxes)

	server.Logger.Level = lgr.Warning
	go server.Start()

	for !server.ShouldStop {
		go player_updates(&server, &players, bounding_boxes, trigger_boxes, interractable_boxes)

		time.Sleep(15090 * time.Microsecond)
	}
}
