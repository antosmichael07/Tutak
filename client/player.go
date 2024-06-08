package main

import (
	"encoding/json"
	"math"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init_player_events(client *tcp.Client, name string, player *rlfp.Player, start_window *bool, players *[]Player) {
	client.On("starter-data", func(data []byte) {
		starter_data := rl.Vector3{}
		err := json.Unmarshal(data, &starter_data)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling starter data: %s", err)
			return
		}
		player.Position = starter_data
		player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
		player.Camera.Target = rl.NewVector3(
			player.Camera.Position.X-float32(math.Cos(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
			player.Camera.Position.Y+float32(math.Sin(float64(player.Rotation.Y)))+(player.Scale.Y/2),
			player.Camera.Position.Z-float32(math.Sin(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
		)

		logger.Log(lgr.Info, "Starter data received: %v", starter_data)

		*start_window = true
	})

	player_update_events(client, player, name, players)

	client.OnConnect(func() {
		logger.Log(lgr.Info, "Connecting to server as a new player: %s", name)
		client.SendData("new-player", []byte(name))
	})
}

func player_update_events(client *tcp.Client, player *rlfp.Player, name string, players *[]Player) {
	client.On("position", func(data []byte) {
		position := PlayerPositionToSend{}
		err := json.Unmarshal(data, &position)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling position data: %s", err)
			return
		}
		if position.Name == name {
			player.Position.X = position.X
			player.Position.Y = position.Y
			player.Position.Z = position.Z
			player.UpdateCameraFirstPerson()
		} else {
			is_exist := false
			for i := range *players {
				if (*players)[i].Name == position.Name {
					(*players)[i].RLFP.Position.X = position.X
					(*players)[i].RLFP.Position.Y = position.Y
					(*players)[i].RLFP.Position.Z = position.Z
					is_exist = true
				}
			}
			if !is_exist {
				new_player := Player{
					Name: position.Name,
					RLFP: rlfp.Player{},
				}
				new_player.RLFP.InitPlayer()
				new_player.RLFP.Position.X = position.X
				new_player.RLFP.Position.Y = position.Y
				new_player.RLFP.Position.Z = position.Z
				*players = append(*players, new_player)
			}
		}
	})

	client.On("rotation", func(data []byte) {
		rotation := PlayerRotationToSend{}
		err := json.Unmarshal(data, &rotation)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling rotation data: %s", err)
			return
		}
		if rotation.Name == name {
			player.Rotation.X = rotation.X
			player.Rotation.Y = rotation.Y
			player.UpdateCameraFirstPerson()
		} else {
			is_exist := false
			for i := range *players {
				if (*players)[i].Name == rotation.Name {
					(*players)[i].RLFP.Rotation.X = rotation.X
					(*players)[i].RLFP.Rotation.Y = rotation.Y
					is_exist = true
				}
			}
			if !is_exist {
				new_player := Player{
					Name: rotation.Name,
					RLFP: rlfp.Player{},
				}
				new_player.RLFP.InitPlayer()
				new_player.RLFP.Rotation.X = rotation.X
				new_player.RLFP.Rotation.Y = rotation.Y
				*players = append(*players, new_player)
			}
		}
	})
}

func player_updates(client *tcp.Client, player *rlfp.Player) {
	input_player(client, player)
}

func input_player(client *tcp.Client, player *rlfp.Player) {
	if rl.IsKeyDown(player.Controls.Forward) || rl.IsKeyDown(player.Controls.Backward) || rl.IsKeyDown(player.Controls.Left) || rl.IsKeyDown(player.Controls.Right) || rl.IsKeyDown(player.Controls.Jump) || rl.IsKeyDown(player.Controls.Crouch) || rl.IsKeyDown(player.Controls.Interact) {
		inputs := Input{}
		if rl.IsKeyDown(player.Controls.Forward) {
			inputs.Forward = true
		}
		if rl.IsKeyDown(player.Controls.Backward) {
			inputs.Backward = true
		}
		if rl.IsKeyDown(player.Controls.Left) {
			inputs.Left = true
		}
		if rl.IsKeyDown(player.Controls.Right) {
			inputs.Right = true
		}
		if rl.IsKeyDown(player.Controls.Jump) {
			inputs.Jump = true
		}
		if rl.IsKeyDown(player.Controls.Crouch) {
			inputs.Crouch = true
		}
		if rl.IsKeyDown(player.Controls.Sprint) {
			inputs.Sprint = true
		}
		if rl.IsKeyDown(player.Controls.Interact) {
			inputs.Interact = true
		}

		to_send, err := json.Marshal(inputs)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling input data: %s", err)
			return
		}
		client.SendData("input", to_send)
	}

	if rl.GetMouseDelta().X != 0 || rl.GetMouseDelta().Y != 0 {
		to_send, err := json.Marshal(player.Rotation)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling rotate data: %s", err)
			return
		}
		client.SendData("rotate", to_send)
	}
}
