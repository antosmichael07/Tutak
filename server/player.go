package main

import (
	"encoding/json"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init_player_events(server *tcp.Server, players *[]Player, bounding_boxes []rl.BoundingBox) {
	server.On("new-player", func(data []byte, conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				return
			}
			if (*players)[i].Name == string(data) {
				server.SendData(conn.Connection, "wrong-name", data)
				return
			}
		}

		new_player := Player{
			Token: conn.Token,
			Name:  string(data),
			RLFP:  PlayerFP{},
		}
		new_player.RLFP.InitPlayer()
		*players = append(*players, new_player)

		logger.Log(lgr.Info, "New player connected: %s", new_player.Name)

		to_send, err := json.Marshal(new_player.RLFP.Position)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling starter data: %s", err)
			return
		}

		server.SendData(conn.Connection, "starter-data", to_send)
	})

	player_update_events(server, players, bounding_boxes)

	server.OnDisconnect(func(conn tcp.Connection) {
		player_name := ""
		player_token := ""
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				player_name = (*players)[i].Name
				player_token = (*players)[i].Token
				break
			}
		}
		for i := range server.Connections {
			if server.Connections[i].Token != player_token {
				server.SendData(server.Connections[i].Connection, "disconnected-player", []byte(player_name))
			}
		}
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				logger.Log(lgr.Info, "Player disconnected: %s", (*players)[i].Name)
				*players = append((*players)[:i], (*players)[i+1:]...)
				break
			}
		}
	})
}

func player_update_events(server *tcp.Server, players *[]Player, bounding_boxes []rl.BoundingBox) {
	server.On("input", func(data []byte, conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				err := json.Unmarshal(data, &(*players)[i].RLFP.Controls)
				if err != nil {
					logger.Log(lgr.Error, "Error unmarshalling input data: %s", err)
					return
				}
				break
			}
		}
	})

	server.On("rotate", func(data []byte, conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				err := json.Unmarshal(data, &(*players)[i].RLFP.Rotation)
				if err != nil {
					logger.Log(lgr.Error, "Error unmarshalling rotate data: %s", err)
					return
				}
				break
			}
		}
	})

	server.On("offset-position", func(data []byte, conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				position := PlayerPositionToSend{}
				err := json.Unmarshal(data, &position)
				if err != nil {
					logger.Log(lgr.Error, "Error unmarshalling offset data: %s", err)
					return
				}
				if (*players)[i].RLFP.Position.X+.5 > position.X || (*players)[i].RLFP.Position.X-.5 < position.X || (*players)[i].RLFP.Position.Y+.5 > position.Y || (*players)[i].RLFP.Position.Y-.5 < position.Y || (*players)[i].RLFP.Position.Z+.5 > position.Z || (*players)[i].RLFP.Position.Z-.5 < position.Z {
					for _, box := range bounding_boxes {
						if !rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(position.X-(*players)[i].RLFP.Scale.X/2, position.Y-(*players)[i].RLFP.Scale.Y/2, position.Z-(*players)[i].RLFP.Scale.Z/2),
							rl.NewVector3(position.X+(*players)[i].RLFP.Scale.X/2, position.Y+(*players)[i].RLFP.Scale.Y/2, position.Z+(*players)[i].RLFP.Scale.Z/2)), box) {
							(*players)[i].RLFP.Position.X = position.X
							(*players)[i].RLFP.Position.Y = position.Y
							(*players)[i].RLFP.Position.Z = position.Z
							(*players)[i].RLFP.UpdateCameraFirstPerson()
						}
					}
				}
				break
			}
		}
	})
}

func player_updates(server *tcp.Server, players *[]Player, bounding_boxes []rl.BoundingBox, trigger_boxes []TriggerBox, interractable_boxes []InteractableBox) {
	for i := range *players {
		(*players)[i].RLFP.UpdatePlayer(bounding_boxes, trigger_boxes, interractable_boxes)

		if (*players)[i].LastPosition != (*players)[i].RLFP.Position {
			to_send, err := json.Marshal(PlayerPositionToSend{
				Name: (*players)[i].Name,
				X:    (*players)[i].RLFP.Position.X,
				Y:    (*players)[i].RLFP.Position.Y,
				Z:    (*players)[i].RLFP.Position.Z,
			})
			if err != nil {
				logger.Log(lgr.Error, "Error marshalling position data: %s", err)
				continue
			}
			server.SendDataToAll("position", to_send)
		}

		if (*players)[i].LastRotation != (*players)[i].RLFP.Rotation {
			to_send, err := json.Marshal(PlayerRotationToSend{
				Name: (*players)[i].Name,
				X:    (*players)[i].RLFP.Rotation.X,
				Y:    (*players)[i].RLFP.Rotation.Y,
			})
			if err != nil {
				logger.Log(lgr.Error, "Error marshalling rotation data: %s", err)
				continue
			}
			server.SendDataToAll("rotation", to_send)
		}

		(*players)[i].LastPosition = (*players)[i].RLFP.Position
		(*players)[i].LastRotation = (*players)[i].RLFP.Rotation
	}
}
