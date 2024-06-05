package main

import (
	"encoding/json"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl_fp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

type Player struct {
	Name   string
	Player rl_fp.Player
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
	trigger_boxes := []rl_fp.TriggerBox{}
	interractable_boxes := []rl_fp.InteractableBox{}
	my_player := rl_fp.Player{}

	client.On("players", func(data []byte) {
		player_package := PlayerPackage{}
		err := json.Unmarshal(data, &player_package)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling players data: %s", err)
		}
		rlfp := rl_fp.Player{}
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

		rl.EndDrawing()
	}
}
