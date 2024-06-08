package main

import (
	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
)

type Player struct {
	Name string
	RLFP rlfp.Player
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

type PlayerPosition struct {
	Name string
	X    float32
	Y    float32
	Z    float32
}