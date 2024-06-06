package main

import (
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerFP struct {
	Speed             Speeds
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
	Timer             float32
	Camera            rl.Camera3D
}

type Speeds struct {
	Normal       float32
	Sprint       float32
	Sneak        float32
	Current      float32
	Acceleration float32
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

func (player *PlayerFP) InitPlayer() {
	player.Speed.Normal = .1
	player.Speed.Sprint = .15
	player.Speed.Sneak = .05
	player.Speed.Current = 0.
	player.Speed.Acceleration = .01
	player.Fovs.Normal = 70.
	player.Fovs.Zoom = 20.
	player.Rotation = rl.NewVector2(0., 0.)
	player.Position = rl.NewVector3(4., .9, 4.)
	player.Scale = rl.NewVector3(.8, 1.8, .8)
	player.ConstScale.Normal = 1.8
	player.ConstScale.Crouch = .9
	player.IsCrouching = false
	player.YVelocity = 0.
	player.Gravity = .0065
	player.JumpPower = .15
	player.LastKeyPressed = ""
	player.FrameTime = 0.
	player.InteractRange = 3.
	player.AlreadyInteracted = false
	player.StepHeight = .5
	player.Stepped = false
	player.Controls.Forward = false
	player.Controls.Backward = false
	player.Controls.Left = false
	player.Controls.Right = false
	player.Controls.Jump = false
	player.Controls.Crouch = false
	player.Controls.Sprint = false
	player.Controls.Interact = false
	player.InitCamera()
}

func (player *PlayerFP) UpdatePlayer(bounding_boxes []rl.BoundingBox, trigger_boxes []TriggerBox, interactable_boxes []InteractableBox) {
	player.FrameTime = float32(time.Now().UnixMicro())*1000000 - player.Timer
	player.Timer = float32(time.Now().UnixMicro()) * 1000000
	player.LastKeyPressedPlayer()
	player.AccelerationPlayer()
	player.ApplyGravityToPlayer(bounding_boxes)
	if player.Speed.Acceleration != 0. {
		player.StepPlayer(bounding_boxes)
		player.MovePlayer(bounding_boxes)
		player.CheckTriggerBoxes(trigger_boxes)
	}
	player.UpdateInteractableBoxes(interactable_boxes)
}

func (player *PlayerFP) LastKeyPressedPlayer() {
	if player.Controls.Forward {
		player.LastKeyPressed = "w"
	}
	if player.Controls.Backward {
		player.LastKeyPressed = "s"
	}
	if player.Controls.Left {
		player.LastKeyPressed = "a"
	}
	if player.Controls.Right {
		player.LastKeyPressed = "d"
	}
}

func (player *PlayerFP) AccelerationPlayer() {
	final_speed := player.Speed.Acceleration * player.FrameTime * 60

	if !player.Controls.Forward && !player.Controls.Backward && !player.Controls.Left && !player.Controls.Right {
		if player.Speed.Current > 0. {
			player.Speed.Current -= final_speed
		} else {
			player.Speed.Current = 0.
		}
	} else if !player.Controls.Sprint && player.Speed.Current > player.Speed.Normal {
		player.Speed.Current -= final_speed
	}
	if player.IsCrouching && player.Speed.Current > player.Speed.Sneak {
		player.Speed.Current -= final_speed
	}

	if player.Speed.Current <= player.Speed.Normal && (player.Controls.Forward || player.Controls.Backward || player.Controls.Left || player.Controls.Right) && !player.Controls.Sprint && !player.Controls.Crouch {
		player.Speed.Current += final_speed
	}
	if player.Controls.Sprint && player.Speed.Current <= player.Speed.Sprint && (player.Controls.Forward || player.Controls.Backward || player.Controls.Left || player.Controls.Right) {
		player.Speed.Current += final_speed
	}
	if player.Controls.Crouch && player.Speed.Current <= player.Speed.Sneak && (player.Controls.Forward || player.Controls.Backward || player.Controls.Left || player.Controls.Right) {
		player.Speed.Current += final_speed
	}
}

func (player *PlayerFP) StepPlayer(bounding_boxes []rl.BoundingBox) {
	player.Stepped = false
	player_tmp := *player
	player_tmp.Position.Y += player.StepHeight + 0.0001
	if !player_tmp.CheckCollisionsYForPlayer(bounding_boxes) && !player_tmp.CheckCollisionsForPlayer(bounding_boxes) && player.CheckCollisionsForPlayer(bounding_boxes) && player.YVelocity == 0. {
		player_position_after_moving := player.GetPlayerPositionAfterMoving()
		player.Position.Y = (player.CheckCollisionsForPlayerAsHighestPoint(bounding_boxes) + player.Scale.Y/2) + 0.0001
		player.Position.X = player_position_after_moving.X
		player.Position.Z = player_position_after_moving.Z
		player.Stepped = true
		return
	}
	collision_x, collision_z := player.CheckCollisionsXZForPlayerWithY(bounding_boxes)
	tmp_collision_x, tmp_collision_z := player_tmp.CheckCollisionsXZForPlayerWithY(bounding_boxes)
	if !player_tmp.CheckCollisionsYForPlayer(bounding_boxes) && !tmp_collision_x && collision_x && player.YVelocity == 0. {
		player.Position.Y = (player.CheckCollisionsXForPlayerAsHighestPoint(bounding_boxes) + player.Scale.Y/2) + 0.0001
		player.Position.X = player.GetPlayerPositionAfterMoving().X
		player.Stepped = true
		return
	}
	if !player_tmp.CheckCollisionsYForPlayer(bounding_boxes) && !tmp_collision_z && collision_z && player.YVelocity == 0. {
		player.Position.Y = (player.CheckCollisionsZForPlayerAsHighestPoint(bounding_boxes) + player.Scale.Y/2) + 0.0001
		player.Position.Z = player.GetPlayerPositionAfterMoving().Z
		player.Stepped = true
		return
	}
}

func (player PlayerFP) CheckCollisionsForPlayerAsHighestPoint(bounding_boxes []rl.BoundingBox) float32 {
	player.Position = player.GetPlayerPositionAfterMoving()

	highest_y := float32(0.)
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			if box.Max.Y > highest_y {
				if box.Min.Y <= box.Max.Y {
					highest_y = box.Max.Y
				} else {
					highest_y = box.Min.Y
				}
			}
		}
	}

	return highest_y
}

