package main

import (
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rlci "github.com/antosmichael07/Raylib-Cube-Image"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init_window() {
	current_monitor := rl.GetCurrentMonitor()
	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Tutak")
	rl.SetExitKey(-1)
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
}

func window_loop(client *tcp.Client, player *rlfp.Player, bounding_boxes []rl.BoundingBox, trigger_boxes []rlfp.TriggerBox, interractable_boxes []rlfp.InteractableBox, players *[]Player) {
	image := rl.LoadImage("image.png")
	cube := rlci.NewCubeImage(image, rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1), rl.White)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(player.Camera)

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

		for i := range *players {
			cube.Position = rl.NewVector3((*players)[i].RLFP.Position.X, (*players)[i].RLFP.Position.Y, (*players)[i].RLFP.Position.Z)
			cube.RotationAngle = (*players)[i].RLFP.Rotation.X
			cube.DrawCubeImage()
		}

		rl.EndMode3D()

		rl.DrawFPS(10, 10)
		position := rl.NewVector3(player.Position.X, player.Position.Y, player.Position.Z)
		player.UpdatePlayer(bounding_boxes, trigger_boxes, interractable_boxes)
		player.Position = rl.NewVector3(position.X, position.Y, position.Z)
		player_updates(client, player)

		rl.EndDrawing()
	}
}
