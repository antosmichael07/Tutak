package main

import (
	"encoding/json"
	"net"
	"time"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var logger = lgr.NewLogger("Tutak")

type Player struct {
	Name string
	RLFP PlayerFP
}

type PlayerPackage struct {
	Name string
	RFLP []byte
}

func main() {
	rl.SetTraceLogLevel(rl.LogError)

	server := tcp.NewServer("localhost:8080")

	players := []Player{{Name: "Test", RLFP: PlayerFP{}}}
	players[0].RLFP.InitPlayer()

	server.On("initialize_player", func(data []byte, conn net.Conn) {
		/*p := Player{RLFP: PlayerFP{}}
		p.RLFP.InitPlayer()
		p.Name = string(data)
		players = append(players, p)
		logger.Log(lgr.Info, "New Player initialized")*/
	})

	server.On("update-rotation", func(data []byte, conn net.Conn) {
		rotation := rl.Vector2{}
		err := json.Unmarshal(data, &rotation)
		if err != nil {
			logger.Log(lgr.Error, "Error unmarshalling player data: %s", err)
		}
		players[0].RLFP.Rotation = rotation
	})

	server.OnConnect(func(conn net.Conn) {
		logger.Log(lgr.Info, "New connection")
	})

	server.OnDisconnect(func(conn net.Conn) {
		logger.Log(lgr.Info, "Connection closed")
	})

	go server.Start()
	server.Logger.Level = lgr.Warning

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

	for {
		for _, p := range players {
			p.RLFP.UpdatePlayer(bounding_boxes, trigger_boxes, interractable_boxes)
		}

		rlfp_data, err := json.Marshal(players[0].RLFP)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling rlfp data: %s", err)
		}
		player_package := PlayerPackage{players[0].Name, rlfp_data}
		data, err := json.Marshal(player_package)
		if err != nil {
			logger.Log(lgr.Error, "Error marshalling player data: %s", err)
		}

		server.SendDataToAll("players", []byte(data))

		time.Sleep(1 * time.Second)
	}
}
