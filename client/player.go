package main

import (
	"encoding/json"
	"math"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (player *Player) InitPlayer() {
	player.RLFP.Speed.Sneak = .04
	player.RLFP.Speed.Normal = .07
	player.RLFP.Speed.Sprint = .1
	player.RLFP.Speed.Acceleration = .0075
	player.RLFP.Gravity = .0075
	player.RLFP.JumpPower = .135
}

func init_player_events(client *tcp.Client, name string, player *Player, start_window *bool, players *[]Player) {
	client.On("starter-data", func(data []byte) {
		starter_data := rl.Vector3{}
		err := json.Unmarshal(data, &starter_data)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling starter data: %s", err)
			return
		}
		player.RLFP.Position = starter_data
		player.RLFP.Camera.Position = rl.NewVector3(player.RLFP.Position.X, player.RLFP.Position.Y+(player.RLFP.Scale.Y/2), player.RLFP.Position.Z)
		player.RLFP.Camera.Target = rl.NewVector3(
			player.RLFP.Camera.Position.X-float32(math.Cos(float64(player.RLFP.Rotation.X)))*float32(math.Cos(float64(player.RLFP.Rotation.Y))),
			player.RLFP.Camera.Position.Y+float32(math.Sin(float64(player.RLFP.Rotation.Y)))+(player.RLFP.Scale.Y/2),
			player.RLFP.Camera.Position.Z-float32(math.Sin(float64(player.RLFP.Rotation.X)))*float32(math.Cos(float64(player.RLFP.Rotation.Y))),
		)

		logger.Log(lgr.Info, "Starter data received: %v", starter_data)

		*start_window = true
	})

	client.On("wrong-name", func(data []byte) {
		logger.Log(lgr.Error, "Name already exists: %s", string(data))
		client.Disconnect()
	})

	client.On("disconnected-player", func(data []byte) {
		player_name := string(data)
		for i := range *players {
			if (*players)[i].Name == player_name {
				*players = append((*players)[:i], (*players)[i+1:]...)
				break
			}
		}
	})

	player_update_events(client, player, name, players)

	client.OnConnect(func() {
		logger.Log(lgr.Info, "Connecting to server as a new player: %s", name)
		client.SendData("new-player", []byte(name))
	})
}

func player_update_events(client *tcp.Client, player *Player, name string, players *[]Player) {
	client.On("position", func(data []byte) {
		logger.Log(lgr.Info, "Position data received: %s", data)
		position := PlayerPositionToSend{}
		err := json.Unmarshal(data, &position)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling position data: %s", err)
			return
		}
		if position.Name == name {
			if player.RLFP.Position.X+.5 < position.X || player.RLFP.Position.X-.5 > position.X || player.RLFP.Position.Y+.5 < position.Y || player.RLFP.Position.Y-.5 > position.Y || player.RLFP.Position.Z+.5 < position.Z || player.RLFP.Position.Z-.5 > position.Z {
				player.RLFP.Position.X = position.X
				player.RLFP.Position.Y = position.Y
				player.RLFP.Position.Z = position.Z
				player.RLFP.UpdateCameraFirstPerson()
			}
			offset_position := PlayerPositionToSend{
				Name: name,
				X:    player.RLFP.Position.X,
				Y:    player.RLFP.Position.Y,
				Z:    player.RLFP.Position.Z,
			}
			to_send, err := json.Marshal(offset_position)
			if err != nil {
				logger.Log(lgr.Error, "Error marshalling offset data: %s", err)
				return
			}
			client.SendData("offset-position", to_send)
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
		if rotation.Name != name {
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

func player_updates(client *tcp.Client, player *Player) {
	input_player(client, player)
}

func input_player(client *tcp.Client, player *Player) {
	inputs := Input{}
	if rl.IsKeyDown(player.RLFP.Controls.Forward) {
		inputs.Forward = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Backward) {
		inputs.Backward = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Left) {
		inputs.Left = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Right) {
		inputs.Right = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Jump) {
		inputs.Jump = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Crouch) {
		inputs.Crouch = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Sprint) {
		inputs.Sprint = true
	}
	if rl.IsKeyDown(player.RLFP.Controls.Interact) {
		inputs.Interact = true
	}

	if player.LastKeysPressed != inputs {
		player.LastKeysPressed = inputs

		to_send, err := json.Marshal(inputs)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling input data: %s", err)
			return
		}
		client.SendData("input", to_send)
	}

	if player.LastRotation != player.RLFP.Rotation {
		to_send, err := json.Marshal(player.RLFP.Rotation)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling rotate data: %s", err)
			return
		}
		player.LastRotation = player.RLFP.Rotation
		client.SendData("rotate", to_send)
	}
}