func (player PlayerFP) CheckCollisionsXForPlayerAsHighestPoint(bounding_boxes []rl.BoundingBox) float32 {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()
	player.Position.X = player_position_after_moving.X
	player.Position.Y = player_position_after_moving.Y

	highest_y := float32(0.)
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			if box.Max.Y > highest_y {
				if box.Min.Y <= box.Max.Y {
					highest_y = box.Max.Y
				} else {
					highest_y = box.Min.Y
				}
			}
		}
	}

	return highest_y
}

func (player PlayerFP) CheckCollisionsZForPlayerAsHighestPoint(bounding_boxes []rl.BoundingBox) float32 {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()
	player.Position.Z = player_position_after_moving.Z
	player.Position.Y = player_position_after_moving.Y

	highest_y := float32(0.)
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			if box.Max.Y > highest_y {
				if box.Min.Y <= box.Max.Y {
					highest_y = box.Max.Y
				} else {
					highest_y = box.Min.Y
				}
			}
		}
	}

	return highest_y
}

func (player *PlayerFP) MovePlayer(bounding_boxes []rl.BoundingBox) {
	half_crouch_scale := player.ConstScale.Crouch / 2

	if player.Controls.Crouch {
		player.Scale.Y = player.ConstScale.Crouch
		if !player.IsCrouching {
			player.Position.Y -= half_crouch_scale
		}
		player.IsCrouching = true
	} else if player.CheckPlayerUncrouch(bounding_boxes) {
		player.Scale.Y = player.ConstScale.Normal
		if player.IsCrouching {
			player.Position.Y += half_crouch_scale
		}
		player.IsCrouching = false
	}
	if player.Controls.Jump && player.YVelocity == 0. && player.CheckIfPlayerOnSurface(bounding_boxes) && !player.IsCrouching {
		player.YVelocity = player.JumpPower
	}

	player.Position.Y += player.YVelocity * (player.FrameTime * 60)

	if !player.Stepped {
		player_position_after_moving := player.GetPlayerPositionAfterMoving()

		collisions_x, collisions_z := player.CheckCollisionsXZForPlayer(bounding_boxes)
		if collisions_x && collisions_z {
			return
		} else if collisions_x {
			player.Position.Z = player_position_after_moving.Z
			return
		} else if collisions_z {
			player.Position.X = player_position_after_moving.X
			return
		}

		if player.CheckCollisionsForPlayer(bounding_boxes) {
			player.Position.X = player_position_after_moving.X
			return
		}

		player.Position.X = player_position_after_moving.X
		player.Position.Z = player_position_after_moving.Z
	}
}

