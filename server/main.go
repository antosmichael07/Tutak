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
	bounding_boxes := []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(-1.5, -.5, -.5), rl.NewVector3(-.5, .5, .5)),
		rl.NewBoundingBox(rl.NewVector3(-2.5, 0., -.5), rl.NewVector3(-1.5, 1., .5)),
		rl.NewBoundingBox(rl.NewVector3(-4.5, .5, -.5), rl.NewVector3(-3.5, 1.5, .5)),
		rl.NewBoundingBox(rl.NewVector3(-5.5, 1., -.5), rl.NewVector3(-4.5, 2., .5)),
	}
	trigger_boxes := []TriggerBox{
		NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(2.5, 1., -.5), rl.NewVector3(3.5, 2., .5))),
		NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(4.5, 2.5, -.5), rl.NewVector3(5.5, 3.5, .5))),
	}
	interractable_boxes := []InteractableBox{
		NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(7.5, 0., -.5), rl.NewVector3(8.5, 1., .5))),
		NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(7.5, .5, -.5), rl.NewVector3(8.5, 1.5, .5))),
	}

	init_player_events(&server, &players, bounding_boxes, trigger_boxes, interractable_boxes)

	server.Logger.Level = lgr.Warning
	server.Start()
}
