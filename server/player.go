package main

import (
	"encoding/json"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init_player_events(server *tcp.Server, players *[]Player, collision_boxes []rl.BoundingBox, trigger_boxes []TriggerBox, interactable_boxes []InteractableBox) {
	server.On("new-player", func(data []byte, conn tcp.Connection) {
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

	player_updates(server, players, collision_boxes, trigger_boxes, interactable_boxes)

	server.OnDisconnect(func(conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				logger.Log(lgr.Info, "Player disconnected: %s", (*players)[i].Name)
				*players = append((*players)[:i], (*players)[i+1:]...)
				break
			}
		}
	})
}

func player_updates(server *tcp.Server, players *[]Player, collision_boxes []rl.BoundingBox, trigger_boxes []TriggerBox, interactable_boxes []InteractableBox) {
	server.On("input", func(data []byte, conn tcp.Connection) {
		for i := range *players {
			if (*players)[i].Token == conn.Token {
				err := json.Unmarshal(data, &(*players)[i].RLFP.Controls)
				if err != nil {
					logger.Log(lgr.Error, "Error unmarshalling input data: %s", err)
					return
				}

				(*players)[i].RLFP.UpdatePlayer(collision_boxes, trigger_boxes, interactable_boxes)

				position := PlayerPositionToSend{
					Name: (*players)[i].Name,
					X:    (*players)[i].RLFP.Position.X,
					Y:    (*players)[i].RLFP.Position.Y,
					Z:    (*players)[i].RLFP.Position.Z,
				}
				to_send, err := json.Marshal(position)
				if err != nil {
					logger.Log(lgr.Error, "Error marshalling position data: %s", err)
					return
				}
				server.SendDataToAll("position", to_send)
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

				(*players)[i].RLFP.UpdateInteractableBoxes(interactable_boxes)
				(*players)[i].RLFP.UpdateCameraFirstPerson()

				rotation := PlayerRotationToSend{
					Name: (*players)[i].Name,
					X:    (*players)[i].RLFP.Rotation.X,
					Y:    (*players)[i].RLFP.Rotation.Y,
				}
				to_send, err := json.Marshal(rotation)
				if err != nil {
					logger.Log(lgr.Error, "Error marshalling rotation data: %s", err)
					return
				}
				server.SendDataToAll("rotation", to_send)
				break
			}
		}
	})
}
