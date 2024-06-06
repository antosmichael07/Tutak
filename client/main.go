package main

import (
	"encoding/json"
	"time"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

type Player struct {
	Name   string
	Player PlayerFP
}

type PlayerPackage struct {
	Name string
	RFLP []byte
}

func main() {
	client := tcp.NewClient("localhost:8080")
	client.Connect()

	name := "Test"
	players := []Player{}
	bounding_boxes := []rl.BoundingBox{}
	trigger_boxes := []TriggerBox{}
	interractable_boxes := []InteractableBox{}
	my_player := PlayerFP{}

	client.On("players", func(data []byte) {
		player_package := PlayerPackage{}
		err := json.Unmarshal(data, &player_package)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling players data: %s", err)
		}
		rlfp := PlayerFP{}
		err = json.Unmarshal(player_package.RFLP, &rlfp)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling rlfp data: %s", err)
		}
		players = []Player{{Name: player_package.Name, Player: rlfp}}
		for i := range players {
			if players[i].Name == name {
				my_player = players[i].Player
			}
		}
	})

	client.OnConnect(func() {
		client.SendData("initialize_player", []byte(name))

		logger.Log(lgr.Info, "Connected to server")
	})

	go client.Listen()
	client.Logger.Level = lgr.Warning

	current_monitor := rl.GetCurrentMonitor()
	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Raylib 3D Custom First Person - Example")
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
	defer rl.CloseWindow()

	time.Sleep(1 * time.Second)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(my_player.Camera)

		rl.DrawGrid(20, 1.)
		for i := range bounding_boxes {
			rl.DrawBoundingBox(bounding_boxes[i], rl.Red)
		}
		for i := range trigger_boxes {
			rl.DrawBoundingBox(trigger_boxes[i].BoundingBox, rl.Green)
		}
		for i := range interractable_boxes {
			rl.DrawBoundingBox(interractable_boxes[i].BoundingBox, rl.Blue)
		}

		rl.EndMode3D()

		rl.DrawFPS(10, 10)
		players[0].Player.UpdatePlayer(bounding_boxes, trigger_boxes, interractable_boxes, client)

		rl.EndDrawing()
	}
}
