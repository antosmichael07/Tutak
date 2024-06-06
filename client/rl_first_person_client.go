package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	lgr "github.com/antosmichael07/Go-Logger"
	tcp "github.com/antosmichael07/Go-TCP-Connection"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerFP struct {
	Speed             Speeds
	Mouse_sensitivity Sensitivities
	Fovs              FOVs
	Position          rl.Vector3
	Rotation          rl.Vector2
	Scale             rl.Vector3
	ConstScale        Scale
	IsCrouching       bool
	YVelocity         float32
	Gravity           float32
	JumpPower         float32
	LastKeyPressed    string
	FrameTime         float32
	InteractRange     float32
	AlreadyInteracted bool
	StepHeight        float32
	Stepped           bool
	Controls          Controls
	Keybinds          Keybinds
	Camera            rl.Camera3D
}

type Speeds struct {
	Normal       float32
	Sprint       float32
	Sneak        float32
	Current      float32
	Acceleration float32
}

type Sensitivities struct {
	Normal float32
	Zoom   float32
}

type FOVs struct {
	Normal float32
	Zoom   float32
}

type Scale struct {
	Normal float32
	Crouch float32
}

type Vector2XZ struct {
	X float32
	Z float32
}

type Controls struct {
	Forward  bool
	Backward bool
	Left     bool
	Right    bool
	Jump     bool
	Crouch   bool
	Sprint   bool
	Interact bool
}

type Keybinds struct {
	Forward  int32
	Backward int32
	Left     int32
	Right    int32
	Jump     int32
	Crouch   int32
	Sprint   int32
	Zoom     int32
	Interact int32
}

type TriggerBox struct {
	BoundingBox rl.BoundingBox
	Triggered   bool
	Triggering  bool
}

type InteractableBox struct {
	BoundingBox  rl.BoundingBox
	Interacted   bool
	Interacting  bool
	RayCollision rl.RayCollision
}

func (player *PlayerFP) UpdatePlayer(bounding_boxes []rl.BoundingBox, trigger_boxes []TriggerBox, interactable_boxes []InteractableBox, client tcp.Client) {
	player.DrawInteractIndicator(interactable_boxes)
	player.RotatePlayer(client)
	player.updateCameraFirstPerson()
}

func (player *PlayerFP) RotatePlayer(client tcp.Client) {
	last_rotationX := player.Rotation.X
	last_rotationY := player.Rotation.Y
	rotationX := player.Rotation.X
	rotationY := player.Rotation.Y
	mouse_delta := rl.GetMouseDelta()
	if rl.IsKeyDown(player.Keybinds.Zoom) {
		rotationX += mouse_delta.X * player.Mouse_sensitivity.Zoom
		rotationY -= mouse_delta.Y * player.Mouse_sensitivity.Zoom
	} else {
		rotationX += mouse_delta.X * player.Mouse_sensitivity.Normal
		rotationY -= mouse_delta.Y * player.Mouse_sensitivity.Normal
	}
	if rotationY > 1.5 {
		rotationY = 1.5
	}
	if rotationY < -1.5 {
		rotationY = -1.5
	}
	to_send, err := json.Marshal(rl.NewVector2(rotationX, rotationY))
	if err != nil {
		logger.Log(lgr.Error, "Error marshalling rotation data: %s", err)
	}
	if last_rotationX != rotationX || last_rotationY != rotationY {
		client.SendData("update-rotation", to_send)
	}
}

func (player PlayerFP) DrawInteractIndicator(interactable_boxes []InteractableBox) {
	for i := range interactable_boxes {
		if interactable_boxes[i].Interacting {
			return
		}
	}
	text := fmt.Sprintf("Press %s to interact", strings.ToUpper(string(player.Keybinds.Interact)))
	text_size := rl.MeasureText(text, 30)
	for i := range interactable_boxes {
		if interactable_boxes[i].RayCollision.Hit && interactable_boxes[i].RayCollision.Distance <= player.InteractRange {
			rl.DrawText(text, int32(rl.GetScreenWidth()/2)-text_size/2, int32(rl.GetScreenHeight()/2)-30, 30, rl.White)
		}
	}
}

func (player *PlayerFP) updateCameraFirstPerson() {
	player.MoveCamera()
	player.RotateCamera()
	player.ZoomCamera()
}

func (player *PlayerFP) MoveCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
}

func (player *PlayerFP) RotateCamera() {
	cos_rotation_y := float32(math.Cos(float64(player.Rotation.Y)))

	player.Camera.Target.X = player.Camera.Position.X - float32(math.Cos(float64(player.Rotation.X)))*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + float32(math.Sin(float64(player.Rotation.Y)))
	player.Camera.Target.Z = player.Camera.Position.Z - float32(math.Sin(float64(player.Rotation.X)))*cos_rotation_y
}

func (player *PlayerFP) ZoomCamera() {
	if rl.IsKeyDown(player.Keybinds.Zoom) {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
