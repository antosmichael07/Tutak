package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Token        string
	Name         string
	RLFP         PlayerFP
	LastPosition rl.Vector3
	LastRotation rl.Vector2
}

type PlayerPositionToSend struct {
	Name string
	X    float32
	Y    float32
	Z    float32
}

type PlayerRotationToSend struct {
	Name string
	X    float32
	Y    float32
}
