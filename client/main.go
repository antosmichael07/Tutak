package main

import (
	"time"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

func main() {
	client := tcp.NewClient("6.tcp.eu.ngrok.io:12544")
	client.Connect()

	go func() {
		for {
			if rl.IsKeyDown(rl.KeyF9) {
				rl.CloseWindow()
				client.Disconnect()
			}
		}
	}()

	start_window := false
	name := "Test"
	players := []Player{}
	bounding_boxes := []rl.BoundingBox{}
	trigger_boxes := []rlfp.TriggerBox{}
	interractable_boxes := []rlfp.InteractableBox{}
	player := rlfp.Player{}
	player.InitPlayer()

	init_player_events(&client, name, &player, &start_window, &players)

	client.Logger.Level = lgr.Warning
	go client.Listen()

	for !start_window {
		time.Sleep(500 * time.Millisecond)
	}
	init_window()
	window_loop(&client, &player, bounding_boxes, trigger_boxes, interractable_boxes, &players)
}