func (player *PlayerFP) ApplyGravityToPlayer(bounding_boxes []rl.BoundingBox) {
	frame_time := player.FrameTime * 60

	player.YVelocity -= player.Gravity * frame_time

	player_y_after_falling := player.Position.Y + player.YVelocity*frame_time
	if player.CheckCollisionsYForPlayer(bounding_boxes) || player_y_after_falling-(player.Scale.Y/2) < 0. {
		player.YVelocity = 0.
		return
	}
}

func (player PlayerFP) CheckCollisionsForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.Position = player.GetPlayerPositionAfterMoving()

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			return true
		}
	}

	return false
}

func (player PlayerFP) CheckCollisionForPlayer(bounding_box rl.BoundingBox) bool {
	return rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
		rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), bounding_box)
}

func (player PlayerFP) CheckCollisionsYForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.Speed.Normal = 0
	player.Speed.Sprint = 0
	player.Speed.Sneak = 0
	player.Speed.Current = 0

	return player.CheckCollisionsForPlayer(bounding_boxes)
}

func (player PlayerFP) CheckCollisionsXZForPlayer(bounding_boxes []rl.BoundingBox) (bool, bool) {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()

	collision_x, collision_z := false, false

	player_position_x := player_position_after_moving.X
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_x+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			collision_x = true
		}
	}

	player_position_z := player_position_after_moving.Z
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player_position_z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player_position_z+player.Scale.Z/2)), box) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player PlayerFP) CheckCollisionsXZForPlayerWithY(bounding_boxes []rl.BoundingBox) (bool, bool) {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()

	collision_x, collision_z := false, false

	player_position_x := player_position_after_moving.X
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-player.Scale.X/2, player_position_after_moving.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_x+player.Scale.X/2, player_position_after_moving.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			collision_x = true
		}
	}

	player_position_z := player_position_after_moving.Z
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player_position_after_moving.Y-player.Scale.Y/2, player_position_z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player_position_after_moving.Y+player.Scale.Y/2, player_position_z+player.Scale.Z/2)), box) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player PlayerFP) GetPlayerPositionAfterMoving() rl.Vector3 {
	frame_time := player.FrameTime * 60

	current_speed := player.Speed.Current

	if player.Speed.Normal == 0. {
		player.Position.Y += player.YVelocity * frame_time
	}

	keys_pressed := 0
	if player.Controls.Forward {
		keys_pressed++
	}
	if player.Controls.Backward {
		keys_pressed++
	}
	if player.Controls.Left {
		keys_pressed++
	}
	if player.Controls.Right {
		keys_pressed++
	}
	if keys_pressed == 2 {
		current_speed = current_speed * .707
	}

	final_speed := current_speed * frame_time

	speeds := Vector2XZ{
		float32(math.Cos(float64(player.Rotation.X))) * final_speed,
		float32(math.Sin(float64(player.Rotation.X))) * final_speed,
	}

	if player.Controls.Forward || player.LastKeyPressed == "w" {
		player.Position.X -= speeds.X
		player.Position.Z -= speeds.Z
	}
	if player.Controls.Backward || player.LastKeyPressed == "s" {
		player.Position.X += speeds.X
		player.Position.Z += speeds.Z
	}
	if player.Controls.Left || player.LastKeyPressed == "a" {
		player.Position.Z += speeds.X
		player.Position.X -= speeds.Z
	}
	if player.Controls.Right || player.LastKeyPressed == "d" {
		player.Position.Z -= speeds.X
		player.Position.X += speeds.Z
	}

	return player.Position
}

