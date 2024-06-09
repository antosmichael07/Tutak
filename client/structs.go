package main

import (
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Name            string
	RLFP            rlfp.Player
	LastRotation    rl.Vector2
	LastKeysPressed Input
}

type Input struct {
	Forward  bool
	Backward bool
	Left     bool
	Right    bool
	Jump     bool
	Crouch   bool
	Sprint   bool
	Interact bool
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