func (player PlayerFP) CheckPlayerUncrouch(bounding_boxes []rl.BoundingBox) bool {
	player.Scale.Y = player.ConstScale.Normal
	player.Position.Y += player.ConstScale.Normal / 2

	return !player.CheckCollisionsForPlayer(bounding_boxes)
}

func (player PlayerFP) CheckIfPlayerOnSurface(bounding_boxes []rl.BoundingBox) bool {
	player.Position.Y -= player.Gravity * (player.FrameTime * 60)
	if player.CheckCollisionsYForPlayer(bounding_boxes) || player.Position.Y-(player.Scale.Y/2) < 0. {
		return true
	}
	return false
}

func (player PlayerFP) CheckTriggerBoxes(trigger_boxes []TriggerBox) {
	for i := range trigger_boxes {
		if !trigger_boxes[i].Triggering {
			trigger_boxes[i].Triggered = player.CheckCollisionForPlayer(trigger_boxes[i].BoundingBox)
		} else {
			trigger_boxes[i].Triggered = false
		}
		trigger_boxes[i].Triggering = player.CheckCollisionForPlayer(trigger_boxes[i].BoundingBox)
	}
}

func NewTriggerBox(box rl.BoundingBox) TriggerBox {
	return TriggerBox{box, false, false}
}

func (player *PlayerFP) UpdateInteractableBoxes(interactable_boxes []InteractableBox) {
	for i := range interactable_boxes {
		interactable_boxes[i].RayCollision = rl.GetRayCollisionBox(rl.GetMouseRay(rl.NewVector2(float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))/2, float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))/2), player.Camera), interactable_boxes[i].BoundingBox)
	}
	player.CheckInteractableBoxes(interactable_boxes)
}

func (player *PlayerFP) CheckInteractableBoxes(interactable_boxes []InteractableBox) {
	for i := range interactable_boxes {
		if player.AlreadyInteracted {
			interactable_boxes[i].Interacted = false
		}
		if player.Controls.Interact && (!player.AlreadyInteracted || interactable_boxes[i].RayCollision.Distance > player.InteractRange) {
			if interactable_boxes[i].RayCollision.Hit && interactable_boxes[i].RayCollision.Distance <= player.InteractRange {
				player.AlreadyInteracted = true
				if !interactable_boxes[i].Interacting {
					interactable_boxes[i].Interacted = true
				} else {
					interactable_boxes[i].Interacted = false
				}
				interactable_boxes[i].Interacting = true
			} else {
				interactable_boxes[i].Interacting = false
				interactable_boxes[i].Interacted = false
			}
		} else if !player.Controls.Interact {
			interactable_boxes[i].Interacting = false
			interactable_boxes[i].Interacted = false
			player.AlreadyInteracted = false
		}
	}
	if player.Controls.Interact {
		player.AlreadyInteracted = true
	} else {
		player.AlreadyInteracted = false
	}
}

func NewInteractableBox(box rl.BoundingBox) InteractableBox {
	return InteractableBox{box, false, false, rl.NewRayCollision(false, 0., rl.NewVector3(0., 0., 0.), rl.NewVector3(0., 0., 0.))}
}

func (player *PlayerFP) InitCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
	player.Camera.Target = rl.NewVector3(
		player.Camera.Position.X-float32(math.Cos(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
		player.Camera.Position.Y+float32(math.Sin(float64(player.Rotation.Y)))+(player.Scale.Y/2),
		player.Camera.Position.Z-float32(math.Sin(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
	)
	player.Camera.Up = rl.NewVector3(0., 1., 0.)
	player.Camera.Fovy = player.Fovs.Normal
	player.Camera.Projection = rl.CameraPerspective
}
